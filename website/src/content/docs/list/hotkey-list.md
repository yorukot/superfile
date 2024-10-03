---
title: Hotkey list
description: Superfile hotkey list
head:
  - tag: title
    content: Hotkey list | superfile
---

:::tip
These are the default hotkeys and you can [change](/configure/custom-hotkeys) them all!
:::

## General

| Function                        | Key              | Variable name    |
| ------------------------------- | ---------------- | ---------------- |
| Open superfile                  | `spf`            |                  |
| Confirm your select or typing   | `enter`, `right` | `confirm_typing` |
| Quit typing, modal or superfile | `esc`, `q`       | `quit`           |
| Cancel typing                   | `ctrl+c`, `esc`  | `cancel_typing`  |
| Open help menu(hotkeylist)      | `?`              | `open_help_menu` |

## Panel navigation

| Function                         | Key                        | Variable name               |
| -------------------------------- | -------------------------- | --------------------------- |
| Create new file panel            | `n`                        | `create_new_file_panel`     |
| Close the focused file panel     | `w`                        | `close_file_panel`          |
| Toggle file preview panel        | `f`                        | `toggle_file_preview_panel` |
| Focus on the next file panel     | `tab`, `L`(shift+l)        | `next_file_panel`           |
| Focus on the previous file panel | `shift+left`, `H`(shift+h) | `previous_file_panel`       |
| Focus on the processbar panel    | `p`                        | `focus_on_process_bar`      |
| Focus on the sidebar             | `s`                        | `focus_on_side_bar`         |
| Focus on the metadata panel      | `m`                        | `focus_on_metadata`         |

## Panel movement

| Function                                           | Key                        | Variable name                                                   |
| -------------------------------------------------- | -------------------------- | --------------------------------------------------------------- |
| Up                                                 | `up`, `k`                  | `list_up`                                                       |
| Down                                               | `down`, `j`                | `list_down`                                                     |
| Return to parent folder                            | `h`, `left`, `backspace`   | `parent_folder`                                                 |
| Select all items in focused file panel             | `A`(shift+a)               | `file_panel_select_all_item` (selection mode only)              |
| Select up with your course                         | `shift+up`, `K`(shift+k)   | `file_panel_select_mode_item_select_up` (selection mode only)   |
| Select down with your course                       | `shift+down`, `J`(shift+j) | `file_panel_select_mode_item_select_down` (selection mode only) |
| Toggle dot file display                            | `.`                        | `toggle_dot_file`                                               |
| Toggle active search bar                           | `/`                        | `search_bar`                                                    |
| Change between selection mode or normal mode       | `v`                        | `change_panel_mode`                                             |
| Pin or Unpin folder to sidebar (can be auto saved) | `P`(shift+p)               | `pinned_folder`                                                 |

## File operations

| Function                                             | Key                | Variable name                                                                          |
| ---------------------------------------------------- | ------------------ | -------------------------------------------------------------------------------------- |
| Create file or folder(/ ends with creating a folder) | `ctrl+n`           | `file_panel_item_create`                                                               |
| Rename file or folder                                | `ctrl+r`           | `file_panel_item_rename`                                                               |
| Copy file or folder (or both)                        | `ctrl+c`           | `copy_single_item` (normal mode) <br> `file_panel_select_mode_item_copy` (select mode) |
| Cut file or folder (or both)                         | `ctrl+x`           | `file_panel_select_mode_item_cut`                                                      |
| Paste all items in your clipboard                    | `ctrl+v`           | `paste_item`                                                                           |
| Delete file or folder (or both)                      | `ctrl+d`, `delete` | `delete_item` (normal mode) <br> `file_panel_select_mode_item_delete` (select mode)    |
| Copy current file or directory path                  | `ctrl+p`           | `copy_path`                                                                            |
| Extract zip file                                     | `ctrl+e`           | `extract_file` (normal mode)                                                           |
| Zip file or folder to .zip file                      | `ctrl+a`           | `compress_file` (normal mode)                                                          |
| Open file with your default editor                   | `e`                | `oepn_file_with_editor` (normal node)                                                  |
| Open current directory with default editor           | `E`(shift+e)       | `current_directory_with_editor` (normal node)                                          |
