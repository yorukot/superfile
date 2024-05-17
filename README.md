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

[You can go here to know how to install](https://superfile.netlify.app/getting-started/installation/)

## Build

You can build the source code yourself by using these steps:

**Requirements**

- [golang](https://go.dev/doc/install)

**Build Steps**

Clone this repository using the following command:

```
git clone https://github.com/yorukot/superfile.git
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
- \[x\] Windows (Not fully supported yet)

## Tutorial

After you install superfile, you can go [here](https://github.com/yorukot/superfile/wiki/Tutorial) to briefly understand how to use superfile!

## Plugins

[Click me to the plugins wiki](https://github.com/yorukot/superfile/wiki/Plugins)

## Themes

### Use an existing theme

You can go to [theme list](https://github.com/yorukot/superfile/blob/main/THEMELIST.md) to find one you like!

> We only have a few themes at the moment, but we will be making more overtime! You can also [submit a pull request](https://github.com/yorukot/superfile/pulls) for your own theme!

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

[**Click me to see the hotkey list**](https://github.com/yorukot/superfile/wiki/Hotkey-list)

> [!TIP]
> If your OS is macOS the file path should be in the `~/Library/Application Support/superfile/hotkeys.toml`

**You can change all hotkeys in** `~/.config/superfile/hotkeys.toml`

> "Normal mode" is the default browsing mode

Global hotkeys cannot conflict with other hotkeys (The only exception is the special hotkey).

The hotkey ranges are found in `hotkeys.toml`

## Troubleshooting

[**Click me to see common problem fix**](https://github.com/yorukot/superfile/wiki/Troubleshooting)

## Contributing

If you want to contribute please follow the [contribution guide](./CONTRIBUTING.md)

## Thanks

### Support

- a Star on my GitHub repository would be nice 🌟
- You can buy a coffee for me 💖

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G1JEGGC)

### Contributors

**Thanks to all the contributors for making this project even greater!**

[![contributors](/asset/contributors.svg)](https://github.com/yorukot/superfile/graphs/contributors)

### Star History

**THANKS FOR All OF YOUR STARS!**
Your stars are my motivation to keep updating!

<a href="https://star-history.com/#yorukot/superfile&Timeline">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=yorukot/superfile&type=Timeline&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=yorukot/superfile&type=Timeline" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=yorukot/superfile&type=Timeline" />
 </picture>
</a>
