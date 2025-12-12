# Hotkey Help Menu Test Plan

## Issue Description
Fixed issue #1047: Missing hotkeys in the help menu

## Changes Made
Added 6 missing hotkeys to the help menu in `src/internal/default_config.go`:
1. **OpenCommandLine** (line ~103) - Open command line with 'c'
2. **OpenSPFPrompt** (line ~108) - Open SPF prompt with ':'
3. **OpenZoxide** (line ~113) - Open zoxide with 'z'
4. **PageUp** (line ~189) - Navigate page up
5. **PageDown** (line ~194) - Navigate page down
6. **CopyPWD** (line ~282) - Copy current working directory with 'W'

## Manual Testing Instructions

### Test 1: Verify Help Menu Display
1. Build the application: `make build` or `go build -o bin/spf ./src/cmd`
2. Run superfile: `./bin/spf`
3. Press '?' to open the help menu
4. Verify the following hotkeys are visible in the help menu:
   - **Command Operations** section:
     - 'c' - Open command line
     - ':' - Open SPF prompt
     - 'z' - Open zoxide
   - **Navigation** section:
     - 'Page Up' - Page up
     - 'Page Down' - Page down
   - **Clipboard** section:
     - 'W' - Copy present working directory

### Test 2: Verify Hotkey Functionality
1. Test Command Operations:
   - Press 'c' and verify command line opens at the bottom
   - Press ':' and verify SPF prompt opens
   - Press 'z' and verify zoxide interface opens (if zoxide is installed)

2. Test Navigation:
   - Press 'Page Up' and verify the list scrolls up by a page
   - Press 'Page Down' and verify the list scrolls down by a page

3. Test Clipboard:
   - Press 'W' and verify the current directory path is copied to clipboard
   - Paste in another application to confirm the path was copied

### Test 3: Regression Testing
1. Verify all existing hotkeys still work correctly
2. Verify no conflicts between hotkeys
3. Check that help menu layout remains readable and organized

## Unit Testing
No unit tests required for this change as it only adds data entries to the help menu configuration. The change is purely presentational and doesn't affect business logic.

## Expected Results
- All 6 missing hotkeys should appear in the help menu
- All hotkeys should function as documented
- No regressions in existing functionality
- Help menu should remain well-organized and readable
