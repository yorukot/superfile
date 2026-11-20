---
title: Navigation and search
description: Panels, sorting, search, zoxide and the mouse.
head:
  - tag: title
    content: Navigation and search | superfile
---

The basic keys are in the [tutorial](/getting-started/tutorial). This page covers how they behave.

## Panels

- `n` opens a new panel in your home folder. `N` opens one in the current folder.
- Each panel has its own folder, search, selection and sort.
- Each panel remembers the cursor position in every folder. Going up puts the cursor on the folder you left.
- The number of panels is limited by the terminal width, up to 10.

To open several panels at start, pass several paths: `spf ~/projects ~/Downloads`. See [Command line](/reference/cli).

## Sort

| Type          | Order                                           |
| ------------- | ----------------------------------------------- |
| Name          | Folders first, A to Z                           |
| Size          | Folders first by item count, then files by size |
| Date Modified | Newest first, folders mixed with files          |
| Type          | Folders first, then files by extension          |
| Natural       | Like Name, but `file2` comes before `file10`    |

`o` and `R` change the focused panel only. The default comes from [`default_sort_type`](/configure/superfile-config#default_sort_type).

## Search

`/` filters the current folder with fuzzy, smart-case matching, like fzf. Press `enter` to keep the filter, or `esc` to clear it.

| Query  | Matches               |
| ------ | --------------------- |
| `'abc` | Contains `abc`        |
| `^abc` | Starts with `abc`     |
| `abc$` | Ends with `abc`       |
| `!abc` | Doesn't contain `abc` |

Hidden dotfiles aren't searched unless you show them with `.`. The filter clears when you change folders.

In the help menu (`?`), `/` searches the hotkeys.

## Zoxide

With [`zoxide_support`](/configure/superfile-config#zoxide_support) on, press `z`, type keywords and press `enter` to jump. Move through the results with the arrow keys; `j` and `k` are typed.

superfile adds folders to zoxide as you open them.

## Mouse

Only the wheel works. It scrolls the focused panel. Hold `shift` to select text in the terminal.

## Good to know

- superfile needs a terminal of at least 60 × 24.
- Whether dotfiles and the footer are shown is remembered between sessions. The preview isn't.
