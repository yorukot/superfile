/// <reference types="vitest" />

import type { AstroConfig } from 'astro';
import { getViteConfig } from 'astro/config';
import { vitePluginStarlightUserConfig } from '../integrations/virtual-user-config';
import { runPlugins, type StarlightUserConfigWithPlugins } from '../utils/plugins';
import { createTestPluginContext } from './test-plugin-utils';

export async function defineVitestConfig(
	{ plugins, ...config }: StarlightUserConfigWithPlugins,
	opts?: {
		build?: Pick<AstroConfig['build'], 'format'>;
		trailingSlash?: AstroConfig['trailingSlash'];
	}
) {
	const root = new URL('./', import.meta.url);
	const srcDir = new URL('./src/', root);
	const build = opts?.build ?? { format: 'directory' };
	const trailingSlash = opts?.trailingSlash ?? 'ignore';

	const { starlightConfig } = await runPlugins(config, plugins, createTestPluginContext());
	return getViteConfig({
		plugins: [
			vitePluginStarlightUserConfig(starlightConfig, { root, srcDir, build, trailingSlash }),
		],
		test: {
			snapshotSerializers: ['./snapshot-serializer-astro-error.ts'],
		},
	});
}
