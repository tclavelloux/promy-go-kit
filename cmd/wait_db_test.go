package cmd

import (
	"context"
	"database/sql/driver"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type timeoutError struct{}

func (timeoutError) Error() string   { return "i/o timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

var errAccessDenied = errors.New("Error 1045 (28000): Access denied for user")

// refusedError returns a genuine ECONNREFUSED dial error on a closed loopback port.
func refusedError(t *testing.T) error {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	addr := ln.Addr().String()
	require.NoError(t, ln.Close())

	conn, err := net.Dial("tcp", addr)
	if err == nil {
		_ = conn.Close()
		t.Fatal("expected the dial to be refused")
	}

	return err
}

func newTestLogger(t *testing.T) (*log.Logger, logtest.TestLogBuffer) {
	t.Helper()

	buffer := logtest.NewDefaultTestLogBuffer()
	logger, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(buffer))
	require.NoError(t, err)

	return logger, buffer
}

// scriptedPing returns the given errors in order, then nil. The last error
// repeats forever when forever is true.
func scriptedPing(errs []error, forever bool) (pingFunc, *int) {
	calls := 0

	return func(context.Context) error {
		i := calls
		calls++

		switch {
		case i < len(errs):
			return errs[i]
		case forever && len(errs) > 0:
			return errs[len(errs)-1]
		default:
			return nil
		}
	}, &calls
}

func recordingSleep() (sleepFunc, *[]time.Duration) {
	sleeps := &[]time.Duration{}

	return func(_ context.Context, d time.Duration) error {
		*sleeps = append(*sleeps, d)

		return nil
	}, sleeps
}

func sum(ds []time.Duration) time.Duration {
	var total time.Duration

	for _, d := range ds {
		total += d
	}

	return total
}

func TestWaitForDatabase_RetriesRefusedThenSucceeds(t *testing.T) {
	t.Parallel()

	logger, buffer := newTestLogger(t)
	refused := refusedError(t)
	ping, calls := scriptedPing([]error{refused, refused, refused}, false)
	sleep, sleeps := recordingSleep()

	err := waitForDatabaseWith(context.Background(), ping, sleep, logger, defaultDBWaitPolicy())

	require.NoError(t, err)
	assert.Equal(t, 4, *calls)
	assert.Equal(t, []time.Duration{time.Second, 2 * time.Second, 4 * time.Second}, *sleeps)

	for attempt := 1; attempt <= 3; attempt++ {
		logtest.AssertContainLogRecord(t, buffer, map[string]interface{}{
			"level":   "warn",
			"message": "database not reachable, retrying",
			"attempt": attempt,
		})
	}

	logtest.AssertContainNotLogRecord(t, buffer, map[string]interface{}{
		"message": "database still unreachable, running migration to report it",
	})
}

func TestWaitForDatabase_GivesUpAfterBudget(t *testing.T) {
	t.Parallel()

	logger, buffer := newTestLogger(t)
	ping, calls := scriptedPing([]error{refusedError(t)}, true)
	sleep, sleeps := recordingSleep()

	err := waitForDatabaseWith(context.Background(), ping, sleep, logger, defaultDBWaitPolicy())

	require.NoError(t, err)
	assert.Equal(t, 14, *calls)
	assert.Equal(t, 87*time.Second, sum(*sleeps))
	assert.LessOrEqual(t, sum(*sleeps), dbWaitBudget)
	logtest.AssertContainLogRecord(t, buffer, map[string]interface{}{
		"level":   "warn",
		"message": "database still unreachable, running migration to report it",
		"attempt": 14,
	})
	logtest.AssertContainLogRecord(t, buffer, map[string]interface{}{
		"message": "database not reachable, retrying",
		"attempt": 13,
	})
}

func TestWaitForDatabase_DoesNotRetryNonConnectionError(t *testing.T) {
	t.Parallel()

	logger, buffer := newTestLogger(t)
	ping, calls := scriptedPing([]error{errAccessDenied}, false)
	sleep, sleeps := recordingSleep()

	err := waitForDatabaseWith(context.Background(), ping, sleep, logger, defaultDBWaitPolicy())

	require.NoError(t, err)
	assert.Equal(t, 1, *calls)
	assert.Empty(t, *sleeps)
	logtest.AssertContainNotLogRecord(t, buffer, map[string]interface{}{
		"message": "database not reachable, retrying",
	})
}

func TestWaitForDatabase_RetriesHandshakeError(t *testing.T) {
	t.Parallel()

	logger, _ := newTestLogger(t)
	ping, calls := scriptedPing([]error{mysql.ErrInvalidConn}, false)
	sleep, sleeps := recordingSleep()

	err := waitForDatabaseWith(context.Background(), ping, sleep, logger, defaultDBWaitPolicy())

	require.NoError(t, err)
	assert.Equal(t, 2, *calls)
	assert.Equal(t, []time.Duration{time.Second}, *sleeps)
}

func TestWaitForDatabase_ParentContextCancelled(t *testing.T) {
	t.Parallel()

	logger, _ := newTestLogger(t)
	ping, calls := scriptedPing([]error{refusedError(t)}, true)
	sleep, sleeps := recordingSleep()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := waitForDatabaseWith(ctx, ping, sleep, logger, defaultDBWaitPolicy())

	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 0, *calls)
	assert.Empty(t, *sleeps)
}

