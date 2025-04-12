package common

import (
	"path/filepath"
	"strings"

	"github.com/yorukot/superfile/src/config/icon"
)

func GetElementIcon(file string, isDir bool, nerdFont bool) icon.Style {
	ext := strings.TrimPrefix(filepath.Ext(file), ".")
	name := file

	if !nerdFont {
		return icon.Style{
			Icon:  "",
			Color: Theme.FilePanelFG,
		}
	}

	if isDir {
		resultIcon := icon.Folders["folder"]
		betterIcon, hasBetterIcon := icon.Folders[name]
		if hasBetterIcon {
			resultIcon = betterIcon
		}
		return resultIcon
	}
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
	fullName := name

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
