package prompt

import "time"

// These could as well be property of prompt Model vs being global consts
// But its fine
const (
	promptHeadlineText = "Superfile Prompt"

	OpenCommand  = "open"
	SplitCommand = "split"
	CdCommand    = "cd"

	// We could later make this configurable. But, not needed now.
	spfPromptChar   = ">"
	shellPromptChar = ":"

	successMessagePrefix = "Success"
	failureMessagePrefix = "Error"

	shellModeString = "(Shell Mode)"
	spfModeString   = "(SPF Mode)"

	// Error message string
	tokenizationError    = "Failed during tokenization"
	splitCommandArgError = "split command should not be given arguments"

	// Timeout for command executed for shell substitution
	shellSubTimeout        = 1000 * time.Millisecond
	shellSubTimeoutInTests = 100 * time.Millisecond

	defaultTestCwd = "/"

	PromptMinWidth  = 10
	PromptMinHeight = 3

	defaultTestWidth     = 100
	defaultTestMaxHeight = 100
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
			command:     OpenCommand,
			usage:       OpenCommand + " <PATH>",
			description: "Open a new panel at a specified path",
		},
		{
			command:     SplitCommand,
			usage:       SplitCommand,
			description: "Open a new panel at a current file panel's path",
		},
		{
			command:     CdCommand,
			usage:       CdCommand + " <PATH>",
			description: "Change directory of current panel",
		},
	}
}
