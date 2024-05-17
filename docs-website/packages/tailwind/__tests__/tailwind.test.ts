import tailwindcss, { type Config } from 'tailwindcss';
import colors from 'tailwindcss/colors';
import postcss from 'postcss';
import { test, expect, describe, vi } from 'vitest';
import StarlightTailwindPlugin from '..';

/** Generate a CSS string based on the passed CSS and HTML content. */
const generatePluginCss = async ({
	css = '@tailwind base;',
	html = '',
	config = {},
}: { css?: string; html?: string; config?: Partial<Config> } = {}): Promise<string> => {
	const result = await postcss(
		tailwindcss({
			// Enable Starlight plugin.
			plugins: [StarlightTailwindPlugin()],
			// Provide content for Tailwind to scan for class names.
			content: [{ raw: html, extension: 'html' }],
			// Spread in any custom Tailwind config.
			...config,
		})
	).process(css, { from: '' });
	return result.css;
};

describe('@tailwind base;', async () => {
	// Generate base CSS with no core Tailwind plugins running to see just Starlight’s output.
	const base = await generatePluginCss({ config: { corePlugins: [] } });

	test('generates Starlight base CSS', async () => {
		expect(base).toMatchInlineSnapshot(`
			"*, ::before, ::after {
			    border-width: 0;
			    border-style: solid;
			    border-color: #e5e7eb;
			}
			::before, ::after {
			    --tw-content: ;
			}
			:root {
			    --sl-font: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
			    --sl-font-mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
			    --sl-color-white: #fff;
			    --sl-color-gray-1: #e5e7eb;
			    --sl-color-gray-2: #d1d5db;
			    --sl-color-gray-3: #9ca3af;
			    --sl-color-gray-4: #4b5563;
			    --sl-color-gray-5: #374151;
			    --sl-color-gray-6: #1f2937;
			    --sl-color-black: #111827;
			    --sl-color-accent-low: #1e1b4b;
			    --sl-color-accent: #4f46e5;
			    --sl-color-accent-high: #c7d2fe;
			}
			:root[data-theme="light"] {
			    --sl-color-white: #111827;
			    --sl-color-gray-1: #1f2937;
			    --sl-color-gray-2: #374151;
			    --sl-color-gray-3: #6b7280;
			    --sl-color-gray-4: #9ca3af;
			    --sl-color-gray-5: #d1d5db;
			    --sl-color-gray-6: #e5e7eb;
			    --sl-color-gray-7: #f3f4f6;
			    --sl-color-black: #fff;
			    --sl-color-accent-low: #c7d2fe;
			    --sl-color-accent: #4f46e5;
			    --sl-color-accent-high: #312e81;
			}"
		`);
	});

	test('configures `--sl-color-*` variables', () => {
		expect(base).includes('--sl-color-gray-1: #e5e7eb;');
		expect(base).includes('--sl-color-accent: #4f46e5;');
	});

	describe('with user theme config', async () => {
		const baseWithConfig = await generatePluginCss({
			config: {
				corePlugins: [],
				theme: { extend: { colors: { accent: colors.amber, gray: colors.slate } } },
			},
		});

		test('generates different CSS from base without user config', () => {
			expect(baseWithConfig).not.toEqual(base);
		});

		test('uses theme values for Starlight colours', () => {
			expect(baseWithConfig).includes('--sl-color-gray-1: #e2e8f0;');
			expect(baseWithConfig).includes('--sl-color-accent: #d97706;');
		});
	});

	test('disables Tailwind preflight', async () => {
		const baseWithDefaultPlugins = await generatePluginCss();
		expect(baseWithDefaultPlugins).not.includes('line-height: 1.5;');
		expect(baseWithDefaultPlugins).includes('--tw-');
		expect(baseWithDefaultPlugins).toMatchInlineSnapshot(`
			"*, ::before, ::after {
			    border-width: 0;
			    border-style: solid;
			    border-color: #e5e7eb;
			}
			::before, ::after {
			    --tw-content: ;
			}
			:root {
			    --sl-font: ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
			    --sl-font-mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
			    --sl-color-white: #fff;
			    --sl-color-gray-1: #e5e7eb;
			    --sl-color-gray-2: #d1d5db;
			    --sl-color-gray-3: #9ca3af;
			    --sl-color-gray-4: #4b5563;
			    --sl-color-gray-5: #374151;
			    --sl-color-gray-6: #1f2937;
			    --sl-color-black: #111827;
			    --sl-color-accent-low: #1e1b4b;
			    --sl-color-accent: #4f46e5;
			    --sl-color-accent-high: #c7d2fe;
			}
			:root[data-theme="light"] {
			    --sl-color-white: #111827;
			    --sl-color-gray-1: #1f2937;
			    --sl-color-gray-2: #374151;
			    --sl-color-gray-3: #6b7280;
			    --sl-color-gray-4: #9ca3af;
			    --sl-color-gray-5: #d1d5db;
			    --sl-color-gray-6: #e5e7eb;
			    --sl-color-gray-7: #f3f4f6;
			    --sl-color-black: #fff;
			    --sl-color-accent-low: #c7d2fe;
			    --sl-color-accent: #4f46e5;
			    --sl-color-accent-high: #312e81;
			}
			*, ::before, ::after {
			    --tw-border-spacing-x: 0;
			    --tw-border-spacing-y: 0;
			    --tw-translate-x: 0;
			    --tw-translate-y: 0;
			    --tw-rotate: 0;
			    --tw-skew-x: 0;
			    --tw-skew-y: 0;
			    --tw-scale-x: 1;
			    --tw-scale-y: 1;
			    --tw-pan-x:  ;
			    --tw-pan-y:  ;
			    --tw-pinch-zoom:  ;
			    --tw-scroll-snap-strictness: proximity;
			    --tw-gradient-from-position:  ;
			    --tw-gradient-via-position:  ;
			    --tw-gradient-to-position:  ;
			    --tw-ordinal:  ;
			    --tw-slashed-zero:  ;
			    --tw-numeric-figure:  ;
			    --tw-numeric-spacing:  ;
			    --tw-numeric-fraction:  ;
			    --tw-ring-inset:  ;
			    --tw-ring-offset-width: 0px;
			    --tw-ring-offset-color: #fff;
			    --tw-ring-color: rgb(59 130 246 / 0.5);
			    --tw-ring-offset-shadow: 0 0 #0000;
			    --tw-ring-shadow: 0 0 #0000;
			    --tw-shadow: 0 0 #0000;
			    --tw-shadow-colored: 0 0 #0000;
			    --tw-blur:  ;
			    --tw-brightness:  ;
			    --tw-contrast:  ;
			    --tw-grayscale:  ;
			    --tw-hue-rotate:  ;
			    --tw-invert:  ;
			    --tw-saturate:  ;
			    --tw-sepia:  ;
			    --tw-drop-shadow:  ;
			    --tw-backdrop-blur:  ;
			    --tw-backdrop-brightness:  ;
			    --tw-backdrop-contrast:  ;
			    --tw-backdrop-grayscale:  ;
			    --tw-backdrop-hue-rotate:  ;
			    --tw-backdrop-invert:  ;
			    --tw-backdrop-opacity:  ;
			    --tw-backdrop-saturate:  ;
			    --tw-backdrop-sepia:  ;
			}
			::backdrop {
			    --tw-border-spacing-x: 0;
			    --tw-border-spacing-y: 0;
			    --tw-translate-x: 0;
			    --tw-translate-y: 0;
			    --tw-rotate: 0;
			    --tw-skew-x: 0;
			    --tw-skew-y: 0;
			    --tw-scale-x: 1;
			    --tw-scale-y: 1;
			    --tw-pan-x:  ;
			    --tw-pan-y:  ;
			    --tw-pinch-zoom:  ;
			    --tw-scroll-snap-strictness: proximity;
			    --tw-gradient-from-position:  ;
			    --tw-gradient-via-position:  ;
			    --tw-gradient-to-position:  ;
			    --tw-ordinal:  ;
			    --tw-slashed-zero:  ;
			    --tw-numeric-figure:  ;
			    --tw-numeric-spacing:  ;
			    --tw-numeric-fraction:  ;
			    --tw-ring-inset:  ;
			    --tw-ring-offset-width: 0px;
			    --tw-ring-offset-color: #fff;
			    --tw-ring-color: rgb(59 130 246 / 0.5);
			    --tw-ring-offset-shadow: 0 0 #0000;
			    --tw-ring-shadow: 0 0 #0000;
			    --tw-shadow: 0 0 #0000;
			    --tw-shadow-colored: 0 0 #0000;
			    --tw-blur:  ;
			    --tw-brightness:  ;
			    --tw-contrast:  ;
			    --tw-grayscale:  ;
			    --tw-hue-rotate:  ;
			    --tw-invert:  ;
			    --tw-saturate:  ;
			    --tw-sepia:  ;
			    --tw-drop-shadow:  ;
			    --tw-backdrop-blur:  ;
			    --tw-backdrop-brightness:  ;
			    --tw-backdrop-contrast:  ;
			    --tw-backdrop-grayscale:  ;
			    --tw-backdrop-hue-rotate:  ;
			    --tw-backdrop-invert:  ;
			    --tw-backdrop-opacity:  ;
			    --tw-backdrop-saturate:  ;
			    --tw-backdrop-sepia:  ;
			}"
		`);
	});
});

describe('@tailwind utilities;', () => {
	test('uses [data-theme="dark"] for dark: utility classes', async () => {
		const utils = await generatePluginCss({
			css: '@tailwind utilities;',
			html: '<div class="dark:text-red-50"></div>',
		});
		expect(utils).includes('[data-theme="dark"] .dark');
		expect(utils).toMatchInlineSnapshot(`
			":is([data-theme="dark"] .dark\\:text-red-50) {
			    --tw-text-opacity: 1;
			    color: rgb(254 242 242 / var(--tw-text-opacity))
			}"
		`);
	});
});

test('warns when a prefix of "sl-" is set', async () => {
	const warn = vi.spyOn(console, 'warn').mockImplementation(() => {});
	await generatePluginCss({ config: { prefix: 'sl-' } });
	expect(warn).toBeCalledTimes(1);
	expect(warn.mock.lastCall?.[0]).toMatchInlineSnapshot(`
		"A Tailwind prefix of "sl-" will clash with Starlight’s built-in styles.
		Please set a different prefix in your Tailwind config file."
	`);
	warn.mockRestore();
});
