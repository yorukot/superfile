package prompt

// This is to generate error objects that can pe nicely printed to UI
type invalidCmdError struct {
	uiMsg        string
	wrappedError error
}

func (e invalidCmdError) Error() string {
	if e.wrappedError == nil {
		return e.uiMsg
	}
	return e.wrappedError.Error()
}

func (e invalidCmdError) Unwrap() error {
	return e.wrappedError
}

func (e invalidCmdError) uiMessage() string {
	return e.uiMsg
}
