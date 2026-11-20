---
title: Metadata
description: The fields the metadata panel shows, and the plugins that add more.
head:
  - tag: title
    content: Metadata | superfile
---

The metadata panel, in the footer, shows details about the item under the cursor. Press `m` to focus it and scroll through the fields.

## Fields

| Field                           | Shows                                                                         |
| ------------------------------- | ----------------------------------------------------------------------------- |
| `Name`, `Size`, `Date Modified` | Basic file information                                                        |
| `Permissions`, `Owner`, `Group` | Access information. `Owner` and `Group` are empty on Windows.                 |
| `Architecture`                  | For programs: `ELF x86-64`, `PE ARM64`, `Mach-O ARM64` and so on              |
| `Attributes`                    | Linux: `lsattr` letters. FreeBSD: `chflags` names. Windows: `RHSAD`.          |
| `MD5Checksum`                   | With [`enable_md5_checksum`](/configure/superfile-config#enable_md5_checksum) |
| Every exiftool field            | With the [`metadata`](/configure/superfile-config#metadata) plugin            |

For a symlink, the panel only shows the target's `Path`.

## Folder size

A folder's size shows `(focus)` until you press `m`. superfile then adds up every file inside it, which can take a while for large folders.

## Plugins

The [`metadata`](/configure/superfile-config#metadata) plugin adds everything [exiftool](https://exiftool.org) reports, such as image size or camera model. If exiftool isn't installed, superfile shows the built-in fields without an error.

[`enable_md5_checksum`](/configure/superfile-config#enable_md5_checksum) reads the whole file every time you move to it. This can be slow for large files.

Restart superfile after changing either setting.
