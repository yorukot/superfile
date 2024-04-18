package components

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/lithammer/shortuuid"
	"github.com/masatana/go-textdistance"
	"github.com/pelletier/go-toml/v2"
	"github.com/rkoesters/xdg/userdirs"
	"github.com/shirou/gopsutil/disk"
)

func getDirectories(height int) []directory {
	directories := []directory{}

	directories = append(directories, getWellKnownDirectories()...)
	if height > 30 {
		directories = append(directories, getPinnedDirectories()...)
		directories = append(directories, getExternalMediaFolders()...)
	}
	return directories
}

func getWellKnownDirectories() []directory {
	directories := []directory{}
	wellKnownDirectories := []directory{
		{location: HomeDir, name: "󰋜 Home"},
		{location: userdirs.Download, name: "󰏔 Downloads"},
		{location: userdirs.Documents, name: "󰈙 Documents"},
		{location: userdirs.Pictures, name: "󰋩 Pictures"},
		{location: userdirs.Videos, name: "󰎁 Videos"},
		{location: userdirs.Music, name: "♬ Music"},
		{location: userdirs.Templates, name: "󰏢 Templates"},
		{location: userdirs.PublicShare, name: " PublicShare"},
	}

	if runtime.GOOS == "darwin" {
		wellKnownDirectories[1].location = HomeDir + "/Downloads/"
		wellKnownDirectories[2].location = HomeDir + "/Documents/"
		wellKnownDirectories[3].location = HomeDir + "/Pictures/"
		wellKnownDirectories[4].location = HomeDir + "/Movies/"
		wellKnownDirectories[5].location = HomeDir + "/Music/"
		wellKnownDirectories[7].location = HomeDir + "/Public/"
	}

	for _, dir := range wellKnownDirectories {
		if _, err := os.Stat(dir.location); !os.IsNotExist(err) {
			// Directory exists
			directories = append(directories, dir)
		}
	}
	return directories
}

func getPinnedDirectories() []directory {
	directories := []directory{}
	var paths []string

	jsonData, err := os.ReadFile(SuperFileDataDir + pinnedFile)
	if err != nil {
		outPutLog("Read superfile data error", err)
	}

	json.Unmarshal(jsonData, &paths)

	for _, path := range paths {
		directoryName := filepath.Base(path)
		directories = append(directories, directory{location: path, name: directoryName})
	}
	return directories
}

func getExternalMediaFolders() (disks []directory) {
	parts, err := disk.Partitions(true)

	if err != nil {
		outPutLog("Error while getting external media: ", err)
	}
	for _, disk := range parts {
		if isExternalDiskPath(disk.Mountpoint) {
			disks = append(disks, directory{
				name:     filepath.Base(disk.Mountpoint),
				location: disk.Mountpoint,
			})
		}
	}
	if err != nil {
		outPutLog("Error while getting external media: ", err)
	}
	return disks
}

func isExternalDiskPath(path string) bool {
	dir := filepath.Dir(path)
	return strings.HasPrefix(dir, "/mnt") ||
		strings.HasPrefix(dir, "/media") ||
		strings.HasPrefix(dir, "/run/media") ||
		strings.HasPrefix(dir, "/Volumes")
}

func returnFocusType(focusPanel focusPanelType) filePanelFocusType {
	if focusPanel == nonePanelFocus {
		return focus
	}
	return secondFocus
}

func returnFolderElement(location string, displayDotFile bool) (folderElement []element) {
	var files []element
	var folders []element

	items, err := os.ReadDir(location)
	if err != nil {
		outPutLog("Return folder element function error", err)
	}

	for _, item := range items {
		fileInfo, err := item.Info()
		if err != nil {
			continue
		}

		if !displayDotFile && strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}
		if fileInfo == nil {
			continue
		}
		newElement := element{
			name:      item.Name(),
			directory: item.IsDir(),
		}
		if location == "/" {
			newElement.location = location + item.Name()
		} else {
			newElement.location = location + "/" + item.Name()
		}

		if item.IsDir() {
			folders = append(folders, newElement)
		} else {
			files = append(files, newElement)
		}
	}

	// Sort folders and files alphabetically
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].name < folders[j].name
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].name < files[j].name
	})

	// Concatenate folders and files
	folderElement = append(folders, files...)

	return folderElement
}

