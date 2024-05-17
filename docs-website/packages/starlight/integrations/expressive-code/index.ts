import {
	astroExpressiveCode,
	type AstroExpressiveCodeOptions,
	type CustomConfigPreprocessors,
} from 'astro-expressive-code';
import { addClassName } from 'astro-expressive-code/hast';
import type { AstroIntegration } from 'astro';
import type { StarlightConfig } from '../../types';
import type { createTranslationSystemFromFs } from '../../utils/translations-fs';
import { pathToLocale } from '../shared/pathToLocale';
import {
	applyStarlightUiThemeColors,
	preprocessThemes,
	type ThemeObjectOrBundledThemeName,
} from './theming';
import { addTranslations } from './translations';

export type StarlightExpressiveCodeOptions = Omit<AstroExpressiveCodeOptions, 'themes'> & {
	/**
	 * Set the themes used to style code blocks.
	 *
	 * See the [Expressive Code `themes` documentation](https://github.com/expressive-code/expressive-code/blob/main/packages/astro-expressive-code/README.md#themes)
	 * for details of the supported theme formats.
	 *
	 * Starlight uses the dark and light variants of Sarah Drasner’s
	 * [Night Owl theme](https://github.com/sdras/night-owl-vscode-theme) by default.
	 *
	 * If you provide at least one dark and one light theme, Starlight will automatically keep
	 * the active code block theme in sync with the current site theme. Configure this behavior
	 * with the [`useStarlightDarkModeSwitch`](#usestarlightdarkmodeswitch) option.
	 *
	 * Defaults to `['starlight-dark', 'starlight-light']`.
	 */
	themes?: ThemeObjectOrBundledThemeName[] | undefined;
	/**
	 * When `true`, code blocks automatically switch between light and dark themes when the
	 * site theme changes.
	 *
	 * When `false`, you must manually add CSS to handle switching between multiple themes.
	 *
	 * **Note**: When setting `themes`, you must provide at least one dark and one light theme
	 * for the Starlight dark mode switch to work.
	 *
	 * Defaults to `true`.
	 */
	useStarlightDarkModeSwitch?: boolean | undefined;
	/**
	 * When `true`, Starlight's CSS variables are used for the colors of code block UI elements
	 * (backgrounds, buttons, shadows etc.), matching the
	 * [site color theme](/guides/css-and-tailwind/#theming).
	 *
	 * When `false`, the colors provided by the active syntax highlighting theme are used for
	 * these elements.
	 *
	 * Defaults to `true` if the `themes` option is not set (= you are using Starlight's
	 * default themes), and `false` otherwise.
	 *
	 * **Note**: When manually setting this to `true` with your custom set of `themes`, you must
	 * provide at least one dark and one light theme to ensure proper color contrast.
	 */
	useStarlightUiThemeColors?: boolean | undefined;
};

type StarlightEcIntegrationOptions = {
	starlightConfig: StarlightConfig;
	useTranslations?: ReturnType<typeof createTranslationSystemFromFs> | undefined;
};

/**
 * Create an Expressive Code configuration preprocessor based on Starlight config.
 * Used internally to set up Expressive Code and by the `<Code>` component.
 */
