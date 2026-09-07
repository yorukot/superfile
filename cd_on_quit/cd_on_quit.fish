function spf
    # Ask the binary where it writes the lastdir file, so this works for every
    # installation method (Homebrew, package managers, manual install, ...)
    set spf_last_dir (command spf path-list --lastdir-file)

    command spf $argv

    if test -f "$spf_last_dir"
        source "$spf_last_dir"
        rm -f -- "$spf_last_dir" >> /dev/null
    end
end
