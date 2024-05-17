import { describe, expect, test } from 'vitest';
import { createHead } from '../../utils/head';

describe('createHead', () => {
	test('merges two <title> tags', () => {
		expect(
			createHead(
				[{ tag: 'title', content: 'Default' }],
				[{ tag: 'title', content: 'Override', attrs: {} }]
			)
		).toEqual([{ tag: 'title', content: 'Override', attrs: {} }]);
	});

	for (const prop of ['name', 'property', 'http-equiv']) {
		test(`merges two <meta> tags with same ${prop} value`, () => {
			expect(
				createHead(
					[{ tag: 'meta', attrs: { [prop]: 'x', content: 'Default' } }],
					[{ tag: 'meta', attrs: { [prop]: 'x', content: 'Test' }, content: '' }]
				)
			).toEqual([{ tag: 'meta', content: '', attrs: { [prop]: 'x', content: 'Test' } }]);
		});
	}

	for (const prop of ['name', 'property', 'http-equiv']) {
		test(`does not merge <meta> tags with different ${prop} values`, () => {
			expect(
				createHead(
					[{ tag: 'meta', attrs: { [prop]: 'x', content: 'X' } }],
					[{ tag: 'meta', attrs: { [prop]: 'y', content: 'Y' }, content: '' }]
				)
			).toEqual([
				{ tag: 'meta', content: '', attrs: { [prop]: 'x', content: 'X' } },
				{ tag: 'meta', content: '', attrs: { [prop]: 'y', content: 'Y' } },
			]);
		});
	}

	test('sorts head by tag importance', () => {
		expect(
			createHead([
				// SEO meta tags
				{ tag: 'meta', attrs: { name: 'x' } },
				// Others
				{ tag: 'link', attrs: { rel: 'stylesheet' } },
				// Important meta tags
				{ tag: 'meta', attrs: { charset: 'utf-8' } },
				{ tag: 'meta', attrs: { name: 'viewport' } },
				{ tag: 'meta', attrs: { 'http-equiv': 'x' } },
				// <title>
				{ tag: 'title', content: 'Title' },
			])
		).toEqual([
			// Important meta tags
			{ tag: 'meta', attrs: { charset: 'utf-8' }, content: '' },
			{ tag: 'meta', attrs: { name: 'viewport' }, content: '' },
			{ tag: 'meta', attrs: { 'http-equiv': 'x' }, content: '' },
			// <title>
			{ tag: 'title', attrs: {}, content: 'Title' },
			// Others
			{ tag: 'link', attrs: { rel: 'stylesheet' }, content: '' },
			// SEO meta tags
			{ tag: 'meta', attrs: { name: 'x' }, content: '' },
		]);
	});

	test('places the default favicon below any user provided icons', () => {
		const defaultFavicon = {
			tag: 'link',
			attrs: {
				rel: 'shortcut icon',
				href: '/favicon.svg',
				type: 'image/svg+xml',
			},
		} as const;

		const userFavicon = {
			tag: 'link',
			attrs: {
				rel: 'icon',
				href: '/favicon.ico',
				sizes: '32x32',
			},
			content: '',
		} as const;

		expect(createHead([defaultFavicon], [userFavicon])).toMatchObject([
			userFavicon,
			defaultFavicon,
		]);
	});
});
