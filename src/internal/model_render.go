package internal

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"image"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/yorukot/ansichroma"

	"github.com/yorukot/superfile/src/config/icon"
)

func (m *model) sidebarRender() string {
	return m.sidebarModel.Render(m.mainPanelHeight, m.focusPanel == sidebarFocus,
		m.fileModel.filePanels[m.filePanelFocusIndex].location)
}

// This also modifies the m.fileModel.filePanels, which it should not
// what modifications we do on this model object are of no consequence.
// Since bubblea passed this 'model' by value in View() function.
func (m *model) filePanelRender() string {
	f := make([]string, len(m.fileModel.filePanels))
	for i, filePanel := range m.fileModel.filePanels {
		// check if cursor or render out of range
		// TODO - instead of this, have a filepanel.validateAndFix(), and log Error
		// This should not ever happen
		if filePanel.cursor > len(filePanel.element)-1 {
			filePanel.cursor = 0
			filePanel.render = 0
		}
		m.fileModel.filePanels[i] = filePanel

		// TODO : Move this to a utility function and clarify the calculation via comments
		// Maybe even write unit tests
		var filePanelWidth int
		if (m.fullWidth-common.Config.SidebarWidth-(4+(len(m.fileModel.filePanels)-1)*2))%len(
			m.fileModel.filePanels,
		) != 0 &&
			i == len(m.fileModel.filePanels)-1 {
			if m.fileModel.filePreview.open {
				filePanelWidth = m.fileModel.width
			} else {
				filePanelWidth = (m.fileModel.width + (m.fullWidth-common.Config.SidebarWidth-
					(4+(len(m.fileModel.filePanels)-1)*2))%len(m.fileModel.filePanels))
			}
		} else {
			filePanelWidth = m.fileModel.width
		}

		f[i] = filePanel.Render(m.mainPanelHeight, filePanelWidth, filePanel.focusType != noneFocus)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, f...)
}

func (panel *filePanel) Render(mainPanelHeight int, filePanelWidth int, focussed bool) string {
	r := ui.FilePanelRenderer(mainPanelHeight+2, filePanelWidth+2, focussed)

	panel.renderTopBar(r, filePanelWidth)
	panel.renderSearchBar(r)
	panel.renderFooter(r)
	panel.renderFileEntries(r, mainPanelHeight, filePanelWidth)

	return r.Render()
}

func (panel *filePanel) renderTopBar(r *rendering.Renderer, filePanelWidth int) {
	// TODO - Add ansitruncate left in renderer and remove truncation here
	truncatedPath := common.TruncateTextBeginning(panel.location, filePanelWidth-4, "...")
	r.AddLines(common.FilePanelTopDirectoryIcon + common.FilePanelTopPathStyle.Render(truncatedPath))
	r.AddSection()
}

func (panel *filePanel) renderSearchBar(r *rendering.Renderer) {
	r.AddLines(" " + panel.searchBar.View())
}

// TODO : Unit test this
func (panel *filePanel) renderFooter(r *rendering.Renderer) {
	sortLabel, sortIcon := panel.getSortInfo()
	modeLabel, modeIcon := panel.getPanelModeInfo()
	cursorStr := panel.getCursorString()

	if common.Config.Nerdfont {
		sortLabel = sortIcon + icon.Space + sortLabel
		modeLabel = modeIcon + icon.Space + modeLabel
	} else {
		// TODO : Figure out if we can set icon.Space to " " if nerdfont is false
		// That would simplify code
		sortLabel = sortIcon + " " + sortLabel
	}

	if common.Config.ShowPanelFooterInfo {
		r.SetBorderInfoItems(sortLabel, modeLabel, cursorStr)
		if r.AreInfoItemsTruncated() {
			r.SetBorderInfoItems(sortIcon, modeIcon, cursorStr)
		}
	} else {
		r.SetBorderInfoItems(cursorStr)
	}
}

