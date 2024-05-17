import config from 'virtual:starlight/user-config';
import { expect, test } from 'vitest';

test('test suite is using correct env', () => {
	expect(config.title).toMatchObject({ fr: 'i18n with root locale' });
});

test('config.isMultilingual is true with multiple locales', () => {
	expect(config.isMultilingual).toBe(true);
	expect(config.locales).keys('root', 'en', 'ar');
});

test('config.defaultLocale is populated from root locale', () => {
	expect(config.defaultLocale.lang).toBe('fr');
	expect(config.defaultLocale.dir).toBe('ltr');
	expect(config.defaultLocale.locale).toBeUndefined();
});
