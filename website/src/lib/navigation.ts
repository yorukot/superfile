export interface NavLeaf {
	label: string;
	slug: string;
	zhLabel?: string;
}

export interface NavGroup {
	label: string;
	zhLabel?: string;
	items: NavLeaf[];
}

export type NavNode = NavLeaf | NavGroup;

export function isGroup(node: NavNode): node is NavGroup {
	return "items" in node;
}

export const navigation: NavNode[] = [
	{ label: "Overview", slug: "overview", zhLabel: "概述" },
	{
		label: "Start Here",
		zhLabel: "開始使用",
		items: [
			{
				label: "Installation",
				slug: "getting-started/installation",
				zhLabel: "安裝"
			},
			{
				label: "Tutorial",
				slug: "getting-started/tutorial",
				zhLabel: "教學"
			},
			{
				label: "Image Preview",
				slug: "getting-started/image-preview",
				zhLabel: "圖片預覽"
			}
		]
	},
	{
		label: "Configure",
		zhLabel: "設定",
		items: [
			{
				label: "Config File Path",
				slug: "configure/config-file-path",
				zhLabel: "設定檔路徑"
			},
			{
				label: "superfile config",
				slug: "configure/superfile-config",
				zhLabel: "superfile 設定"
			},
			{
				label: "Custom Hotkeys",
				slug: "configure/custom-hotkeys",
				zhLabel: "自訂快捷鍵"
			},
			{
				label: "Custom Theme",
				slug: "configure/custom-theme",
				zhLabel: "自訂主題"
			},
			{
				label: "Enable Plugin",
				slug: "configure/enable-plugin",
				zhLabel: "啟用插件"
			}
		]
	},
	{
		label: "List",
		zhLabel: "清單",
		items: [
			{
				label: "Hotkey List",
				slug: "list/hotkey-list",
				zhLabel: "快捷鍵清單"
			},
			{
				label: "Theme List",
				slug: "list/theme-list",
				zhLabel: "主題清單"
			},
			{
				label: "Plugin List",
				slug: "list/plugin-list",
				zhLabel: "插件清單"
			}
		]
	},
	{
		label: "Contribute",
		zhLabel: "貢獻",
		items: [
			{
				label: "How to Contribute",
				slug: "contribute/how-to-contribute",
				zhLabel: "如何貢獻"
			},
			{
				label: "File Structure",
				slug: "contribute/file-struct",
				zhLabel: "檔案結構"
			},
			{
				label: "Implementation Info",
				slug: "contribute/implementation-info",
				zhLabel: "實作資訊"
			}
		]
	},
	{ label: "Troubleshooting", slug: "troubleshooting", zhLabel: "疑難排解" },
	{ label: "Special Thanks", slug: "special-thanks", zhLabel: "特別感謝" },
	{ label: "NOTICE", slug: "notice", zhLabel: "NOTICE" },
	{ label: "Changelog", slug: "changelog", zhLabel: "更新紀錄" }
];

/** Flatten nav tree to ordered list for prev/next. */
export function flattenNav(): NavLeaf[] {
	const result: NavLeaf[] = [];
	for (const node of navigation) {
		if (isGroup(node)) {
			result.push(...node.items);
		} else {
			result.push(node);
		}
	}
	return result;
}

/** Get prev/next entries for a given base slug (no locale prefix). */
export function getPrevNext(slug: string): {
	prev?: NavLeaf;
	next?: NavLeaf;
} {
	const flat = flattenNav();
	const index = flat.findIndex(e => e.slug === slug);
	if (index === -1) return {};
	return {
		prev: flat[index - 1],
		next: flat[index + 1]
	};
}
