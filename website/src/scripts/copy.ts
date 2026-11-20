// Copy buttons: a `[data-copy]` button copies the command or the code block it sits in.
const CHECK_ICON = '<path d="M20 6 9 17l-5-5"/>';

document.addEventListener("click", async event => {
	const button = (event.target as Element).closest<HTMLElement>("[data-copy]");
	if (!button || button.classList.contains("is-copied")) return;
	const source = button.closest(".sf-cmd")?.querySelector<HTMLElement>("[data-cmd-text]") ?? button.closest(".sf-code")?.querySelector("pre");
	if (!source) return;
	const label = button.querySelector("[data-copy-label]")!;

	try {
		await navigator.clipboard.writeText(source.innerText.trimEnd());
	} catch {
		// The browser refused (no permission, or not https): select the text, so copying by hand is one key press.
		getSelection()?.selectAllChildren(source);
		label.textContent = "press ctrl+c";
		setTimeout(() => (label.textContent = "copy"), 1600);
		return;
	}

	// The copy icon becomes a check and the word becomes "copied" for a moment.
	const icon = button.querySelector("svg")!;
	const copyIcon = icon.innerHTML;
	button.classList.add("is-copied");
	icon.innerHTML = CHECK_ICON;
	label.textContent = "copied";
	setTimeout(() => {
		button.classList.remove("is-copied");
		icon.innerHTML = copyIcon;
		label.textContent = "copy";
	}, 1600);
});
