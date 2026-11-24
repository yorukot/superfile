spf() {
    # Each session gets its own file, so concurrent sessions can't cd into
    # each other's directory
    local spf_cwd_file
    spf_cwd_file="$(mktemp "${TMPDIR:-/tmp}/spf-cwd.XXXXXX")" || {
        command spf "$@"
        return
    }

    command spf --cwd-file "$spf_cwd_file" "$@"

    if [ -s "$spf_cwd_file" ]; then
        cd -- "$(cat -- "$spf_cwd_file")"
    fi
    rm -f -- "$spf_cwd_file" > /dev/null
}
