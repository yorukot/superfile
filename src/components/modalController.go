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
