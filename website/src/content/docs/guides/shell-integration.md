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

On quit, superfile writes a `cd` command to the lastdir file, and the shell function runs it. Run `spf pl --lastdir-file` to print where that file is.

:::caution

The fish function always reads `~/.local/state/superfile/lastdir`. If you set `XDG_STATE_HOME`, change the path in the function.

:::

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

To pick several files, press `v`, select them, then press `e`.

The file holds full paths, one per line, with no trailing newline. Quitting with `q` writes nothing.
