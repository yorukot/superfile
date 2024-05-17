import { type HeadConfig, HeadConfigSchema, type HeadUserConfig } from '../schemas/head';

const HeadSchema = HeadConfigSchema();

/** Create a fully parsed, merged, and sorted head entry array from multiple sources. */
export function createHead(defaults: HeadUserConfig, ...heads: HeadConfig[]) {
	let head = HeadSchema.parse(defaults);
	for (const next of heads) {
		head = mergeHead(head, next);
	}
	return sortHead(head);
}

/**
 * Test if a head config object contains a matching `<title>` or `<meta>` tag.
 *
 * For example, will return true if `head` already contains
 * `<meta name="description" content="A">` and the passed `tag`
 * is `<meta name="description" content="B">`. Tests against `name`,
 * `property`, and `http-equiv` attributes for `<meta>` tags.
 */
function hasTag(head: HeadConfig, entry: HeadConfig[number]): boolean {
	switch (entry.tag) {
		case 'title':
			return head.some(({ tag }) => tag === 'title');
		case 'meta':
			return hasOneOf(head, entry, ['name', 'property', 'http-equiv']);
		default:
			return false;
	}
}

/**
 * Test if a head config object contains a tag of the same type
 * as `entry` and a matching attribute for one of the passed `keys`.
 */
function hasOneOf(head: HeadConfig, entry: HeadConfig[number], keys: string[]): boolean {
	const attr = getAttr(keys, entry);
	if (!attr) return false;
	const [key, val] = attr;
	return head.some(({ tag, attrs }) => tag === entry.tag && attrs[key] === val);
}

/** Find the first matching key–value pair in a head entry’s attributes. */
function getAttr(
	keys: string[],
	entry: HeadConfig[number]
): [key: string, value: string | boolean] | undefined {
	let attr: [string, string | boolean] | undefined;
	for (const key of keys) {
		const val = entry.attrs[key];
		if (val) {
			attr = [key, val];
			break;
		}
	}
	return attr;
}

/** Merge two heads, overwriting entries in the first head that exist in the second. */
function mergeHead(oldHead: HeadConfig, newHead: HeadConfig) {
	return [...oldHead.filter((tag) => !hasTag(newHead, tag)), ...newHead];
}

/** Sort head tags to place important tags first and relegate “SEO” meta tags. */
function sortHead(head: HeadConfig) {
	return head.sort((a, b) => {
		const aImportance = getImportance(a);
		const bImportance = getImportance(b);
		return aImportance > bImportance ? -1 : bImportance > aImportance ? 1 : 0;
	});
}

/** Get the relative importance of a specific head tag. */
function getImportance(entry: HeadConfig[number]) {
	// 1. Important meta tags.
	if (
		entry.tag === 'meta' &&
		('charset' in entry.attrs || 'http-equiv' in entry.attrs || entry.attrs.name === 'viewport')
	) {
		return 100;
	}
	// 2. Page title
	if (entry.tag === 'title') return 90;
	// 3. Anything that isn’t an SEO meta tag.
	if (entry.tag !== 'meta') {
		// The default favicon should be below any extra icons that the user may have set
		// because if several icons are equally appropriate, the last one is used and we
		// want to use the SVG icon when supported.
		if (entry.tag === 'link' && 'rel' in entry.attrs && entry.attrs.rel === 'shortcut icon') {
			return 70;
		}
		return 80;
	}
	// 4. SEO meta tags.
	return 0;
}
