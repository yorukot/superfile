---
title: Learn how to use tutorial
description: Quickly get started with Superfile
head:
  - tag: title
    content: Tutorial | superfile
---

This tutorial will help you learn to use superfile step by step.

:::caution
If you didn't install superfile yet please [click here](/getting-started/installation) to install.
:::

:::tip
If you just want to know about Hotkey without so much supplementary content, please go to [here](/list/hotkey-list)
:::

## Hotkeys tutorial
First, if you want to open `superfile` by opening a terminal and typing `spf`.

Afterwards if you want to exit just press `q` or `esc`.

![demo](https://github.com/yorukot/superfile/assets/107802416/ddd9f05c-b39b-4f55-838b-d248c845a589)

### Navigation
You can put focus on the sidebar by pressing `s`.

Press `p` to focus the processbar.

Press `m` to focus on metadata.

:::tip
The size of the folder will only be obtained when you focus on metadata.

If you want to get more detailed metadata, you can install the metadata plugin.
:::
If you want to return to the file panel, just press again to remove focus.

![demo](https://github.com/yorukot/superfile/assets/107802416/ec7062ce-1884-4395-b68b-e0546c8a02de)

### File panel navigation
Now you might be thinking that a file panel is not enough. 

Therefore, you can press `n` to create a new file panel and `w` to close the file panel.

Now you know how to create and close file panels.

But how to switch to the next or previous archive panel?

You can press `tab` or `L` (shift+l) to move to the next file panel.

Then press `shift+left` or `H` (shift+h) to move to the previous archive panel.

![demo](https://github.com/yorukot/superfile/assets/107802416/2c2a7632-c5c0-43b6-80a7-d6e21fcf63b1)

### File panel movement

Now let us introduce how to operate the File panel

First of all, if you don’t want to see dotfiles, you can press `.` which will hide all dotfiles

Then if you think you will use this folder frequently, you can also put it on the sidebar, just press `P` to pin or unpin

When you focus on the file panel you can press `up` or `k` to up
press `down` or `j` to down

After navigation to the file or folder you want, you can press `enter` or `l` to confirm. The file will be opened using your default Application (if not, there will be no response) and the folder will be entered.

Press `h` or `backspace` will return to the parent directory.

Pressing `o` will bring up a menu to choose how you would like the panel to sort the files. You can choose between `Name`, `Size`, or `Date Modified`. `enter` to select, and `esc`, `o`, or `ctrl+c` to cancel.

To reverse the order of the sort, press `R` (`shift+r`)

If you have a large number of files, you can also use `/` to search,After entering the Key you want, you can press `/` again or `enter`

If you want to clear the current search, you can press `ctrl+c` or `esc`

![demo](https://github.com/yorukot/superfile/assets/107802416/f6fd9e4e-f73f-4848-a113-416732abf126)

### File selection mode movement

You might be thinking what is selection mode?

That's really easy!This mode is similar to Vim's Visual.But select file or folder instead of code.

After selecting, you will be able to perform [file operations](#file-operations) on all selected files or folders.

First to enter this mode you can press `v`. To return to browse mode the same is done by pressing `v`

:::tip
The following operations can only be performed in Select mode,You can see the current mode in the lower right corner of the file panel
:::

After entering you can use the same shortcut as [file panel movement](#file-panel-movement) to move.

Now you may have moved to the file or folder you want to select. You can press `enter` or `L` (shift+l) to select and press again to deselect.

But this will be a bit slow when you want to select a large number of files

So here are some faster ways

You can press `shift+up` or `K` (shift+k) to select all files or folders passed by the cursor when it goes up.

Of course, the same is true for `shift+down` or `J` (shift+j)

You can also press `A` to select all folders in the current directory

![demo](https://github.com/yorukot/superfile/assets/107802416/4306fd31-04e0-456c-b1f2-3923e8d932e1)

### File operations

:::note
Only copy, cut and delete can be used in selection mode
:::

You have learned how to use superfile to browse files and select files. Let's learn how to perform file operations!

First, let me teach you how to create a file. You can press `ctrl+n` to create a file or folder, if you want to create folder you need add `/` in the end.

Then if you want to rename it, press `ctrl+r` and it will name the location of your cursor.

If you want to copy, you can press `ctrl+c` and the copied file list will be displayed in the clipboard (lower right corner).

If you want to cut, press `ctrl+x`.

:::tip
The copy here will actually be copied to the clipboard of your system.
:::

Your copy process will be displayed in the processbar (lower left corner).

You can press `ctrl+d` to delete file (The deletion here is not direct deletion but will be placed in the trash can.). But when you use an external hard drive, it will be deleted directly.

If you want to decompress or compress you can press `ctrl+a` to compress and `ctrl+e` to decompress.

To open a file with an editor, press `e`.

To open the current directory with an editor, press `E`.

To change the default editor, you can set the `EDITOR` environment variable in your terminal. For example:

```bash
EDITOR=nvim
```


This will set Neovim as your default editor. After setting this, the specified editor will be used when opening files with the `e` or `E` key bindings.

:::caution
If your editor does not support opening the current directory with an editor, you may encounter an error when pressing `E`.
:::

(Sorry, this video has a little bit of lag)
[demo video](https://github.com/yorukot/superfile/assets/107802416/d0770b3f-025e-40c9-ad3f-8b2adaf1c6c5)

