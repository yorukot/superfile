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
- [Install](#install)
  - [Font](#font)
- [Support system](#support-system)
- [Themes](#themes)
  - [Use an existing theme](#use-an-existing-theme)
  - [Completely customize your theme](#completely-customize-your-theme)
- [Hotkey](#hotkey)
- [Contribute](#contribute)
  - [Share your idea](#share-your-idea)
  - [Bug report](#bug-report)
  - [Share your themes](#share-your-themes)
- [Todo list](#todo-list)
- [Star History](#star-history)
## Features

- Fancy gui
- Fully customizable
- Vim's selection mode
- Easy to use 
- Trash can

## Install

> I am still try to make more install method! Like `HomeBrew` or `snap`

You can go to [latest release](https://github.com/MHNightCat/superfile/releases/latest) and download binary file

cd to download and move binary to bin after that please install [font](#font)
```bash
cd ~/Download
chmod +x ./spf
sudo mv ./spf /bin/
```

### Font

You **MUST** install [nerd font](https://www.nerdfonts.com/font-downloads)

If after install it still not working
Please check your terminal preference setting

## Support system

- [X] Linux  
- [X] Macos - Not tested, but should work
- [ ] Windows - Not support

## Themes

### Use an existing theme
You can go to [theme list](https://github.com/MHNightCat/superfile/blob/main/THEMELIST.md) to find which is you liked!

> We only have a few themes at the moment, we will be making more over the next time! [Or you can public your own theme](https://github.com/MHNightCat/superfile/pulls)!

and editor `~/.superfile/config/config.json`

Edit config.json using `Nano`:

```
nano ~/.superfile/config/config.json
```

Edit config.json using `Vim`:

```
vim ~/.superfile/config/config.json
```

change 

```json
"theme": "gruvbox",
```
to

```json
"theme": "you like theme name",
```

### Completely customize your theme

If you want to customize your own theme, you can go to `~/.superfile/theme/YOUR_THEME_NAME.json`
and copy the existing themes json to your own theme file

Now you can customize your own theme!!

And if you complete your theme you can change

```json
"theme": "gruvbox",
```
to

```json
"theme": "YOUR_THEME_NAME",
```

[If you are satisfied with your theme, you might as well put it into the default theme list!](#contribute)

## Hotkey

[**Click me to watch the hotkey list**](https://github.com/MHNightCat/superfile/blob/main/HOTKEY.md)

**You can change the all hotkey in** `~/.superfile/config/config.json`

Edit config.json using `Nano`:

```
nano ~/.superfile/config/config.json
```

Edit config.json using `Vim`:

```
vim ~/.superfile/config/config.json
```

> normal mode mean browser mode

All global hotkeys cannot conflict with other hotkeys(Except special hotkey).

The hotkey ranges I wrote in config.json

## Contribute

[**Click me to learn how to contribute**](https://docs.github.com/en/get-started/exploring-projects-on-github/contributing-to-a-project)

### Share your idea
[**I have some idea but i don't know how to contribute**](https://github.com/MHNightCat/superfile/discussions)

> We welcome anyone with any ideas about superfile!!

### Bug report

[**Bug report in here~**](https://github.com/MHNightCat/superfile/issues)

### Share your themes

Same as contribution. Just put your own theme into `/themes`
and create a pull request!

If you really want to share your theme but you don't know how to do it
You can go to [here](https://github.com/MHNightCat/superfile/discussions/new?category=theme) create a discussion and i will help you(if i have time)


## Todo list

**I hope i can complete all this todo list haha**

- [x] Auto init config file
- [ ] Compress and decompress files
- [ ] Can cancel the progress of the process bar
- [ ] Undo function
- [ ] Auto clear trash can
- [ ] AES encryption and decryption

## Star History

**THANKS FOR All OF YOUR STAR!**

## Star History

<a href="https://star-history.com/#MHNightCat/superfile&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=MHNightCat/superfile&type=Date" />
 </picture>
</a>
