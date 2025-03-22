---
title: Config file path
description: All superfile config file path
head:
  - tag: title
    content: Config file path | superfile
---

:::tip
If you want to get the set path you can try `spf pl` which will print out the file locations of all superfile.
:::

## Directories

#### Config directory

|         Linux         |              macOS              |          Windows           |
| :-------------------: | :-----------------------------: | :------------------------: |
| `~/.config/superfile` | `~/Library/Application Support/superfile` | `%LOCALAPPDATA%/superfile` |

#### Theme directory

|            Linux            |                      macOS                      |             Windows              |
| :-------------------------: | :---------------------------------------------: | :------------------------------: |
| `~/.config/superfile/theme` | `~/Library/Application Support/superfile/theme` | `%LOCALAPPDATA%/superfile/theme` |

#### Data directory

|           Linux            |                   macOS                    |          Windows           |
| :------------------------: | :----------------------------------------: | :------------------------: |
| `~/.local/share/superfile` | `~/Library/Application Support/superfile/` | `%LOCALAPPDATA%/superfile` |

### Changing Config File Path

You can use the `-c` or `--config-file` flag to specify a different path for the `config.toml` file:

```bash
spf -c /path/to/your/config.toml
```

You can use the `--hotkey-file` flag to specify a different path for the `hotkey.toml` file:

```bash
spf --hotkey-file /path/to/your/hotkey.toml
```

#### Log directory

|           Linux            |                   macOS                   |          Windows           |
| :------------------------: | :---------------------------------------: | :------------------------: |
| `~/.local/state/superfile` | `~/Library/Application Support/superfile` | `%LOCALAPPDATA%/superfile` |

---

## All config file path

#### Config

|               Linux               |                    macOS                    |                Windows                 |
| :-------------------------------: | :-----------------------------------------: | :------------------------------------: |
| `~/.config/superfile/config.toml` | `~/Library/Application Support/superfile/config.toml` | `%LOCALAPPDATA%/superfile/config.toml` |

#### Hotkeys

|               Linux                |                    macOS                     |                 Windows                 |
| :--------------------------------: | :------------------------------------------: | :-------------------------------------: |
| `~/.config/superfile/hotkeys.toml` | `~/Library/Application Support/superfile/hotkeys.toml` | `%LOCALAPPDATA%/superfile/hotkeys.toml` |

#### Log file

|                  Linux                   |                          macOS                          |                 Windows                  |
| :--------------------------------------: | :-----------------------------------------------------: | :--------------------------------------: |
| `~/.local/state/superfile/superfile.log` | `~/Library/Application Support/superfile/superfile.log` | `%LOCALAPPDATA%/superfile/superfile.log` |
