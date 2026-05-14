package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type target struct {
	goos   string
	goarch string
}

func main() {
	targets := []target{
		{goos: "linux", goarch: "amd64"},
		{goos: "linux", goarch: "arm64"},
		{goos: "darwin", goarch: "amd64"},
		{goos: "darwin", goarch: "arm64"},
		{goos: "windows", goarch: "amd64"},
		{goos: "windows", goarch: "arm64"},
	}
	sections := map[string]string{}

	for _, target := range targets {
		report, err := runGoLicenses(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "generate notice for %s/%s: %v\n", target.goos, target.goarch, err)
			os.Exit(1)
		}

		for name, section := range parseSections(report) {
			existing, ok := sections[name]
			if !ok || shouldReplace(existing, section) {
				sections[name] = section
			}
		}
	}

	names := make([]string, 0, len(sections))
	for name := range sections {
		names = append(names, name)
	}
	sort.Strings(names)

	var output bytes.Buffer
	output.WriteString("# Notices\n\n\n")
	for index, name := range names {
		if index > 0 {
			output.WriteString("\n\n")
		}
		output.WriteString(strings.TrimRight(sections[name], "\n"))
	}
	output.WriteByte('\n')

	if err := os.WriteFile("NOTICE.md", output.Bytes(), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "write NOTICE.md: %v\n", err)
		os.Exit(1)
	}
}

func runGoLicenses(target target) (string, error) {
	cmd := exec.Command("go-licenses", "report", "./", "--template=notice.tmpl")
	cmd.Env = append(os.Environ(), "GOOS="+target.goos, "GOARCH="+target.goarch)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func parseSections(report string) map[string]string {
	report = strings.ReplaceAll(report, "\r\n", "\n")
	lines := strings.Split(report, "\n")
	sections := map[string]string{}

	var currentName string
	var current []string

	flush := func() {
		if currentName == "" {
			return
		}
		sections[currentName] = strings.Join(current, "\n")
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			flush()
			currentName = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			current = []string{line}
			continue
		}

		if currentName != "" {
			current = append(current, line)
		}
	}
	flush()

	return sections
}

func shouldReplace(existing, candidate string) bool {
	existingUnknown := strings.Contains(existing, "](Unknown)")
	candidateUnknown := strings.Contains(candidate, "](Unknown)")

	if existingUnknown && !candidateUnknown {
		return true
	}
	if !existingUnknown && candidateUnknown {
		return false
	}

	return len(candidate) > len(existing)
}
