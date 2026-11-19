package common

import (
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"

	"github.com/yorukot/superfile/src/config/icon"
)

func getFileIcon(file string, isLink bool) icon.Style {
	if isLink {
		return icon.Icons["link_file"]
	}
	ext := strings.TrimPrefix(filepath.Ext(file), ".")
	// default icon for all files. try to find a better one though...
	resultIcon := icon.Icons["file"]
	// resolve aliased extensions
	extKey := strings.ToLower(ext)
	alias, hasAlias := icon.Aliases[extKey]
	if hasAlias {
		extKey = alias
	}

	// see if we can find a better icon based on extension alone
	betterIcon, hasBetterIcon := icon.Icons[extKey]
	if hasBetterIcon {
		resultIcon = betterIcon
	}

	// now look for icons based on full names
	fullName := file

	fullName = strings.ToLower(fullName)
	fullAlias, hasFullAlias := icon.Aliases[fullName]
	if hasFullAlias {
		fullName = fullAlias
	}
	bestIcon, hasBestIcon := icon.Icons[fullName]
	if hasBestIcon {
		resultIcon = bestIcon
	}
	if resultIcon.Color == "NONE" {
		return icon.Style{
			Icon:  resultIcon.Icon,
			Color: Theme.FilePanelFG,
		}
	}
	return resultIcon
}

func GetElementIcon(file string, isDir bool, isLink bool, nerdFont bool) icon.Style {
	if !nerdFont {
		return icon.Style{
			Icon:  "",
			Color: Theme.FilePanelFG,
		}
	}

	if isDir {
		if isLink {
			return icon.Folders["link_folder"]
		}
		resultIcon := icon.Folders["folder"]
		betterIcon, hasBetterIcon := icon.Folders[file]
		if hasBetterIcon {
			resultIcon = betterIcon
		}
		return resultIcon
	}

	return getFileIcon(file, isLink)
}

func GetDirectoryIcon(path string, name string, nerdFont bool) string {
	if !nerdFont {
		return ""
	}

	switch path {
	case xdg.Home:
		return icon.Home
	case xdg.UserDirs.Desktop:
		return icon.Desktop
	case xdg.UserDirs.Download:
		return icon.Download
	case xdg.UserDirs.Documents:
		return icon.Documents
	case xdg.UserDirs.Pictures:
		return icon.Pictures
	case xdg.UserDirs.Videos:
		return icon.Videos
	case xdg.UserDirs.Music:
		return icon.Music
	case xdg.UserDirs.Templates:
		return icon.Templates
	case xdg.UserDirs.PublicShare:
		return icon.PublicShare
	default:
		result, exists := icon.Folders[name]
		if !exists {
			return icon.Directory
		}

		return result.Icon
	}
}
