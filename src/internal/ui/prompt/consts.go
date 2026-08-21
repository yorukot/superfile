package prompt

import "time"

// These could as well be property of prompt Model vs being global consts
// But its fine
const (
	promptHeadlineText = "superfile Prompt"

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
	shellSubTimeout = 1000 * time.Millisecond

	// Budget for substitutions that tests expect to succeed. Those commands
	// return immediately, so the budget only has to cover spawning a shell.
	// Being generous costs nothing on a healthy run, and keeps a CI runner
	// that is slow to spawn one from being reported as a substitution failure.
	shellSubSuccessTimeoutInTests = 30 * time.Second

	// Budget for the substitution that tests expect to expire. It is paired
	// with a command that runs for far longer, so a slow spawn can never look
	// like a command that finished in time.
	shellSubTimeoutInTests = 100 * time.Millisecond

	defaultTestCwd = "/"

	PromptMinWidth  = 10
	PromptMinHeight = 3

	defaultTestWidth     = 100
	defaultTestMaxHeight = 100

	// UI dimension constants for prompt modal
	// promptInputPadding is total padding for prompt input fields
	promptInputPadding = 6 // 2 + 1 + 2 + 1 (borders and spacing)

	// expectedArgCount is the expected number of prompt arguments
	expectedArgCount = 2
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
			description: "Open a new panel at the current file panel's path",
		},
		{
			command:     CdCommand,
			usage:       CdCommand + " <PATH>",
			description: "Change directory of current panel",
		},
	}
}
