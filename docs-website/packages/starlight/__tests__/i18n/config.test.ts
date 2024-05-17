import config from 'virtual:starlight/user-config';
import { expect, test } from 'vitest';

test('test suite is using correct env', () => {
	expect(config.title).toMatchObject({ 'en-US': 'i18n with no root locale' });
});

test('config.isMultilingual is true with multiple locales', () => {
	expect(config.isMultilingual).toBe(true);
	expect(config.locales).keys('fr', 'en', 'ar', 'pt-br');
});

test('config.defaultLocale is populated from the user’s chosen default', () => {
	expect(config.defaultLocale.locale).toBe('en');
	expect(config.defaultLocale.lang).toBe('en-US');
	expect(config.defaultLocale.dir).toBe('ltr');
});
