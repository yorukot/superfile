// Theme colours are read at build time from the theme files superfile ships.
import fs from "node:fs";
import path from "node:path";

export interface Theme {
	name: string;
	/** CSS custom properties that repaint the home page's preview in this theme. */
	style: string;
}

const themeDir = path.join(process.env.GITHUB_WORKSPACE ?? "..", "src/superfile_config/theme");
const featured = ["catppuccin-mocha", "catppuccin-latte", "dracula", "gruvbox", "nord", "rose-pine", "tokyonight"];

function readTheme(name: string): Theme {
	const source = fs.readFileSync(path.join(themeDir, `${name}.toml`), "utf-8");
	const colors: Record<string, string> = {};
	for (const [, key, value] of source.matchAll(/^(\w+)\s*=\s*"(#[0-9a-fA-F]{3,8})"/gm)) colors[key] = value;
	const pick = (...keys: string[]) => keys.map(key => colors[key]).find(Boolean);
	const vars = {
		"--term-bg": pick("file_panel_bg", "full_screen_bg"),
		"--term-ink": pick("file_panel_fg", "full_screen_fg"),
		"--term-muted": pick("file_panel_fg", "full_screen_fg"),
		"--term-line": pick("file_panel_border"),
		"--term-teal": pick("file_panel_top_path", "sidebar_title"),
		"--term-signal": pick("file_panel_border_active"),
		"--signal": pick("file_panel_border_active"),
		"--term-amber": pick("sidebar_title", "hint"),
		"--theme-cursor": pick("file_panel_item_selected_fg", "cursor")
	};
	const style = Object.entries(vars)
		.filter(([, value]) => value)
		.map(([key, value]) => `${key}:${value}`)
		.join(";");
	return { name, style };
}

/** How many themes superfile ships. */
export const themeCount = fs.readdirSync(themeDir).filter(file => file.endsWith(".toml")).length;

/** The themes the home page can preview. */
export const featuredThemes = featured.map(readTheme);
