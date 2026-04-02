#!/bin/bash
# Runs tests with full-module coverage and outputs uncovered line ranges.
# Subprocess coverage (from TestMain dispatching) is merged into the profile.
#
# Usage: ./cover.sh [file-filter]
# Example: ./cover.sh ex.go
# Example: ./cover.sh refactor/
set -euo pipefail

filter="${1:-}"

coverpkgs="rsc.io/rf,rsc.io/rf/refactor,rsc.io/rf/diff"

# Create a directory for subprocess coverage data.
# GOCOVERDIR passes through to child processes automatically.
subcoverdir=$(mktemp -d)
trap 'rm -rf "$subcoverdir"' EXIT

no_proxy="" NO_PROXY="" GOCOVERDIR="$subcoverdir" \
    go test -coverprofile=cover.out -coverpkg="$coverpkgs" -count=1 \
    rsc.io/rf rsc.io/rf/diff rsc.io/rf/refactor

# Merge subprocess coverage data into the profile, if any was collected.
if ls "$subcoverdir"/*.* >/dev/null 2>&1; then
    go tool covdata textfmt -i="$subcoverdir" -o=subcover.out
    # Append subprocess profile lines (skipping its mode header) to cover.out.
    tail -n +2 subcover.out >> cover.out
    rm -f subcover.out
fi

echo ""
"$(dirname "$0")/uncovered.sh" cover.out "$filter"