func returnFolderElementBySearchString(location string, displayDotFile bool, searchString string) (folderElement []element) {

	items, err := os.ReadDir(location)
	if err != nil {
		outPutLog("Return folder element function error", err)
	}

	for _, item := range items {
		fileInfo, _ := item.Info()
		if !displayDotFile && strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}
		if fileInfo == nil {
			continue
		}
		newElement := element{
			name:      item.Name(),
			directory: item.IsDir(),
			matchRate: textdistance.JaroWinklerDistance(item.Name(), searchString),
		}
		if location == "/" {
			newElement.location = location + item.Name()
		} else {
			newElement.location = location + "/" + item.Name()
		}
		if newElement.matchRate > 0 {
			folderElement = append(folderElement, newElement)
		}
	}

	// Sort folders and files by match rate
	sort.Slice(folderElement, func(i, j int) bool {
		return folderElement[i].matchRate > folderElement[j].matchRate
	})

	return folderElement
}

func panelElementHeight(mainPanelHeight int) int {
	return mainPanelHeight - 3
}

func bottomElementHight(bottomElementHight int) int {
	return bottomElementHight - 5
}

func arrayContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func outPutLog(values ...interface{}) {
	log.SetOutput(logOutput)
	for _, value := range values {
		log.Println(value)
	}
}

func removeElementByValue(slice []string, value string) []string {
	newSlice := []string{}
	for _, v := range slice {
		if v != value {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

func renameIfDuplicate(destination string) (string, error) {
	info, err := os.Stat(destination)
	if os.IsNotExist(err) {
		return destination, nil
	} else if err != nil {
		return "", err
	}

	if info.IsDir() {
		match := regexp.MustCompile(`\((\d+)\)$`).FindStringSubmatch(info.Name())
		if len(match) > 1 {
			number, _ := strconv.Atoi(match[1])
			for {
				number++
				newDirName := fmt.Sprintf("%s(%d)", info.Name()[:len(info.Name())-len(match[0])], number)
				newPath := filepath.Join(filepath.Dir(destination), newDirName)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					return newPath, nil
				}
			}
		} else {
			for i := 1; ; i++ {
				newDirName := fmt.Sprintf("%s(%d)", info.Name(), i)
				newPath := filepath.Join(filepath.Dir(destination), newDirName)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					return newPath, nil
				}
			}
		}
	} else {
		baseName := filepath.Base(destination)
		ext := filepath.Ext(baseName)
		fileName := baseName[:len(baseName)-len(ext)]
		match := regexp.MustCompile(`\((\d+)\)$`).FindStringSubmatch(fileName)
		if len(match) > 1 {
			number, _ := strconv.Atoi(match[1])
			for {
				number++
				newFileName := fmt.Sprintf("%s(%d)%s", fileName[:len(fileName)-len(match[0])], number, ext)
				newPath := filepath.Join(filepath.Dir(destination), newFileName)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					return newPath, nil
				}
			}
		} else {
			for i := 1; ; i++ {
				newFileName := fmt.Sprintf("%s(%d)%s", fileName, i, ext)
				newPath := filepath.Join(filepath.Dir(destination), newFileName)
				if _, err := os.Stat(newPath); os.IsNotExist(err) {
					return newPath, nil
				}
			}
		}
	}
}

func pasteFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		outPutLog("Paste file function open file error", err)
	}
	defer srcFile.Close()

	dst, err = renameIfDuplicate(dst)
	if err != nil {
		outPutLog("Paste file function rename error", err)
	}
	dstFile, err := os.Create(dst)
	if err != nil {
		outPutLog("Paste file function create file error", err)
	}
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		outPutLog("Paste file function copy file error", err)
	}
	if err != nil {
		return err
	}
	return nil
}

func pasteDir(src, dst string, id string, m model) (model, error) {
	// Check if destination directory already exists
	dst, err := renameIfDuplicate(dst)
	if err != nil {
		return m, err
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		newPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			newPath, err = renameIfDuplicate(newPath)
			if err != nil {
				return err
			}
			err = os.MkdirAll(newPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			p := m.processBarModel.process[id]
			if m.copyItems.cut {
				p.name = "󰆐 " + filepath.Base(path)
			} else {
				p.name = "󰆏 " + filepath.Base(path)
			}

			if len(channel) < 5 {
				channel <- channelMessage{
					messageId:       id,
					processNewState: p,
				}
			}

			err := pasteFile(path, newPath)
			if err != nil {
				p.state = failure
				channel <- channelMessage{
					messageId:       id,
					processNewState: p,
				}
				return err
			}
			p.done++
			if len(channel) < 5 {
				channel <- channelMessage{
					messageId:       id,
					processNewState: p,
				}
			}
			m.processBarModel.process[id] = p
		}

		return nil
	})

	if err != nil {
		return m, err
	}

	return m, nil
}

