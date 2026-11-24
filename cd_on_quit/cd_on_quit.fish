function spf
    # Each session gets its own file, so concurrent sessions can't cd into
    # each other's directory
    set -l spf_cwd_file (mktemp)

    command spf --cwd-file $spf_cwd_file $argv

    if test -s "$spf_cwd_file"
        cd -- (cat -- "$spf_cwd_file")
    end
    rm -f -- "$spf_cwd_file" >> /dev/null
end
