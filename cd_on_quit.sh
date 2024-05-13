spf() {
	export SPF_LAST_DIR="${XDG_STATE_HOME:-$HOME/.local/state}/superfile/lastdir"

	command spf "$@"
	
	[ ! -f "$SPF_LAST_DIR" ] || {
		. "$SPF_LAST_DIR"
		rm -f -- "$SPF_LAST_DIR" > /dev/null
	}
}