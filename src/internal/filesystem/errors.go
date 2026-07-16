package filesystem

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupported  = errors.New("unsupported filesystem operation")
	ErrPermission   = errors.New("filesystem permission denied")
	ErrNotFound     = errors.New("filesystem path not found")
	ErrCanceled     = errors.New("filesystem operation canceled")
	ErrDisconnected = errors.New("filesystem session disconnected")
	ErrConflict     = errors.New("filesystem conflict")
)

type OperationError struct {
	Kind      error
	Provider  ProviderKind
	Operation Operation
	Path      Path
	Message   string
}

func (e *OperationError) Error() string {
	if e == nil {
		return ""
	}
	message := e.Kind.Error()
	if e.Message != "" {
		message = e.Message
	}
	return fmt.Sprintf("%s: provider=%s operation=%s path=%s", message, e.Provider, e.Operation, e.Path.String())
}

func (e *OperationError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Kind
}

func (e *OperationError) Is(target error) bool {
	return e != nil && errors.Is(e.Kind, target)
}

func NewUnsupportedError(provider ProviderKind, operation Operation, path Path, message string) error {
	return &OperationError{Kind: ErrUnsupported, Provider: provider, Operation: operation, Path: path, Message: message}
}

func NewPermissionError(provider ProviderKind, operation Operation, path Path, message string) error {
	return &OperationError{Kind: ErrPermission, Provider: provider, Operation: operation, Path: path, Message: message}
}

func NewNotFoundError(provider ProviderKind, operation Operation, path Path, message string) error {
	return &OperationError{Kind: ErrNotFound, Provider: provider, Operation: operation, Path: path, Message: message}
}

func NewCanceledError(provider ProviderKind, operation Operation, path Path, message string) error {
	return &OperationError{Kind: ErrCanceled, Provider: provider, Operation: operation, Path: path, Message: message}
}

func NewDisconnectedError(provider ProviderKind, operation Operation, path Path, message string) error {
	return &OperationError{
		Kind:      ErrDisconnected,
		Provider:  provider,
		Operation: operation,
		Path:      path,
		Message:   message,
	}
}

func NewConflictError(provider ProviderKind, operation Operation, path Path, message string) error {
	return &OperationError{Kind: ErrConflict, Provider: provider, Operation: operation, Path: path, Message: message}
}
