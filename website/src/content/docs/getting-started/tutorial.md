---
title: Learn how to use tutorial
description: Quickly get started with superfile
head:
  - tag: title
    content: Tutorial | superfile
---

This tutorial will teach you how to use superfile step by step.

:::caution
If you haven't installed superfile yet, please [click here](/getting-started/installation).
:::

:::tip
A full list of hotkeys are available [here](/list/hotkey-list)
:::

## Hotkeys tutorial

Let's start by running superfile! Open a terminal, type `spf` and press `enter`.

To exit, press `q` or `esc`.

![hotkeys-demo](../../../assets/demo/hotkeys-demo.gif)

### Panel navigation

Once superfile is running, it displays five panels:

- sidebar
- file
- processes
- metadata
- clipboard
- command execution bar

The file panel is the focused view by default. You can change focus onto three other panels.

Press `s` to focus on the sidebar.

Press `p` to focus on the processes.

Press `m` to focus on the metadata.

Press `:` to open command execution bar.

To return focus back onto the file panel, press the same hotkey again.

> For command execution bar you need press `esc` or `ctrl+c`

You can also press `f` to show or hide the preview window.

Also press `F` to hide or show all footer panel.

![panel-navigation-demo](../../../assets/demo/panel-navigation-demo.gif)

:::tip
The size of the folder will only be shown when you focus on the metadata.

For more detailed metadata, [click here](/configure/enable-plugin) to install the metadata plugin.
:::

To create more file panels, press `n`. Press `w` to close the focused file panel.

To move through multiple file panels, press `tab` or `L` (shift+l). To move to the previous panel, press `shift`+`left` or `H` (shift+h).

![multiple-panels-demo](../../../assets/demo/multiple-panels-demo.gif)

### Panel movement

superfile provides multiple hotkeys to move through directories. The angle bracket cursor `>` tells you where you are.

While focused on the file panel, move the cursor up with `up` or `k` and down with `down` or `j`.

After navigating to the your file/folder, press `enter` or `l` to confirm your selection. Files are opened with your default application (if none set, there will be no response) and folders are opened for viewing. Press `h` or `backspace` to return to the parent directory.

![panel-movement-demo](../../../assets/demo/panel-movement-demo.gif)

Folders can be pinned to the sidebar panel. Navigate to and open your folder. Press `P` (shift+p) to pin or unpin it.

Press `o` to bring up the sort options menu. You can sort by:

- `Name`
- `Size`
- `Date Modified`

Press `enter` to confirm your sort option. Press `esc`, `o`, or `ctrl`+`c` to cancel. To reverse the order of the sort, press `R` (shift+r).

Press `/` to bring up the search bar. Type the name (you may need to first delete the `/` if it auto-populates). superfile searches in the current directory and dynamically displays the results. To exit the search bar, press `ctrl`+`c` or `esc`.

Press `.` to show or hide dotfiles.

#### Selection mode

Use selection mode for bulk operations. If you are familiar with Vim, selection mode is similar to Vim's [visual mode](https://vimhelp.org/visual.txt.html#Visual).

Press `v` to toggle between selection mode and normal (browser) mode.

Once in selection mode, you can perform [file operations](#file-operations) on all selected files/folders. [Panel movement](#panel-movement) hotkeys also work in selection mode.

:::tip
The following operations can only be performed while in selection mode. Your current mode is displayed in the lower-right corner of the file panel (Select or Browser).
:::

To make selections, navigate to your file/folder and press `enter` or `L` (shift+l). Press the same key again to deselect.

This may become tedious when you have a large number of items. Instead, you can press `shift`+`up` or `K` (shift+k) to select everything above the cursor. Press `shift`+`down` or `J` (shift+j) to select everything below the cursor.

You can also press `A` (shift+a) to select everything in the current directory.

![selection-mode-demo](../../../assets/demo/selection-mode-demo.gif)

### File operations

:::note
Only copy, cut and delete can be used in selection mode.
:::

Now let's learn how to perform file operations.

Create a new file with `ctrl`+`n`. Type your new file's name and press `enter`. To create a new folder, add `/` to the end of the name.

:::tip
You can create a directory, subdirectory and file in one string. For example:

`directory/subdirectory/filename`
:::

To rename, point your cursor at a file/folder and press `ctrl`+`r`.

To copy, you can press `ctrl`+`c`.

To cut, you can press `ctrl`+`x`.

Both cut and copied items are shown in the clipboard panel (lower-right corner). The progress of your operations is displayed in the processes panel (lower-left corner).

To paste, you can press `ctrl`+`v`.

:::note
In some terminals, for example Windows Powershell, `ctrl`+`v` pastes input from clipboard to terminal. So, `ctrl`+`v` might not work for paste. Either you can add `ctrl`+`w` hotkey for paste, or override default behaviour of `ctrl`+`v` on your terminal.
:::

To delete, you can press `ctrl`+`d`

:::note
The deletion here is not direct deletion, but will be placed in the trash can. However, when you use an external hard drive, it will be deleted directly.
:::

To compress, press `ctrl`+`a`. To decompress, press `ctrl`+`e`.

To open a file with an editor, press `e`.

To open the current directory with an editor, press `E` (shift+e).

To change the default file editor, you can set the `EDITOR` environment variable in your terminal or you can use the `editor` config option (take priority over `EDITOR` environment variable). 
To change the default directory editor, you can use the `dir_editor` config option.
For example:

```bash
EDITOR=nvim
```

This will set Neovim as your default editor. After setting this, Neovim will be used when opening files with the `e` key bindings.

```
editor = "nano"
dir_editor = "vi"
```

These are changes in config file. See [superfile-config](/configure/superfile-config) for more info.
This will set `nano` as your default editor, and `vi` as your default directory editor. After setting this, `nano` will be used when opening files with the `e` key bindings, and `vi` will be used to open current directory with `E` key bindings.

:::caution
If your directory editor does not support opening the current directory with an editor, you may encounter an error when pressing `E`.
:::

![file-operations-demo](../../../assets/demo/file-operations-demo.gif)

### SPF Prompt
#### Shell Mode
Press `:` to open the prompt in shell mode, and execute any shell command in the current directory.
![Prompt-Shell-Mode](../../../assets/git-assets/prompt_shell_mode.png)

:::note
You won't receive any stdout outputs.
For now, this is meant for executing more complex file manipulations via the shell,
rather than handling interactive outputs.
You will be able to see the exit code of the command.
:::

#### SPF Mode
Press `>` to open the prompt in SPF mode. 
![Prompt-SPF-Mode](../../../assets/git-assets/prompt_spf_mode.png)

In this mode, you can execute these spf commands :
- `split` - Open a new panel at a current file panel's path.
- `open <PATH>` - Open a new panel at a specified path.
- `cd <PATH>` - Change directory of current panel.

In this mode, You can substitute shell environment variables via `${}`, shell commands via `$()` and prefix path with `~` to get substituted to home directory 
For example 
- `cd ${HOME}` or `cd ~/xyz`
- `open $(dirname $(which bash))`

Press `esc` or `ctrl`+`c` to exit Prompt.