func (panel *filePanel) renderFileEntries(r *rendering.Renderer, mainPanelHeight, filePanelWidth int) {
	if len(panel.element) == 0 {
		r.AddLines(common.FilePanelNoneText)
		return
	}

	end := min(panel.render+panelElementHeight(mainPanelHeight), len(panel.element))

	for i := panel.render; i < end; i++ {
		// TODO : Fix this, this is O(n^2) complexity. Considered a file panel with 200 files, and 100 selected
		// We will be doing a search in 100 item slice for all 200 files.
		isSelected := arrayContains(panel.selected, panel.element[i].location)

		if panel.renaming && i == panel.cursor {
			r.AddLines(panel.rename.View())
			continue
		}

		cursor := " "
		if i == panel.cursor && !panel.searchBar.Focused() {
			cursor = icon.Cursor
		}

		// Performance TODO: Remove or cache this if not needed at render time
		// This will unnecessarily slow down rendering. There should be a way to avoid this at render
		_, err := os.ReadDir(panel.element[i].location)
		dirExists := err == nil || panel.element[i].directory

		renderedName := common.PrettierName(
			panel.element[i].name,
			filePanelWidth-5,
			dirExists,
			isSelected,
			common.FilePanelBGColor,
		)

		r.AddLines(common.FilePanelCursorStyle.Render(cursor+" ") + renderedName)
	}
}

func (panel *filePanel) getSortInfo() (string, string) {
	opts := panel.sortOptions.data
	selected := opts.options[opts.selected]
	label := selected
	if selected == string(sortingDateModified) {
		label = "Date"
	}

	iconStr := icon.SortAsc

	if opts.reversed {
		iconStr = icon.SortDesc
	}
	return label, iconStr
}

func (panel *filePanel) getPanelModeInfo() (string, string) {
	switch panel.panelMode {
	case browserMode:
		return "Browser", icon.Browser
	case selectMode:
		return "Select", icon.Select
	default:
		return "", ""
	}
}

func (panel *filePanel) getCursorString() string {
	cursor := panel.cursor
	if len(panel.element) > 0 {
		cursor++ // Convert to 1-based
	}
	return fmt.Sprintf("%d/%d", cursor, len(panel.element))
}

func (m *model) processBarRender() string {
	return m.processBarModel.Render(m.focusPanel == processBarFocus)
}

func (m *model) clipboardRender() string {
	// render
	var bottomWidth int
	if m.fullWidth%3 != 0 {
		bottomWidth = utils.FooterWidth(m.fullWidth + m.fullWidth%3 + 2)
	} else {
		bottomWidth = utils.FooterWidth(m.fullWidth)
	}
	r := ui.ClipboardRenderer(m.footerHeight+2, bottomWidth+2)
	if len(m.copyItems.items) == 0 {
		// TODO move this to a string
		r.AddLines("", " "+icon.Error+"  No content in clipboard")
	} else {
		for i := 0; i < len(m.copyItems.items) && i < m.footerHeight; i++ {
			if i == m.footerHeight-1 && i != len(m.copyItems.items)-1 {
				// Last Entry we can render, but there are more that one left
				r.AddLines(strconv.Itoa(len(m.copyItems.items)-i) + " item left....")
			} else {
				fileInfo, err := os.Stat(m.copyItems.items[i])
				if err != nil {
					slog.Error("Clipboard render function get item state ", "error", err)
				}
				if !os.IsNotExist(err) {
					// TODO : There is an inconsistency in parameter that is being passed,
					// and its name in ClipboardPrettierName function
					r.AddLines(common.ClipboardPrettierName(m.copyItems.items[i],
						utils.FooterWidth(m.fullWidth)-3, fileInfo.IsDir(), false))
				}
			}
		}
	}
	return r.Render()
}

func (m *model) terminalSizeWarnRender() string {
	fullWidthString := strconv.Itoa(m.fullWidth)
	fullHeightString := strconv.Itoa(m.fullHeight)
	minimumWidthString := strconv.Itoa(common.MinimumWidth)
	minimumHeightString := strconv.Itoa(common.MinimumHeight)
	if m.fullHeight < common.MinimumHeight {
		fullHeightString = common.TerminalTooSmall.Render(fullHeightString)
	}
	if m.fullWidth < common.MinimumWidth {
		fullWidthString = common.TerminalTooSmall.Render(fullWidthString)
	}
	fullHeightString = common.TerminalCorrectSize.Render(fullHeightString)
	fullWidthString = common.TerminalCorrectSize.Render(fullWidthString)

	heightString := common.MainStyle.Render(" Height = ")
	return common.FullScreenStyle(m.fullHeight, m.fullWidth).Render(`Terminal size too small:` + "\n" +
		"Width = " + fullWidthString +
		heightString + fullHeightString + "\n\n" +
		"Needed for current config:" + "\n" +
		"Width = " + common.TerminalCorrectSize.Render(minimumWidthString) +
		heightString + common.TerminalCorrectSize.Render(minimumHeightString))
}

