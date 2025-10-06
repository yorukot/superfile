package internal

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
)

func (m *model) panelBulkRename() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]

	if panel.panelMode != selectMode || len(panel.selected) == 0 {
		return
	}

	m.bulkRenameModal.open = true
	m.bulkRenameModal.renameType = 0
	m.bulkRenameModal.cursor = 0
	m.bulkRenameModal.startNumber = 1
	m.bulkRenameModal.caseType = 0
	m.bulkRenameModal.errorMessage = ""
	m.firstTextInput = true

	m.bulkRenameModal.findInput = common.GenerateBulkRenameTextInput("Find text")
	m.bulkRenameModal.replaceInput = common.GenerateBulkRenameTextInput("Replace with")
	m.bulkRenameModal.prefixInput = common.GenerateBulkRenameTextInput("Add prefix")
	m.bulkRenameModal.suffixInput = common.GenerateBulkRenameTextInput("Add suffix")
}

func (m *model) generateBulkRenamePreview() []renamePreview {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	previews := make([]renamePreview, 0, len(panel.selected))

	for i, itemPath := range panel.selected {
		preview := m.createRenamePreview(itemPath, i)
		previews = append(previews, preview)
	}

	return previews
}

func (m *model) createRenamePreview(itemPath string, index int) renamePreview {
	oldName := filepath.Base(itemPath)
	newName := m.applyBulkRenameTransformation(oldName, index)

	validation := renameValidation{
		oldName:  oldName,
		newName:  newName,
		itemPath: itemPath,
	}

	return renamePreview{
		oldName: oldName,
		newName: newName,
		error:   validateRename(validation),
	}
}

func (m *model) applyBulkRenameTransformation(oldName string, index int) string {
	modal := &m.bulkRenameModal
	
	transformers := map[int]func() string{
		0: func() string { return applyFindReplace(oldName, modal.findInput.Value(), modal.replaceInput.Value()) },
		1: func() string { return applyPrefix(oldName, modal.prefixInput.Value()) },
		2: func() string { return applySuffix(oldName, modal.suffixInput.Value()) },
		3: func() string { return applyNumbering(oldName, modal.startNumber+index) },
		4: func() string { return applyCaseConversion(oldName, modal.caseType) },
	}
	
	if transformer, exists := transformers[modal.renameType]; exists {
		return transformer()
	}
	return oldName
}

func validateRename(v renameValidation) string {
	if v.newName == "" {
		return "Empty filename"
	}
	if v.newName == v.oldName {
		return "No change"
	}
	return checkRenameConflict(v)
}

func checkRenameConflict(v renameValidation) string {
	newPath := filepath.Join(filepath.Dir(v.itemPath), v.newName)
	if _, statErr := os.Stat(newPath); statErr == nil && newPath != v.itemPath {
		return "File already exists"
	}
	return ""
}

type renameTransformer struct {
	find    string
	replace string
	prefix  string
	suffix  string
	number  int
	caseOp  caseOperation
}

type caseOperation int

const (
	toLower caseOperation = iota
	toUpper
	toTitle
)

func applyFindReplace(filename, find, replace string) string {
	if find == "" {
		return filename
	}
	return strings.ReplaceAll(filename, find, replace)
}

func applyPrefix(filename, prefix string) string {
	if prefix == "" {
		return filename
	}
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return prefix + nameWithoutExt + ext
}

func applySuffix(filename, suffix string) string {
	if suffix == "" {
		return filename
	}
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return nameWithoutExt + suffix + ext
}

func applyNumbering(filename string, number int) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return nameWithoutExt + "_" + strconv.Itoa(number) + ext
}

func applyCaseConversion(filename string, caseType int) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	converters := map[int]func(string) string{
		0: strings.ToLower,
		1: strings.ToUpper,
		2: toTitleCase,
	}

	if converter, exists := converters[caseType]; exists {
		return converter(nameWithoutExt) + ext
	}
	return filename
}

func toTitleCase(text string) string {
	words := strings.Fields(strings.ToLower(text))
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
