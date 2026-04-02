#!/bin/bash
# Reads a Go coverage profile and prints uncovered line ranges per file.
# Usage: ./uncovered.sh [coverprofile] [file-filter]
# Example: ./uncovered.sh c.out ex.go
#
# A line is only reported as uncovered if it appears in at least one
# count=0 block and does NOT appear in any count>0 block.

profile="${1:-c.out}"
filter="${2:-}"

awk '
NR == 1 { next }  # skip mode line
{
    split($0, parts, " ")
    count = parts[3] + 0
    spec = parts[1]
    ci = index(spec, ":")
    file = substr(spec, 1, ci - 1)
    rest = substr(spec, ci + 1)
    split(rest, lr, ",")
    split(lr[1], sl, ".")
    split(lr[2], el, ".")
    startLine = sl[1] + 0
    endLine = el[1] + 0

    for (l = startLine; l <= endLine; l++) {
        if (count > 0) {
            covered[file, l] = 1
        } else {
            uncov[file, l] = 1
        }
        files[file] = 1
    }
}
END {
    for (file in files) {
        # Collect truly uncovered lines
        n = 0
        for (key in uncov) {
            split(key, kp, SUBSEP)
            if (kp[1] != file) continue
            l = kp[2] + 0
            if (!covered[file, l]) {
                lines[++n] = l
            }
        }
        if (n == 0) continue

        # Insertion sort
        for (i = 2; i <= n; i++) {
            key = lines[i]
            j = i - 1
            while (j >= 1 && lines[j] > key) {
                lines[j+1] = lines[j]
                j--
            }
            lines[j+1] = key
        }

        # Build ranges
        ranges = ""
        rstart = lines[1]; rend = lines[1]
        for (i = 2; i <= n; i++) {
            if (lines[i] == rend + 1) {
                rend = lines[i]
            } else {
                if (ranges != "") ranges = ranges ","
                ranges = ranges (rstart == rend ? rstart : rstart "-" rend)
                rstart = lines[i]; rend = lines[i]
            }
        }
        if (ranges != "") ranges = ranges ","
        ranges = ranges (rstart == rend ? rstart : rstart "-" rend)

        print file ": " ranges
        delete lines
    }
}
' "$profile" | sort | if [ -n "$filter" ]; then grep "$filter"; else cat; fi
