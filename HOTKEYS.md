## General

| Function        | Key        | Variable name |
| --------------- | ---------- | ------------- |
| Start SuperFile | `spf`      |               |
| Reload          | `ctrl+R`   | `reload`      |
| Quit            | `esc`, `q` | `quit`        |

## Panel navigation

| Function                                                   | Key               | Variable name        |
| ---------------------------------------------------------- | ----------------- | -------------------- |
| Pin or Unpin folder to sidebar (can be auto saved) | `ctrl+p`          | `pinnedFolder`       |
| Create new file panel                                      | `ctrl+n`          | `createNewFilePanel` |
| Close the focused file panel                               | `ctrl+w`          | `closeFilePanel`     |
| Focus on the file panel                                    | `tab`             | `nextFilePanel`      |
| Focus on the previous file panel                           | `shift+left`, `q` | `previousFilePanel`  |
| Focus on the processbar panel                              | `p`               | `focusOnProcessBar`  |
| Focus on the sidebar                                       | `b`               | `focusOnSideBar`     |
| Focus on the metadata panel                                | `m`               | `focusOnMetaData`    |

## Panel movement

| Function                                            | Key                        | Variable name                                                      |
| --------------------------------------------------- | -------------------------- | ------------------------------------------------------------------ |
| Change between selection mode or normal mode             | `v`                        | `changePanelMode`                                                  |
| Up                                                  | `up`, `k`                  | `listUp`                                                           |
| Down                                                | `down`, `j`                | `listDown`                                                         |
| Go to folder                                        | `enter`, `l`               | `selectItem`                                                       |
| Return to parent folder                             | `h`, `backspace`           | `parentFolder`                                                     |
| Select all items in focused file panel               | `ctrl+a`                   | `filePanelSelectAllItem`(only works in selection mode)              |
| Select with your course                             | `shift+up`, `K`(shift+k)   | `filePanelSelectModeItemSelectUp`(only works in selection mode)     |
| Select with your course                             | `shift+left`, `j`(shift+j) | `filePanelSelectModeItemSelectDown`(only works in selection mode)   |
| Select the item where the current cursor is located | `enter`, `l`               | `filePanelSelectModeItemSingleSelect`(only works in selection mode) |

## File operations

| Function                         | Key      | Variable name                                                                 |
| -------------------------------- | -------- | ----------------------------------------------------------------------------- |
| Create a new folder              | `f`      | `filePanelFolderCreate`                                                       |
| Create a new file                | `c`      | `filePanelFileCreate`                                                         |
| Rename file or folder            | `r`      | `filePanelItemRename`                                                         |
| Delete file or folder (or both)    | `ctrl+d` | `deleteItem`(normal mode) <br> `filePanelSelectModeItemDelete`(select mode)   |
| Copy file or folder (or both)      | `ctrl+c` | `copySingleItem`(normal mode) <br> `filePanelSelectModeItemCopy`(select mode) |
| Paste all items in your clipboard | `ctrl+v` | `pasteItem`                                                                   |

## Special

| Function                      | Key             | Variable name |
| ----------------------------- | --------------- | ------------- |
| Confirm rename or create item | `enter`         | `confirm`     |
| Cancel rename or create item  | `ctrl+c`, `esc` | `cancel`      |
