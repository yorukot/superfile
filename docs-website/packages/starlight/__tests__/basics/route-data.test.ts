import { expect, test, vi } from 'vitest';
import { generateRouteData } from '../../utils/route-data';
import { routes } from '../../utils/routing';

vi.mock('astro:content', async () =>
	(await import('../test-utils')).mockedAstroContent({
		docs: [
			['index.mdx', { title: 'Home Page' }],
			['getting-started.mdx', { title: 'Splash', template: 'splash' }],
			// @ts-expect-error — Using a slug not present in Starlight docs site
			['showcase.mdx', { title: 'ToC Disabled', tableOfContents: false }],
			['environmental-impact.md', { title: 'Explicit update date', lastUpdated: new Date() }],
		],
	})
);

test('adds data to route shape', () => {
	const route = routes[0]!;
	const data = generateRouteData({
		props: { ...route, headings: [{ depth: 1, slug: 'heading-1', text: 'Heading 1' }] },
		url: new URL('https://example.com'),
	});
	expect(data.hasSidebar).toBe(true);
	expect(data).toHaveProperty('lastUpdated');
	expect(data.toc).toMatchInlineSnapshot(`
		{
		  "items": [
		    {
		      "children": [],
		      "depth": 2,
		      "slug": "_top",
		      "text": "Overview",
		    },
		  ],
		  "maxHeadingLevel": 3,
		  "minHeadingLevel": 2,
		}
	`);
	expect(data.pagination).toMatchInlineSnapshot(`
		{
		  "next": {
		    "attrs": {},
		    "badge": undefined,
		    "href": "/environmental-impact/",
		    "isCurrent": false,
		    "label": "Explicit update date",
		    "type": "link",
		  },
		  "prev": undefined,
		}
	`);
	expect(data.sidebar.map((entry) => entry.label)).toMatchInlineSnapshot(`
		[
		  "Home Page",
		  "Explicit update date",
		  "Splash",
		  "ToC Disabled",
		]
	`);
});

test('disables table of contents for splash template', () => {
	const route = routes[1]!;
	const data = generateRouteData({
		props: { ...route, headings: [{ depth: 1, slug: 'heading-1', text: 'Heading 1' }] },
		url: new URL('https://example.com/getting-started/'),
	});
	expect(data.toc).toBeUndefined();
});

test('disables table of contents if frontmatter includes `tableOfContents: false`', () => {
	const route = routes[2]!;
	const data = generateRouteData({
		props: { ...route, headings: [{ depth: 1, slug: 'heading-1', text: 'Heading 1' }] },
		url: new URL('https://example.com/showcase/'),
	});
	expect(data.toc).toBeUndefined();
});

test('uses explicit last updated date from frontmatter', () => {
	const route = routes[3]!;
	const data = generateRouteData({
		props: { ...route, headings: [{ depth: 1, slug: 'heading-1', text: 'Heading 1' }] },
		url: new URL('https://example.com/showcase/'),
	});
	expect(data.lastUpdated).toBeInstanceOf(Date);
	expect(data.lastUpdated).toEqual(route.entry.data.lastUpdated);
});

test('includes localized labels', () => {
	const route = routes[0]!;
	const data = generateRouteData({
		props: { ...route, headings: [{ depth: 1, slug: 'heading-1', text: 'Heading 1' }] },
		url: new URL('https://example.com'),
	});
	expect(data.labels).toBeDefined();
	expect(data.labels['skipLink.label']).toBe('Skip to content');
});
