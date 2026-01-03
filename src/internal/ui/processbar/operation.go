package processbar

import "github.com/yorukot/superfile/src/config/icon"

type OperationType int

const (
	OpCopy OperationType = iota
	OpCut
	OpDelete
	OpCompress
	OpExtract
)

// GetIcon returns the appropriate icon for the operation type
func (op OperationType) GetIcon() string {
	switch op {
	case OpCopy:
		return icon.Copy
	case OpCut:
		return icon.Cut
	case OpDelete:
		return icon.Delete
	case OpCompress:
		return icon.CompressFile
	case OpExtract:
		return icon.ExtractFile
	default:
		return icon.InOperation
	}
}

// GetVerb returns the present tense verb for the operation
func (op OperationType) GetVerb() string {
	switch op {
	case OpCopy:
		return "Copying"
	case OpCut:
		return "Moving"
	case OpDelete:
		return "Deleting"
	case OpCompress:
		return "Compressing"
	case OpExtract:
		return "Extracting"
	default:
		return "Processing"
	}
}

// GetPastVerb returns the past tense verb for the operation
func (op OperationType) GetPastVerb() string {
	switch op {
	case OpCopy:
		return "Copied"
	case OpCut:
		return "Moved"
	case OpDelete:
		return "Deleted"
	case OpCompress:
		return "Compressed"
	case OpExtract:
		return "Extracted"
	default:
		return "Processed"
	}
}
