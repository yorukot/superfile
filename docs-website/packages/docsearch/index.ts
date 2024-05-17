import type { StarlightPlugin } from '@astrojs/starlight/types';
import type { AstroUserConfig, ViteUserConfig } from 'astro';
import { z } from 'astro/zod';

/** Config options users must provide for DocSearch to work. */
const DocSearchConfigSchema = z.object({
	appId: z.string(),
	apiKey: z.string(),
	indexName: z.string(),
});
export type DocSearchConfig = z.input<typeof DocSearchConfigSchema>;

/** Starlight DocSearch plugin. */
export default function starlightDocSearch(userConfig: DocSearchConfig): StarlightPlugin {
	const opts = DocSearchConfigSchema.parse(userConfig);
	return {
		name: 'starlight-docsearch',
		hooks: {
			setup({ addIntegration, config, logger, updateConfig }) {
				// If the user has already has a custom override for the Search component, don't override it.
				if (config.components?.Search) {
					logger.warn(
						'It looks like you already have a `Search` component override in your Starlight configuration.'
					);
					logger.warn(
						'To render `@astrojs/starlight-docsearch`, remove the override for the `Search` component.\n'
					);
				} else {
					// Otherwise, add the Search component override to the user's configuration.
					updateConfig({
						pagefind: false,
						components: {
							...config.components,
							Search: '@astrojs/starlight-docsearch/DocSearch.astro',
						},
					});
				}

				// Add an Astro integration that injects a Vite plugin to expose
				// the DocSearch config via a virtual module.
				addIntegration({
					name: 'starlight-docsearch',
					hooks: {
						'astro:config:setup': ({ updateConfig }) => {
							updateConfig({
								vite: {
									plugins: [vitePluginDocSearch(opts)],
								},
							} satisfies AstroUserConfig);
						},
					},
				});
			},
		},
	};
}

/** Vite plugin that exposes the DocSearch config via virtual modules. */
function vitePluginDocSearch(config: DocSearchConfig): VitePlugin {
	const moduleId = 'virtual:starlight/docsearch-config';
	const resolvedModuleId = `\0${moduleId}`;
	const moduleContent = `export default ${JSON.stringify(config)}`;

	return {
		name: 'vite-plugin-starlight-docsearch-config',
		load(id) {
			return id === resolvedModuleId ? moduleContent : undefined;
		},
		resolveId(id) {
			return id === moduleId ? resolvedModuleId : undefined;
		},
	};
}

type VitePlugin = NonNullable<ViteUserConfig['plugins']>[number];
