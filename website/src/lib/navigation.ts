export interface NavLeaf {
	label: string;
	slug: string;
}

export interface NavGroup {
	label: string;
	items: NavLeaf[];
}

/** The docs sidebar. Group labels are rendered as `// label`. */
export const navigation: NavGroup[] = [
	{
		label: "get started",
		items: [
			{ label: "overview", slug: "overview" },
			{ label: "install", slug: "getting-started/installation" },
			{ label: "tutorial", slug: "getting-started/tutorial" }
		]
	},
	{
		label: "guides",
		items: [
			{ label: "navigation and search", slug: "guides/navigation" },
			{ label: "file operations", slug: "guides/file-operations" },
			{ label: "file preview", slug: "getting-started/image-preview" },
			{ label: "metadata", slug: "guides/metadata" },
			{ label: "sidebar", slug: "guides/sidebar" },
			{ label: "prompt", slug: "guides/prompt" },
			{ label: "shell integration", slug: "guides/shell-integration" }
		]
	},
	{
		label: "customize",
		items: [
			{ label: "config file path", slug: "configure/config-file-path" },
			{ label: "config", slug: "configure/superfile-config" },
			{ label: "hotkeys", slug: "configure/custom-hotkeys" },
			{ label: "themes", slug: "configure/custom-theme" },
			{ label: "plugins", slug: "configure/enable-plugin" }
		]
	},
	{
		label: "reference",
		items: [
			{ label: "command line", slug: "reference/cli" },
			{ label: "dependencies", slug: "reference/dependencies" },
			{ label: "hotkey list", slug: "list/hotkey-list" },
			{ label: "theme list", slug: "list/theme-list" },
			{ label: "plugin list", slug: "list/plugin-list" },
			{ label: "troubleshooting", slug: "troubleshooting" },
			{ label: "changelog", slug: "changelog" }
		]
	},
	{
		label: "contribute",
		items: [
			{ label: "how to contribute", slug: "contribute/how-to-contribute" },
			{ label: "file structure", slug: "contribute/file-struct" },
			{ label: "implementation info", slug: "contribute/implementation-info" }
		]
	},
	{
		label: "project",
		items: [
			{ label: "special thanks", slug: "special-thanks" },
			{ label: "notice", slug: "notice" }
		]
	}
];

/**
 * The sidebar of one set of docs: only the pages that set has. Pages it has that today's sidebar
 * does not know (they were renamed or removed since) go into a last group, named after their file.
 */
export function sidebarFor(pages: string[]): NavGroup[] {
	const known = new Set(navigation.flatMap(group => group.items.map(item => item.slug)));
	const groups = navigation.map(group => ({ ...group, items: group.items.filter(item => pages.includes(item.slug)) })).filter(group => group.items.length);
	const others = pages.filter(page => !known.has(page)).sort();
	if (others.length) groups.push({ label: "more", items: others.map(page => ({ slug: page, label: page.split("/").pop()!.replaceAll("-", " ") })) });
	return groups;
}

/** The entries before and after a page, in the order of its sidebar. */
export function getPrevNext(sidebar: NavGroup[], slug: string): { prev?: NavLeaf; next?: NavLeaf } {
	const flat = sidebar.flatMap(group => group.items);
	const index = flat.findIndex(item => item.slug === slug);
	return { prev: flat[index - 1], next: flat[index + 1] };
}
