---
title: File operations
description: How superfile copies, moves, deletes, compresses, extracts and opens files.
head:
  - tag: title
    content: File operations | superfile
---

In normal mode, an operation acts on the item under the cursor. In selection mode, it acts on the selected items. See [which keys work in which mode](/list/hotkey-list#which-keys-work-in-which-mode).

## Create and rename

Press `ctrl+n` and type a name. End it with `/` (`\` on Windows) to create a folder. `docs/notes/todo.txt` creates the missing folders too.

If the name is taken, superfile creates `name(1).ext` instead.

Press `ctrl+r` to rename. If the new name is taken, superfile warns you; press `enter` to overwrite the existing item.

## Copy, cut and paste

`ctrl+c` and `ctrl+x` put items on superfile's own clipboard, shown in the clipboard panel. `ctrl+v` pastes them into the current folder.

- Pasting a copy never overwrites. A name clash makes `name(1).ext`.
- After a copy, the clipboard keeps its items, so you can paste again. After a cut, it is emptied.
- Symlinks are copied as links.

:::caution

A cut doesn't rename on a name clash. A moved file replaces the existing file, and a moved folder is merged into the existing folder.

:::

## Delete

`ctrl+d` moves items to the trash. `D` deletes them permanently. Both ask you to confirm.

| System  | Trash                                                                                       |
| ------- | ------------------------------------------------------------------------------------------- |
| Linux   | `~/.local/share/Trash`. Items on other drives go to `.Trash-<uid>` at the top of the drive. |
| macOS   | The system Trash                                                                            |
| Windows | The Recycle Bin                                                                             |

When no trash is available, superfile deletes permanently, and the confirmation says "completely delete". This happens inside a trash folder, on a drive superfile can't write to, and on systems other than the three above.

superfile can't restore items from the trash.

## Compress and extract

`ctrl+a` compresses into a `.zip` file in the current folder. In selection mode, all selected items go into one archive. `photo.jpg` becomes `photo.zip`.

`ctrl+e` extracts the archive under the cursor into a new folder next to it. `photos.tar.gz` extracts into `photos/`.

| Can extract                                                    | Can't extract yet                 |
| -------------------------------------------------------------- | --------------------------------- |
| `.zip`, `.tar`, `.gz`, `.tar.gz`, `.bz`, `.7z`, `.rar`, `.iso` | `.tar.bz2`, `.tgz`, `.xz`, `.zst` |

## Copy paths

`ctrl+p` copies the paths of the item under the cursor, or of all selected items, to your system clipboard. `c` copies the current folder's path.

On Linux, this needs `wl-copy`, `xclip` or `xsel`. See [Dependencies](/reference/dependencies).

## Progress and errors

The processes panel lists every operation of the session with its progress. Press `p` to focus it. Operations can't be cancelled.

If an item fails, superfile shows an error. Press `enter` to skip it and continue, or `esc` to stop.

## Open files

`enter` opens a file with your system's default application: `xdg-open` on Linux, `open` on macOS. To choose a program per extension, use [`open_with`](/configure/superfile-config#open_with).

`e` opens the item in your [`editor`](/configure/superfile-config#editor), and `E` opens the current folder in your [`dir_editor`](/configure/superfile-config#dir_editor). superfile waits until the editor exits.
