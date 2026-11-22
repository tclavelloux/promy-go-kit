setup:
	pre-commit install --install-hooks

# Crash loop after a go.mod Go bump: "go: go.mod requires go >= 1.26.x (running go 1.26.y; GOTOOLCHAIN=local)".
# Cause: `up` never builds or pulls, so the cached golang base image stays older than go.mod.
# Fix: `make build-up-recreate` (keeps volumes, MySQL data survives).
# `fresh` does not pull either (--no-cache skips the layer cache only); `full-fresh` (-v) deletes named volumes.
up:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose up -d

# Re-pull the golang base image, rebuild, recreate containers; keeps volumes. Fixes the `up` crash loop above.
build-up-recreate:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose build --pull
	docker compose up -d --force-recreate

down:
	docker compose down

fresh:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose down --remove-orphans
	docker compose build --no-cache
	docker compose up -d --build -V

full-fresh:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi
	docker compose down --remove-orphans -v
	docker compose build --no-cache
	docker compose up -d --build -V

logs:
	docker compose logs -f

test:
	CGO_ENABLED=0 go test -v -race -cover -count=1 -failfast ./...

lint:
	scripts/golangci-lint.sh run --timeout=5m

lint-install:
	scripts/golangci-lint.sh --version

# Pinned to the version promy-github-workflows' go-coverage.yml runs in CI; bump both together.
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@v2.21.0

# Runs the race suite and gates it on .testcoverage.yml — the same thresholds CI enforces.
# Wired to the pre-push stage in .pre-commit-config.yaml. `make test` is deliberately NOT
# a separate pre-push hook: this target already runs the full suite.
check-coverage: install-go-test-coverage
	go test -race -failfast -coverpkg=./... -covermode=atomic -coverprofile=coverage.out -count=1 ./...
	$$(go env GOPATH)/bin/go-test-coverage --config=./.testcoverage.yml

migrate:
	docker-compose up migration-runner

migrate-status:
	docker-compose exec migration-runner go run main.go migrate-database status

migrate-down:
	docker-compose exec migration-runner go run main.go migrate-database down

migrate-reset:
	docker-compose exec migration-runner go run main.go migrate-database reset

# Git workflow helpers
git-status:
	@echo "=== Git Status with Branch Info ==="
	@git status -sb
	@echo "\n=== Uncommitted Files by Domain ==="
	@git status --short | grep -E "^ M|^M |^A |^ A" | awk '{print $$2}' | grep "internal/domain/" | sed 's|internal/domain/\([^/]*\)/.*|\1|' | sort -u | sed 's/^/  - /'

git-log:
	@git log --oneline --graph --decorate -20

git-diff:
	@echo "=== Staged Changes ==="
	@git diff --cached --stat
	@echo "\n=== Unstaged Changes ==="
	@git diff --stat

git-check:
	@echo "=== Pre-Commit Checklist ==="
	@echo "[ ] One domain/layer per commit?"
	@echo "[ ] Message follows Conventional Commits?"
	@echo "[ ] Files related to same feature?"
	@echo "[ ] < 10 files, < 500 lines?"
	@echo ""
	@echo "=== Running Tests ==="
	@go test ./... -short
	@echo ""
	@echo "=== Running Linter ==="
	@scripts/golangci-lint.sh run --timeout=5m

# Same pinned govulncheck + Go toolchain as CI's `vuln` job; blocks on any called vuln.
vuln:
	pre-commit run go-vuln --hook-stage manual
