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

	// Error message string
	tokenizationError    = "Failed during tokenization"
	splitCommandArgError = "split command should not be given arguments"

	// Timeout for command executed for shell substitution
	shellSubTimeoutMsec = 1000
)

func modeString(shellMode bool) string {
	if shellMode {
		return shellModeString
	}
	return spfModeString
}

func shellPrompt(shellMode bool) string {
	if shellMode {
		return shellPromptChar
	}
	return spfPromptChar
}

func defaultCommandSlice() []promptCommand {
	return []promptCommand{
		{
			command:     openCommand,
			usage:       openCommand + " <PATH>",
			description: "Open a new panel at a specified path",
		},
		{
			command:     splitCommand,
			usage:       splitCommand,
			description: "Open a new panel at a current file panel's path",
		},
		{
			command:     cdCommand,
			usage:       cdCommand + " <PATH>",
			description: "Change directory of current panel",
		},
	}
}
