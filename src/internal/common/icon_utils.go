package common

import (
	"path/filepath"
	"strings"

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
