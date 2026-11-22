package cmd

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"net"
	"time"

	"github.com/ankorstore/yokai/log"
	"github.com/go-sql-driver/mysql"
)

const (
	dbWaitInitialBackoff = 1 * time.Second
	dbWaitMaxBackoff     = 8 * time.Second
	// dbWaitBudget is the total of the sleeps between attempts.
	dbWaitBudget      = 90 * time.Second
	dbWaitPingTimeout = 5 * time.Second
)

type pingFunc func(ctx context.Context) error

type sleepFunc func(ctx context.Context, d time.Duration) error

type dbWaitPolicy struct {
	initial, max, budget, pingTimeout time.Duration
}

func defaultDBWaitPolicy() dbWaitPolicy {
	return dbWaitPolicy{
		initial:     dbWaitInitialBackoff,
		max:         dbWaitMaxBackoff,
		budget:      dbWaitBudget,
		pingTimeout: dbWaitPingTimeout,
	}
}

// needsDatabaseWait reports whether a goose command connects to the database.
func needsDatabaseWait(command string) bool {
	switch command {
	case "create", "fix", "validate":
		return false
	default:
		return true
	}
}

// waitForDatabase blocks until the database answers, the error is not a
// connection error, or the retry budget is spent. It only fails when ctx is
// done: otherwise the migration runs next and reports any remaining error.
func waitForDatabase(ctx context.Context, db *sql.DB, logger *log.Logger) error {
	return waitForDatabaseWith(ctx, db.PingContext, sleepContext, logger, defaultDBWaitPolicy())
}

func waitForDatabaseWith(
	ctx context.Context,
	ping pingFunc,
	sleep sleepFunc,
	logger *log.Logger,
	policy dbWaitPolicy,
) error {
	wait := policy.initial

	var waited time.Duration

	for attempt := 1; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := pingOnce(ctx, ping, policy.pingTimeout)
		if err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !isConnectionError(err) {
			return nil
		}

		if waited+wait > policy.budget {
			logger.Warn().Err(err).Int("attempt", attempt).
				Msg("database still unreachable, running migration to report it")

			return nil
		}

		logger.Warn().Err(err).Int("attempt", attempt).Dur("wait", wait).
			Msg("database not reachable, retrying")

		if sleepErr := sleep(ctx, wait); sleepErr != nil {
			return sleepErr
		}

		waited += wait
		wait = nextBackoff(wait, policy.max)
	}
}

func pingOnce(ctx context.Context, ping pingFunc, timeout time.Duration) error {
	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return ping(pingCtx)
}

// isConnectionError reports whether err means the server cannot be reached
// yet. Server-side errors (auth, unknown database, SQL) are not.
func isConnectionError(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) && opErr.Op == "dial" {
		return true
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, driver.ErrBadConn) ||
		errors.Is(err, mysql.ErrInvalidConn) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF)
}

func nextBackoff(current, maxBackoff time.Duration) time.Duration {
	return min(current*2, maxBackoff)
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