export function getStarlightEcConfigPreprocessor({
	starlightConfig,
	useTranslations,
}: StarlightEcIntegrationOptions): CustomConfigPreprocessors['preprocessAstroIntegrationConfig'] {
	return (input): AstroExpressiveCodeOptions => {
		const astroConfig = input.astroConfig;
		const ecConfig = input.ecConfig as StarlightExpressiveCodeOptions;

		const {
			themes: themesInput,
			customizeTheme,
			styleOverrides: { textMarkers: textMarkersStyleOverrides, ...otherStyleOverrides } = {},
			useStarlightDarkModeSwitch,
			useStarlightUiThemeColors = ecConfig.themes === undefined,
			plugins = [],
			...rest
		} = ecConfig;

		// Handle the `themes` option
		const themes = preprocessThemes(themesInput);
		if (useStarlightUiThemeColors === true && themes.length < 2) {
			console.warn(
				`*** Warning: Using the config option "useStarlightUiThemeColors: true" ` +
					`with a single theme is not recommended. For better color contrast, ` +
					`please provide at least one dark and one light theme.\n`
			);
		}

		// Add the `not-content` class to all rendered blocks to prevent them from being affected
		// by Starlight's default content styles
		plugins.push({
			name: 'Starlight Plugin',
			hooks: {
				postprocessRenderedBlock: ({ renderData }) => {
					addClassName(renderData.blockAst, 'not-content');
				},
			},
		});

		// Add Expressive Code UI translations (if any) for all defined locales
		if (useTranslations) addTranslations(starlightConfig, useTranslations);

		return {
			themes,
			customizeTheme: (theme) => {
				if (useStarlightUiThemeColors) {
					applyStarlightUiThemeColors(theme);
				}
				if (customizeTheme) {
					theme = customizeTheme(theme) ?? theme;
				}
				return theme;
			},
			defaultLocale: starlightConfig.defaultLocale?.lang ?? starlightConfig.defaultLocale?.locale,
			themeCssSelector: (theme, { styleVariants }) => {
				// If one dark and one light theme are available, and the user has not disabled it,
				// generate theme CSS selectors compatible with Starlight's dark mode switch
				if (useStarlightDarkModeSwitch !== false && styleVariants.length >= 2) {
					const baseTheme = styleVariants[0]?.theme;
					const altTheme = styleVariants.find((v) => v.theme.type !== baseTheme?.type)?.theme;
					if (theme === baseTheme || theme === altTheme) return `[data-theme='${theme.type}']`;
				}
				// Return the default selector
				return `[data-theme='${theme.name}']`;
			},
			styleOverrides: {
				borderRadius: '0px',
				borderWidth: '1px',
				codePaddingBlock: '0.75rem',
				codePaddingInline: '1rem',
				codeFontFamily: 'var(--__sl-font-mono)',
				codeFontSize: 'var(--sl-text-code)',
				codeLineHeight: 'var(--sl-line-height)',
				uiFontFamily: 'var(--__sl-font)',
				textMarkers: {
					lineDiffIndicatorMarginLeft: '0.25rem',
					defaultChroma: '45',
					backgroundOpacity: '60%',
					...textMarkersStyleOverrides,
				},
				...otherStyleOverrides,
			},
			getBlockLocale: ({ file }) => pathToLocale(file.path, { starlightConfig, astroConfig }),
			plugins,
			...rest,
		};
	};
}

export const starlightExpressiveCode = ({
	starlightConfig,
	useTranslations,
}: StarlightEcIntegrationOptions): AstroIntegration[] => {
	// If Expressive Code is disabled, add a shim to prevent build errors and provide
	// a helpful error message in case the user tries to use the `<Code>` component
	if (starlightConfig.expressiveCode === false) {
		const modules: Record<string, string> = {
			'virtual:astro-expressive-code/api': 'export default {}',
			'virtual:astro-expressive-code/config': 'export default {}',
			'virtual:astro-expressive-code/preprocess-config': `throw new Error("Starlight's
				Code component requires Expressive Code, which is disabled in your Starlight config.
				Please remove \`expressiveCode: false\` from your config or import Astro's built-in
				Code component from 'astro:components' instead.")`.replace(/\s+/g, ' '),
		};

		return [
			{
				name: 'astro-expressive-code-shim',
				hooks: {
					'astro:config:setup': ({ updateConfig }) => {
						updateConfig({
							vite: {
								plugins: [
									{
										name: 'vite-plugin-astro-expressive-code-shim',
										enforce: 'post',
										resolveId: (id) => (id in modules ? `\0${id}` : undefined),
										load: (id) => (id?.[0] === '\0' ? modules[id.slice(1)] : undefined),
									},
								],
							},
						});
					},
				},
			},
		];
	}

	const configArgs =
		typeof starlightConfig.expressiveCode === 'object'
			? (starlightConfig.expressiveCode as AstroExpressiveCodeOptions)
			: {};
	return [
		astroExpressiveCode({
			...configArgs,
			customConfigPreprocessors: {
				preprocessAstroIntegrationConfig: getStarlightEcConfigPreprocessor({
					starlightConfig,
					useTranslations,
				}),
				preprocessComponentConfig: `
					import starlightConfig from 'virtual:starlight/user-config'
					import { useTranslations, getStarlightEcConfigPreprocessor } from '@astrojs/starlight/internal'

					export default getStarlightEcConfigPreprocessor({ starlightConfig, useTranslations })
				`,
			},
		}),
	];
};
