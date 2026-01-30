#!/bin/bash
#
# Script to capture all required screenshots for the Kanagawa theme PR
# Usage: ./scripts/capture-kanagawa-screenshots.sh
#
# Requirements:
# - Go toolchain working (no version mismatch)
# - superfile built with Kanagawa theme
# - Terminal with window size of at least 120x40
#

set -e

THEME_NAME="kanagawa"
OUTPUT_DIR="asset/theme-screenshots"
BINARY="./bin/spf"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== Kanagawa Theme Screenshot Capture Script ===${NC}"
echo ""

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    echo -e "${YELLOW}Binary not found at $BINARY${NC}"
    echo "Building superfile..."
    make build
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"
echo -e "${GREEN}✓ Output directory: $OUTPUT_DIR${NC}"
echo ""

# Instructions for each screenshot
echo -e "${BLUE}Screenshots needed:${NC}"
echo ""
echo "1. Full view - Show all panels (sidebar, file previewer, process panel, metadata panel, clipboard)"
echo "   - Make sure file previewer has content"
echo "   - Process panel should show at least one process"
echo "   - Clipboard should have at least one entry"
echo ""
echo "2. Sidebar focused - Press Tab to focus sidebar, note the border color"
echo ""
echo "3. Process panel focused - Press Tab until process panel is focused"
echo ""
echo "4. Help menu - Press ? to open help"
echo ""
echo "5. New file popup - Press Ctrl+n to create new file"
echo ""
echo "6. Image preview - Navigate to an image file and preview it"
echo ""
echo "7. Successful shell command - Run a command that succeeds (e.g., :echo hello)"
echo ""
echo "8. Failed shell command - Run a command that fails (e.g., :false)"
echo ""

# Screenshot function using macOS screencapture
capture_screenshot() {
    local name="$1"
    local description="$2"

    echo -e "${YELLOW}Ready to capture: $description${NC}"
    echo "Press Enter when you have superfile showing the correct view..."
    read -r

    local timestamp=$(date +%Y%m%d-%H%M%S)
    local filename="${OUTPUT_DIR}/${name}-${timestamp}.png"

    if command -v screencapture &> /dev/null; then
        screencapture -x "$filename"
        echo -e "${GREEN}✓ Captured: $filename${NC}"
    else
        echo -e "${YELLOW}Warning: screencapture not found. Please take manual screenshot.${NC}"
        echo "  Save as: $filename"
    fi
    echo ""
}

# Run superfile and capture screenshots
echo -e "${BLUE}Launching superfile with Kanagawa theme...${NC}"
echo "Configure superfile to use theme = \"kanagawa\" in config first"
echo ""
echo "Once superfile is running, follow the prompts to capture each screenshot"
echo ""
echo "Press Enter to start..."
read -r

# Capture each screenshot
capture_screenshot "01-full-view" "Full view of superfile with all panels"
capture_screenshot "02-sidebar-focused" "Sidebar focused (border color)"
capture_screenshot "03-process-panel-focused" "Process panel focused (border color)"
capture_screenshot "04-help-menu" "Help menu (press ?)"
capture_screenshot "05-new-file-popup" "New file popup (press Ctrl+n)"
capture_screenshot "06-image-preview" "Image being previewed"
capture_screenshot "07-successful-command" "Successful shell command"
capture_screenshot "08-failed-command" "Failed shell command"

echo -e "${GREEN}=== All screenshots captured! ===${NC}"
echo ""
echo "Screenshots saved to: $OUTPUT_DIR/"
echo ""
echo "Next steps:"
echo "1. Review the screenshots"
echo "2. Copy the best full-view screenshot to asset/theme/kanagawa.png"
echo "3. Add all screenshots to the PR description"
