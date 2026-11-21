---
title: Shell integration
description: Change your shell's folder when superfile quits, and use superfile as a file picker.
head:
  - tag: title
    content: Shell integration | superfile
---

superfile can hand the last folder back to your shell, and the files you pick back to another program.

## cd on quit

Your shell moves to the folder you were browsing when superfile quits.

1. Set [`cd_on_quit`](/configure/superfile-config#cd_on_quit) to `true`.
2. Add the `spf` shell function from the [config page](/configure/superfile-config#cd_on_quit).

Press `Q` to quit and cd once, even when `cd_on_quit` is `false`.

On quit, the shell function reads the last folder from a temporary file it passes with `--cwd-file`, so each shell session has its own file, and runs `cd` to it. `--cwd-file <path>` writes only the folder path, and only when superfile quits with cd on quit. Superfile still writes the older `cd` command to the lastdir file, which `spf path-list --lastdir-file` prints, so a function from an earlier release keeps working. To use the new behavior, replace it with the current one from the [config page](/configure/superfile-config#cd_on_quit) and reload your shell.

`spf --print-last-dir` prints the last folder on exit. It doesn't work in `$(...)`, because superfile draws its interface on standard output too.

## File picker

Start superfile with `--chooser-file` to pick files for an editor or a script. Instead of opening a file, superfile writes its path to the chooser file and quits.

```bash
chosen=$(mktemp)
spf --chooser-file "$chosen"
cat "$chosen"
```

| Key               | Writes                                                          |
| ----------------- | --------------------------------------------------------------- |
| `enter` on a file | The file's path                                                 |
| `e`               | The selected items, or the item under the cursor (also folders) |
| `E`               | The current folder                                              |

To pick several files, press `v`, mark each item with `enter`, then press `e`. The paths appear in the same order as the items in the panel.

The file holds full paths, one per line, with no trailing newline. Quitting with `q` writes nothing.
