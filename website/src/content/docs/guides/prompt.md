---
title: Prompt
description: Run shell commands and superfile commands from the built-in prompt.
head:
  - tag: title
    content: Prompt | superfile
---

Run shell commands or superfile commands without leaving superfile.

## Open and close

| Key             | Variable            | Action             |
| --------------- | ------------------- | ------------------ |
| `:`             | `open_command_line` | Open in shell mode |
| `>`             | `open_spf_prompt`   | Open in SPF mode   |
| `enter`         | `confirm_typing`    | Run the command    |
| `esc`, `ctrl+c` | `cancel_typing`     | Close the prompt   |

Type `:` or `>` in an empty prompt to switch modes.

## Shell mode

Runs the command in the focused panel's directory and shows the exit status and output.

```
: git status
```

- Uses `/bin/sh` (PowerShell on Windows), not your login shell.
- Shows stdout and stderr together.

:::caution

- superfile waits until a command finishes. It stops the shell after 5 seconds, but programs the shell started can keep running, and superfile waits for them too.
- Commands get no input, so interactive programs like `sudo` do not work.
- Long output is cut off. Redirect it to a file instead.

:::

## SPF mode

| Command       | Action                             |
| ------------- | ---------------------------------- |
| `split`       | New panel in the current directory |
| `open <PATH>` | New panel at `PATH`                |
| `cd <PATH>`   | Change the focused panel to `PATH` |

```
cd ~/projects
open "My Documents"
cd ${HOME}/Downloads
open $(dirname $(which bash))
```

- Relative paths start from the focused panel's directory.
- Quote paths with spaces, or escape them with `\`.
- `${VAR}` inserts an environment variable. `$VAR` is not expanded.
- `$(cmd)` inserts command output. Its shell is stopped after 1 second, but programs it started can make superfile wait longer.

:::note

Do not quote `$(cmd)`. Its trailing newline becomes part of the path.

:::

## Close after success

Set [shell_close_on_success](/configure/superfile-config#shell_close_on_success) to `true` to close the prompt after a successful command in either mode. You will not see the output.
