spf() {
    # Ask the binary where it writes the lastdir file, so this works for every
    # installation method (Homebrew, package managers, manual install, ...)
    # assigned inside `if`: `export VAR="$(...)"` always reports success, so a
    # failed lookup would go unnoticed (SC2155)
    if SPF_LAST_DIR="$(command spf path-list --lastdir-file)"; then
        export SPF_LAST_DIR
    else
        SPF_LAST_DIR=''
    fi

    command spf "$@"

    [ ! -f "$SPF_LAST_DIR" ] || {
        . "$SPF_LAST_DIR"
        rm -f -- "$SPF_LAST_DIR" > /dev/null
    }
}
