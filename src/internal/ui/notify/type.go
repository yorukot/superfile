package notify

type ConfirmActionType int

const (
	RenameAction ConfirmActionType = iota
	DeleteAction
	NoAction
)
