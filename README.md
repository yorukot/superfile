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

- [Installation](#installation)
- [Build](#build)
- [Supported Systems](#supported-systems)
- [Tutorial](#tutorial)
- [Plugins](#plugins)
- [Themes](#themes)
- [Hotkeys](#hotkeys)
- [Contributing](#contributing)
- [Troubleshooting](#troubleshooting)
- [Thanks](#thanks)
  - [Support](#Support)
  - [Contributors](#contributors)
  - [Star History](#star-history)

## Installation

Quick install (Support MacOs and linux)

```bash
bash -c "$(wget -qO- https://raw.githubusercontent.com/yorukot/superfile/main/install.sh)"
```

### More installation methods
[Click me to check on how to install](https://superfile.netlify.app/getting-started/installation/)

## Build

You can build the source code yourself by using these steps:

**Requirements**

- [golang](https://go.dev/doc/install)

**Build Steps**

Clone this repository using the following command:

```
git clone https://github.com/yorukot/superfile.git --depth=1
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
sudo mv ./bin/spf /usr/local/bin
```

## Supported Systems

- \[x\] Linux
- \[x\] MacOS
- \[x\] Windows (Not fully supported yet)

## Tutorial

After you install superfile, you can go [here](https://superfile.netlify.app/getting-started/tutorial/) to briefly understand how to use superfile!

## Plugins

[Click me to the plugins wiki](https://github.com/yorukot/superfile/wiki/Plugins)

## Themes

[Click me to the theme wiki](https://superfile.netlify.app/configure/custom-theme/)

## Hotkeys

> [!WARNING]
> If you are vim/nvim user please change your default hotkeys config to vim version!

[**Click me to see the hotkey wiki**](https://superfile.netlify.app/configure/custom-hotkeys/)

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
