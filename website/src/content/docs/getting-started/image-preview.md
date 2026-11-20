---
title: File preview
description: What superfile can preview, the tools each type needs, and how images are drawn.
head:
  - tag: title
    content: File preview | superfile
---

The preview panel shows the item under the cursor. Press `f` to show or hide it.

## Supported files

| File                                                 | Preview                  | Needs              |
| ---------------------------------------------------- | ------------------------ | ------------------ |
| Text and code                                        | First lines, highlighted | Nothing            |
| `.jpg` `.png` `.gif` `.webp` `.bmp` `.tiff`          | The image                | Nothing            |
| `.pdf`                                               | First page               | `pdftoppm`         |
| `.ps` `.eps`                                         | First page               | `gs` (Ghostscript) |
| `.mp4` `.mkv` `.mov` `.webm` `.avi` and other videos | A frame                  | `ffmpeg`           |
| Folders                                              | The folder's contents    | Nothing            |

superfile can't preview SVG or ICO files yet. A GIF shows its first frame only. See [Dependencies](/reference/dependencies) to install the tools.

[`show_image_preview`](/configure/superfile-config#show_image_preview) turns off images, PDFs and videos together.

## Image rendering

superfile uses the Kitty graphics protocol when `TERM` or `TERM_PROGRAM` is one of these:

```
kitty  xterm-kitty  ghostty  WezTerm  iTerm2  Konsole  WarpTerminal
```

Everywhere else, including inside tmux, images are drawn with colored half-block characters. Sixel is not supported.

:::note

Half-block images need your theme's `file_panel_bg` to be a full hex color such as `#1e1e2e`. Otherwise you see "Error converting image to ANSI".

:::

## Code highlighting

Code is highlighted by file name, with the style set by `code_syntax_highlight` in your [theme](/configure/custom-theme#theme-settings).

To use [bat](https://github.com/sharkdp/bat) instead, set [`code_previewer = "bat"`](/configure/superfile-config#code_previewer). superfile also finds it as `batcat`. If bat is missing, code files show an error.

Files that aren't UTF-8 text, such as binaries, show "Unsupported formats".

## Refreshing

The preview doesn't watch for changes. For text, code and folders, move the cursor away and back to see the latest content.

Images are cached. A changed image can keep showing its old version for several minutes. Resize the terminal or restart superfile to load the new one. A changed PDF, PostScript or video file shows its old thumbnail until you restart superfile.
