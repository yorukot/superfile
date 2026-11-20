// Copy buttons: a `[data-copy]` button copies the command or the code block it sits in.
// Lucide "check" and "x".
const ICONS = { copied: '<path d="M20 6 9 17l-5-5"/>', failed: '<path d="M18 6 6 18"/><path d="m6 6 12 12"/>' };

/** Shows the outcome on the button for a moment: its icon, its colour and a word. */
function flash(button: HTMLElement, outcome: keyof typeof ICONS, word: string) {
	const icon = button.querySelector("svg")!;
	const label = button.querySelector("[data-copy-label]")!;
	const copyIcon = icon.innerHTML;
	button.classList.add(`is-${outcome}`);
	icon.innerHTML = ICONS[outcome];
	label.textContent = word;
	setTimeout(() => {
		button.classList.remove(`is-${outcome}`);
		icon.innerHTML = copyIcon;
		label.textContent = "copy";
	}, 1600);
}

document.addEventListener("click", async event => {
	const button = (event.target as Element).closest<HTMLElement>("[data-copy]");
	if (!button || button.matches(".is-copied, .is-failed")) return;
	const source = button.closest(".sf-cmd")?.querySelector<HTMLElement>("[data-cmd-text]") ?? button.closest(".sf-code")?.querySelector("pre");
	if (!source) return;
	try {
		await navigator.clipboard.writeText(source.innerText.trimEnd());
		flash(button, "copied", "copied");
	} catch {
		// The browser refused (no permission, or not https): select the text, so copying by hand is one step.
		getSelection()?.selectAllChildren(source);
		flash(button, "failed", "press ctrl+c");
	}
});
