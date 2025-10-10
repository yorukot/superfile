# zoxide package
This is for the Zoxide navigation modal of superfile

Handles user input for zoxide queries, integrates with the go-zoxide library, and returns navigation actions to the model.

## Features

- Interactive zoxide directory search and navigation
- Real-time suggestions with scores from zoxide database  
- Keyboard navigation with standard superfile hotkeys
- Integration with existing file panel navigation system

## Usage

The zoxide modal is opened by pressing the `z` hotkey and allows users to:
1. Type directory names to search zoxide's database
2. See top 5 matching directories with relevance scores
3. Navigate to selected directory in the current file panel
4. Close the modal with Escape or successful navigation

## Architecture

- `Model`: Main zoxide modal state and behavior
- `HandleUpdate()`: Processes keyboard input and zoxide queries
- `Render()`: Displays search interface and suggestions
- Integration with `*zoxidelib.Client` for zoxide database queries

## Coverage

Current test coverage: **0%**

No tests have been implemented yet for this package. The package is functional and integrated, but lacks unit test coverage.

```bash
cd /path/to/ui/zoxide
# Basic coverage (when tests exist)
go test -cover

# HTML report (when tests exist)
go test -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
```