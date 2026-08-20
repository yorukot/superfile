package prompt

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	tea "charm.land/bubbletea/v2"
)

// completionsReadyMsg carries shell completion results back to the prompt.
type completionsReadyMsg struct {
	completions []string
	prefix      string
}

func (msg completionsReadyMsg) String() string {
	return fmt.Sprintf("completionsReadyMsg{n=%d, prefix=%q}", len(msg.completions), msg.prefix)
}

// fetchCompletions returns a tea.Cmd that asynchronously fetches shell
// completions for the last token in the input string.
func fetchCompletions(input string, cwd string) tea.Cmd {
	prefix, _ := getLastToken(input)
	if prefix == "" {
		return nil
	}
	return func() tea.Msg {
		completions := getShellCompletions(prefix, cwd)
		return completionsReadyMsg{completions: completions, prefix: prefix}
	}
}

// getShellCompletions queries bash's compgen for file and command completions.
// Falls back to manual PATH + directory listing if compgen is unavailable.
func getShellCompletions(prefix string, cwd string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	escaped := shellEscape(prefix)
	// -c: commands, -f: files, combined gives both
	cmd := exec.CommandContext(ctx, "bash", "-c",
		fmt.Sprintf("shopt -s nocaseglob 2>/dev/null; compgen -cf -- %s 2>/dev/null", escaped))
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			// Deduplicate while preserving order
			seen := make(map[string]bool, len(lines))
			uniq := lines[:0]
			for _, l := range lines {
				key := strings.ToLower(l)
				if !seen[key] {
					seen[key] = true
					uniq = append(uniq, l)
				}
			}
			return uniq
		}
	}

	// Fallback: manual file + command completion
	return manualCompletions(prefix, cwd)
}

// manualCompletions provides a basic fallback when compgen is not available.
// It lists matching files in cwd and matching executables from $PATH.
func manualCompletions(prefix string, cwd string) []string {
	var results []string
	seen := make(map[string]bool)

	// File completions from cwd
	entries, err := os.ReadDir(cwd)
	if err == nil {
		for _, e := range entries {
			name := e.Name()
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				// Skip hidden files unless prefix starts with '.'
				if !strings.HasPrefix(prefix, ".") && strings.HasPrefix(name, ".") {
					continue
				}
				if e.IsDir() {
					name += string(os.PathSeparator)
				}
				results = append(results, name)
				seen[name] = true
			}
		}
	}

	// Command completions from $PATH (only for first token or when prefix looks like a command)
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			name := e.Name()
			if !strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) || seen[name] {
				continue
			}
			seen[name] = true
			results = append(results, name)
		}
	}

	return results
}

// getLastToken extracts the last whitespace-delimited token from the input.
// Respects single and double quotes. Returns the token and its start index.
func getLastToken(input string) (token string, startIdx int) {
	trimmed := strings.TrimRight(input, " \t")
	if trimmed == "" {
		return "", len(input)
	}

	inQuote := false
	var quote rune
	runes := []rune(trimmed)
	lastSpace := -1

	for i, r := range runes {
		if inQuote {
			if r == quote {
				inQuote = false
			}
		} else {
			if r == '"' || r == '\'' {
				inQuote = true
				quote = r
			} else if unicode.IsSpace(r) {
				lastSpace = i
			}
		}
	}

	startIdx = lastSpace + 1
	token = string(runes[startIdx:])
	return token, startIdx
}

// applyCompletion replaces the last token in input with the given completion.
func applyCompletion(input string, completion string) string {
	_, startIdx := getLastToken(input)
	runes := []rune(input)
	// Trim trailing whitespace from input before applying completion
	trimLen := len(runes)
	for trimLen > 0 && unicode.IsSpace(runes[trimLen-1]) {
		trimLen--
	}
	return string(runes[:startIdx]) + completion
}

// commonPrefix returns the longest common prefix among all strings.
func commonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}

// shellEscape wraps a string in single quotes for safe shell usage.
func shellEscape(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
