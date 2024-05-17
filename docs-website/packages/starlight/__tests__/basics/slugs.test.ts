import { describe, expect, test, vi } from 'vitest';
import {
	localeToLang,
	localizedId,
	localizedSlug,
	slugToLocaleData,
	slugToParam,
	slugToPathname,
	urlToSlug,
} from '../../utils/slugs';

describe('slugToLocaleData', () => {
	test('returns an undefined locale for root locale slugs', () => {
		expect(slugToLocaleData('test').locale).toBeUndefined();
		expect(slugToLocaleData('dir/test').locale).toBeUndefined();
	});
	test('returns default "en" lang when no locale config is set', () => {
		expect(slugToLocaleData('test').lang).toBe('en');
		expect(slugToLocaleData('dir/test').lang).toBe('en');
	});
	test('returns default "ltr" dir when no locale config is set', () => {
		expect(slugToLocaleData('test').dir).toBe('ltr');
		expect(slugToLocaleData('dir/test').dir).toBe('ltr');
	});
});

describe('slugToParam', () => {
	test('returns undefined for empty slug (index)', () => {
		expect(slugToParam('')).toBeUndefined();
	});
	test('returns undefined for root index', () => {
		expect(slugToParam('index')).toBeUndefined();
	});
	test('strips index from end of nested slug', () => {
		expect(slugToParam('dir/index')).toBe('dir');
		expect(slugToParam('dir/index/sub-dir/index')).toBe('dir/index/sub-dir');
	});
	test('returns other slugs unchanged', () => {
		expect(slugToParam('slug')).toBe('slug');
		expect(slugToParam('dir/page')).toBe('dir/page');
		expect(slugToParam('dir/sub-dir/page')).toBe('dir/sub-dir/page');
	});
});

describe('slugToPathname', () => {
	test('returns "/" for empty slug', () => {
		expect(slugToPathname('')).toBe('/');
	});
	test('returns "/" for root index', () => {
		expect(slugToPathname('index')).toBe('/');
	});
	test('strips index from end of nested slug', () => {
		expect(slugToPathname('dir/index')).toBe('/dir/');
		expect(slugToPathname('dir/index/sub-dir/index')).toBe('/dir/index/sub-dir/');
	});
	test('returns slugs with leading and trailing slashes added', () => {
		expect(slugToPathname('slug')).toBe('/slug/');
		expect(slugToPathname('dir/page')).toBe('/dir/page/');
		expect(slugToPathname('dir/sub-dir/page')).toBe('/dir/sub-dir/page/');
	});
});

describe('localeToLang', () => {
	test('returns lang for root locale', () => {
		expect(localeToLang(undefined)).toBe('en');
	});
});

describe('localizedId', () => {
	test('returns unchanged when no locales are set', () => {
		expect(localizedId('test.md', undefined)).toBe('test.md');
	});
});

describe('localizedSlug', () => {
	test('returns unchanged when no locales are set', () => {
		expect(localizedSlug('test', undefined)).toBe('test');
	});
});

describe('urlToSlug', () => {
	test('returns slugs with `build.output: "directory"`', () => {
		expect(urlToSlug(new URL('https://example.com'))).toBe('');
		expect(urlToSlug(new URL('https://example.com/slug'))).toBe('slug');
		expect(urlToSlug(new URL('https://example.com/dir/page/'))).toBe('dir/page');
		expect(urlToSlug(new URL('https://example.com/dir/sub-dir/page/'))).toBe('dir/sub-dir/page');
	});

	test('returns slugs with `build.output: "file"`', () => {
		expect(urlToSlug(new URL('https://example.com/index.html'))).toBe('');
		expect(urlToSlug(new URL('https://example.com/slug.html'))).toBe('slug');
		expect(urlToSlug(new URL('https://example.com/dir/page/index.html'))).toBe('dir/page');
		expect(urlToSlug(new URL('https://example.com/dir/sub-dir/page.html'))).toBe(
			'dir/sub-dir/page'
		);
	});

	test('returns slugs with a custom `base` option', () => {
		vi.stubEnv('BASE_URL', '/base/');
		expect(urlToSlug(new URL('https://example.com/base'))).toBe('');
		expect(urlToSlug(new URL('https://example.com/base/slug'))).toBe('slug');
		expect(urlToSlug(new URL('https://example.com/base/dir/page/'))).toBe('dir/page');
		expect(urlToSlug(new URL('https://example.com/base/dir/sub-dir/page/'))).toBe(
			'dir/sub-dir/page'
		);
		vi.unstubAllEnvs();
	});
});
