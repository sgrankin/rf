# CLAUDE.md

## Project

rf is a Go refactoring tool by Russ Cox. Module path: `rsc.io/rf`.

Commands: `mv`, `ex`, `inline`, `add`, `rm`, `key`, `cp` (unimplemented).
See `doc.go` for full documentation.

## Build & Test

```
go test ./...                          # run all tests
go test -run TestRun/mv_func.txt       # run a single testdata case
go test -run TestRun -u                # update expected test output
go test -coverprofile=c.out ./...      # coverage
```

## Network

This environment routes through an egress proxy. Go toolchain downloads
and module fetches require clearing `no_proxy` to avoid bypassing the proxy
for domains that have no direct DNS path:

```
no_proxy="" NO_PROXY="" go mod tidy
no_proxy="" NO_PROXY="" go test ./...
```

## Testing Style

Tests use Russ Cox's **txtar** format (`golang.org/x/tools/txtar`).
Each `.txt` file in `testdata/` is a self-contained test case.

### txtar file structure

```
rf-command args
-- file.go --
package m
...
-- stdout --
expected diff output
-- stderr --
expected error output (optional)
```

- The **comment section** (before the first `-- file --` marker) contains
  the rf script to run. Lines starting with `-` are flags.
- **File sections** define the input source files for a temporary Go module.
- `-- stdout --` holds the expected unified diff output.
- `-- stderr --` holds expected error messages. Omit if no errors expected.
- The test harness creates a temporary directory with `go.mod` (`module m`)
  and writes all file sections into it, then runs the rf script with `-diff`.

### Naming conventions

Test files are named `{command}_{variant}.txt`:
- `mv_func.txt`, `mv_file.txt`, `mv_stmts.txt` — mv command variants
- `ex_type.txt`, `ex_import.txt` — ex command variants
- `add.txt`, `add_error.txt` — add command, including error cases
- `inline_const.txt`, `inline_const_rm.txt` — inline command

Error cases use `_error` suffix or include a `-- stderr --` section.

### Writing new tests

1. Create `testdata/{command}_{description}.txt`
2. Write the rf command on the first line
3. Add input `.go` files as txtar sections
4. Add `-- stdout --` with expected diff output
5. Or run `go test -run TestRun/{name} -u` to auto-generate expected output

### Unit tests

Small pure functions use table-driven tests with struct slices
(see `ex_test.go`, `refactor/addr_test.go`). No test frameworks — just
`testing.T` with `t.Errorf`.

## Code Style

- No test frameworks, no assertion libraries.
- No unnecessary abstractions. Three similar lines are better than a premature helper.
- Comments only where the logic is non-obvious.
- `types.Unalias()` must be used when comparing `types.Type` values,
  since Go 1.24+ represents type aliases as `*types.Alias`.
