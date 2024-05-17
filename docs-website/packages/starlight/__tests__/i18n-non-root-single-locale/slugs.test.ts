import { describe, expect, test } from 'vitest';
import { localeToLang, localizedId, localizedSlug, slugToLocaleData } from '../../utils/slugs';

describe('slugToLocaleData', () => {
	test('returns default "fr" locale', () => {
		expect(slugToLocaleData('fr/test').locale).toBe('fr');
		expect(slugToLocaleData('fr/dir/test').locale).toBe('fr');
	});
	test('returns default locale "fr" lang', () => {
		expect(slugToLocaleData('fr/test').lang).toBe('fr-CA');
		expect(slugToLocaleData('fr/dir/test').lang).toBe('fr-CA');
	});
	test('returns default locale "ltr" dir', () => {
		expect(slugToLocaleData('fr/test').dir).toBe('ltr');
		expect(slugToLocaleData('fr/dir/test').dir).toBe('ltr');
	});
});

describe('localeToLang', () => {
	test('returns lang for default locale', () => {
		expect(localeToLang('fr')).toBe('fr-CA');
	});
});

describe('localizedId', () => {
	test('returns unchanged for default locale', () => {
		expect(localizedId('fr/test.md', 'fr')).toBe('fr/test.md');
	});
});

describe('localizedSlug', () => {
	test('returns unchanged for default locale', () => {
		expect(localizedSlug('fr', 'fr')).toBe('fr');
		expect(localizedSlug('fr/test', 'fr')).toBe('fr/test');
		expect(localizedSlug('fr/dir/test', 'fr')).toBe('fr/dir/test');
	});
});
