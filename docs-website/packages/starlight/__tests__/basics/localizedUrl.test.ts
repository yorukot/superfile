import { describe, expect, test } from 'vitest';
import { localizedUrl } from '../../utils/localizedUrl';

describe('with `build.output: "directory"`', () => {
	test('it has no effect in a monolingual project', () => {
		const url = new URL('https://example.com/en/guide/');
		expect(localizedUrl(url, undefined).href).toBe(url.href);
	});

	test('has no effect on index route in a monolingual project', () => {
		const url = new URL('https://example.com/');
		expect(localizedUrl(url, undefined).href).toBe(url.href);
	});
});

describe('with `build.output: "file"`', () => {
	test('it has no effect in a monolingual project', () => {
		const url = new URL('https://example.com/en/guide.html');
		expect(localizedUrl(url, undefined).href).toBe(url.href);
	});

	test('has no effect on index route in a monolingual project', () => {
		const url = new URL('https://example.com/index.html');
		expect(localizedUrl(url, undefined).href).toBe(url.href);
	});
});
