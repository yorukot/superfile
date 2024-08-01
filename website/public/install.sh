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


temp_dir=$(mktemp -d)
if [ $? -ne 0 ]; then
    echo -e "${red}❌ Fail install superfile: ${yellow}Unable to create temporary directory${nc}"
    exit 1
fi

package=superfile
version=1.1.4
arch=$(uname -m)
os=$(uname -s)

cd "${temp_dir}"

if [[ "$arch" == "x86_64" ]]; then
    arch="amd64"
elif [[ "$arch" == "arm"* ]]; then
    arch="arm64"
else
    echo -e "${red}❌ Fail install superfile: ${yellow}Unsupported architecture${nc}"
    exit 1
fi

if [[ "$os" == "Linux" ]]; then
    os="linux"
elif [[ "$os" == "Darwin" ]]; then
    os="darwin"
else
    echo -e "${red}❌ Fail install superfile: ${yellow}Unsupported operating system${nc}"
    exit 1
fi

file_name=${package}-${os}-v${version}-${arch}

url="https://github.com/yorukot/superfile/releases/download/v${version}/${file_name}.tar.gz"

if command -v curl &> /dev/null; then
    echo -e "${bright_yellow}Downloading ${cyan}${package} v${version} for ${os} (${arch})...${nc}"
    curl -sLO "$url"
else
    echo -e "${bright_yellow}Downloading ${cyan}${package} v${version} for ${os} (${arch})...${nc}"
    wget -q "$url"
fi

echo -e "${bright_yellow}Extracting ${cyan}${package}...${nc}"
tar -xzf "${file_name}.tar.gz"

echo -e "${bright_yellow}Installing ${cyan}${package}...${nc}"
cd ./dist/${file_name}
chmod +x ./spf
echo -e "${yellow}Press ctrl+C to not install as sudo and try locally.${nc}"
if ! sudo mv ./spf /usr/local/bin/; then
  echo -e "${yellow}Unable to move binary to /usr/local/bin. Do you have sudo permissions?${nc}"
  mkdir -p ~/.local/bin
  if ! mv ./spf ~/.local/bin/; then
    echo -e "${red}❌ Failed to install superfile: Unable to move to ~/.local/bin as well.${nc}"
  else
    if ! [[ ":$PATH:" == *":$HOME/.local/bin:"* ]]; then
      shell_found_and_not_bash=1
      case $SHELL in
        */bash)
          echo 'export PATH="${HOME}/.local/bin":${PATH}' >> ~/.bashrc
          shell_found_and_not_bash=0
          ;;
        */zsh)
          echo 'export PATH="${HOME}/.local/bin":${PATH}' >> ~/.zshrc
          ;;
        */fish)
          echo 'fish_add_path "${HOME}/.local/bin"' >> ~/.config/fish/config.fish 
          ;;
        */ksh)
          echo 'export PATH="${HOME}/.local/bin":${PATH}' >> ~/.kshrc
          ;;
        */xonsh)
          echo '$PATH.prepend("${HOME}/.local/bin")' >> ~/.xonshrc
          ;;
        */csh)
          echo 'setenv PATH "${HOME}/.local/bin":${PATH}' >> ~/.cshrc
          ;;
        */tcsh)
          echo 'setenv PATH "${HOME}/.local/bin":${PATH}' >> ~/.tshrc
          ;;
        *)
          echo -e "${red}Unsupported shell: ${SHELL}. Please add ${white}\"${bright_cyan}\${HOME}/.local/bin${white}\" ${red}to PATH in your shell's config file.${red}"
          shell_found_and_not_bash=0
          ;;
      esac
      if [ $shell_found_and_not_bash == 1 ]; then
        echo -e "${white}\"${bright_purple}${HOME}/.local/bin${white}\"${yellow} has been added to your PATH.${nc}"
        echo -e "${yellow}Please source your config file/relogin.${nc}"
      fi
    fi
    echo -e "🎉 ${bright_cyan}Local ${bright_green}Installation complete!${nc}"
    echo -e "${bright_cyan}You can type ${white}\"${bright_yellow}spf${white}\" ${bright_cyan}to start!${nc}"
  fi
else
  echo -e "🎉 ${bright_green}Installation complete!${nc}"
  echo -e "${bright_cyan}You can type ${white}\"${bright_yellow}spf${white}\" ${bright_cyan}to start!${nc}"
fi

rm -rf "$temp_dir"
