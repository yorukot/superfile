<div align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/asset/superfilelogowhite.png" />
  <source media="(prefers-color-scheme: light)" srcset="/asset/superfilelogoblack.png" />
  <img alt="superfile LOGO" src="/asset/superfilelogowhite.png" />
</picture>

![](/asset/demo.png)

</div>

## Demo

| Perform common operations |
| ------------------------- |
| ![](/asset/demo.gif)      |

## Content

- [Installation](#install)
  - [Homebrew](#homebrew)
  - [Linux](#linux)
  - [Font](#font)
- [Build](#build)
- [Supported Systems](#supported-systems)
- [Tutorial](#tutorial)
- [Plugins](#plugins)
- [Themes](#themes)
  - [Use an existing theme](#use-an-existing-theme)
  - [Create your own theme](#create-your-own-theme)
- [Hotkeys](#hotkeys)
- [Contributing](#contributing)
- [Todo list](#todo-list)
- [Star History](#star-history)

## Install

> I am still working on different installation methods like `homebrew` and `snap`

**Requirements**

- Any [`Nerd Font`](#font)

### Homebrew

If you want to use homebrew please install `go` first!

```bash
brew tap mhnightcat/superfile https://github.com/MHNightCat/homebrew-superfile.git && brew install superfile
```

### Linux

You can go to the [latest release](https://github.com/MHNightCat/superfile/releases/latest) and download the binary file. Once it is downloaded please excrate the file after that enter the following in your terminal:

```bash
cd ~/Download
chmod +x ./spf
sudo mv ./spf /bin/
```

<h4>
     NixOS
</h4>

<details><summary>Click to expand</summary>
<p>

Add superfile to your flake inputs:

```nix
inputs = {
  superfile = {
    url = "github:MHNightCat/superfile";
  };
  # ...
};
```

Then you can add it to your packages:

```nix
let
  system = "x86_64-linux";
in {
  environment.systemPackages = with pkgs; [
    # ...
    inputs.superfile.packages.${system}.default  ];
}
```

</details>

### Font

> [!WARNING] 
> This is a reminder that you must use a [Nerd font](https://www.nerdfonts.com/font-downloads)

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
./build.sh
```

Add the binary file to your $PATH, e.g. in `/usr/local/bin`:

```bash
mv ./bin/spf /usr/local/bin
```

## Supported Systems

- \[x\] Linux
- \[x\] MacOS
- \[ \] Windows

## Tutorial

After you install superfile, you can go [here](https://github.com/MHNightCat/superfile/wiki/Tutorial) to briefly understand how to use superfile!

## Plugins

[Click me to the plugins wiki](https://github.com/MHNightCat/superfile/wiki/Plugins)

## Themes

### Use an existing theme

You can go to [theme list](https://github.com/MHNightCat/superfile/blob/main/THEMELIST.md) to find one you like!

> We only have a few themes at the moment, but we will be making more over time! You can also [submit a pull request](https://github.com/MHNightCat/superfile/pulls) for your own theme!

copy `theme_name` in:

```
Theme name: theme_name
```

Edit `config.json` using your preferred editor:

```
$EDITOR ~/.config/superfile/config.json
```

and change:

```toml
theme = "gruvbox"
```

to:

```toml
theme = "theme-name"
```

### Create your own theme

If you want to customize your own theme, you can go to `~/.config/superfile/theme/YOUR_THEME_NAME.json` and copy the existing theme's json to your own theme file

Don't forget to change the `theme` variable in `config.json` to your theme name.

[If you are satisfied with your theme, you might as well put it into the default theme list!](#contribute)

## Hotkeys

[**Click me to see the hotkey list**](https://github.com/MHNightCat/superfile/wiki/Hotkey-list)

**You can change all hotkeys in** `~/.config/superfile/hotkeys.toml`

> "Normal mode" is the default browsing mode

Global hotkeys cannot conflict with other hotkeys (The only exception is the special hotkey).

The hotkey ranges are found in `hotkeys.toml`

## Contributing

If you want to contribute please follow the [contribution guide](./CONTRIBUTING.md)

## Todo list

See the todo list in [here](https://github.com/MHNightCat/superfile/blob/main/TODOLIST.md)

## Star History

**THANKS FOR All OF YOUR STARS!**

<a href="https://star-history.com/#MHNightCat/superfile&Timeline">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Timeline&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Timeline" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Timeline" />
 </picture>
</a>
