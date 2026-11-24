---
title: Command line
description: Every command, flag and argument of spf.
head:
  - tag: title
    content: Command line | superfile
---

```
spf [COMMAND] [OPTIONS] [PATH]...
```

## Paths

Each path opens in its own panel.

```bash
spf ~/projects /tmp    # two panels
spf ~/notes/todo.md    # opens ~/notes with the cursor on todo.md
```

Without a path, superfile opens [`default_directory`](/configure/superfile-config#default_directory). A path that doesn't exist is looked up with zoxide if [`zoxide_support`](/configure/superfile-config#zoxide_support) is on, otherwise superfile opens your home folder.

## Options

| Option                  | Short   | Description                                                                                          |
| ----------------------- | ------- | ---------------------------------------------------------------------------------------------------- |
| `--config-file <path>`  | `-c`    | Use another config file                                                                              |
| `--with-config-override <key=value>` | `--wco` | Override a config setting for this run. Repeatable. Same syntax as config.toml |
| `--hotkey-file <path>`  | `--hf`  | Use another hotkey file                                                                              |
| `--chooser-file <path>` | `--cf`  | Write the chosen file paths here and quit. See [file picker](/guides/shell-integration#file-picker). |
| `--print-last-dir`      | `--pld` | Print the last folder on exit                                                                        |
| `--fix-config-file`     | `--fch` | Add missing settings to the config file                                                              |
| `--fix-hotkeys`         | `--fh`  | Add missing hotkeys to the hotkey file                                                               |
| `--debug-info`          | `--di`  | Print versions, paths and installed tools, for bug reports                                           |
| `--version`             | `-v`    | Print the version                                                                                    |
| `--help`                | `-h`    | Print the help                                                                                       |

Short options also work with one dash, such as `-pld`.

:::note

`--fix-config-file` and `--fix-hotkeys` back up your file as `<name>.bak-<random>`, then rewrite it. Your comments are lost.

:::

## Config overrides

Use `--with-config-override` (or `--wco`) to override settings from the loaded config file for the current run without modifying the file. Each value uses TOML `key=value` syntax. Repeat the flag to override multiple settings.

```bash
spf --wco 'debug=true' \
    --wco 'theme="gruvbox"' \
    --wco 'sidebar_sections=["home","pinned"]'
```

## Commands

| Command       | Description                                                                        |
| ------------- | ---------------------------------------------------------------------------------- |
| `spf pl`      | Print the paths of the config, hotkey, log and data files                          |
| `spf pl --ld` | Print the path of the lastdir file used by [cd on quit](/guides/shell-integration) |
| `spf help`    | Print the help                                                                     |
