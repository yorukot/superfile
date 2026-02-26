package internal

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Cancel typing modal e.g. create file or directory
func (m *model) cancelTypingModal() {
	m.typingModal.textInput.Blur()
	m.typingModal.open = false
}

// Confirm to create file or directory
func (m *model) createItem() {
	if err := checkFileNameValidity(m.typingModal.textInput.Value()); err != nil {
		m.typingModal.errorMesssage = err.Error()
		slog.Error("Errow while createItem during item creation", "error", err)

		return
	}

	defer func() {
		m.typingModal.errorMesssage = ""
		m.typingModal.open = false
		m.typingModal.textInput.Blur()
	}()

	path := filepath.Join(m.typingModal.location, m.typingModal.textInput.Value())
	if !strings.HasSuffix(m.typingModal.textInput.Value(), string(filepath.Separator)) {
		path, _ = renameIfDuplicate(path)
		if err := os.MkdirAll(filepath.Dir(path), utils.UserDirPerm); err != nil {
			slog.Error("Error while createItem during directory creation", "error", err)
			return
		}
		f, err := os.Create(path)
		if err != nil {
			slog.Error("Error while createItem during file creation", "error", err)
			return
		}
		defer f.Close()
	} else {
		err := os.MkdirAll(path, utils.UserDirPerm)
		if err != nil {
			slog.Error("Error while createItem during directory creation", "error", err)
			return
		}
	}
}

// Cancel rename file or directory
func (m *model) cancelRename() {
	panel := m.getFocusedFilePanel()
	panel.Rename.Blur()
	panel.Renaming = false
	m.fileModel.Renaming = false
}

// Connfirm rename file or directory
func (m *model) confirmRename() {
	panel := m.getFocusedFilePanel()

	// Although we dont expect this to happen based on our current flow
	// Just adding it here to be safe
	if panel.Empty() {
		slog.Error("confirmRename called on empty panel")
		return
	}

	oldPath := panel.GetFocusedItem().Location
	newPath := filepath.Join(panel.Location, panel.Rename.Value())

	// Rename the file
	err := os.Rename(oldPath, newPath)
	if err != nil {
		slog.Error("Error while confirmRename during rename", "error", err)
		// Dont return. We have to also reset the panel and model information
	}
	m.fileModel.Renaming = false
	panel.Rename.Blur()
	panel.Renaming = false
}

func (m *model) confirmSortOptions() {
	panel := m.getFocusedFilePanel()
	panel.SortKind = m.sortModal.GetSelectedKind()
	m.sortModal.Close()
}

// Cancel search, this will clear all searchbar input
func (m *model) cancelSearch() {
	panel := m.getFocusedFilePanel()
	panel.SearchBar.Blur()
	panel.SearchBar.SetValue("")
}

// Confirm search. This will exit the search bar and filter the files
func (m *model) confirmSearch() {
	panel := m.getFocusedFilePanel()
	panel.SearchBar.Blur()
}

func (m *model) getFocusedFilePanel() *filepanel.Model {
	return m.fileModel.GetFocusedFilePanel()
}
