package prompt

import "fmt"

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

type envVarNotFoundError struct {
	varName string
}

func (e envVarNotFoundError) Error() string {
	return fmt.Sprintf("env var %s not found", e.varName)
}

type paranthesisMatchError struct {
	openChar  rune
	closeChar rune
}

func (p paranthesisMatchError) Error() string {
	return fmt.Sprintf("could not find matching %v for %v", p.closeChar, p.openChar)
}

func bracketParMatchError() paranthesisMatchError {
	return paranthesisMatchError{openChar: '(', closeChar: ')'}
}

func curlyBracketParMatchError() paranthesisMatchError {
	return paranthesisMatchError{openChar: '{', closeChar: '}'}
}
