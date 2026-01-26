package sortmodel

func New() Model {
	return Model{
		Height: sortOptionsDefaultHeight,
		Width:  sortOptionsDefaultWidth,
		Cursor: 0,
		open:   false,
	}
}
