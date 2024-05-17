import { describe, expect, test, vi } from 'vitest';
import { flattenSidebar, getPrevNextLinks, getSidebar } from '../../utils/navigation';

vi.mock('astro:content', async () =>
	(await import('../test-utils')).mockedAstroContent({
		docs: [
			['index.mdx', { title: 'Home Page' }],
			['environmental-impact.md', { title: 'Eco-friendly docs' }],
			['guides/authoring-content.md', { title: 'Authoring Markdown' }],
			['reference/frontmatter.md', { title: 'Frontmatter Reference', sidebar: { hidden: true } }],
			['guides/components.mdx', { title: 'Components' }],
		],
	})
);

describe('getSidebar', () => {
	test('returns an array of sidebar entries', () => {
		expect(getSidebar('/', undefined)).toMatchInlineSnapshot(`
			[
			  {
			    "attrs": {},
			    "badge": undefined,
			    "href": "/",
			    "isCurrent": true,
			    "label": "Home Page",
			    "type": "link",
			  },
			  {
			    "attrs": {},
			    "badge": undefined,
			    "href": "/environmental-impact/",
			    "isCurrent": false,
			    "label": "Eco-friendly docs",
			    "type": "link",
			  },
			  {
			    "badge": undefined,
			    "collapsed": false,
			    "entries": [
			      {
			        "attrs": {},
			        "badge": undefined,
			        "href": "/guides/authoring-content/",
			        "isCurrent": false,
			        "label": "Authoring Markdown",
			        "type": "link",
			      },
			      {
			        "attrs": {},
			        "badge": undefined,
			        "href": "/guides/components/",
			        "isCurrent": false,
			        "label": "Components",
			        "type": "link",
			      },
			    ],
			    "label": "guides",
			    "type": "group",
			  },
			]
		`);
	});

	test.each(['/', '/environmental-impact/', '/guides/authoring-content/'])(
		'marks current path with isCurrent: %s',
		(currentPath) => {
			const items = flattenSidebar(getSidebar(currentPath, undefined));
			const currentItems = items.filter((item) => item.type === 'link' && item.isCurrent);
			expect(currentItems).toHaveLength(1);
			const currentItem = currentItems[0];
			if (currentItem?.type !== 'link') throw new Error('Expected current item to be link');
			expect(currentItem.href).toBe(currentPath);
		}
	);

	test('ignore trailing slashes when marking current path with isCurrent', () => {
		const pathWithTrailingSlash = '/environmental-impact/';
		const items = flattenSidebar(getSidebar(pathWithTrailingSlash, undefined));
		const currentItems = items.filter((item) => item.type === 'link' && item.isCurrent);
		expect(currentItems).toMatchInlineSnapshot(`
			[
			  {
			    "attrs": {},
			    "badge": undefined,
			    "href": "/environmental-impact/",
			    "isCurrent": true,
			    "label": "Eco-friendly docs",
			    "type": "link",
			  },
			]
		`);
	});

	test('nests files in subdirectory in group when autogenerating', () => {
		const sidebar = getSidebar('/', undefined);
		expect(sidebar.every((item) => item.type === 'group' || !item.href.startsWith('/guides/')));
		const guides = sidebar.find((item) => item.type === 'group' && item.label === 'guides');
		expect(guides?.type).toBe('group');
		// @ts-expect-error — TypeScript doesn’t know we know we’re in a group.
		expect(guides.entries).toHaveLength(2);
	});

	test('uses page title as label when autogenerating', () => {
		const sidebar = getSidebar('/', undefined);
		const homeLink = sidebar.find((item) => item.type === 'link' && item.href === '/');
		expect(homeLink?.label).toBe('Home Page');
	});
});

describe('flattenSidebar', () => {
	test('flattens nested sidebar array', () => {
		const sidebar = getSidebar('/', undefined);
		const flattened = flattenSidebar(sidebar);
		// Sidebar should include some nested group items.
		expect(sidebar.some((item) => item.type === 'group')).toBe(true);
		// Flattened sidebar should only include link items.
		expect(flattened.every((item) => item.type === 'link')).toBe(true);

		expect(flattened).toMatchInlineSnapshot(`
			[
			  {
			    "attrs": {},
			    "badge": undefined,
			    "href": "/",
			    "isCurrent": true,
			    "label": "Home Page",
			    "type": "link",
			  },
			  {
			    "attrs": {},
			    "badge": undefined,
			    "href": "/environmental-impact/",
			    "isCurrent": false,
			    "label": "Eco-friendly docs",
			    "type": "link",
			  },
			  {
			    "attrs": {},
			    "badge": undefined,
			    "href": "/guides/authoring-content/",
			    "isCurrent": false,
			    "label": "Authoring Markdown",
			    "type": "link",
			  },
			  {
			    "attrs": {},
			    "badge": undefined,
			    "href": "/guides/components/",
			    "isCurrent": false,
			    "label": "Components",
			    "type": "link",
			  },
			]
		`);
	});
});

