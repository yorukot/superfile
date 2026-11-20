---
title: Troubleshooting
description: Have you encountered any problems? Come here and take a look.
head:
  - tag: title
    content: Troubleshooting | superfile
---

## My superfile icon doesn't display correctly

Try these things below:

- Make sure you have already installed [nerdfont](https://www.nerdfonts.com/font-downloads) (You can choose whatever font you like!)
- Apply this font to your terminal. This may require different settings depending on the terminal. You can check how to set it up!
- **WSL users:** installing the font inside the Linux filesystem (e.g. `~/.fonts` + `fc-cache`) has no effect. The terminal is rendered by the Windows host app, so you need to install the Nerd Font on Windows and set it in your terminal's settings (Windows Terminal, VS Code, etc.).

## Help! My superfile's rendering is all messed up!

Try these things below:

- Set your locale to utf-8
- chcp 65001 ( If that's an option for your shell )
- Set environment variable RUNEWIDTH_EASTASIAN to 0 (`RUNEWIDTH_EASTASIAN=0`)

## Images show as colored blocks

superfile only uses the Kitty graphics protocol in [supported terminals](/getting-started/image-preview#image-rendering). Try these things below:

- Check your terminal with `echo $TERM $TERM_PROGRAM`
- Run superfile outside tmux

## PDF or video files show "Unsupported formats"

Install the tool for that file type (`pdftoppm` for PDF, `ffmpeg` for video), then restart superfile. See [Dependencies](/reference/dependencies).

## Reporting a bug

Run `spf --debug-info` and include its output in your issue.
