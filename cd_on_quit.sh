spf() {
    # Linux
    if [ "$(uname)" == "Linux" ]; then
        export SPF_LAST_DIR="${XDG_STATE_HOME:-$HOME/.local/state}/superfile/lastdir"
    fi
    
    # macOS
    if [ "$(uname)" == "Darwin" ]; then
        export SPF_LAST_DIR="$HOME/Library/Application Support/superfile/lastdir"
    fi
    
    # Windows
    if [ "$(uname)" == "Windows" ]; then
        export SPF_LAST_DIR="$LOCALAPPDATA/superfile/lastdir"
    fi
    
    command spf "$@"
    
    [ ! -f "$SPF_LAST_DIR" ] || {
        . "$SPF_LAST_DIR"
        rm -f -- "$SPF_LAST_DIR" > /dev/null
    }
}