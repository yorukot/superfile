---
title: Install superfile
description: Let's install superfile to your computer..
head:
  - tag: title
    content: Install superfile | superfile
---

## Before install

First make sure you have the following tools installed on your machine:

- [Any Nerd-font ](https://www.nerdfonts.com/font-downloads), and set the font for your terminal application to use the installed Nerd-font

:::tip
If you don't install `Nerd font`, superfile will still work, but the UI may look a bit off. It's recommended to disable the Nerd font option to avoid this issue.
:::

## Installation Scripts

Copy and paste the following one-line command into your machine's terminal.

### Linux / MacOs

With `curl`:

```bash
bash -c "$(curl -sLo- https://superfile.netlify.app/install.sh)"
```

Or with `wget`:
```bash
bash -c "$(wget -qO- https://superfile.netlify.app/install.sh)"
```

Use `SPF_INSTALL_VERSION` to specify a version :

```bash
SPF_INSTALL_VERSION=1.2.1 bash -c "$(curl -sLo- https://superfile.netlify.app/install.sh)"
```

### Windows

With `powershell`:

```bash
powershell -ExecutionPolicy Bypass -Command "Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://superfile.netlify.app/install.ps1'))"
```

:::note
To uninstall, run the above `powershell` command with the modified URL:

`https://superfile.netlify.app/uninstall.ps1`
:::

Use `SPF_INSTALL_VERSION` to specify a version :

```bash
powershell -ExecutionPolicy Bypass -Command "$env:SPF_INSTALL_VERSION=1.2.1; Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://superfile.netlify.app/install.ps1'))"
```

With [Winget](https://winget.run/):

```powershell
winget install superfile
``````

With [Scoop](https://scoop.sh/):

```bash
scoop install superfile
```

## Community maintained packages

[![Packaging status](https://repology.org/badge/vertical-allrepos/superfile.svg)](https://repology.org/project/superfile/versions)

> Sort by letter

### Arch

###### Builds package from sources

```bash
sudo pacman -S superfile
```

###### Builds most recent version from GitHub

```bash
yay -S superfile-git
```

### Homebrew

Install [Homebrew](https://brew.sh/) and then run the following command:

```bash
brew install superfile
```

### NixOS

###### Install with nix command-line

```bash
nix profile install github:yorukot/superfile#superfile
```

###### Install with flake

Add superfile to your flake inputs:

```nix
inputs = {
  superfile = {
    url = "github:yorukot/superfile";
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

### Pixi

Install [Pixi](https://pixi.sh/latest/) and then run the following command:

```bash
pixi global install superfile
```

### X-CMD

[x-cmd](https://www.x-cmd.com/) is a **toolbox for Posix Shell**, offering a lightweight package manager built using shell and awk.
```sh
x env use superfile
```

## Start superfile

After completing the installation, you can restart the terminal (if necessary).

Run `spf` to start superfile

```bash
spf
```

## Next steps

- [Tutorial](/getting-started/tutorial)
- [Hotkey list](/list/hotkey-list)
