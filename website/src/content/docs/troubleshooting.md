---
title: Troubleshooting
description: Have you encountered any problems? Come here and take a look.
head:
  - tag: title
    content: Troubleshooting | superfile
---

## My superfile icon doesn't display correctly

Try these things below:

- Make sure you already install [nerdfont](https://www.nerdfonts.com/font-downloads) (You can choose whatever font you like!)
- Apply this font to your terminal,This may require different settings depending on the terminal.You can check how to set it up!

## Help! My superfile's rendering is all messed up!

Try these things below:

- Set your locale to utf-8  
- chcp 65001 ( If that's an option for your shell )  
- Set environment variable RUNEWIDTH_EASTASIAN to 0 (`RUNEWIDTH_EASTASIAN=0`)

## superfile chooser/save mode is not working through my portal wrapper

If you launch superfile through `xdg-desktop-portal-termfilechooser` or another wrapper:

- Use `--chooser-file` for open-file selection output.
- Use `--save-file` for save-target selection output.
- `--chooser-file` now supports multi-select and writes newline-delimited absolute paths.
- `--save-file` uses superfile's save flow, where `e` confirms the focused file or ghost and `E` confirms the current directory plus the ghost name.
- The `xdg-desktop-portal-termfilechooser` superfile wrapper must call `spf --save-file="$out" "$path"` for save requests. Older wrappers that always call `--chooser-file` will not enter save mode.
