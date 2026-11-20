---
title: Sidebar
description: Jump to common folders, pinned folders and disks.
head:
  - tag: title
    content: Sidebar | superfile
---

The sidebar lists common folders, pinned folders and disks. Press `s` to focus it, then `enter` to open an entry in the focused panel.

## Sections

| Section | Lists                                                                                      |
| ------- | ------------------------------------------------------------------------------------------ |
| Home    | Home, Desktop, Downloads, Documents, Pictures, Videos, Music, Templates, Trash (Linux)     |
| Pinned  | Folders you pinned                                                                         |
| Disks   | `/` and drives under `/mnt`, `/media`, `/run/media` and `/Volumes`. All drives on Windows. |

Folders that don't exist are hidden. Choose and order the sections with [`sidebar_sections`](/configure/superfile-config#sidebar_sections), or hide the sidebar with [`sidebar_width = 0`](/configure/superfile-config#sidebar_width).

## Pinned folders

`P` pins or unpins the folder you are in, not the item under the cursor.

- With the sidebar focused, `ctrl+r` renames a pin's label. The folder itself isn't renamed.
- Pins are saved in `pinned.json` in the [data directory](/configure/config-file-path#data-directory). Edit it to reorder them.
- Pins to folders that no longer exist are removed.

## Filter

With the sidebar focused, press `/` to filter the entries. `enter` keeps the filter and `esc` clears it.
