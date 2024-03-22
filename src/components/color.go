package components

func FilePanelFocusColor(focusType filePanelFocusType) string {
	if focusType == noneFocus {
		return "#A4A2A2"
	} else if focusType == secondFocus {
		return "#656565"
	} else {
		return "#D79921"
	}
}
