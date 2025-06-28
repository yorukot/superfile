#!/usr/bin/env bash

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
RUN_TESTSUITE=false
SKIP_TESTS=false
VERBOSE=false

# Function to print colored output
print_step() {
    echo "${BLUE}==>${NC} $1"
}

print_success() {
    echo "${GREEN}✓${NC} $1"
}

print_warning() {
    echo "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo "${RED}✗${NC} $1"
}

# Function to show usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "A comprehensive script for formatting, testing, and building superfile"
    echo ""
    echo "OPTIONS:"
    echo "  -t, --testsuite     Run integration testsuite after unit tests"
    echo "  -s, --skip-tests    Skip unit tests (only format, lint, and build)"
    echo "  -v, --verbose       Enable verbose output"
    echo "  -h, --help          Show this help message"
    echo ""
    echo "STEPS PERFORMED:"
    echo "  1. Tidy Go modules"
    echo "  2. Format code with 'go fmt'"
    echo "  3. Run golangci-lint"
    echo "  4. Run unit tests (unless --skip-tests)"
    echo "  5. Run integration testsuite (if --testsuite)"
    echo "  6. Build spf binary"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--testsuite)
            RUN_TESTSUITE=true
            shift
            ;;
        -s|--skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Set verbose flag for commands if requested
VERBOSE_FLAG=""
if [ "$VERBOSE" = true ]; then
    VERBOSE_FLAG="-v"
fi

echo "${BLUE}🚀 Starting superfile development workflow${NC}"
echo ""

# Step 1: Tidy up the go mod
print_step "Tidying Go modules..."
if go mod tidy $VERBOSE_FLAG; then
    print_success "Go modules tidied"
else
    print_error "Failed to tidy Go modules"
    exit 1
fi

# Step 2: Format the code
print_step "Formatting Go code..."
if go fmt ./...; then
    print_success "Code formatted"
else
    print_error "Failed to format code"
    exit 1
fi

# Step 3: Run the linter
print_step "Running golangci-lint..."
if golangci-lint run; then
    print_success "Linting passed"
else
    print_error "Linting failed"
    exit 1
fi

# Step 4: Run unit tests (unless skipped)
if [ "$SKIP_TESTS" = false ]; then
    print_step "Running unit tests..."
    if [ "$VERBOSE" = true ]; then
        if go test -v ./...; then
            print_success "Unit tests passed"
        else
            print_error "Unit tests failed"
            exit 1
        fi
    else
        if go test ./...; then
            print_success "Unit tests passed"
        else
            print_error "Unit tests failed"
            exit 1
        fi
    fi
else
    print_warning "Skipping unit tests"
fi

# Step 5: Run integration testsuite (if requested)
if [ "$RUN_TESTSUITE" = true ]; then
    print_step "Running integration testsuite..."

    # Check if Python is available
    if ! command -v python3 &> /dev/null; then
        print_error "Python3 is required for testsuite but not found"
        exit 1
    fi

    # Check if testsuite requirements are installed
    if [ ! -f "testsuite/requirements.txt" ]; then
        print_error "testsuite/requirements.txt not found"
        exit 1
    fi

    cd testsuite

    # Install requirements if needed (you might want to do this manually)
    print_step "Installing testsuite requirements..."
    if pip3 install -r requirements.txt > /dev/null 2>&1; then
        print_success "Testsuite requirements installed"
    else
        print_warning "Failed to install testsuite requirements - continuing anyway"
    fi

    # Run the testsuite
    if [ "$VERBOSE" = true ]; then
        if python3 main.py --debug; then
            print_success "Integration testsuite passed"
        else
            print_error "Integration testsuite failed"
            cd ..
            exit 1
        fi
    else
        if python3 main.py; then
            print_success "Integration testsuite passed"
        else
            print_error "Integration testsuite failed"
            cd ..
            exit 1
        fi
    fi

    cd ..
fi

# Step 6: Build the app
print_step "Building spf binary..."
if CGO_ENABLED=0 go build -o ./bin/spf; then
    print_success "Build completed successfully"
else
    print_error "Build failed"
    exit 1
fi

echo ""
echo "${GREEN}🎉 All steps completed successfully!${NC}"
echo "${BLUE}Binary location:${NC} ./bin/spf"

# Show binary info
if [ -f "./bin/spf" ]; then
    BINARY_SIZE=$(du -h ./bin/spf | cut -f1)
    echo "${BLUE}Binary size:${NC} $BINARY_SIZE"
fi
