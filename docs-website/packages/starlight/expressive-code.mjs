/**
 * @file This file is exported by Starlight as `@astrojs/starlight/expressive-code`.
 *
 * It is required by the `<Code>` component to access the same configuration preprocessor
 * function as the one used by the integration.
 *
 * It also provides access to all of the Expressive Code classes and functions without having
 * to install `astro-expressive-code` as an additional dependency into a user's project
 * (and thereby risiking version conflicts).
 *
 * Note: This file is intentionally not a TypeScript module to allow access to all exported
 * functionality even if TypeScript is not available, e.g. from the `ec.config.mjs` file
 * that does not get processed by Vite.
 */

export * from 'astro-expressive-code';

// @ts-ignore - Types are provided by the separate `expressive-code.d.ts` file
export function defineEcConfig(config) {
	return config;
}
