import { describe, expect, test } from 'vitest';
import { localizedUrl } from '../../utils/localizedUrl';

describe('with `build.output: "directory"`', () => {
	test('it has no effect if locale matches', () => {
		const url = new URL('https://example.com/en/guide/');
		expect(localizedUrl(url, 'en').href).toBe(url.href);
	});

	test('it has no effect if locale matches for index', () => {
		const url = new URL('https://example.com/en/');
		expect(localizedUrl(url, 'en').href).toBe(url.href);
	});

	test('it changes locale to requested locale', () => {
		const url = new URL('https://example.com/en/guide/');
		expect(localizedUrl(url, 'fr').href).toBe('https://example.com/fr/guide/');
	});

	test('it changes locale to requested locale for index', () => {
		const url = new URL('https://example.com/en/');
		expect(localizedUrl(url, 'fr').href).toBe('https://example.com/fr/');
	});
});

describe('with `build.output: "file"`', () => {
	test('it has no effect if locale matches', () => {
		const url = new URL('https://example.com/en/guide.html');
		expect(localizedUrl(url, 'en').href).toBe(url.href);
	});

	test('it has no effect if locale matches for index', () => {
		const url = new URL('https://example.com/en.html');
		expect(localizedUrl(url, 'en').href).toBe(url.href);
	});

	test('it changes locale to requested locale', () => {
		const url = new URL('https://example.com/en/guide.html');
		expect(localizedUrl(url, 'fr').href).toBe('https://example.com/fr/guide.html');
	});

	test('it changes locale to requested locale for index', () => {
		const url = new URL('https://example.com/en.html');
		expect(localizedUrl(url, 'fr').href).toBe('https://example.com/fr.html');
	});
});