func (m *model) terminalSizeWarnAfterFirstRender() string {
	minimumWidthInt := common.Config.SidebarWidth + 20*len(m.fileModel.filePanels) + 20 - 1
	minimumWidthString := strconv.Itoa(minimumWidthInt)
	fullWidthString := strconv.Itoa(m.fullWidth)
	fullHeightString := strconv.Itoa(m.fullHeight)
	minimumHeightString := strconv.Itoa(common.MinimumHeight)

	if m.fullHeight < common.MinimumHeight {
		fullHeightString = common.TerminalTooSmall.Render(fullHeightString)
	}
	if m.fullWidth < minimumWidthInt {
		fullWidthString = common.TerminalTooSmall.Render(fullWidthString)
	}
	fullHeightString = common.TerminalCorrectSize.Render(fullHeightString)
	fullWidthString = common.TerminalCorrectSize.Render(fullWidthString)

	heightString := common.MainStyle.Render(" Height = ")
	return common.FullScreenStyle(m.fullHeight, m.fullWidth).Render(`You change your terminal size too small:` + "\n" +
		"Width = " + fullWidthString +
		heightString + fullHeightString + "\n\n" +
		"Needed for current config:" + "\n" +
		"Width = " + common.TerminalCorrectSize.Render(minimumWidthString) +
		heightString + common.TerminalCorrectSize.Render(minimumHeightString))
}

func (m *model) typineModalRender() string {
	previewPath := filepath.Join(m.typingModal.location, m.typingModal.textInput.Value())

	fileLocation := common.FilePanelTopDirectoryIconStyle.Render(" "+icon.Directory+icon.Space) +
		common.FilePanelTopPathStyle.Render(
			common.TruncateTextBeginning(previewPath, common.ModalWidth-4, "..."),
		) + "\n"

	confirm := common.ModalConfirm.Render(" (" + common.Hotkeys.ConfirmTyping[0] + ") Create ")
	cancel := common.ModalCancel.Render(" (" + common.Hotkeys.CancelTyping[0] + ") Cancel ")

	tip := confirm +
		lipgloss.NewStyle().Background(common.ModalBGColor).Render("           ") +
		cancel

	var err string
	if m.typingModal.errorMesssage != "" {
		err = "\n\n" + common.ModalErrorStyle.Render(m.typingModal.errorMesssage)
	}
	// TODO : Move this all to rendering package to avoid specifying newlines manually
	return common.ModalBorderStyle(common.ModalHeight, common.ModalWidth).
		Render(fileLocation + "\n" + m.typingModal.textInput.View() + "\n\n" + tip + err)
}

func (m *model) introduceModalRender() string {
	title := common.SidebarTitleStyle.Render(" Thanks for using superfile!!") +
		common.ModalStyle.Render("\n You can read the following information before starting to use it!")
	vimUserWarn := common.ProcessErrorStyle.Render("  ** Very importantly ** If you are a Vim/Nvim user, go to:\n" +
		"  https://superfile.netlify.app/configure/custom-hotkeys/ to change your hotkey settings!")
	subOne := common.SidebarTitleStyle.Render("  (1)") +
		common.ModalStyle.Render(" If this is your first time, make sure you read:\n"+
			"      https://superfile.netlify.app/getting-started/tutorial/")
	subTwo := common.SidebarTitleStyle.Render("  (2)") +
		common.ModalStyle.Render(" If you forget the relevant keys during use,\n"+
			"      you can press \"?\" (shift+/) at any time to query the keys!")
	subThree := common.SidebarTitleStyle.Render("  (3)") +
		common.ModalStyle.Render(" For more customization you can refer to:\n"+
			"      https://superfile.netlify.app/")
	subFour := common.SidebarTitleStyle.Render("  (4)") +
		common.ModalStyle.Render(" Thank you again for using superfile.\n"+
			"      If you have any questions, please feel free to ask at:\n"+
			"      https://github.com/yorukot/superfile\n"+
			"      Of course, you can always open a new issue to share your idea \n"+
			"      or report a bug!")
	return common.FirstUseModal(m.helpMenu.height, m.helpMenu.width).
		Render(title + "\n\n" + vimUserWarn + "\n\n" + subOne + "\n\n" +
			subTwo + "\n\n" + subThree + "\n\n" + subFour + "\n\n")
}

func (m *model) promptModalRender() string {
	return m.promptModal.Render()
}

