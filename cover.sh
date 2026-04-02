#!/bin/bash
# Runs tests with full-module coverage and outputs uncovered line ranges.
#
# Usage: ./cover.sh [file-filter]
# Example: ./cover.sh ex.go
# Example: ./cover.sh refactor/
set -euo pipefail

filter="${1:-}"

coverpkgs="rsc.io/rf,rsc.io/rf/refactor,rsc.io/rf/diff"

no_proxy="" NO_PROXY="" go test -coverprofile=cover.out -coverpkg="$coverpkgs" -count=1 \
    rsc.io/rf rsc.io/rf/diff rsc.io/rf/refactor

echo ""
"$(dirname "$0")/uncovered.sh" cover.out "$filter"
