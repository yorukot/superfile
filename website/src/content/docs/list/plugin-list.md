---
title: Plugin List
description: Complete list of available superfile plugins
head:
  - tag: title
    content: Plugin List | superfile
---

Superfile supports various plugins to extend its functionality. Below is a complete list of available plugins and their requirements.

### Metadata

- **Description:** Show more detailed metadata for files and directories

- **Requirements:** [`exiftool`](https://exiftool.org)

- **Config name:** `metadata`

### MD5 Checksum

- **Description:** Show MD5 checksums for regular files in the metadata panel

- **Requirements:** None

- **Config name:** `enable_md5_checksum`

- **Note:** Calculating checksums reads the selected file, so it may be slow for large files.

### Zoxide

- **Description:** Smart directory jumping integration with zoxide. Navigate to frequently used directories quickly with a searchable modal interface.

- **Requirements:** [`zoxide`](https://github.com/ajeetdsouza/zoxide)

- **Config name:** `zoxide_support`

- **Usage:** Press `z` to open the zoxide navigation modal. Start typing to search directories, use arrow keys to navigate results, and press Enter to jump to a directory.

### Find file/folder

- **Description:** Recursive file and folder search using fd. Search from a modal, then jump to a match — selecting a file moves to its parent directory and places the cursor on it; selecting a folder moves into it.

- **Requirements:** [`fd`](https://github.com/sharkdp/fd)

- **Config name:** `find_file_support`

- **Usage:** Press `ctrl+f` (default) or `F` (vim hotkeys) to open the find modal. Start typing to filter results (empty query lists all files and folders), use arrow keys to navigate results, and press Enter to move to the selected file or folder.