func (m *model) helpMenuRender() string {
	maxKeyLength := 0

	for _, data := range m.helpMenu.data {
		totalKeyLen := 0
		for _, key := range data.hotkey {
			totalKeyLen += len(key)
		}
		saprateLen := len(data.hotkey) - 1*3
		if data.subTitle == "" && totalKeyLen+saprateLen > maxKeyLength {
			maxKeyLength = totalKeyLen + saprateLen
		}
	}

	valueLength := m.helpMenu.width - maxKeyLength - 2
	if valueLength < m.helpMenu.width/2 {
		valueLength = m.helpMenu.width/2 - 2
	}

	totalTitleCount := 0
	cursorBeenTitleCount := 0

	for i, data := range m.helpMenu.data {
		if data.subTitle != "" {
			if i < m.helpMenu.cursor {
				cursorBeenTitleCount++
			}
			totalTitleCount++
		}
	}

	renderHotkeyLength := m.getRenderHotkeyLengthHelpmenuModal()
	helpMenuContent := m.getHelpMenuContent(renderHotkeyLength, valueLength)

	bottomBorder := common.GenerateFooterBorder(fmt.Sprintf("%s/%s",
		strconv.Itoa(m.helpMenu.cursor+1-cursorBeenTitleCount),
		strconv.Itoa(len(m.helpMenu.data)-totalTitleCount)), m.helpMenu.width-2)

	return common.HelpMenuModalBorderStyle(m.helpMenu.height, m.helpMenu.width, bottomBorder).Render(helpMenuContent)
}

func (m *model) getRenderHotkeyLengthHelpmenuModal() int {
	renderHotkeyLength := 0
	for i := m.helpMenu.renderIndex; i < m.helpMenu.height+m.helpMenu.renderIndex && i < len(m.helpMenu.data); i++ {
		hotkey := ""

		if m.helpMenu.data[i].subTitle != "" {
			continue
		}

		for i, key := range m.helpMenu.data[i].hotkey {
			if i != 0 {
				hotkey += " | "
			}
			hotkey += key
		}

		renderHotkeyLength = max(renderHotkeyLength, len(common.HelpMenuHotkeyStyle.Render(hotkey)))
	}
	return renderHotkeyLength
}

func (m *model) getHelpMenuContent(renderHotkeyLength int, valueLength int) string {
	helpMenuContent := ""
	for i := m.helpMenu.renderIndex; i < m.helpMenu.height+m.helpMenu.renderIndex && i < len(m.helpMenu.data); i++ {
		if i != m.helpMenu.renderIndex {
			helpMenuContent += "\n"
		}

		if m.helpMenu.data[i].subTitle != "" {
			helpMenuContent += common.HelpMenuTitleStyle.Render(" " + m.helpMenu.data[i].subTitle)
			continue
		}

		hotkey := ""
		description := common.TruncateText(m.helpMenu.data[i].description, valueLength, "...")

		for i, key := range m.helpMenu.data[i].hotkey {
			if i != 0 {
				hotkey += " | "
			}
			hotkey += key
		}

		cursor := "  "
		if m.helpMenu.cursor == i {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor + " ")
		}
		helpMenuContent += cursor + common.ModalStyle.Render(fmt.Sprintf("%*s%s", renderHotkeyLength,
			common.HelpMenuHotkeyStyle.Render(hotkey+" "), common.ModalStyle.Render(description)))
	}
	return helpMenuContent
}

func (m *model) sortOptionsRender() string {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	sortOptionsContent := common.ModalTitleStyle.Render(" Sort Options") + "\n\n"
	for i, option := range panel.sortOptions.data.options {
		cursor := " "
		if i == panel.sortOptions.cursor {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor)
		}
		sortOptionsContent += cursor + common.ModalStyle.Render(" "+option) + "\n"
	}
	bottomBorder := common.GenerateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(panel.sortOptions.cursor+1),
		strconv.Itoa(len(panel.sortOptions.data.options))), panel.sortOptions.width-2)

	return common.SortOptionsModalBorderStyle(panel.sortOptions.height, panel.sortOptions.width,
		bottomBorder).Render(sortOptionsContent)
}

func readFileContent(filepath string, maxLineLength int, previewLine int) (string, error) {
	var resultBuilder strings.Builder
	file, err := os.Open(filepath)
	if err != nil {
		return resultBuilder.String(), err
	}
	defer file.Close()

	reader := transform.NewReader(file, unicode.BOMOverride(unicode.UTF8.NewDecoder()))
	scanner := bufio.NewScanner(reader)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = ansi.Truncate(line, maxLineLength, "")
		resultBuilder.WriteString(line)
		resultBuilder.WriteRune('\n')
		lineCount++
		if previewLine > 0 && lineCount >= previewLine {
			break
		}
	}
	// returns the first non-EOF error that was encountered by the [Scanner]
	return resultBuilder.String(), scanner.Err()
}

