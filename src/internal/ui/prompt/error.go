package prompt

import "fmt"

// This is to generate error objects that can be nicely printed to UI
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

type envVarNotFoundError struct {
	varName string
}

func (e envVarNotFoundError) Error() string {
	return fmt.Sprintf("env var %s not found", e.varName)
}

type bracketMatchError struct {
	openChar  rune
	closeChar rune
}

func (p bracketMatchError) Error() string {
	return fmt.Sprintf("could not find matching %c for %c", p.closeChar, p.openChar)
}

func roundBracketMatchError() bracketMatchError {
	return bracketMatchError{openChar: '(', closeChar: ')'}
}

func curlyBracketMatchError() bracketMatchError {
	return bracketMatchError{openChar: '{', closeChar: '}'}
}
