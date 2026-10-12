# AGENTS.md — superfile fork (`feat/mouse-external-drop`)

This is a personal fork of superfile with a custom mouse feature set:

- mouse drag-and-drop that moves files by dropping them onto a directory
- a confirmation prompt before a drag drop moves anything (same modal flow as the
  delete confirmation)
- external file drops (terminal paste) that import paths as a copy
- right-click context menu and other pointer support

Do not "fix", simplify, or route around this behavior. It is intentional and was
built over several days. When in doubt, preserve existing behavior and ask.

## Build

`golangci-lint` is NOT installed on this machine, so `make build` always fails at
the lint step even when the code is fine.

Build directly with Go instead (the `main` package is at the repo root, not in
`src/cmd` — building `./src/cmd` produces a library archive, not an executable):

```sh
go build -o ./bin/spf .
```

The produced binary goes in gitignored `./bin/`. The installed copy the user
runs is `~/Work/bin/superfile-mouse` (a plain copy of the same binary) and
`/usr/local/bin/spf` (installed with sudo, root-owned).

## Tests

```sh
go test ./src/...
```

Known pre-existing failure (unrelated to this fork, also fails on a clean upstream
tree): `TestFileDelete/Move_to_trash`. Treat every other failure in the
`*mouse*`, `*drag*`, or `*context_menu*` tests as a regression of this fork.

## Merge/update care

The mouse work is uncommitted-upstream; upstream superfile does not contain it.
Any fetch/merge/rebase of upstream can break this fork in three ways:

1. **Merge conflicts** in shared files that this fork touches, especially:
   - `src/internal/key_function.go` (notify confirm/cancel dispatch, panel creation)
   - `src/internal/model.go`
   - `src/internal/type.go` (`moveToConfirm`, `drag` state)
   - `src/internal/ui/notify/type.go` (`MoveAction`)
   - `src/internal/common/predefined_variable.go` (`DragMoveWarn*` constants)
   - `src/internal/ui/filemodel/dimensions.go`
   - `src/internal/ui/notify/model.go` / `handle_file_operations.go` if upstream
     reworks the modal flow
2. **Dependency bumps** — `go mod tidy` on a `charm.land/bubbletea/v2` upgrade can
   silently change mouse-event delivery. Rebuild and retest after any such bump.
3. **Modal flow drift** — if upstream changes the delete/rename/notify key dispatch,
   the `MoveAction` hook stops firing with no compile error.

After any merge/update, always rebuild and run `go test ./src/...` before calling
it done.