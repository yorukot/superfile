import { describe, expect, test } from 'vitest';
import { FaviconSchema } from '../../schemas/favicon';
import { TitleTransformConfigSchema } from '../../schemas/site-title';

describe('FaviconSchema', () => {
	test('returns the proper href and type attributes', () => {
		const icon = '/custom-icon.jpg';

		const favicon = FaviconSchema().parse(icon);

		expect(favicon.href).toBe(icon);
		expect(favicon.type).toBe('image/jpeg');
	});

	test('throws on invalid favicon extensions', () => {
		expect(() => FaviconSchema().parse('/favicon.pdf')).toThrow();
	});
});

describe('TitleTransformConfigSchema', () => {
	test('title can be a string', () => {
		const title = 'My Site';
		const defaultLang = 'en';

		const siteTitle = TitleTransformConfigSchema(defaultLang).parse(title);

		expect(siteTitle).toEqual({
			en: title,
		});
	});

	test('title can be an object', () => {
		const title = {
			en: 'My Site',
			es: 'Mi Sitio',
		};
		const defaultLang = 'en';

		const siteTitle = TitleTransformConfigSchema(defaultLang).parse(title);

		expect(siteTitle).toEqual(title);
	});

	test('throws on missing default language key', () => {
		const title = {
			es: 'Mi Sitio',
		};
		const defaultLang = 'en';

		expect(() => TitleTransformConfigSchema(defaultLang).parse(title)).toThrow();
	});
});
