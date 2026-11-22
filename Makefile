.PHONY: setup
setup:
	pre-commit install --install-hooks

.PHONY: test
test:
	go test -v -race -cover -count=1 ./...

.PHONY: lint
lint:
	scripts/golangci-lint.sh run --timeout=5m

.PHONY: lint-install
lint-install:
	scripts/golangci-lint.sh --version

.PHONY: coverage
coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Pinned to the version promy-github-workflows' go-coverage.yml runs in CI; bump both together.
.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@v2.23.0

# Runs the race suite and gates it on .testcoverage.yml, the same thresholds CI enforces.
# Wired to the pre-push stage in .pre-commit-config.yaml.
.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test -race -failfast -coverpkg=./... -covermode=atomic -coverprofile=coverage.out -count=1 ./...
	$$(go env GOPATH)/bin/go-test-coverage --config=./.testcoverage.yml

# Same pinned govulncheck + Go toolchain as CI's `vuln` job; blocks on any called vuln.
.PHONY: vuln
vuln:
	pre-commit run go-vuln --hook-stage manual
