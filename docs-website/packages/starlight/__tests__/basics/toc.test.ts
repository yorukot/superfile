import { expect, test } from 'vitest';
import { generateToC } from '../../utils/generateToC';

const defaultOpts = { minHeadingLevel: 2, maxHeadingLevel: 3, title: 'Overview' };

test('generates an overview entry with no headings available', () => {
	const toc = generateToC([], defaultOpts);
	expect(toc).toHaveLength(1);
	expect(toc[0]).toEqual({ children: [], depth: 2, slug: '_top', text: 'Overview' });
});

test('generates entries from heading array', () => {
	const toc = generateToC([{ text: 'One', slug: 'one', depth: 2 }], defaultOpts);
	expect(toc).toHaveLength(2);
	expect(toc[0]).toEqual({ children: [], depth: 2, slug: '_top', text: 'Overview' });
	expect(toc[1]).toEqual({ children: [], depth: 2, slug: 'one', text: 'One' });
});

test('nests lower-level headings in children array h2 => h3', () => {
	const toc = generateToC(
		[
			{ text: 'One', slug: 'one', depth: 2 },
			{ text: 'Two', slug: 'two', depth: 3 },
		],
		defaultOpts
	);
	expect(toc).toHaveLength(2);
	expect(toc[0]).toEqual({ children: [], depth: 2, slug: '_top', text: 'Overview' });
	expect(toc[1]).toEqual({
		children: [{ children: [], depth: 3, slug: 'two', text: 'Two' }],
		depth: 2,
		slug: 'one',
		text: 'One',
	});
});

test('nests lower-level headings in children array h2 => h4', () => {
	const toc = generateToC(
		[
			{ text: 'One', slug: 'one', depth: 2 },
			{ text: 'Two', slug: 'two', depth: 4 },
		],
		{ ...defaultOpts, maxHeadingLevel: 6 }
	);
	expect(toc).toHaveLength(2);
	expect(toc[0]).toEqual({ children: [], depth: 2, slug: '_top', text: 'Overview' });
	expect(toc[1]).toEqual({
		children: [{ children: [], depth: 4, slug: 'two', text: 'Two' }],
		depth: 2,
		slug: 'one',
		text: 'One',
	});
});

test('nests lower-level headings deeply h2 => h4 => h6', () => {
	const toc = generateToC(
		[
			{ text: 'One', slug: 'one', depth: 2 },
			{ text: 'Two', slug: 'two', depth: 4 },
			{ text: 'Three', slug: 'three', depth: 6 },
		],
		{ ...defaultOpts, maxHeadingLevel: 6 }
	);
	expect(toc).toHaveLength(2);
	expect(toc[0]).toEqual({ children: [], depth: 2, slug: '_top', text: 'Overview' });
	expect(toc[1]).toMatchInlineSnapshot(`
		{
		  "children": [
		    {
		      "children": [
		        {
		          "children": [],
		          "depth": 6,
		          "slug": "three",
		          "text": "Three",
		        },
		      ],
		      "depth": 4,
		      "slug": "two",
		      "text": "Two",
		    },
		  ],
		  "depth": 2,
		  "slug": "one",
		  "text": "One",
		}
	`);
});

test('adds higher-level headings sequentially h6 => h4 => h2', () => {
	const toc = generateToC(
		[
			{ text: 'One', slug: 'one', depth: 6 },
			{ text: 'Two', slug: 'two', depth: 4 },
			{ text: 'Three', slug: 'three', depth: 2 },
		],
		{ ...defaultOpts, maxHeadingLevel: 6 }
	);
	expect(toc).toHaveLength(2);
	expect(toc).toMatchInlineSnapshot(`
		[
		  {
		    "children": [
		      {
		        "children": [],
		        "depth": 6,
		        "slug": "one",
		        "text": "One",
		      },
		      {
		        "children": [],
		        "depth": 4,
		        "slug": "two",
		        "text": "Two",
		      },
		    ],
		    "depth": 2,
		    "slug": "_top",
		    "text": "Overview",
		  },
		  {
		    "children": [],
		    "depth": 2,
		    "slug": "three",
		    "text": "Three",
		  },
		]
	`);
});

test('filters out higher-level headings than minHeadingLevel', () => {
	const toc = generateToC([{ text: 'One', slug: 'one', depth: 1 }], defaultOpts);
	expect(toc).toHaveLength(1);
	expect(toc[0]).toEqual({ children: [], depth: 2, slug: '_top', text: 'Overview' });
});

test('filters out lower-level headings than maxHeadingLevel', () => {
	const toc = generateToC([{ text: 'One', slug: 'one', depth: 4 }], defaultOpts);
	expect(toc).toHaveLength(1);
	expect(toc[0]).toEqual({ children: [], depth: 2, slug: '_top', text: 'Overview' });
});
