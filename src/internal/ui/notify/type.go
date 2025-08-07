package notify

type ConfirmActionType int

const (
	RenameAction ConfirmActionType = iota
	DeleteAction
	QuitAction
	NoAction
)
