package prompt

// These could as well be property of prompt Model vs being global consts
// But its fine
const (
	promptHeadlineText = "Superfile Prompt"

	openCommand  = "open"
	splitCommand = "split"
	cdCommand    = "cd"

	// We could later make this configurable. But, not needed now.
	spfPromptChar   = ">"
	shellPromptChar = ":"

	successMessagePrefix = "Success"
	failureMessagePrefix = "Error"

	shellModeString = "(Shell Mode)"
	spfModeString   = "(Prompt Mode)"

	shellSubTimeoutMsec = 1000
)
