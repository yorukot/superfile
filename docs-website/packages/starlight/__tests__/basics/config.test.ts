import config from 'virtual:starlight/user-config';
import { expect, test } from 'vitest';

test('test suite is using correct env', () => {
	expect(config.title).toMatchObject({ en: 'Basics' });
});

test('isMultilingual is false when no locales configured ', () => {
	expect(config.locales).toBeUndefined();
	expect(config.isMultilingual).toBe(false);
});

test('default locale is set when no locales configured', () => {
	expect(config.defaultLocale).not.toBeUndefined();
	expect(config.defaultLocale.lang).toBe('en');
	expect(config.defaultLocale.label).toBe('English');
	expect(config.defaultLocale.dir).toBe('ltr');
});

test('lastUpdated defaults to false', () => {
	expect(config.lastUpdated).toBe(false);
});

test('favicon defaults to the provided SVG icon', () => {
	expect(config.favicon.href).toBe('/favicon.svg');
	expect(config.favicon.type).toBe('image/svg+xml');
});
