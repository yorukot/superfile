import { describe, expect, test } from 'vitest';
import { getSitemapConfig, starlightSitemap } from '../../integrations/sitemap';
import type { StarlightConfig } from '../../types';
import { StarlightConfigSchema, type StarlightUserConfig } from '../../utils/user-config';

describe('starlightSitemap', () => {
	test('returns @astrojs/sitemap integration', () => {
		const integration = starlightSitemap({} as StarlightConfig);
		expect(integration.name).toBe('@astrojs/sitemap');
	});
});

describe('getSitemapConfig', () => {
	test('configures i18n config', () => {
		const config = getSitemapConfig(
			StarlightConfigSchema.parse({
				title: 'i18n test',
				locales: { root: { lang: 'en', label: 'English' }, fr: { label: 'French' } },
			} satisfies StarlightUserConfig)
		);
		expect(config).toMatchInlineSnapshot(`
			{
			  "i18n": {
			    "defaultLocale": "root",
			    "locales": {
			      "fr": "fr",
			      "root": "en",
			    },
			  },
			}
		`);
	});

	test('no config for monolingual sites', () => {
		const config = getSitemapConfig(
			StarlightConfigSchema.parse({ title: 'i18n test' } satisfies StarlightUserConfig)
		);
		expect(config).toMatchInlineSnapshot('{}');
	});
});
