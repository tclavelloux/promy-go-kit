#!/bin/bash
# scripts/setup-new-repo.sh
# Run this script after creating a new repo from the template

set -e

echo "🚀 Promy Repository Setup"
echo ""

# `make setup` runs `pre-commit install --install-hooks`, so pre-commit must already be on PATH.
if ! command -v pre-commit >/dev/null 2>&1; then
  echo "❌ pre-commit is not installed — 'make setup' needs it."
  echo "   Install it with: brew install pre-commit"
  exit 1
fi
echo "✅ pre-commit found ($(pre-commit --version))"

# The template's own identity. TEMPLATE_REPO is the bare repo name used as
# release-please's package-name; TEMPLATE_MODULE is the import path every .go file carries.
TEMPLATE_REPO="promy-template-go"
TEMPLATE_MODULE="github.com/tclavelloux/$TEMPLATE_REPO"

# Prompt for package name
read -p "Enter new package name (e.g., promy-notification): " PACKAGE_NAME

if [ -z "$PACKAGE_NAME" ]; then
  echo "❌ Package name cannot be empty"
  exit 1
fi

read -p "Enter GitHub org or user [tclavelloux]: " GITHUB_ORG
GITHUB_ORG="${GITHUB_ORG:-tclavelloux}"

NEW_MODULE="github.com/$GITHUB_ORG/$PACKAGE_NAME"

echo ""
echo "📝 Updating package name from TEMPLATE_PACKAGE_NAME to $PACKAGE_NAME..."
echo "📝 Updating module path from $TEMPLATE_MODULE to $NEW_MODULE..."

# Replace in all files. Three substitutions, applied in this order per line:
#   1. the full module path — a real, already-resolving import in every .go file and go.mod
#   2. the bare repo name — release-please's package-name, which carries no github.com prefix
#   3. TEMPLATE_PACKAGE_NAME — the placeholder left in the docs
# Order matters: the module path must go first, or step 2 would shred it.
find . -type f ! -path "./.git/*" \
  \( -name "*.go" -o -name "*.mod" -o -name "*.json" -o -name "*.yaml" -o -name "*.yml" -o -name "*.md" \) \
  -exec sed -i '' \
    -e "s|$TEMPLATE_MODULE|$NEW_MODULE|g" \
    -e "s/$TEMPLATE_REPO/$PACKAGE_NAME/g" \
    -e "s/TEMPLATE_PACKAGE_NAME/$PACKAGE_NAME/g" \
    -e "s/TEMPLATE_ORG/$GITHUB_ORG/g" {} +

# Fail loudly rather than leaving a half-renamed repo that only breaks at build time.
# This script is excluded: it holds TEMPLATE_REPO literally and is deleted below.
LEFTOVERS=$(grep -rlE "$TEMPLATE_REPO|TEMPLATE_PACKAGE_NAME|TEMPLATE_ORG" . \
  --exclude-dir=.git --exclude="setup-new-repo.sh" || true)
if [ -n "$LEFTOVERS" ]; then
  echo "❌ Rename incomplete — template identifiers still present in:"
  echo "$LEFTOVERS"
  exit 1
fi

echo "✅ Package name and module path updated"
echo ""

echo "🔧 Installing git hooks (pre-commit, commit-msg, pre-push)..."
make setup
echo "✅ Hooks installed in .git/hooks (managed by pre-commit)"
echo "   If a pre-existing .git/hooks/pre-commit.legacy appeared, delete it."
echo ""

# .golangci-version pins the linter for local hooks and `make lint`. It must stay in
# sync with the version installed in .github/workflows/ci.yml, or local and CI disagree.
echo ""
echo "🔍 Linter pin: .golangci-version = $(cat .golangci-version 2>/dev/null || echo 'MISSING')"
echo "   Keep it identical to the golangci-lint version in .github/workflows/ci.yml."

echo ""
echo "🗑️  Cleaning up template files..."
rm -f TEMPLATE_INSTRUCTION.md
# Drops the README provenance line, which points at the now-deleted TEMPLATE_INSTRUCTION.md.
sed -i '' '/^Scaffolded from /d' README.md
rm -f scripts/setup-new-repo.sh

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Enable GitHub Actions permissions in repo settings"
echo "2. Set the coverage threshold (threshold.total) in .testcoverage.yml"
echo "3. Review the diff, stage by explicit path, then commit"
echo "4. Push to GitHub: git push origin main"
echo "5. Start developing with feature branches!"
echo ""
echo "Happy coding! 🎉"
