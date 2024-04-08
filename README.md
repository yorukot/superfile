<div align="center">

# SUPERFILE

![](/asset/demo.png)

</div>

## Demo

| Perform common operations |
| ------------------------- |
| ![](/asset/demo.gif)      |

## Content

- [Features](#features)
- [Installation](#install)
  - [Homebrew](#homebrew)
  - [Linux](#linux)
  - [Font](#font)
- [Build](#build)
- [Supported Systems](#supported-systems)
- [Themes](#themes)
  - [Use an existing theme](#use-an-existing-theme)
  - [Create your own theme](#create-your-own-theme)
- [Hotkeys](#hotkeys)
- [Contribute](#contribute)
  - [Share your idea](#share-your-idea)
  - [Bug report](#bug-report)
- [Todo list](#todo-list)
- [Star History](#star-history)

## Features

- Fancy gui
- Fully customizable
- Vim's selection mode
- Easy to use
- Trash can
- Metadata detil
- Copy file to the clipboard
- Copy and paste file
- Auto rename file or folder when duplicate
- Rename files in a modern way
- Open file with default app
- Open terminal with current path

## Install

> I am still working on different installation methods like `homebrew` and `snap`

**Requirements**

- [`Exiftool`](#exiftool)
- Any [`Nerd Font`](#font)

### Homebrew

Download [this homebrew file](https://github.com/MHNightCat/superfile/blob/main/superfile.rb) and enter the following in your terminal:

```bash
brew install ~/Download/superfile.rb
```

### Linux

You can go to the [latest release](https://github.com/MHNightCat/superfile/releases/latest) and download the binary file. Once it is downloaded enter the following in your terminal:

```bash
cd ~/Download
chmod +x ./spf
sudo mv ./spf /bin/
```

### Exiftool

[`exiftool`](https://github.com/exiftool/exiftool) is a tool used to obtain file metadata. If it is not installed, it will cause errors.

**Installation:**

```bash
# Homebrew:
brew install exiftool

# Fedora:
sudo dnf install perl-Image-ExifTool

# Ubuntu:
sudo apt install exiftool

# Archlinux:
sudo pacman -S perl-image-exiftool
```

### Font

> WARNING: This is a reminder that you must use a [Nerd font](https://www.nerdfonts.com/font-downloads)

Once the font is installed if `superfile` isn't working make sure to update your terminal preferences to use the font.

## Build

You can build the source code yourself by using these steps:

**Requirements**

- [golang](https://go.dev/doc/install)

**Build Steps**

Clone this repo using the following command:

```
git clone https://github.com/MHNightCat/superfile.git
```

Enter the downloaded directory:

```bash
cd superfile
```

Run the `build.sh` file:

```bash
sh build.sh
```

Move the binary file to /bin (on Linux):

```bash
mv ./bin/spf /bin
```

or on OSX:

```bash
mv ./bin/spf /usr/local/bin
```

## Supported Systems

- \[x\] Linux
- \[x\] MacOS
- \[ \] Windows

## Themes

### Use an existing theme

You can go to [theme list](https://github.com/MHNightCat/superfile/blob/main/THEMELIST.md) to find one you like!

> We only have a few themes at the moment, but we will be making more over time! You can also [submit a pull request](https://github.com/MHNightCat/superfile/pulls) for your own theme!

Edit config.json using `Nano`:

```
nano ~/.superfile/config/config.json
```

Edit config.json using `Vim`:

```
vim ~/.superfile/config/config.json
```

then change:

```json
"theme": "gruvbox",
```

to:

```json
"theme": "theme_name",
```

### Create your own theme

If you want to customize your own theme, you can go to `~/.superfile/theme/YOUR_THEME_NAME.json` and copy the existing theme's json to your own theme file

Now you can customize your own theme!!

And if you complete your theme you can change:

```json
"theme": "gruvbox",
```

to:

```json
"theme": "YOUR_THEME_NAME",
```

[If you are satisfied with your theme, you might as well put it into the default theme list!](#contribute)

## Hotkeys

[**Click me to see the hotkey list**](https://github.com/MHNightCat/superfile/blob/main/HOTKEYS.md)

**You can change all hotkeys in** `~/.superfile/config/config.json`

Edit config.json using `Nano`:

```
nano ~/.superfile/config/config.json
```

Edit config.json using `Vim`:

```
vim ~/.superfile/config/config.json
```

> "Normal mode" is the default browsing mode

Global hotkeys cannot conflict with other hotkeys (The only exception is the special hotkey).

The hotkey ranges are found in `config.json`

## Contribute

[**Click me to learn how to contribute**](https://docs.github.com/en/get-started/exploring-projects-on-github/contributing-to-a-project)

> For example, add your custom themes to `/themes` and submit a pull request

### Share your idea

[**I have some ideas but i don't know how to contribute**](https://github.com/MHNightCat/superfile/discussions)

> We welcome anyone with any ideas about superfile!!

### Bug report

[**Submit a bug report here**](https://github.com/MHNightCat/superfile/issues)

## Todo list 

**I will do my best to complete this list haha**

- \[x\] Auto init config file
- \[ \] Extract files
- \[ \] Open terminal in the focused file panel location
- \[ \] Open file with enter key
- \[ \] File panel search / filter
- \[ \] Add help bar down below bottom bar
- \[ \] Compress files
- \[ \] Can cancel the progress of the process bar
- \[ \] Undo function
- \[ \] Auto clear trash can
- \[ \] AES encryption and decryption
- \[ \] Add more theme

#### 1.2

- [ ] AES encryption and decryption
- [ ] Auto clear trash can
- [ ] Can cancel the progress of the process bar

## Star History

**THANKS FOR All OF YOUR STARS!**

<a href="https://star-history.com/#MHNightCat/superfile&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Date" />
 </picture>
</a>
