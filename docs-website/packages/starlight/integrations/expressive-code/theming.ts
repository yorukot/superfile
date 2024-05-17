import { ExpressiveCodeTheme, type ThemeObjectOrShikiThemeName } from 'astro-expressive-code';
import nightOwlDark from './themes/night-owl-dark.jsonc?raw';
import nightOwlLight from './themes/night-owl-light.jsonc?raw';

export type BundledThemeName = 'starlight-dark' | 'starlight-light';

export type ThemeObjectOrBundledThemeName = ThemeObjectOrShikiThemeName | BundledThemeName;

/**
 * Converts the Starlight `themes` config option into a format understood by Expressive Code,
 * loading any bundled themes and using the Starlight defaults if no themes were provided.
 */
export function preprocessThemes(
	themes: ThemeObjectOrBundledThemeName[] | undefined
): ThemeObjectOrShikiThemeName[] {
	// Try to gracefully handle cases where the user forgot to use an array in the config
	themes = themes && !Array.isArray(themes) ? [themes] : themes;
	// If no themes were provided, use our bundled default themes
	if (!themes || !themes.length) themes = ['starlight-dark', 'starlight-light'];

	return themes.map((theme) => {
		// If the current entry is the name of a bundled theme, load it
		if (theme === 'starlight-dark' || theme === 'starlight-light') {
			const bundledTheme = theme === 'starlight-dark' ? nightOwlDark : nightOwlLight;
			return customizeBundledTheme(ExpressiveCodeTheme.fromJSONString(bundledTheme));
		}
		// Otherwise, just pass it through
		return theme;
	});
}

/**
 * Customizes some settings of the bundled theme to make it fit better with Starlight.
 */
function customizeBundledTheme(theme: ExpressiveCodeTheme) {
	theme.colors['titleBar.border'] = theme.colors['tab.activeBackground'];
	theme.colors['editorGroupHeader.tabsBorder'] = theme.colors['tab.activeBackground'];

	// Add underline font style to link syntax highlighting tokens
	// to match the new GitHub theme link style
	theme.settings.forEach((s) => {
		if (s.name?.includes('Link')) s.settings.fontStyle = 'underline';
	});

	return theme;
}

/**
 * Modifies the given theme by applying Starlight's CSS variables to the colors of UI elements
 * (backgrounds, buttons, shadows etc.). This ensures that code blocks match the site's theme.
 */
export function applyStarlightUiThemeColors(theme: ExpressiveCodeTheme) {
	const isDark = theme.type === 'dark';
	const neutralMinimal = isDark ? '#ffffff17' : '#0000001a';
	const neutralDimmed = isDark ? '#ffffff40' : '#00000055';

	// Make borders slightly transparent
	const borderColor = 'color-mix(in srgb, var(--sl-color-gray-5), transparent 25%)';
	theme.colors['titleBar.border'] = borderColor;
	theme.colors['editorGroupHeader.tabsBorder'] = borderColor;

	// Use the same color for terminal title bar background and editor tab bar background
	const backgroundColor = isDark ? 'var(--sl-color-black)' : 'var(--sl-color-gray-6)';
	theme.colors['titleBar.activeBackground'] = backgroundColor;
	theme.colors['editorGroupHeader.tabsBackground'] = backgroundColor;

	// Use the same color for terminal titles and tab titles
	theme.colors['titleBar.activeForeground'] = 'var(--sl-color-text)';
	theme.colors['tab.activeForeground'] = 'var(--sl-color-text)';

	// Set tab border colors
	const activeBorderColor = isDark ? 'var(--sl-color-accent-high)' : 'var(--sl-color-accent)';
	theme.colors['tab.activeBorder'] = 'transparent';
	theme.colors['tab.activeBorderTop'] = activeBorderColor;

	// Use neutral colors for scrollbars
	theme.colors['scrollbarSlider.background'] = neutralMinimal;
	theme.colors['scrollbarSlider.hoverBackground'] = neutralDimmed;

	// Set theme `bg` color property for contrast calculations
	theme.bg = isDark ? '#23262f' : '#f6f7f9';
	// Set actual background color to the appropriate Starlight CSS variable
	const editorBackgroundColor = isDark ? 'var(--sl-color-gray-6)' : 'var(--sl-color-gray-7)';

	theme.styleOverrides.frames = {
		// Use the same color for editor background, terminal background and active tab background
		editorBackground: editorBackgroundColor,
		terminalBackground: editorBackgroundColor,
		editorActiveTabBackground: editorBackgroundColor,
		terminalTitlebarDotsForeground: borderColor,
		terminalTitlebarDotsOpacity: '0.75',
		inlineButtonForeground: 'var(--sl-color-text)',
		frameBoxShadowCssValue: 'none',
	};

	// Use neutral, semi-transparent colors for default text markers
	// to avoid conflicts with the user's chosen background color
	theme.styleOverrides.textMarkers = {
		markBackground: neutralMinimal,
		markBorderColor: neutralDimmed,
	};

	return theme;
}
