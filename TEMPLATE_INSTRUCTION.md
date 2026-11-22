# New Repository Setup Instructions

You've created a new repository from the Promy template! 🎉

## Quick Setup (5 minutes)

### 1. Update Package Names

Run the scaffold script — the only supported rename path. It prompts for the package
name and GitHub org, then rewrites both the `TEMPLATE_PACKAGE_NAME` placeholder and the
template's Go module path (`github.com/tclavelloux/promy-template-go`) across the repo:

```bash
./scripts/setup-new-repo.sh
```

Files it rewrites:
- `go.mod` - module name
- `main.go`, `cmd/`, `internal/` - import paths
- `release-please-config.json`, `.github/workflows/release-please.yml` - package-name
- `README.md`, `HOWTO.md` - service name

Then edit by hand — the script cannot infer these:
- [ ] `README.md` - one-line role, Overview, Boundaries tables
- [ ] `docker-compose.yaml` - `container_name` per service
- [ ] `.testcoverage.yml` - `threshold.total`
- [ ] `.golangci-version` - local golangci-lint pin, must match the version CI resolves

Verify the rename took: `grep -rn promy-template-go . --exclude-dir=.git` must print nothing.

### 2. Enable GitHub Actions

1. Go to repository Settings → Actions → General
2. Select "Read repository contents and packages permissions" (read-only default token; each workflow requests extra scopes in its own `permissions:` block)
3. Check "Allow GitHub Actions to create and approve pull requests"
4. Save

### 3. Set Up Git Hooks

Install `pre-commit` first. `make setup` runs `pre-commit install --install-hooks` and fails
without it:

```bash
brew install pre-commit
make setup
```

`make setup` writes three hooks into `.git/hooks/` (`core.hooksPath` stays unset), as listed in
`default_install_hook_types` in [`.pre-commit-config.yaml`](.pre-commit-config.yaml):
- `pre-commit` - `no-direct-commit-to-main` (shared hook, blocks commits on `main`), formatting, lint, gitleaks, `go mod tidy`
- `commit-msg` - Conventional Commits
- `pre-push` - `check-coverage` (`make check-coverage`: race suite, then `go-test-coverage` against `.testcoverage.yml`, the same threshold CI enforces)

If `.git/hooks/pre-commit` already existed, `pre-commit install` moves it to
`pre-commit.legacy` and keeps chaining to it. Delete `pre-commit.legacy` after setup.

### 4. Initialize Git History

```bash
# Delete this instruction file
rm TEMPLATE_INSTRUCTION.md

# Initial commit — stage by explicit path, never `git add .`
git add <paths>
git commit -S -m "chore: initialize from promy-template-go

- Configure Release Please for automated releases
- Add Makefile helpers
- Configure Yokai framework structure
- Wire pre-commit hooks"

git push origin main
```

### 5. Verify Setup

- [ ] Git hooks installed (`ls .git/hooks` lists `pre-commit`, `commit-msg` and `pre-push`)
- [ ] Make commands work: `make help`
- [ ] Release Please workflow exists in Actions
- [ ] README reflects your project
- [ ] `pre-commit run --all-files` passes (lint, format, gitleaks, go-mod-tidy)
- [ ] `pre-commit run --hook-stage pre-push --all-files` passes (coverage gate)

### 6. Start Developing!

```bash
# Create your first feature — the `no-direct-commit-to-main` hook blocks commits on main
git checkout -b feature/initial-setup
```

---

Delete this file after setup is complete!
