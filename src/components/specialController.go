package components

import (
	"os"
	"path/filepath"
)

func CancelModal(m model) model {
	m.createNewItem.textInput.Blur()
	m.createNewItem.open = false
	return m
}

func CreateItem(m model) model {
	if m.createNewItem.itemType == newFile {
		path := m.createNewItem.location + "/" + m.createNewItem.textInput.Value()
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			OutputLog(err)
		}
		f, err := os.Create(path)
		if err != nil {
			OutputLog(err)
		}
		defer f.Close()
	} else {
		path := m.createNewItem.location + "/" + m.createNewItem.textInput.Value()
		err := os.MkdirAll(path, 0755)
		if err != nil {
			OutputLog(err)
		}
	}
	m.createNewItem.open = false
	m.createNewItem.textInput.Blur()
	return m
}

func CancelReanem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.rename.Blur()
	panel.renaming = false
	m.fileModel.renaming = false
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

func ConfirmRename(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	oldPath := panel.element[panel.cursor].location
    newPath := panel.location + "/" + panel.rename.Value()        
    
    // Rename the file
    err := os.Rename(oldPath, newPath)
    if err != nil {
		OutputLog(err)
    }
	m.fileModel.renaming = false
	panel.rename.Blur()
	panel.renaming = false
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}