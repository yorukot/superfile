---
title: Install superfile
description: Let's install superfile to your computer..
head:
  - tag: title
    content: Install superfile | superfile
---

## Requirements

First make sure you have the following tools installed on your machine:

- [Any Nerd-font ](https://www.nerdfonts.com/font-downloads)

:::tip
If you don't install `Nerd font` superfile it will still work, but the UI may be a bit ugly.
:::

## Installation Script

### Linux / MacOs
Just copy and paste this one-line command:

```bash
bash -c "$(curl -sLo- https://superfile.netlify.app/install.sh)"
```
Or wget:
```bash
bash -c "$(wget -qO- https://superfile.netlify.app/install.sh)"
```

## Windows

It actually supports windows! Well.. sort of.

Use powershell to run this command:

```bash
powershell -ExecutionPolicy Bypass -Command "Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://superfile.netlify.app/install.ps1'))"
```
:::note
For uninstall do the same but uninstall.ps1
:::


## Community maintained packages

[![Packaging status](https://repology.org/badge/vertical-allrepos/superfile.svg)](https://repology.org/project/superfile/versions)

> Sort by letter

### Arch

###### Builds package from sources

```bash
sudo pacman -S superfile
```

###### Fetches prebuilt binaries from github

```bash
sudo pacman -S superfile-bin
```

### Homebrew

Install homebrew and execute the following commands

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

### X-CMD

[x-cmd](https://www.x-cmd.com/) is a **toolbox for Posix Shell**, offering a lightweight package manager built using shell and awk.
```sh
x env use superfile
```

## Start superfile

After completing the installation, you can restart the terminal (if necessary)

You can use `spf` to start superfile

```bash
spf
```

## Next-step

- [Tutorial](/getting-started/tutorial)
- [Hotkey list](/list/hotkey-list)
