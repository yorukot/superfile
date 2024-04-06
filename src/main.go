package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/MHNightCat/superfile/components"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	currentVersion      = "v0.1.0-beta"
	latestVersionURL    = "https://api.github.com/repos/MHNightCat/superfile/releases/latest"
	latestVersionGithub = "github.com//MHNightCat/superfile/releases/latest"
	dir                 = "./.superfile/data/lastCheckVersion"
)

func main() {

	p := tea.NewProgram(components.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	CheckForUpdates()
}

func CheckForUpdates() {
	lastTime, err := readFromFile(dir)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error reading from file:", err)
		return
	}

	currentTime := time.Now()

	if lastTime.IsZero() || currentTime.Sub(lastTime) >= 24*time.Hour {
		resp, err := http.Get(latestVersionURL)
		if err != nil {
			fmt.Println("Error checking for updates:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		var release GitHubRelease
		if err := json.Unmarshal(body, &release); err != nil {
			return
		}
		if release.TagName != currentVersion {
			fmt.Printf("A new version %s is available.\nPlease update.\n┏\n\n        %s\n\n                                                               ┛\n", release.TagName, latestVersionGithub)
		}

		timeStr := currentTime.Format(time.RFC3339)
		err = writeToFile(dir, timeStr)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func readFromFile(filename string) (time.Time, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return time.Time{}, err
	}
	if len(content) == 0 {
		return time.Time{}, nil
	}
	lastTime, err := time.Parse(time.RFC3339, string(content))
	if err != nil {
		return time.Time{}, err
	}

	return lastTime, nil
}

func writeToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func downloadFile(url string, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), os.ModePerm)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}