func TestWaitForDatabase_SleepInterrupted(t *testing.T) {
	t.Parallel()

	logger, _ := newTestLogger(t)
	ping, _ := scriptedPing([]error{refusedError(t)}, true)

	sleep := func(context.Context, time.Duration) error { return context.Canceled }

	err := waitForDatabaseWith(context.Background(), ping, sleep, logger, defaultDBWaitPolicy())

	require.ErrorIs(t, err, context.Canceled)
}

func TestSleepContext(t *testing.T) {
	t.Parallel()

	require.NoError(t, sleepContext(context.Background(), time.Millisecond))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.ErrorIs(t, sleepContext(ctx, time.Minute), context.Canceled)
}

func TestIsConnectionError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"dial op error", &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}, true},
		{"refused on loopback", refusedError(t), true},
		{"dns error", &net.DNSError{Err: "no such host", Name: "db"}, true},
		{"timeout net error", timeoutError{}, true},
		{"deadline exceeded", context.DeadlineExceeded, true},
		{"bad conn", driver.ErrBadConn, true},
		{"mysql invalid conn", mysql.ErrInvalidConn, true},
		{"eof", io.EOF, true},
		{"unexpected eof", io.ErrUnexpectedEOF, true},
		{"wrapped eof", errors.Join(errors.New("handshake"), io.EOF), true},
		{"access denied", errAccessDenied, false},
		{"mysql server error", &mysql.MySQLError{Number: 1049, Message: "Unknown database"}, false},
		{"context canceled", context.Canceled, false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isConnectionError(tt.err))
		})
	}
}

func TestNextBackoff(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 2*time.Second, nextBackoff(time.Second, dbWaitMaxBackoff))
	assert.Equal(t, 8*time.Second, nextBackoff(4*time.Second, dbWaitMaxBackoff))
	assert.Equal(t, 8*time.Second, nextBackoff(8*time.Second, dbWaitMaxBackoff))
}

func TestNeedsDatabaseWait(t *testing.T) {
	t.Parallel()

	for command, want := range map[string]bool{
		"create":   false,
		"fix":      false,
		"validate": false,
		"up":       true,
		"status":   true,
		"down":     true,
		"reset":    true,
	} {
		assert.Equal(t, want, needsDatabaseWait(command), command)
	}
}

func TestMigrateOptions(t *testing.T) {
	t.Parallel()

	assert.Len(t, migrateOptions([]string{"up"}), len(migrateOptions([]string{"create"}))+1)
}