func (m *model) filePreviewPanelRender() string {
	return m.fileModel.filePreview.content
}

// Helper function to handle file info errors
func (m *model) renderFileInfoError(r *rendering.Renderer, _ lipgloss.Style, err error) string {
	slog.Error("Error get file info", "error", err)
	return r.Render()
}

// Helper function to handle unsupported formats
func (m *model) renderUnsupportedFormat(box lipgloss.Style) string {
	return box.Render(common.FilePreviewUnsupportedFormatText)
}

// Helper function to handle unsupported mode
func (m *model) renderUnsupportedFileMode(r *rendering.Renderer) string {
	r.AddLines(common.FilePreviewUnsupportedFileMode)
	return r.Render()
}

// Helper function to handle directory preview
func (m *model) renderDirectoryPreview(r *rendering.Renderer, itemPath string, previewHeight int) string {
	files, err := os.ReadDir(itemPath)
	if err != nil {
		slog.Error("Error render directory preview", "error", err)
		r.AddLines(common.FilePreviewDirectoryUnreadableText)
		return r.Render()
	}

	if len(files) == 0 {
		r.AddLines(common.FilePreviewEmptyText)
		return r.Render()
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		}
		if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return files[i].Name() < files[j].Name()
	})

	for i := 0; i < previewHeight && i < len(files); i++ {
		file := files[i]
		style := common.GetElementIcon(file.Name(), file.IsDir(), common.Config.Nerdfont)
		res := lipgloss.NewStyle().Foreground(lipgloss.Color(style.Color)).Background(common.FilePanelBGColor).
			Render(style.Icon+" ") + common.FilePanelStyle.Render(file.Name())
		r.AddLines(res)
	}
	return r.Render()
}

// Helper function to handle image preview
func (m *model) renderImagePreview(box lipgloss.Style, itemPath string, previewWidth,
	previewHeight int, sideAreaWidth int) string {
	if !m.fileModel.filePreview.open {
		return box.Render("\n --- Preview panel is closed ---")
	}

	if !common.Config.ShowImagePreview {
		return box.Render("\n --- Image preview is disabled ---")
	}

	// Use the new auto-detection function to choose the best renderer
	imageRender, err := m.imagePreviewer.ImagePreview(itemPath, previewWidth, previewHeight,
		common.Theme.FilePanelBG, sideAreaWidth)
	if errors.Is(err, image.ErrFormat) {
		return box.Render("\n --- " + icon.Error + " Unsupported image formats ---")
	}

	if err != nil {
		slog.Error("Error convert image to ansi", "error", err)
		return box.Render("\n --- " + icon.Error + " Error convert image to ansi ---")
	}

	// Check if this looks like Kitty protocol output (starts with escape sequences)
	// For Kitty protocol, avoid using lipgloss alignment to prevent layout drift
	if strings.HasPrefix(imageRender, "\x1b_G") {
		rendered := common.FilePreviewBox(previewHeight, previewWidth).Render(imageRender)
		return rendered
	}

	// For ANSI output, we can safely use vertical alignment
	return box.AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center).Render(imageRender)
}

// Helper function to handle text file preview
func (m *model) renderTextPreview(r *rendering.Renderer, box lipgloss.Style, itemPath string,
	previewWidth, previewHeight int) string {
	format := lexers.Match(filepath.Base(itemPath))
	if format == nil {
		isText, err := common.IsTextFile(itemPath)
		if err != nil {
			slog.Error("Error while checking text file", "error", err)
			return box.Render(common.FilePreviewError)
		} else if !isText {
			return box.Render(common.FilePreviewUnsupportedFormatText)
		}
	}

	fileContent, err := readFileContent(itemPath, previewWidth, previewHeight)
	if err != nil {
		slog.Error("Error open file", "error", err)
		return box.Render(common.FilePreviewError)
	}

	if fileContent == "" {
		return box.Render(common.FilePreviewEmptyText)
	}

	if format != nil {
		background := ""
		if !common.Config.TransparentBackground {
			background = common.Theme.FilePanelBG
		}
		if common.Config.CodePreviewer == "bat" {
			if batCmd == "" {
				return box.Render("\n --- " + icon.Error +
					" 'bat' is not installed or not found. ---\n --- Cannot render file preview. ---")
			}
			fileContent, err = getBatSyntaxHighlightedContent(itemPath, previewHeight, background)
		} else {
			fileContent, err = ansichroma.HightlightString(fileContent, format.Config().Name,
				common.Theme.CodeSyntaxHighlightTheme, background)
		}
		if err != nil {
			slog.Error("Error render code highlight", "error", err)
			return box.Render("\n" + common.FilePreviewError)
		}
	}

	r.AddLines(fileContent)
	return r.Render()
}
func (m *model) filePreviewPanelRenderWithDimensions() string {
	panel := m.getFocusedFilePanel()
	if len(panel.element) == 0 {
		return m.fileModel.filePreview.RenderText("")
	}
	itemPath := panel.getSelectedItem().location
	return m.filePreviewPanelRenderWithPath(itemPath)
}

