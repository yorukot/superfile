## General

| Function        | Key                      | Variable name |
| --------------- | ------------------------ | ------------- |
| Start SuperFile | `spf`                    |               |
| Reload          | press any key can reload |               |
| Quit            | `esc`, `q`               | `quit`        |

## Panel navigation

| Function                                           | Key               | Variable name           |
| -------------------------------------------------- | ----------------- | ----------------------- |
| Pin or Unpin folder to sidebar (can be auto saved) | `ctrl+p`          | `pinned_folder`         |
| Create new file panel                              | `ctrl+n`          | `create_new_file_panel` |
| Close the focused file panel                       | `ctrl+w`          | `close_file_panel`      |
| Focus on the file panel                            | `tab`             | `next_file_panel`       |
| Focus on the previous file panel                   | `shift+left`, `q` | `previous_file_panel`   |
| Focus on the processbar panel                      | `p`               | `focus_on_process_bar`  |
| Focus on the sidebar                               | `b`               | `focus_on_side_bar`     |
| Focus on the metadata panel                        | `m`               | `focus_on_meta_data`    |

## Panel movement

| Function                                            | Key                        | Variable name                                                             |
| --------------------------------------------------- | -------------------------- | ------------------------------------------------------------------------- |
| Change between selection mode or normal mode        | `v`                        | `change_panel_mode`                                                       |
| Up                                                  | `up`, `k`                  | `list_up`                                                                 |
| Down                                                | `down`, `j`                | `list_down`                                                               |
| Go to folder                                        | `enter`, `l`               | `select_item`                                                             |
| Return to parent folder                             | `h`, `backspace`           | `parent_folder`                                                           |
| Select all items in focused file panel              | `ctrl+a`                   | `file_panel_select_all_item`(only works in selection mode)                |
| Select with your course                             | `shift+up`, `K`(shift+k)   | `file_panel_select_mode_item_select_up`(only works in selection mode)     |
| Select with your course                             | `shift+left`, `J`(shift+j) | `file_panel_select_mode_items_elect_down`(only works in selection mode)   |
| Select the item where the current cursor is located | `enter`, `l`               | `file_panel_select_mode_item_single_select`(only works in selection mode) |
| Toggle dot file display                             | `ctrl+h`                   | `toggle_dot_file`                                                         |
| Toggle active search bar                            | `ctrl+f`                   | `search_bar`                                                              |

## File operations

| Function                                | Key          | Variable name                                                                        |
| --------------------------------------- | ------------ | ------------------------------------------------------------------------------------ |
| Create a new folder                     | `f`          | `file_panel_folder_create`                                                           |
| Create a new file                       | `c`          | `file_panel_file_create`                                                             |
| Rename file or folder                   | `r`          | `file_panel_item_rename`                                                             |
| Extract zip file                        | `ctrl+e`     | `extract_file`(normal mode)                                                          |
| Zip file or folder to .zip file         | `ctrl+r`     | `compress_file`(normal mode)                                                         |
| Delete file or folder (or both)         | `ctrl+d`     | `delete_item`(normal mode) <br> `file_panel_select_mode_item_delete`(select mode)    |
| Copy file or folder (or both)           | `ctrl+c`     | `copy_single_item`(normal mode) <br> `file_panel_select_mode_item_copy`(select mode) |
| Paste all items in your clipboard       | `ctrl+v`     | `paste_item`                                                                         |
| Open file with your default application | `enter`, `l` | `select_item`(select mode)                                                           |

## Special

| Function                                                                   | Key             | Variable name |
| -------------------------------------------------------------------------- | --------------- | ------------- |
| Confirm rename or create item or exit search bar                           | `enter`         | `confirm`     |
| Cancel rename or create item or exit search bar and clear search bar value | `ctrl+c`, `esc` | `cancel`      |
