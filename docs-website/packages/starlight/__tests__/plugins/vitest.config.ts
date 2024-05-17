import { defineVitestConfig } from '../test-config';

export default defineVitestConfig({
	title: 'Plugins',
	sidebar: [{ label: 'Getting Started', link: 'getting-started' }],
	plugins: [
		{
			name: 'test-plugin-1',
			hooks: {
				setup({ config, updateConfig }) {
					updateConfig({
						title: `${config.title} - Custom`,
						description: 'plugin 1',
						/**
						 * The configuration received by a plugin should be the user provided configuration as-is
						 * befor any Zod `transform`s are applied.
						 * To test this, we use this plugin to update the `favicon` value to a specific value if
						 * the `favicon` config value is an object, which would mean that the associated Zod
						 * `transform` was applied.
						 */
						favicon: typeof config.favicon === 'object' ? 'invalid.svg' : 'valid.svg',
					});
				},
			},
		},
		{
			name: 'test-plugin-2',
			hooks: {
				setup({ config, updateConfig }) {
					updateConfig({
						description: `${config.description} - plugin 2`,
						sidebar: [{ label: 'Showcase', link: 'showcase' }],
					});
				},
			},
		},
		{
			name: 'test-plugin-3',
			hooks: {
				async setup({ config, updateConfig }) {
					await Promise.resolve();
					updateConfig({
						description: `${config.description} - plugin 3`,
					});
				},
			},
		},
	],
});
