#!/usr/bin/env bash

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if colors should be disabled
if [ "$FORCE_COLOR" != "1" ] && ([ -n "$MAKEFLAGS" ] || [ "$TERM" = "dumb" ] || [ ! -t 1 ]); then
    # Disable colors when running under Make or non-interactive
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    NC=''
fi

# Default values
RUN_TESTSUITE=false
SKIP_TESTS=false
VERBOSE=false
USE_GLOBAL_ENV=false

# Function to print colored output
print_step() {
    printf "${BLUE}==>${NC} %s\n" "$1"
}

print_success() {
    printf "${GREEN}✓${NC} %s\n" "$1"
}

print_warning() {
    printf "${YELLOW}⚠${NC} %s\n" "$1"
}

print_error() {
    printf "${RED}✗${NC} %s\n" "$1"
}

# Function to setup Python virtual environment
setup_venv() {
    local venv_path="$1"
    
    # Remove existing incomplete virtual environment if activate script is missing or pip is broken
    if [ -d "$venv_path" ]; then
        if [ ! -f "$venv_path/bin/activate" ]; then
            print_warning "Removing incomplete virtual environment (missing activate)..."
            rm -rf "$venv_path"
        else
            # Test if pip works in the existing virtual environment
            if ! (source "$venv_path/bin/activate" && python -m pip --version > /dev/null 2>&1); then
                print_warning "Removing broken virtual environment (pip not working)..."
                rm -rf "$venv_path"
            fi
        fi
    fi
    
    if [ ! -d "$venv_path" ]; then
        print_step "Creating Python virtual environment..."
        if python3 -m venv "$venv_path" --upgrade-deps; then
            print_success "Virtual environment created at $venv_path"
        else
            print_error "Failed to create virtual environment"
            return 1
        fi
    else
        print_step "Using existing virtual environment at $venv_path"
    fi
    
    # Check if activate script exists and has proper permissions
    if [ ! -f "$venv_path/bin/activate" ]; then
        print_error "Virtual environment activate script not found at $venv_path/bin/activate"
        return 1
    fi
    
    # Ensure activate script has execution permissions
    chmod +x "$venv_path/bin/activate"
    
    # Activate virtual environment
    source "$venv_path/bin/activate"
    
    # Verify that we're in the virtual environment
    if [ -z "$VIRTUAL_ENV" ]; then
        print_error "Failed to activate virtual environment"
        return 1
    fi
    
    # Upgrade pip to latest version
    print_step "Upgrading pip in virtual environment..."
    if python -m pip install --upgrade pip > /dev/null 2>&1; then
        print_success "Pip upgraded successfully"
    else
        print_warning "Failed to upgrade pip - continuing anyway"
    fi
    
    return 0
}

# Function to cleanup virtual environment
cleanup_venv() {
    if [ -n "$VIRTUAL_ENV" ]; then
        deactivate 2>/dev/null || true
    fi
}

# Setup trap for cleanup on exit/interruption
trap cleanup_venv EXIT INT TERM

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
    echo "  --use-global-env    Use global Python environment instead of virtual environment"
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
        --use-global-env)
            USE_GLOBAL_ENV=true
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

printf "${BLUE}🚀 Starting superfile development workflow${NC}\n"
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

    # Use virtual environment by default, global environment if requested
    if [ "$USE_GLOBAL_ENV" = true ]; then
        # Install requirements globally
        print_step "Installing testsuite requirements globally..."
        print_warning "Using global Python environment - consider removing --use-global-env flag to use virtual environment"
        if python3 -m pip install -r requirements.txt > /dev/null 2>&1; then
            print_success "Testsuite requirements installed globally"
        else
            print_warning "Failed to install testsuite requirements - continuing anyway"
        fi
    else
        # Setup virtual environment (default behavior)
        VENV_PATH="./venv"
        
        if ! setup_venv "$VENV_PATH"; then
            print_error "Failed to setup virtual environment"
            cd ..
            exit 1
        fi
        
        # Install requirements in virtual environment
        print_step "Installing testsuite requirements in virtual environment..."
        if python -m pip install -r requirements.txt > /dev/null 2>&1; then
            print_success "Testsuite requirements installed in virtual environment"
        else
            print_error "Failed to install testsuite requirements in virtual environment"
            cd ..
            exit 1
        fi
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
printf "${GREEN}🎉 All steps completed successfully!${NC}\n"
printf "${BLUE}Binary location:${NC} ./bin/spf\n"

# Show binary info
if [ -f "./bin/spf" ]; then
    BINARY_SIZE=$(du -h ./bin/spf | cut -f1)
    printf "${BLUE}Binary size:${NC} $BINARY_SIZE\n"
fi
