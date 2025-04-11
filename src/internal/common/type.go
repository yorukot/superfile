package common

import (
	"path/filepath"
	"strings"

	"github.com/yorukot/superfile/src/config/icon"
)

// Placeholder inteface for now, might later move 'model' type to commons and have
// and add an execute(model) function to this
type ModelAction interface {
	String() string
}

type NoAction struct {
}

func (n NoAction) String() string {
	return "NoAction"
}

type ShellCommandAction struct {
	Command string
}

func (s ShellCommandAction) String() string {
	return "ShellCommandAction for command " + s.Command
}

// We could later move 'model' type to commons and have
// these actions implement an execute(model) interface
type SplitPanelAction struct{}

func (s SplitPanelAction) String() string {
	return "SplitPanelAction"
}

type CDCurrentPanelAction struct {
	Location string
}

func (c CDCurrentPanelAction) String() string {
	return "CDCurrentPanelAction to " + c.Location
}

type OpenPanelAction struct {
	Location string
}

func (o OpenPanelAction) String() string {
	return "OpenPanelAction at " + o.Location
}

func GetElementIcon(file string, isDir bool) icon.Style {
	ext := strings.TrimPrefix(filepath.Ext(file), ".")
	name := file

	if !Config.Nerdfont {
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
