import { describe, expect, test } from 'vitest';
import { localizedUrl } from '../../utils/localizedUrl';

describe('with `build.output: "directory"`', () => {
	test('it has no effect in a monolingual project with a non-root single locale', () => {
		const url = new URL('https://example.com/fr/guide/');
		expect(localizedUrl(url, 'fr').href).toBe(url.href);
	});

	test('has no effect on index route in a monolingual project with a non-root single locale', () => {
		const url = new URL('https://example.com/fr/');
		expect(localizedUrl(url, 'fr').href).toBe(url.href);
	});
});

describe('with `build.output: "file"`', () => {
	test('it has no effect in a monolingual project with a non-root single locale', () => {
		const url = new URL('https://example.com/fr/guide.html');
		expect(localizedUrl(url, 'fr').href).toBe(url.href);
	});

	test('has no effect on index route in a monolingual project with a non-root single locale', () => {
		const url = new URL('https://example.com/fr.html');
		expect(localizedUrl(url, 'fr').href).toBe(url.href);
	});
});