func returnMetaData(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	cursor := panel.cursor
	LastTimeCursorMove = [2]int{int(time.Now().UnixMicro()), cursor}
	time.Sleep(150 * time.Millisecond)
	if LastTimeCursorMove[1] != cursor && m.focusPanel != metadataFocus {
		return m
	}
	m.fileMetaData.metaData = m.fileMetaData.metaData[:0]
	id := shortuuid.New()
	if len(panel.element) == 0 {
		channel <- channelMessage{
			messageId:    id,
			loadMetadata: true,
			metadata:     m.fileMetaData.metaData,
		}
		return m
	}
	if len(panel.element[panel.cursor].metaData) != 0 && m.focusPanel != metadataFocus {
		m.fileMetaData.metaData = panel.element[panel.cursor].metaData
		channel <- channelMessage{
			messageId:    id,
			loadMetadata: true,
			metadata:     m.fileMetaData.metaData,
		}
		return m
	}
	filePath := panel.element[panel.cursor].location

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"Link file is broken!(you can only delete this file)", ""})
		return m
	}
	if err != nil {
		outPutLog("Return meta data function get file state error", err)
	}
	if fileInfo.IsDir() {
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderName", fileInfo.Name()})
		if m.focusPanel == metadataFocus {
			m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderSize", formatFileSize(dirSize(filePath))})
		}
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"FolderModifyDate", fileInfo.ModTime().String()})
		channel <- channelMessage{
			messageId:    id,
			loadMetadata: true,
			metadata:     m.fileMetaData.metaData,
		}
		return m
	}

	if Config.Metadata {
		fileInfos := et.ExtractMetadata(filePath)

		for _, fileInfo := range fileInfos {
			if fileInfo.Err != nil {
				outPutLog("Return meta data function error", fileInfo, fileInfo.Err)
				continue
			}

			for k, v := range fileInfo.Fields {
				temp := [2]string{k, fmt.Sprintf("%v", v)}
				m.fileMetaData.metaData = append(m.fileMetaData.metaData, temp)
			}
		}
	} else {
		fileName := [2]string{"FileName", fileInfo.Name()}
		fileSize := [2]string{"FileSize", formatFileSize(fileInfo.Size())}
		fileModifyData := [2]string{"FileModifyDate", fileInfo.ModTime().String()}
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, fileName, fileSize, fileModifyData)
	}

	channel <- channelMessage{
		messageId:    id,
		loadMetadata: true,
		metadata:     m.fileMetaData.metaData,
	}

	panel.element[panel.cursor].metaData = m.fileMetaData.metaData
	return m
}

func formatFileSize(size int64) string {
	if size == 0 {
		return "0B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	unitIndex := int(math.Floor(math.Log(float64(size)) / math.Log(1024)))
	adjustedSize := float64(size) / math.Pow(1024, float64(unitIndex))

	return fmt.Sprintf("%.2f %s", adjustedSize, units[unitIndex])
}

func dirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			outPutLog("Dir size function error", err)
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size
}

func countFiles(dirPath string) (int, error) {
	count := 0

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})

	return count, err
}

