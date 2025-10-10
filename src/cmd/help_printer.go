package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

// CustomHelpPrinter provides cargo-style colored help output for superfile CLI
func CustomHelpPrinter(w io.Writer, templ string, data interface{}) {
	// Define color styles matching superfile's aesthetic
	titleColor := color.New(color.FgGreen, color.Bold)
	flagColor := color.New(color.FgCyan, color.Bold)
	commandColor := color.New(color.FgBlue, color.Bold)
	accentColor := color.New(color.FgMagenta, color.Bold)

	switch v := data.(type) {
	case *cli.Command:
		// Get the actual binary name from os.Args[0]
		binaryName := filepath.Base(os.Args[0])
		printUsage(w, titleColor, accentColor, binaryName, v)

		printCommands(w, titleColor, commandColor, v)

		printFlags(w, titleColor, flagColor, v)
		// Print version info if available
		if v.Version != "" {
			fmt.Printf("Version: ")
			accentColor.Fprintf(w, "%s\n\n", v.Version)
		}

		// Print help footer using the actual binary name
		fmt.Fprint(w, "Use \"")
		accentColor.Fprintf(w, "%s", binaryName)
		fmt.Fprint(w, " [COMMAND] --help\" for more information about a command.\n")

	default:
		// Fallback to default template rendering for other cases
		cli.HelpPrinterCustom(w, templ, data, nil)
	}
}

func printUsage(w io.Writer, titleColor *color.Color, accentColor *color.Color, binaryName string, v *cli.Command) {
	titleColor.Fprintf(w, "Usage:")
	fmt.Fprint(w, " ")
	accentColor.Fprintf(w, "%s", binaryName)
	if len(v.Commands) > 0 {
		fmt.Fprint(w, " [COMMAND]")
	}
	if len(v.Flags) > 0 {
		fmt.Fprint(w, " [OPTIONS]")
	}
	if v.ArgsUsage != "" {
		fmt.Fprintf(w, " %s", v.ArgsUsage)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w)
	if v.Description != "" {
		fmt.Fprintf(w, "%s\n\n", strings.TrimSpace(v.Description))
	}
}

func printCommands(w io.Writer, titleColor *color.Color, commandColor *color.Color, v *cli.Command) {
	if len(v.Commands) == 0 {
		return
	}
	titleColor.Fprintf(w, "Commands:\n")
	for _, cmd := range v.Commands {
		// Format command name with aliases
		cmdDisplay := cmd.Name
		if len(cmd.Aliases) > 0 {
			cmdDisplay = fmt.Sprintf("%s, %s", cmd.Name, strings.Join(cmd.Aliases, ", "))
		}

		commandColor.Fprintf(w, "  %-20s", cmdDisplay)
		fmt.Fprintf(w, " %s\n", cmd.Usage)
	}
	fmt.Fprintln(w)
}

func printFlags(w io.Writer, titleColor *color.Color, flagColor *color.Color, v *cli.Command) {
	if len(v.Flags) == 0 {
		return
	}
	titleColor.Fprintf(w, "Options:\n")
	for _, flag := range v.Flags {
		names := flag.Names()

		// Format flag names with proper prefixes and aliases
		var flagParts []string
		var valuePlaceholder string
		var usage string

		// Determine flag type, value placeholder, and usage in one switch
		switch f := flag.(type) {
		case *cli.BoolFlag:
			// Boolean flags don't need values
			valuePlaceholder = ""
			usage = f.Usage
		case *cli.StringFlag:
			valuePlaceholder = " <value>"
			usage = f.Usage
			if f.Value != "" {
				usage += fmt.Sprintf(" (default: %q)", f.Value)
			}
		case *cli.StringSliceFlag:
			valuePlaceholder = " <value>..."
			usage = f.Usage
		case *cli.IntFlag:
			valuePlaceholder = " <number>"
			usage = f.Usage
			if f.Value != 0 {
				usage += fmt.Sprintf(" (default: %d)", f.Value)
			}
		default:
			valuePlaceholder = " <value>"
			usage = "No description available"
		}

		for _, name := range names {
			if len(name) == 1 {
				flagParts = append(flagParts, "-"+name)
			} else {
				flagParts = append(flagParts, "--"+name)
			}
		}
		flagStr := strings.Join(flagParts, ", ") + valuePlaceholder

		flagColor.Fprintf(w, "  %-30s", flagStr)
		fmt.Fprintf(w, " %s\n", usage)
	}
	fmt.Fprintln(w)
}
