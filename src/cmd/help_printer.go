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
	titleColor := color.New(color.FgHiGreen, color.Bold)
	flagColor := color.New(color.FgHiYellow)
	commandColor := color.New(color.FgHiBlue)
	descColor := color.New(color.FgWhite)
	accentColor := color.New(color.FgHiCyan)

	switch v := data.(type) {
	case *cli.Command:
		// Get the actual binary name from os.Args[0]
		binaryName := filepath.Base(os.Args[0])

		// Print usage section
		titleColor.Fprintf(w, "Usage:")
		fmt.Fprintf(w, " %s", binaryName)
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

		// Print description if available
		if v.Description != "" {
			descColor.Fprintf(w, "%s\n\n", strings.TrimSpace(v.Description))
		}

		// Print commands section
		if len(v.Commands) > 0 {
			titleColor.Fprintf(w, "Commands:\n")
			for _, cmd := range v.Commands {
				// Format command name with aliases
				cmdDisplay := cmd.Name
				if len(cmd.Aliases) > 0 {
					cmdDisplay = fmt.Sprintf("%s, %s", cmd.Name, strings.Join(cmd.Aliases, ", "))
				}

				commandColor.Fprintf(w, "  %-20s", cmdDisplay)
				descColor.Fprintf(w, " %s\n", cmd.Usage)
			}
			fmt.Fprintln(w)
		}

		// Print global options section
		if len(v.Flags) > 0 {
			titleColor.Fprintf(w, "Options:\n")
			for _, flag := range v.Flags {
				names := flag.Names()

				// Format flag names with proper prefixes and aliases
				var flagParts []string
				for _, name := range names {
					if len(name) == 1 {
						flagParts = append(flagParts, "-"+name)
					} else {
						flagParts = append(flagParts, "--"+name)
					}
				}
				flagStr := strings.Join(flagParts, ", ")

				flagColor.Fprintf(w, "  %-24s", flagStr)

				// Get usage text from different flag types
				var usage string
				switch f := flag.(type) {
				case *cli.BoolFlag:
					usage = f.Usage
				case *cli.StringFlag:
					usage = f.Usage
					if f.Value != "" {
						usage += fmt.Sprintf(" (default: %q)", f.Value)
					}
				case *cli.StringSliceFlag:
					usage = f.Usage
				case *cli.IntFlag:
					usage = f.Usage
					if f.Value != 0 {
						usage += fmt.Sprintf(" (default: %d)", f.Value)
					}
				default:
					usage = "No description available"
				}

				descColor.Fprintf(w, " %s\n", usage)
			}
			fmt.Fprintln(w)
		}

		// Print version info if available
		if v.Version != "" {
			accentColor.Fprintf(w, "Version: %s\n\n", v.Version)
		}

		// Print help footer using the actual binary name
		descColor.Fprintf(w, "Use \"%s [COMMAND] --help\" for more information about a command.\n", binaryName)

	default:
		// Fallback to default template rendering for other cases
		cli.HelpPrinterCustom(w, templ, data, nil)
	}
}
