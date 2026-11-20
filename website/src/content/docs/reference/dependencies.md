---
title: Dependencies
description: Optional tools that add features to superfile.
head:
  - tag: title
    content: Dependencies | superfile
---

superfile runs without any of these tools. Each one adds a feature. Run `spf --debug-info` to see which ones superfile can find.

| Tool                                              | Adds                                                                                          |
| ------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| A [Nerd Font](https://www.nerdfonts.com)          | File icons. Otherwise set `nerdfont = false`.                                                 |
| `pdftoppm` (from Poppler)                         | [PDF preview](/getting-started/image-preview)                                                 |
| `gs` (Ghostscript)                                | PostScript preview                                                                            |
| `ffmpeg`                                          | Video preview                                                                                 |
| [`bat`](https://github.com/sharkdp/bat)           | Code highlighting with [`code_previewer = "bat"`](/configure/superfile-config#code_previewer) |
| [`exiftool`](https://exiftool.org)                | The [`metadata`](/configure/superfile-config#metadata) plugin                                 |
| [`zoxide`](https://github.com/ajeetdsouza/zoxide) | The [`zoxide_support`](/configure/superfile-config#zoxide_support) plugin                     |
| `wl-copy`, `xclip` or `xsel` (Linux)              | Copying paths with `ctrl+p` and `c`                                                           |
| `xdg-open` (Linux)                                | Opening files with `enter`                                                                    |

superfile looks for these tools when it starts. Restart it after installing one.
