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

### Zoxide

- **Description:** Smart directory jumping integration with zoxide. Navigate to frequently used directories quickly with a searchable modal interface.

- **Requirements:** [`zoxide`](https://github.com/ajeetdsouza/zoxide)

- **Config name:** `zoxide`

- **Usage:** Press `z` to open the zoxide navigation modal. Start typing to search directories, use arrow keys to navigate results, and press Enter to jump to a directory.