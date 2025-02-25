<div align="center">

<p>
  <H4>
    <a href="https://ko-fi.com/yorukot">superfile is supported by the community.</a>
  </H3>
</p>
<sup>Special thanks to our sponsor:</sup>
<br>
<br>
<a href="https://www.warp.dev/?utm_source=github&utm_medium=referral&utm_campaign=superfile">
  <div>
    <picture>
      <img alt="Warp" width="500" src="/asset/warp.png">
    </picture>
  </div>
  <b>Warp, the intelligent terminal</b>
  <div>
    <sup>superfile's Terminal of Choice<br>
Visit warp.dev to learn more
    </sup>
  </div>
</a>
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

### MacOS and Linux

```bash
bash -c "$(curl -sLo- https://superfile.netlify.app/install.sh)"
```
If you want to inspect the script, see : [install.sh](./website/public/install.sh)

### Windows

#### Powershell
```powershell
powershell -ExecutionPolicy Bypass -Command "Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://superfile.netlify.app/install.ps1'))"
```
If you want to inspect the script, see : [install.ps1](./website/public/install.ps1)

#### [Scoop](https://scoop.sh/)
```
scoop install superfile
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

## Start Superfile

```bash
spf
```

## Supported Systems

- \[x\] Linux
- \[x\] MacOS
- \[x\] Windows (Not fully supported yet)

## Tutorial

After you install superfile, you can go [here](https://superfile.netlify.app/getting-started/tutorial/) to briefly understand how to use superfile!

## Plugins

[Click me to the plugins wiki](https://superfile.netlify.app/list/plugin-list/)

## Themes

[Click me to the theme wiki](https://superfile.netlify.app/configure/custom-theme/)

## Hotkeys

> [!WARNING]
> If you are vim/nvim user please change your default hotkeys config to vim version!

[**Click me to see the hotkey wiki**](https://superfile.netlify.app/configure/custom-hotkeys/)

## Troubleshooting

[**Click me to see common problem fix**](https://superfile.netlify.app/troubleshooting/)

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
powershell -ExecutionPolicy Bypass -Command "Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://superfile.netlify.app/uninstall.ps1'))"
```

## Contributing

If you want to contribute please follow the [contribution guide](./CONTRIBUTING.md)

## Thanks

### Support

- a Star on my GitHub repository would be nice 🌟
- You can buy a coffee for me 💖

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G1JEGGC)

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
