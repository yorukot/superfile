package common

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// Standard folder labels shared by navigation and the sidebar.
const (
	FolderHome        = "Home"
	FolderDownloads   = "Downloads"
	FolderDesktop     = "Desktop"
	FolderDocuments   = "Documents"
	FolderPictures    = "Pictures"
	FolderMusic       = "Music"
	FolderVideos      = "Videos"
	FolderTemplates   = "Templates"
	FolderPublicShare = "PublicShare"
	FolderTrash       = "Trash"
)

// FolderShortcut describes a standard folder or a one-based filtered pin slot.
type FolderShortcut struct {
	Name string
	Keys []string
	Pin  int
}

// FolderShortcuts returns effective bindings in precedence order. Existing
// actions always win collisions, including actions in other panel modes.
func (h HotkeysType) FolderShortcuts() []FolderShortcut {
	shortcuts := []FolderShortcut{
		{Name: FolderHome, Keys: h.GoToHome},
		{Name: FolderDownloads, Keys: h.GoToDownloads},
		{Name: FolderDesktop, Keys: h.GoToDesktop},
		{Name: FolderDocuments, Keys: h.GoToDocuments},
		{Name: FolderPictures, Keys: h.GoToPictures},
		{Name: FolderMusic, Keys: h.GoToMusic},
		{Name: FolderVideos, Keys: h.GoToVideos},
		{Name: FolderTemplates, Keys: h.GoToTemplates},
		{Name: FolderPublicShare, Keys: h.GoToPublicShare},
		{Name: FolderTrash, Keys: h.GoToTrash},
	}
	pins := [][]string{h.GoToPin1, h.GoToPin2, h.GoToPin3, h.GoToPin4, h.GoToPin5,
		h.GoToPin6, h.GoToPin7, h.GoToPin8, h.GoToPin9}
	for i, keys := range pins {
		shortcuts = append(
			shortcuts,
			FolderShortcut{Name: "Pinned folder " + strconv.Itoa(i+1), Keys: keys, Pin: i + 1},
		)
	}

	used := h.existingFolderShortcutConflicts()
	for i := range shortcuts {
		keys := make([]string, 0, len(shortcuts[i].Keys))
		for _, key := range shortcuts[i].Keys {
			if key != "" && !used[key] {
				keys = append(keys, key)
				used[key] = true
			}
		}
		shortcuts[i].Keys = keys
	}
	return shortcuts
}

func (h HotkeysType) existingFolderShortcutConflicts() map[string]bool {
	used := make(map[string]bool)
	val := reflect.ValueOf(h)
	for i := range val.NumField() {
		name := val.Type().Field(i).Tag.Get("toml")
		if strings.HasPrefix(name, "go_to_") || name == "confirm_typing" || name == "cancel_typing" {
			continue
		}
		for j := range val.Field(i).Len() {
			used[val.Field(i).Index(j).String()] = true
		}
	}
	return used
}

// AltHint omits the held modifier and ignores non-Alt bindings.
func (s FolderShortcut) AltHint() string {
	for _, key := range s.Keys {
		parts := strings.Split(key, "+")
		if i := slices.Index(parts, "alt"); i >= 0 {
			return strings.Join(slices.Delete(parts, i, i+1), "+")
		}
	}
	return ""
}