func loadConfigFile(dir string) (toggleDotFileBool bool, firstFilePanelDir string) {
	var err error

	logOutput, err = os.OpenFile(SuperFileCacheDir+logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while opening superfile.log file: %v", err)
	}

	data, err := os.ReadFile(SuperFileMainDir + configFile)
	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}
	err = toml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Error decoding config json( your config file may have misconfigured ): %v", err)
	}

	data, err = os.ReadFile(SuperFileMainDir + hotkeysFile)
	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}
	err = toml.Unmarshal(data, &hotkeys)
	if err != nil {
		log.Fatalf("Error decoding config json( your config file may have misconfigured ): %v", err)
	}

	data, err = os.ReadFile(SuperFileMainDir + themeFolder + "/" + Config.Theme + ".toml")
	if err != nil {
		log.Fatalf("Theme file doesn't exist: %v", err)
	}

	err = toml.Unmarshal(data, &theme)
	if err != nil {
		log.Fatalf("Error while decoding theme json( Your theme file may have errors ): %v", err)
	}
	toggleDotFileData, err := os.ReadFile(SuperFileDataDir + toggleDotFile)
	if err != nil {
		outPutLog("Error while reading toggleDotFile data error:", err)
	}
	if string(toggleDotFileData) == "true" {
		toggleDotFileBool = true
	} else if string(toggleDotFileData) == "false" {
		toggleDotFileBool = false
	}
	LoadThemeConfig()
	
	if Config.Metadata {
		et, err = exiftool.NewExiftool()
		if err != nil {
			outPutLog("Initial model function init exiftool error", err)
		}
	}

	firstFilePanelDir = HomeDir
	if dir != "" {
		firstFilePanelDir, err = filepath.Abs(dir)
		if err != nil {
			firstFilePanelDir = HomeDir
		}
	}
	return toggleDotFileBool, firstFilePanelDir
}

func unzip(src, dest string) error {
	id := shortuuid.New()
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	totalFiles := len(r.File)
	// progessbar
	prog := progress.New(generateGradientColor())
	prog.PercentageStyle = footerStyle
	// channel message
	p := process{
		name:     "unzip file",
		progress: prog,
		state:    inOperation,
		total:    totalFiles,
		done:     0,
	}
	if _, err := os.Stat(filepath.Join(dest, filepath.Base(src))); os.IsExist(err) {
		p.state = failure
		p.name = "󰛫 Directory already exist"
		channel <- channelMessage{
			messageId:       id,
			processNewState: p,
		}
		return nil
	}
	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)

			if err != nil {

				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		p.name = "󰛫 " + f.Name
		if len(channel) < 3 {
			channel <- channelMessage{
				messageId:       id,
				processNewState: p,
			}
		}
		err := extractAndWriteFile(f)
		if err != nil {
			p.state = failure
			channel <- channelMessage{
				messageId:       id,
				processNewState: p,
			}
			return err
		}
		p.done++
		if len(channel) < 3 {
			channel <- channelMessage{
				messageId:       id,
				processNewState: p,
			}
		}
	}

	p.total = totalFiles
	p.state = successful
	channel <- channelMessage{
		messageId:       id,
		processNewState: p,
	}

	return nil
}

func zipSource(source, target string) error {
	id := shortuuid.New()
	prog := progress.New()
	prog.PercentageStyle = footerStyle

	totalFiles, err := countFiles(source)

	if err != nil {
		outPutLog("Zip file count files error: ", err)
	}

	p := process{
		name:     "zip files",
		progress: prog,
		state:    inOperation,
		total:    totalFiles,
		done:     0,
	}

	_, err = os.Stat(target)
	if os.IsExist(err) {
		p.name = "󰗄 File already exist"
		channel <- channelMessage{
			messageId:       id,
			processNewState: p,
		}
		return nil
	}

	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		p.name = "󰗄 " + filepath.Base(path)
		if len(channel) < 5 {
			channel <- channelMessage{
				messageId:       id,
				processNewState: p,
			}
		}

		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		if err != nil {
			return err
		}
		p.done++
		if len(channel) < 5 {
			channel <- channelMessage{
				messageId:       id,
				processNewState: p,
			}
		}
		return nil
	})

	if err != nil {
		outPutLog("Error while zip file:", err)
		p.state = failure
		channel <- channelMessage{
			messageId:       id,
			processNewState: p,
		}
	}
	p.state = successful
	p.done = totalFiles
	channel <- channelMessage{
		messageId:       id,
		processNewState: p,
	}

	return nil
}

func generateSearchBar() textinput.Model {
	ti := textinput.New()
	ti.Cursor.Style = footerCursorStyle
	ti.Cursor.TextStyle = footerStyle
	ti.TextStyle = filePanelStyle
	ti.Prompt = filePanelTopDirectoryIconStyle.Render(" ")
	ti.Cursor.Blink = true
	ti.PlaceholderStyle = filePanelStyle
	ti.Placeholder = "(" + hotkeys.SearchBar[0] + ") Type something"
	ti.Blur()
	ti.CharLimit = 156
	return ti
}
