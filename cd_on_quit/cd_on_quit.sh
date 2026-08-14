spf() {
    # Ask the binary where it writes the lastdir file, so this works for every
    # installation method (Homebrew, package managers, manual install, ...)
    export SPF_LAST_DIR="$(command spf path-list --lastdir-file)"

    command spf "$@"

    [ ! -f "$SPF_LAST_DIR" ] || {
        . "$SPF_LAST_DIR"
        rm -f -- "$SPF_LAST_DIR" > /dev/null
    }
}
