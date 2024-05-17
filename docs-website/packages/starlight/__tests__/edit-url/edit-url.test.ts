import { expect, test, vi } from 'vitest';
import { generateRouteData } from '../../utils/route-data';
import { routes } from '../../utils/routing';

vi.mock('astro:content', async () =>
	(await import('../test-utils')).mockedAstroContent({
		docs: [
			['index.mdx', { title: 'Home Page' }],
			['getting-started.mdx', { title: 'Getting Started' }],
			[
				// @ts-expect-error — Using a slug not present in Starlight docs site
				'showcase.mdx',
				{ title: 'Custom edit link', editUrl: 'https://example.com/custom-edit?link' },
			],
		],
	})
);

test('synthesizes edit URL using file location and `editLink.baseUrl`', () => {
	{
		const route = routes[0]!;
		const data = generateRouteData({
			props: { ...route, headings: [] },
			url: new URL('https://example.com'),
		});
		expect(data.editUrl?.href).toBe(
			'https://github.com/withastro/starlight/edit/main/docs/src/content/docs/index.mdx'
		);
	}
	{
		const route = routes[1]!;
		const data = generateRouteData({
			props: { ...route, headings: [] },
			url: new URL('https://example.com'),
		});
		expect(data.editUrl?.href).toBe(
			'https://github.com/withastro/starlight/edit/main/docs/src/content/docs/getting-started.mdx'
		);
	}
});

test('uses frontmatter `editUrl` if defined', () => {
	const route = routes[2]!;
	const data = generateRouteData({
		props: { ...route, headings: [] },
		url: new URL('https://example.com'),
	});
	expect(data.editUrl?.href).toBe('https://example.com/custom-edit?link');
});
