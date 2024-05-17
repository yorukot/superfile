import plugin from 'tailwindcss/plugin';

/**
 * Starlight Tailwind Plugin
 *
 * - Disables Tailwind Preflight.
 * - Configures `dark:` variants for Starlight dark mode.
 * - Links Starlight’s colors to `gray` and `accent` in Tailwind theme settings.
 * - Links Starlight’s fonts to `sans` and `mono` in Tailwind theme settings.
 *
 * @example
 * // tailwind.config.mjs
 * import colors from 'tailwindcss/colors';
 * import starlightPlugin from '@astrojs/starlight-tailwind';
 *
 * export default {
 *   plugins: [
 *     // Add Starlight’s Tailwind plugin
 *     starlightPlugin(),
 *   ],
 *   theme: {
 *     extend: {
 *       colors: {
 *         // Set an accent color for Astro to use. Indigo is closest to Astro’s defaults.
 *         accent: colors.yellow,
 *         // Configure your preferred gray scale. Zinc is closest to Astro’s defaults.
 *         gray: colors.zinc,
 *       },
 *     },
 *   },
 * }
 */
const StarlightTailwindPlugin = () =>
	plugin(
		({ addBase, theme, config }) => {
			if (config<string>('prefix') === 'sl-') {
				console.warn(
					'A Tailwind prefix of "sl-" will clash with Starlight’s built-in styles.\n' +
						'Please set a different prefix in your Tailwind config file.'
				);
			}

			/** Utility to apply accent colors based on a user’s theme config. */
			const themeAccent = (
				shade: 50 | 100 | 200 | 300 | 400 | 500 | 600 | 700 | 800 | 900 | 950,
				fallback: string
			) =>
				shade === 950
					? theme(`colors.accent.${shade}`, theme(`colors.accent.900`, fallback))
					: theme(`colors.accent.${shade}`, fallback);

			let white: string = theme('colors.white');
			if (typeof white !== 'string') {
				console.warn(
					`Expected \`colors.white\` in Tailwind theme to be a string, received ${typeof white}.\n` +
						`Try setting a single value, for example \`white: '#fafaf9'\` or \`white: colors.stone[50]\`.`
				);
				// Ensure a usable value for white if the user-configured one is wrong.
				white = '#fff';
			}

			addBase({
				// Restore crucial styles from Tailwind Preflight: https://tailwindcss.com/docs/preflight
				// Allow adding a border to an element by just adding a border-width. (https://github.com/tailwindcss/tailwindcss/pull/116)
				'*, ::before, ::after': {
					borderWidth: '0',
					borderStyle: 'solid',
					borderColor: theme('borderColor.DEFAULT', 'currentColor'),
				},
				'::before, ::after': { '--tw-content': '' },

				// Wire up Starlight theme to use Tailwind config.
				':root': {
					// Use Tailwind-configured font families.
					'--sl-font': theme('fontFamily.sans'),
					'--sl-font-mono': theme('fontFamily.mono'),
					// Dark mode Starlight theme variables.
					'--sl-color-white': white,
					'--sl-color-gray-1': theme('colors.gray.200'),
					'--sl-color-gray-2': theme('colors.gray.300'),
					'--sl-color-gray-3': theme('colors.gray.400'),
					'--sl-color-gray-4': theme('colors.gray.600'),
					'--sl-color-gray-5': theme('colors.gray.700'),
					'--sl-color-gray-6': theme('colors.gray.800'),
					'--sl-color-black': theme('colors.gray.900'),
					'--sl-color-accent-low': themeAccent(950, '#1e1b4b'),
					'--sl-color-accent': themeAccent(600, '#4f46e5'),
					'--sl-color-accent-high': themeAccent(200, '#c7d2fe'),
					// Light mode Starlight theme variables
					'&[data-theme="light"]': {
						'--sl-color-white': theme('colors.gray.900'),
						'--sl-color-gray-1': theme('colors.gray.800'),
						'--sl-color-gray-2': theme('colors.gray.700'),
						'--sl-color-gray-3': theme('colors.gray.500'),
						'--sl-color-gray-4': theme('colors.gray.400'),
						'--sl-color-gray-5': theme('colors.gray.300'),
						'--sl-color-gray-6': theme('colors.gray.200'),
						'--sl-color-gray-7': theme('colors.gray.100'),
						'--sl-color-black': white,
						'--sl-color-accent-low': themeAccent(200, '#c7d2fe'),
						'--sl-color-accent': themeAccent(600, '#4f46e5'),
						'--sl-color-accent-high': themeAccent(900, '#312e81'),
					},
				},
			});
		},
		{
			// Starlight uses a `data-theme` attribute to power its dark mode.
			darkMode: ['class', '[data-theme="dark"]'],
			corePlugins: {
				// Disable Tailwind’s default reset styles which conflict with Starlight.
				preflight: false,
			},
		}
	);

export default StarlightTailwindPlugin;
