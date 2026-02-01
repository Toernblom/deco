#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="/workspace/deco"
cd "$REPO_DIR"

if ! command -v go >/dev/null 2>&1; then
  echo "Go not found. Enable Go in the environment or install it before running this script."
  exit 1
fi

export GOBIN="$HOME/.local/bin"
mkdir -p "$GOBIN"

# Install bd (Beads)
go install github.com/steveyegge/beads/cmd/bd@latest

# Install deco from this repo
go install ./cmd/deco

# Optional: prefetch modules for faster builds/tests
go mod download

# Persist PATH for future shells
if ! grep -q 'export PATH="$HOME/.local/bin:$PATH"' "$HOME/.bashrc" 2>/dev/null; then
  echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
fi
