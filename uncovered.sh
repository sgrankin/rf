#!/bin/bash
# Reads a Go coverage profile and prints uncovered line ranges per file.
# Usage: ./uncovered.sh [coverprofile] [file-filter]
# Example: ./uncovered.sh c.out ex.go

profile="${1:-c.out}"
filter="${2:-}"

awk '
NR == 1 { next }  # skip mode line
{
    # Format: pkg/file.go:startLine.startCol,endLine.endCol numStmts count
    # Split on space first to get the block spec, numStmts, count
    split($0, parts, " ")
    count = parts[3] + 0

    if (count != 0) next

    # Parse "pkg/file.go:startLine.startCol,endLine.endCol"
    spec = parts[1]
    # Find the colon separating file from line info
    ci = index(spec, ":")
    file = substr(spec, 1, ci - 1)
    rest = substr(spec, ci + 1)
    # rest is "startLine.startCol,endLine.endCol"
    split(rest, lr, ",")
    split(lr[1], sl, ".")
    split(lr[2], el, ".")
    startLine = sl[1] + 0
    endLine = el[1] + 0

    # Store uncovered ranges per file
    n = file_count[file]++
    starts[file, n] = startLine
    ends[file, n] = endLine
}
END {
    for (file in file_count) {
        n = file_count[file]
        # Sort ranges by start line (simple insertion sort)
        for (i = 1; i < n; i++) {
            s = starts[file, i]; e = ends[file, i]
            j = i - 1
            while (j >= 0 && starts[file, j] > s) {
                starts[file, j+1] = starts[file, j]
                ends[file, j+1] = ends[file, j]
                j--
            }
            starts[file, j+1] = s
            ends[file, j+1] = e
        }
        # Merge overlapping/adjacent ranges
        m = 0
        ms[0] = starts[file, 0]; me[0] = ends[file, 0]
        for (i = 1; i < n; i++) {
            if (starts[file, i] <= me[m] + 1) {
                if (ends[file, i] > me[m]) me[m] = ends[file, i]
            } else {
                m++
                ms[m] = starts[file, i]; me[m] = ends[file, i]
            }
        }
        # Format output
        ranges = ""
        for (i = 0; i <= m; i++) {
            if (ranges != "") ranges = ranges ","
            if (ms[i] == me[i])
                ranges = ranges ms[i]
            else
                ranges = ranges ms[i] "-" me[i]
        }
        print file ": " ranges
    }
}
' "$profile" | sort | if [ -n "$filter" ]; then grep "$filter"; else cat; fi
