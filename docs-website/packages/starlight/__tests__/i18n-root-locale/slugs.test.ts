import { describe, expect, test } from 'vitest';
import { localeToLang, localizedId, localizedSlug, slugToLocaleData } from '../../utils/slugs';

describe('slugToLocaleData', () => {
	test('returns an undefined locale for root locale slugs', () => {
		expect(slugToLocaleData('test').locale).toBeUndefined();
		expect(slugToLocaleData('dir/test').locale).toBeUndefined();
	});
	test('returns a locale for localized slugs', () => {
		expect(slugToLocaleData('en/test').locale).toBe('en');
		expect(slugToLocaleData('ar/test').locale).toBe('ar');
	});
	test('returns default locale lang for root locale slugs', () => {
		expect(slugToLocaleData('test').lang).toBe('fr');
		expect(slugToLocaleData('dir/test').lang).toBe('fr');
	});
	test('returns langs for localized slugs', () => {
		expect(slugToLocaleData('ar/test').lang).toBe('ar');
		expect(slugToLocaleData('en/dir/test').lang).toBe('en-US');
	});
	test('returns default locale dir for root locale slugs', () => {
		expect(slugToLocaleData('test').dir).toBe('ltr');
		expect(slugToLocaleData('dir/test').dir).toBe('ltr');
	});
	test('returns configured dir for localized slugs', () => {
		expect(slugToLocaleData('ar/test').dir).toBe('rtl');
		expect(slugToLocaleData('en/dir/test').dir).toBe('ltr');
	});
});

describe('localeToLang', () => {
	test('returns lang for root locale', () => {
		expect(localeToLang(undefined)).toBe('fr');
	});
	test('returns lang for non-root locales', () => {
		expect(localeToLang('en')).toBe('en-US');
		expect(localeToLang('ar')).toBe('ar');
	});
});

describe('localizedId', () => {
	test('returns unchanged when already in requested locale', () => {
		expect(localizedId('test.md', undefined)).toBe('test.md');
		expect(localizedId('dir/test.md', undefined)).toBe('dir/test.md');
		expect(localizedId('en/test.md', 'en')).toBe('en/test.md');
		expect(localizedId('en/dir/test.md', 'en')).toBe('en/dir/test.md');
		expect(localizedId('ar/test.md', 'ar')).toBe('ar/test.md');
		expect(localizedId('ar/dir/test.md', 'ar')).toBe('ar/dir/test.md');
	});
	test('returns localized id for requested locale', () => {
		expect(localizedId('test.md', 'en')).toBe('en/test.md');
		expect(localizedId('dir/test.md', 'en')).toBe('en/dir/test.md');
		expect(localizedId('en/test.md', 'ar')).toBe('ar/test.md');
		expect(localizedId('en/test.md', undefined)).toBe('test.md');
	});
});

describe('localizedSlug', () => {
	test('returns unchanged when already in requested locale', () => {
		expect(localizedSlug('', undefined)).toBe('');
		expect(localizedSlug('test', undefined)).toBe('test');
		expect(localizedSlug('dir/test', undefined)).toBe('dir/test');
		expect(localizedSlug('en', 'en')).toBe('en');
		expect(localizedSlug('en/test', 'en')).toBe('en/test');
		expect(localizedSlug('en/dir/test', 'en')).toBe('en/dir/test');
	});
	test('returns localized slug for requested locale', () => {
		expect(localizedSlug('', 'en')).toBe('en');
		expect(localizedSlug('test', 'en')).toBe('en/test');
		expect(localizedSlug('dir/test', 'en')).toBe('en/dir/test');
		expect(localizedSlug('en', undefined)).toBe('');
		expect(localizedSlug('en/test', undefined)).toBe('test');
		expect(localizedSlug('en/dir/test', undefined)).toBe('dir/test');
		expect(localizedSlug('en', 'ar')).toBe('ar');
		expect(localizedSlug('en/test', 'ar')).toBe('ar/test');
		expect(localizedSlug('en/dir/test', 'ar')).toBe('ar/dir/test');
	});
});
