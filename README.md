<div align="center">

<p>
  <h4>
    <a href="https://ko-fi.com/yorukot">superfile is supported by the community.</a>
  </h4>
<div align="center" markdown="1">
   <sup>Special thanks to:</sup>
   <br>
   <br>
   <a href="https://www.warp.dev/?utm_source=github&utm_medium=referral&utm_campaign=superfile">
      <img alt="Warp sponsorship" width="300" src="/asset/warp.png">
   </a>

### [Warp, the AI terminal for developers](https://www.warp.dev/?utm_source=github&utm_medium=referral&utm_campaign=superfile)
[Available for MacOS, Linux, & Windows](https://www.warp.dev/?utm_source=github&utm_medium=referral&utm_campaign=superfile)<br>

</div>
<hr>

</div>

<div align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="/asset/superfilelogowhite.png" />
  <source media="(prefers-color-scheme: light)" srcset="/asset/superfilelogoblack.png" />
  <img alt="superfile LOGO" src="/asset/superfilelogowhite.png" />
</picture>

[![Go Report Card](https://goreportcard.com/badge/github.com/yorukot/superfile)](https://goreportcard.com/report/github.com/yorukot/superfile)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/yorukot/superfile/refs/heads/main/LICENSE)
[![Discord Link](https://img.shields.io/discord/1338415256875307110?label=discord&logo=discord&logoColor=white)](https://discord.gg/YYtJ23Du7B)
[![Release](https://img.shields.io/github/v/release/yorukot/superfile.svg?style=flat-square)](https://github.com/yorukot/superfile/releases/latest)
[![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/yorukot/superfile?utm_source=oss&utm_medium=github&utm_campaign=yorukot%2Fsuperfile&labelColor=171717&color=FF570A&&label=CodeRabbit+Reviews)](https://www.coderabbit.ai/)

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
- [Notes](#notes)
- [Contributing](#contributing)
- [Troubleshooting](#troubleshooting)
- [Thanks](#thanks)
  - [Support](#Support)
  - [Core maintainer](#core-maintainer)
  - [Contributors](#contributors)
  - [Star History](#star-history)

## Installation

### MacOS and Linux

```bash
bash -c "$(curl -sLo- https://superfile.dev/install.sh)"
```
If you want to inspect the script, see : [install.sh](./website/public/install.sh)

### Windows

#### Powershell
```powershell
powershell -ExecutionPolicy Bypass -Command "Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://superfile.dev/install.ps1'))"
```
If you want to inspect the script, see : [install.ps1](./website/public/install.ps1)

#### [Winget](https://winget.run/)
```powershell
winget install --id yorukot.superfile
```

#### [Scoop](https://scoop.sh/)
```
scoop install superfile
```

### More installation methods
[Click me to check on how to install](https://superfile.dev/getting-started/installation/)

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

### For MacOS/Linux
Run the `build.sh` file:

```bash
./build.sh
```

Add the binary file to your $PATH, e.g., in `/usr/local/bin`:

```bash
sudo mv ./bin/spf /usr/local/bin
```

### For Windows

```bash
go build -o bin/spf.exe
```

Edit System Environment Variables and add superfile repo's `bin` directory to your PATH  

## Start superfile

```bash
spf
```

## Supported Systems

- \[x\] Linux
- \[x\] MacOS
- \[x\] Windows (Not fully supported yet)

## Tutorial

After you install superfile, you can go [here](https://superfile.dev/getting-started/tutorial/) to briefly understand how to use superfile!

## Plugins

[Click me to the plugins wiki](https://superfile.dev/list/plugin-list/)

## Themes

[Click me to the theme wiki](https://superfile.dev/configure/custom-theme/)

## Hotkeys

> [!WARNING]
> If you are vim/nvim user please change your default hotkeys config to vim version!

[**Click me to see the hotkey wiki**](https://superfile.dev/configure/custom-hotkeys/)

## Notes

We have an auto update functionality, that fetches superfile's latest released version from github (if last timestamp of last version check was less than 24 hours) and prints a prompt to user, if there is a newer version available.

You can turn this off, by setting `auto_check_update` to false in superfile config. [**Click me to see the config wiki**](https://superfile.dev/configure/superfile-config/) 

## Troubleshooting

[**Click me to see common problem fix**](https://superfile.dev/troubleshooting/)

## Uninstalling

### MacOS and Linux

On MacOS and Linux, you can uninstall superfile by simply removing the binary. If you installed superfile with sudo, runw

```bash
sudo rm /usr/local/bin/spf
```

If you installed superfile without sudo, run

```bash
rm ~/.local/bin/spf
```

If you don't rember, just try removing both.


### Window

To uninstall superfile on Windows, use this powershell script.

```powershell
powershell -ExecutionPolicy Bypass -Command "Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://superfile.dev/uninstall.ps1'))"
```

## Contributing

If you want to contribute please follow the [contribution guide](./CONTRIBUTING.md)

[**Click me to see changelog**](https://superfile.dev/changelog)

## Thanks

### Support

- a Star on my GitHub repository would be nice 🌟
- You can buy a coffee for me 💖

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G1JEGGC)

### Core maintainer

> We welcome anyone who wants to become a core maintainer. Feel free to reach out!

- **[@yorukot](https://github.com/yorukot)** - Original author and maintainer
- **[@lazysegtree](https://github.com/lazysegtree)** - Core maintainer

### Contributors

**Thanks to all the contributors for making this project even greater!**

<a href="https://github.com/yorukot/superfile/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=yorukot/superfile" />
</a>

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


<div align="center">

## ༼ つ ◕_◕ ༽つ  Please share.

</div>
