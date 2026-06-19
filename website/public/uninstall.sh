#!/bin/bash

green='\033[0;32m'
red='\033[0;31m'
yellow='\033[0;33m'
blue='\033[0;34m'
purple='\033[0;35m'
cyan='\033[0;36m'
white='\033[0;37m'
bright_red='\033[1;31m'
bright_green='\033[1;32m'
bright_yellow='\033[1;33m'
bright_blue='\033[1;34m'
bright_purple='\033[1;35m'
bright_cyan='\033[1;36m'
bright_white='\033[1;37m'
nc='\033[0m' # No Color

echo -e '
\033[0;31m                                                    ______   __  __
\033[1;31m                                                   /      \ /  |/  |
\033[0;33m  _______  __    __   ______    ______    ______  /$$$$$$  |$$/ $$ |  ______
\033[1;33m /       |/  |  /  | /      \  /      \  /      \ $$ |_ $$/ /  |$$ | /      \
\033[0;32m/$$$$$$$/ $$ |  $$ |/$$$$$$  |/$$$$$$  |/$$$$$$  |$$   |    $$ |$$ |/$$$$$$  |
\033[1;32m$$      \ $$ |  $$ |$$ |  $$ |$$    $$ |$$ |  $$/ $$$$/     $$ |$$ |$$    $$ |
\033[0;34m $$$$$$  |$$ \__$$ |$$ |__$$ |$$$$$$$$/ $$ |      $$ |      $$ |$$ |$$$$$$$$/
\033[1;34m/     $$/ $$    $$/ $$    $$/ $$       |$$ |      $$ |      $$ |$$ |$$       |
\033[0;35m$$$$$$$/   $$$$$$/  $$$$$$$/   $$$$$$$/ $$/       $$/       $$/ $$/  $$$$$$$/
\033[1;35m                    $$ |
\033[0;31m                    $$ |
\033[1;31m                    $$/
'

found=0
failed=0

if [ -z "${HOME:-}" ] || [ "$HOME" = "/" ]; then
    echo -e "${red}❌ Refusing to run uninstall with an unsafe HOME value.${nc}"
    exit 1
fi

# Resolve XDG dirs, falling back to spec-defined defaults
XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"
XDG_DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}"
XDG_CACHE_HOME="${XDG_CACHE_HOME:-$HOME/.cache}"
XDG_STATE_HOME="${XDG_STATE_HOME:-$HOME/.local/state}"

# Remove binary from /usr/local/bin
if [ -f /usr/local/bin/spf ]; then
    found=1
    echo -e "${bright_yellow}Removing ${cyan}/usr/local/bin/spf${bright_yellow}...${nc}"
    if ! sudo rm /usr/local/bin/spf; then
        echo -e "${red}❌ Failed to remove ${white}/usr/local/bin/spf${red}. Do you have sudo permissions?${nc}"
        failed=1
    else
        echo -e "${bright_green}✔ Removed ${white}/usr/local/bin/spf${nc}"
    fi
fi

# Remove binary from ~/.local/bin
if [ -f "$HOME/.local/bin/spf" ]; then
    found=1
    echo -e "${bright_yellow}Removing ${cyan}~/.local/bin/spf${bright_yellow}...${nc}"
    if ! rm "$HOME/.local/bin/spf"; then
        echo -e "${red}❌ Failed to remove ${white}~/.local/bin/spf${nc}"
        failed=1
    else
        echo -e "${bright_green}✔ Removed ${white}~/.local/bin/spf${nc}"
    fi
fi

# Remove config directory
if [ -d "$XDG_CONFIG_HOME/superfile" ]; then
    found=1
    echo -e "${bright_yellow}Removing ${cyan}$XDG_CONFIG_HOME/superfile${bright_yellow}...${nc}"
    if ! rm -rf "$XDG_CONFIG_HOME/superfile"; then
        echo -e "${red}❌ Failed to remove ${white}$XDG_CONFIG_HOME/superfile${nc}"
        failed=1
    else
        echo -e "${bright_green}✔ Removed ${white}$XDG_CONFIG_HOME/superfile${nc}"
    fi
fi

# Remove data directory
if [ -d "$XDG_DATA_HOME/superfile" ]; then
    found=1
    echo -e "${bright_yellow}Removing ${cyan}$XDG_DATA_HOME/superfile${bright_yellow}...${nc}"
    if ! rm -rf "$XDG_DATA_HOME/superfile"; then
        echo -e "${red}❌ Failed to remove ${white}$XDG_DATA_HOME/superfile${nc}"
        failed=1
    else
        echo -e "${bright_green}✔ Removed ${white}$XDG_DATA_HOME/superfile${nc}"
    fi
fi

# Remove cache directory
if [ -d "$XDG_CACHE_HOME/superfile" ]; then
    found=1
    echo -e "${bright_yellow}Removing ${cyan}$XDG_CACHE_HOME/superfile${bright_yellow}...${nc}"
    if ! rm -rf "$XDG_CACHE_HOME/superfile"; then
        echo -e "${red}❌ Failed to remove ${white}$XDG_CACHE_HOME/superfile${nc}"
        failed=1
    else
        echo -e "${bright_green}✔ Removed ${white}$XDG_CACHE_HOME/superfile${nc}"
    fi
fi

# Remove state directory
if [ -d "$XDG_STATE_HOME/superfile" ]; then
    found=1
    echo -e "${bright_yellow}Removing ${cyan}$XDG_STATE_HOME/superfile${bright_yellow}...${nc}"
    if ! rm -rf "$XDG_STATE_HOME/superfile"; then
        echo -e "${red}❌ Failed to remove ${white}$XDG_STATE_HOME/superfile${nc}"
        failed=1
    else
        echo -e "${bright_green}✔ Removed ${white}$XDG_STATE_HOME/superfile${nc}"
    fi
fi

# Remove Application Support directory (macOS)
if [ -d "$HOME/Library/Application Support/superfile" ]; then
    found=1
    echo -e "${bright_yellow}Removing ${cyan}~/Library/Application Support/superfile${bright_yellow}...${nc}"
    if ! rm -rf "$HOME/Library/Application Support/superfile"; then
        echo -e "${red}❌ Failed to remove ${white}~/Library/Application Support/superfile${nc}"
        failed=1
    else
        echo -e "${bright_green}✔ Removed ${white}~/Library/Application Support/superfile${nc}"
    fi
fi

# Remove cache directory (macOS)
if [ -d "$HOME/Library/Caches/superfile" ]; then
    found=1
    echo -e "${bright_yellow}Removing ${cyan}~/Library/Caches/superfile${bright_yellow}...${nc}"
    if ! rm -rf "$HOME/Library/Caches/superfile"; then
        echo -e "${red}❌ Failed to remove ${white}~/Library/Caches/superfile${nc}"
        failed=1
    else
        echo -e "${bright_green}✔ Removed ${white}~/Library/Caches/superfile${nc}"
    fi
fi

if [ "$found" -eq 0 ]; then
    echo -e "${yellow}No superfile installation found. Nothing to remove.${nc}"
elif [ "$failed" -eq 1 ]; then
    echo -e "\n${red}⚠ Uninstall completed with errors. Please review the messages above.${nc}"
    exit 1
else
    echo -e "\n👋 ${bright_green}superfile has been uninstalled.${nc}"
fi