describe('getPrevNextLinks', () => {
	test('returns stable previous/next values', () => {
		const sidebar = getSidebar('/environmental-impact/', undefined);
		const links = getPrevNextLinks(sidebar, true, {});
		expect(links).toMatchInlineSnapshot(`
			{
			  "next": {
			    "attrs": {},
			    "badge": undefined,
			    "href": "/guides/authoring-content/",
			    "isCurrent": false,
			    "label": "Authoring Markdown",
			    "type": "link",
			  },
			  "prev": {
			    "attrs": {},
			    "badge": undefined,
			    "href": "/",
			    "isCurrent": false,
			    "label": "Home Page",
			    "type": "link",
			  },
			}
		`);
	});

	test('returns no links when pagination is disabled', () => {
		const sidebar = getSidebar('/environmental-impact/', undefined);
		const links = getPrevNextLinks(sidebar, false, {});
		expect(links).toEqual({ prev: undefined, next: undefined });
	});

	test('returns no previous link for first item', () => {
		const sidebar = getSidebar('/', undefined);
		const links = getPrevNextLinks(sidebar, true, {});
		expect(links.prev).toBeUndefined();
	});

	test('returns no next link for last item', () => {
		const sidebar = getSidebar('/guides/components/', undefined);
		const links = getPrevNextLinks(sidebar, true, {});
		expect(links.next).toBeUndefined();
	});

	test('final parameter can disable prev/next', () => {
		const sidebar = getSidebar('/environmental-impact/', undefined);
		expect(getPrevNextLinks(sidebar, true, { prev: true }).prev).toBeDefined();
		expect(getPrevNextLinks(sidebar, true, { prev: false }).prev).toBeUndefined();
		expect(getPrevNextLinks(sidebar, true, { next: true }).next).toBeDefined();
		expect(getPrevNextLinks(sidebar, true, { next: false }).next).toBeUndefined();
	});

	test('final parameter can set custom link label with string', () => {
		const sidebar = getSidebar('/environmental-impact/', undefined);
		const withDefaultLabels = getPrevNextLinks(sidebar, true, {});
		const withCustomLabels = getPrevNextLinks(sidebar, true, { prev: 'x', next: 'y' });
		expect(withCustomLabels.prev?.label).toBe('x');
		expect(withCustomLabels.prev?.label).not.toBe(withDefaultLabels.prev?.label);
		expect(withCustomLabels.next?.label).toBe('y');
		expect(withCustomLabels.next?.label).not.toBe(withDefaultLabels.next?.label);
	});

	test('final parameter can set custom link label with object', () => {
		const sidebar = getSidebar('/environmental-impact/', undefined);
		const withDefaultLabels = getPrevNextLinks(sidebar, true, {});
		const withCustomLabels = getPrevNextLinks(sidebar, true, {
			prev: { label: 'x' },
			next: { label: 'y' },
		});
		expect(withCustomLabels.prev?.label).toBe('x');
		expect(withCustomLabels.prev?.label).not.toBe(withDefaultLabels.prev?.label);
		expect(withCustomLabels.next?.label).toBe('y');
		expect(withCustomLabels.next?.label).not.toBe(withDefaultLabels.next?.label);
	});

	test('final parameter can set custom link destination', () => {
		const sidebar = getSidebar('/environmental-impact/', undefined);
		const withDefaults = getPrevNextLinks(sidebar, true, {});
		const withCustomLinks = getPrevNextLinks(sidebar, true, {
			prev: { link: '/x' },
			next: { link: '/y' },
		});
		expect(withCustomLinks.prev?.href).toBe('/x');
		expect(withCustomLinks.prev?.href).not.toBe(withDefaults.prev?.href);
		expect(withCustomLinks.prev?.label).toBe(withDefaults.prev?.label);
		expect(withCustomLinks.next?.href).toBe('/y');
		expect(withCustomLinks.next?.href).not.toBe(withDefaults.next?.href);
		expect(withCustomLinks.next?.label).toBe(withDefaults.next?.label);
	});

	test('final parameter can set custom link even if no default link existed', () => {
		const sidebar = getSidebar('/', undefined);
		const withDefaults = getPrevNextLinks(sidebar, true, {});
		const withCustomLinks = getPrevNextLinks(sidebar, true, {
			prev: { link: 'x', label: 'X' },
		});
		expect(withDefaults.prev).toBeUndefined();
		expect(withCustomLinks.prev).toEqual({
			type: 'link',
			href: '/x',
			label: 'X',
			isCurrent: false,
			attrs: {},
		});
	});

	test('final parameter can override global pagination toggle', () => {
		const sidebar = getSidebar('/environmental-impact/', undefined);
		const withDefaults = getPrevNextLinks(sidebar, true, {});
		const withOverrides = getPrevNextLinks(sidebar, false, { prev: true, next: true });
		expect(withOverrides).toEqual(withDefaults);
	});
});