// TODO: Avoid accessing the model here.
// There is used in tea.Cmd for async rendering. Avoid access inside async
// command closures to reduce data-race risk. Require immutable values as parameter
// (e.g., previewW, previewH, sideAreaWidth, imagePreviewer, theme colors) and pass
// them to a pure renderer that doesn’t touch m.
func (m *model) filePreviewPanelRenderWithPath(itemPath string) string {
	previewHeight := m.fileModel.filePreview.height
	previewWidth := m.fileModel.filePreview.width

	box := common.FilePreviewBox(previewHeight, previewWidth)
	r := ui.FilePreviewPanelRenderer(previewHeight, previewWidth)
	clearCmd := m.imagePreviewer.ClearKittyImages()

	// This could create errors if panel.cursor ever becomes negative, or goes out of bounds
	// We should have a panel validation function in our View() function
	// Panel is a full fledged object with own state, its accessed and modified so many times.
	// Ideally we dont should never access data from it via directly accessing its variables
	// TODO : Instead we should have helper functions for panel object and access data that way
	// like panel.GetCurrentSelectedElem() . This abstration of implemetation of panel is needed.
	// Now this lack of abstraction has caused issues ( See PR#730 ) . And now
	// someone needs to scan through the entire codebase to figure out which access of panel
	// data is causing crash.

	fileInfo, infoErr := os.Stat(itemPath)
	if infoErr != nil {
		return m.renderFileInfoError(r, box, infoErr) + clearCmd
	}
	slog.Debug("Attempting to render preview", "itemPath", itemPath,
		"mode", fileInfo.Mode().String(), "isRegular", fileInfo.Mode().IsRegular())

	// For non regular files which are not directories Dont try to read them
	// See Issue
	if !fileInfo.Mode().IsRegular() && (fileInfo.Mode()&fs.ModeDir) == 0 {
		return m.renderUnsupportedFileMode(r) + clearCmd
	}

	ext := filepath.Ext(itemPath)
	if slices.Contains(common.UnsupportedPreviewFormats, ext) {
		return m.renderUnsupportedFormat(box) + clearCmd
	}

	if fileInfo.IsDir() {
		return m.renderDirectoryPreview(r, itemPath, previewHeight) + clearCmd
	}

	if isImageFile(itemPath) {
		return m.renderImagePreview(box, itemPath, previewWidth, previewHeight, m.fullWidth-previewWidth+1)
	}

	return m.renderTextPreview(r, box, itemPath, previewWidth, previewHeight) + clearCmd
}

func getBatSyntaxHighlightedContent(itemPath string, previewLine int, background string) (string, error) {
	// --plain: use the plain style without line numbers and decorations
	// --force-colorization: force colorization for non-interactive shell output
	// --line-range <:m>: only read from line 1 to line "m"
	batArgs := []string{itemPath, "--plain", "--force-colorization", "--line-range", fmt.Sprintf(":%d", previewLine-1)}

	// set timeout for the external command execution to 500ms max
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, batCmd, batArgs...)

	fileContentBytes, err := cmd.Output()
	if err != nil {
		slog.Error("Error render code highlight", "error", err)
		return "", err
	}

	fileContent := string(fileContentBytes)
	if !common.Config.TransparentBackground {
		fileContent = setBatBackground(fileContent, background)
	}
	return fileContent, nil
}

func setBatBackground(input string, background string) string {
	tokens := strings.Split(input, "\x1b[0m")
	backgroundStyle := lipgloss.NewStyle().Background(lipgloss.Color(background))
	for idx, token := range tokens {
		tokens[idx] = backgroundStyle.Render(token)
	}
	return strings.Join(tokens, "\x1b[0m")
}
