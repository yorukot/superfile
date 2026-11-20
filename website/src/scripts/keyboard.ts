/** True when a key press belongs to a text field or an open dialog, not to the page. */
export function isTyping(event: KeyboardEvent): boolean {
	if (event.metaKey || event.altKey) return true;
	return event.target instanceof Element && !!event.target.closest("input, textarea, select, [contenteditable='true'], dialog[open]");
}
