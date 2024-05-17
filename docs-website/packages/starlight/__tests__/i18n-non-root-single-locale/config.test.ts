import config from 'virtual:starlight/user-config';
import { expect, test } from 'vitest';

test('test suite is using correct env', () => {
	expect(config.title).toMatchObject({ 'fr-CA': 'i18n with a non-root single locale' });
});

test('config.isMultilingual is false with a single locale', () => {
	expect(config.isMultilingual).toBe(false);
	expect(config.locales).keys('fr');
});

test('config.defaultLocale is populated from default locale', () => {
	expect(config.defaultLocale.lang).toBe('fr-CA');
	expect(config.defaultLocale.dir).toBe('ltr');
	expect(config.defaultLocale.locale).toBe('fr');
});
