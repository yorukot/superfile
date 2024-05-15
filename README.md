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
  - [Install pre-built binaries](#install-pre-built-binaries)
  - [Windows](#Windows)
  - [NixOs](#nixos)
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
- [Troubleshooting](#troubleshooting)
- [Thanks](#thanks)
  - [Support](#Support)
  - [Contributors](#contributors)
  - [Star History](#star-history)

## Installation

**Requirements**

- Any [`Nerd Font`](#font)

### Homebrew

Install homebrew and execute the following commands

```bash
brew install superfile
```

### Install pre-built binaries
**Just copy and paste this one-line command:**

```bash
bash -c "$(curl -sLo- https://raw.githubusercontent.com/MHNightCat/superfile/main/install.sh)"
```

Or wget:

```bash
bash -c "$(wget -qO- https://raw.githubusercontent.com/MHNightCat/superfile/main/install.sh)"
```

### Windows

It actually supports windows! Well.. sort of.

Use powershell to run this command:

```powershell
powershell -ExecutionPolicy Bypass -Command "Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/MHNightCat/superfile/main/install.ps1'))"

```
For uninstall do the same but uninstall.ps1

### NixOS

<details><summary>Click to expand</summary>
<p>

#### Install with nix command-line:

```bash
nix profile install github:MHNightCat/superfile#superfile
```

#### Install with flake:

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

After installed, type "spf" to open superfile.

## Build

You can build the source code yourself by using these steps:

**Requirements**

- [golang](https://go.dev/doc/install)

**Build Steps**

Clone this repository using the following command:

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

Add the binary file to your $PATH, e.g., in `/usr/local/bin`:

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

> We only have a few themes at the moment, but we will be making more overtime! You can also [submit a pull request](https://github.com/MHNightCat/superfile/pulls) for your own theme!

copy `theme_name` in:

```
Theme name: theme_name
```

Edit `config.toml` using your preferred editor:

> [!TIP]
> If your OS is macOS the file path should be in the `~/Library/Application Support/superfile/config.toml`

```
$EDITOR ~/.config/superfile/config.toml
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

If you want to customize your own theme, you can go to `~/.config/superfile/theme/YOUR_THEME_NAME.toml` and copy the existing theme's json to your own theme file

Don't forget to change the `theme` variable in `config.toml` to your theme name.

[If you are satisfied with your theme, you might as well put it into the default theme list!](#contribute)

## Hotkeys

[**Click me to see the hotkey list**](https://github.com/MHNightCat/superfile/wiki/Hotkey-list)

> [!TIP]
> If your OS is macOS the file path should be in the `~/Library/Application Support/superfile/hotkeys.toml`

**You can change all hotkeys in** `~/.config/superfile/hotkeys.toml`

> "Normal mode" is the default browsing mode

Global hotkeys cannot conflict with other hotkeys (The only exception is the special hotkey).

The hotkey ranges are found in `hotkeys.toml`

## Troubleshooting

[**Click me to see common problem fix**](https://github.com/MHNightCat/superfile/wiki/Troubleshooting)

## Contributing

If you want to contribute please follow the [contribution guide](./CONTRIBUTING.md)

## Thanks

### Support

- a Star on my GitHub repository would be nice 🌟
- You can buy a coffee for me 💖

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G1JEGGC)

### Contributors

**Thanks to all the contributors for making this project even greater!**

[![contributors](/asset/contributors.svg)](https://github.com/mhnightcat/superfile/graphs/contributors)

### Star History

**THANKS FOR All OF YOUR STARS!**
Your stars are my motivation to keep updating!

<a href="https://star-history.com/#MHNightCat/superfile&Timeline">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Timeline&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Timeline" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Timeline" />
 </picture>
</a>
