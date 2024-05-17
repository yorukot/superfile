import config from 'virtual:starlight/user-config';
import type { Badge } from '../schemas/badge';
import type { PrevNextLinkConfig } from '../schemas/prevNextLink';
import type {
	AutoSidebarGroup,
	LinkHTMLAttributes,
	SidebarItem,
	SidebarLinkItem,
} from '../schemas/sidebar';
import { createPathFormatter } from './createPathFormatter';
import { formatPath } from './format-path';
import { pickLang } from './i18n';
import { ensureLeadingSlash, ensureTrailingSlash, stripLeadingAndTrailingSlashes } from './path';
import { getLocaleRoutes, type Route } from './routing';
import { localeToLang, slugToPathname } from './slugs';

const DirKey = Symbol('DirKey');
const SlugKey = Symbol('SlugKey');

export interface Link {
	type: 'link';
	label: string;
	href: string;
	isCurrent: boolean;
	badge: Badge | undefined;
	attrs: LinkHTMLAttributes;
}

interface Group {
	type: 'group';
	label: string;
	entries: (Link | Group)[];
	collapsed: boolean;
	badge: Badge | undefined;
}

export type SidebarEntry = Link | Group;

/**
 * A representation of the route structure. For each object entry:
 * if it’s a folder, the key is the directory name, and value is the directory
 * content; if it’s a route entry, the key is the last segment of the route, and value
 * is the full entry.
 */
interface Dir {
	[DirKey]: undefined;
	[SlugKey]: string;
	[item: string]: Dir | Route;
}

/** Create a new directory object. */
function makeDir(slug: string): Dir {
	const dir = {} as Dir;
	// Add DirKey and SlugKey as non-enumerable properties so that `Object.entries(dir)` ignores them.
	Object.defineProperty(dir, DirKey, { enumerable: false });
	Object.defineProperty(dir, SlugKey, { value: slug, enumerable: false });
	return dir;
}

/** Test if the passed object is a directory record.  */
function isDir(data: Record<string, unknown>): data is Dir {
	return DirKey in data;
}

/** Convert an item in a user’s sidebar config to a sidebar entry. */
function configItemToEntry(
	item: SidebarItem,
	currentPathname: string,
	locale: string | undefined,
	routes: Route[]
): SidebarEntry {
	if ('link' in item) {
		return linkFromConfig(item, locale, currentPathname);
	} else if ('autogenerate' in item) {
		return groupFromAutogenerateConfig(item, locale, routes, currentPathname);
	} else {
		return {
			type: 'group',
			label: pickLang(item.translations, localeToLang(locale)) || item.label,
			entries: item.items.map((i) => configItemToEntry(i, currentPathname, locale, routes)),
			collapsed: item.collapsed,
			badge: item.badge,
		};
	}
}

/** Autogenerate a group of links from a user’s sidebar config. */
function groupFromAutogenerateConfig(
	item: AutoSidebarGroup,
	locale: string | undefined,
	routes: Route[],
	currentPathname: string
): Group {
	const { collapsed: subgroupCollapsed, directory } = item.autogenerate;
	const localeDir = locale ? locale + '/' + directory : directory;
	const dirDocs = routes.filter(
		(doc) =>
			// Match against `foo.md` or `foo/index.md`.
			stripExtension(doc.id) === localeDir ||
			// Match against `foo/anything/else.md`.
			doc.id.startsWith(localeDir + '/')
	);
	const tree = treeify(dirDocs, localeDir);
	return {
		type: 'group',
		label: pickLang(item.translations, localeToLang(locale)) || item.label,
		entries: sidebarFromDir(tree, currentPathname, locale, subgroupCollapsed ?? item.collapsed),
		collapsed: item.collapsed,
		badge: item.badge,
	};
}

/** Check if a string starts with one of `http://` or `https://`. */
const isAbsolute = (link: string) => /^https?:\/\//.test(link);

/** Create a link entry from a user config object. */
function linkFromConfig(
	item: SidebarLinkItem,
	locale: string | undefined,
	currentPathname: string
) {
	let href = item.link;
	if (!isAbsolute(href)) {
		href = ensureLeadingSlash(href);
		// Inject current locale into link.
		if (locale) href = '/' + locale + href;
	}
	const label = pickLang(item.translations, localeToLang(locale)) || item.label;
	return makeLink(href, label, currentPathname, item.badge, item.attrs);
}

/** Create a link entry. */
function makeLink(
	href: string,
	label: string,
	currentPathname: string,
	badge?: Badge,
	attrs?: LinkHTMLAttributes
): Link {
	if (!isAbsolute(href)) {
		href = formatPath(href);
	}

	const isCurrent = pathsMatch(encodeURI(href), currentPathname);

	return { type: 'link', label, href, isCurrent, badge, attrs: attrs ?? {} };
}

/** Test if two paths are equivalent even if formatted differently. */
function pathsMatch(pathA: string, pathB: string) {
	const format = createPathFormatter({ trailingSlash: 'never' });
	return format(pathA) === format(pathB);
}

/** Get the segments leading to a page. */
function getBreadcrumbs(path: string, baseDir: string): string[] {
	// Strip extension from path.
	const pathWithoutExt = stripExtension(path);
	// Index paths will match `baseDir` and don’t include breadcrumbs.
	if (pathWithoutExt === baseDir) return [];
	// Ensure base directory ends in a trailing slash.
	baseDir = ensureTrailingSlash(baseDir);
	// Strip base directory from path if present.
	const relativePath = pathWithoutExt.startsWith(baseDir)
		? pathWithoutExt.replace(baseDir, '')
		: pathWithoutExt;

	return relativePath.split('/');
}

/** Turn a flat array of routes into a tree structure. */
function treeify(routes: Route[], baseDir: string): Dir {
	const treeRoot: Dir = makeDir(baseDir);
	routes
		// Remove any entries that should be hidden
		.filter((doc) => !doc.entry.data.sidebar.hidden)
		// Sort by depth, to build the tree depth first.
		.sort((a, b) => b.id.split('/').length - a.id.split('/').length)
		// Build the tree
		.forEach((doc) => {
			const parts = getBreadcrumbs(doc.id, baseDir);
			let currentNode = treeRoot;

			parts.forEach((part, index) => {
				const isLeaf = index === parts.length - 1;

				// Handle directory index pages by renaming them to `index`
				if (isLeaf && currentNode.hasOwnProperty(part)) {
					currentNode = currentNode[part] as Dir;
					part = 'index';
				}

				// Recurse down the tree if this isn’t the leaf node.
				if (!isLeaf) {
					const path = currentNode[SlugKey];
					currentNode[part] ||= makeDir(stripLeadingAndTrailingSlashes(path + '/' + part));
					currentNode = currentNode[part] as Dir;
				} else {
					currentNode[part] = doc;
				}
			});
		});

	return treeRoot;
}

/** Create a link entry for a given content collection entry. */
function linkFromRoute(route: Route, currentPathname: string): Link {
	return makeLink(
		slugToPathname(route.slug),
		route.entry.data.sidebar.label || route.entry.data.title,
		currentPathname,
		route.entry.data.sidebar.badge,
		route.entry.data.sidebar.attrs
	);
}

/**
 * Get the sort weight for a given route or directory. Lower numbers rank higher.
 * Directories have the weight of the lowest weighted route they contain.
 */
function getOrder(routeOrDir: Route | Dir): number {
	return isDir(routeOrDir)
		? Math.min(...Object.values(routeOrDir).flatMap(getOrder))
		: // If no order value is found, set it to the largest number possible.
		  routeOrDir.entry.data.sidebar.order ?? Number.MAX_VALUE;
}

/** Sort a directory’s entries by user-specified order or alphabetically if no order specified. */
function sortDirEntries(dir: [string, Dir | Route][]): [string, Dir | Route][] {
	const collator = new Intl.Collator(localeToLang(undefined));
	return dir.sort(([_keyA, a], [_keyB, b]) => {
		const [aOrder, bOrder] = [getOrder(a), getOrder(b)];
		// Pages are sorted by order in ascending order.
		if (aOrder !== bOrder) return aOrder < bOrder ? -1 : 1;
		// If two pages have the same order value they will be sorted by their slug.
		return collator.compare(isDir(a) ? a[SlugKey] : a.slug, isDir(b) ? b[SlugKey] : b.slug);
	});
}

/** Create a group entry for a given content collection directory. */
function groupFromDir(
	dir: Dir,
	fullPath: string,
	dirName: string,
	currentPathname: string,
	locale: string | undefined,
	collapsed: boolean
): Group {
	const entries = sortDirEntries(Object.entries(dir)).map(([key, dirOrRoute]) =>
		dirToItem(dirOrRoute, `${fullPath}/${key}`, key, currentPathname, locale, collapsed)
	);
	return {
		type: 'group',
		label: dirName,
		entries,
		collapsed,
		badge: undefined,
	};
}

/** Create a sidebar entry for a directory or content entry. */
function dirToItem(
	dirOrRoute: Dir[string],
	fullPath: string,
	dirName: string,
	currentPathname: string,
	locale: string | undefined,
	collapsed: boolean
): SidebarEntry {
	return isDir(dirOrRoute)
		? groupFromDir(dirOrRoute, fullPath, dirName, currentPathname, locale, collapsed)
		: linkFromRoute(dirOrRoute, currentPathname);
}

/** Create a sidebar entry for a given content directory. */
function sidebarFromDir(
	tree: Dir,
	currentPathname: string,
	locale: string | undefined,
	collapsed: boolean
) {
	return sortDirEntries(Object.entries(tree)).map(([key, dirOrRoute]) =>
		dirToItem(dirOrRoute, key, key, currentPathname, locale, collapsed)
	);
}

/** Get the sidebar for the current page. */
export function getSidebar(pathname: string, locale: string | undefined): SidebarEntry[] {
	const routes = getLocaleRoutes(locale);
	if (config.sidebar) {
		return config.sidebar.map((group) => configItemToEntry(group, pathname, locale, routes));
	} else {
		const tree = treeify(routes, locale || '');
		return sidebarFromDir(tree, pathname, locale, false);
	}
}

/** Turn the nested tree structure of a sidebar into a flat list of all the links. */
export function flattenSidebar(sidebar: SidebarEntry[]): Link[] {
	return sidebar.flatMap((entry) =>
		entry.type === 'group' ? flattenSidebar(entry.entries) : entry
	);
}

/** Get previous/next pages in the sidebar or the ones from the frontmatter if any. */
export function getPrevNextLinks(
	sidebar: SidebarEntry[],
	paginationEnabled: boolean,
	config: {
		prev?: PrevNextLinkConfig;
		next?: PrevNextLinkConfig;
	}
): {
	/** Link to previous page in the sidebar. */
	prev: Link | undefined;
	/** Link to next page in the sidebar. */
	next: Link | undefined;
} {
	const entries = flattenSidebar(sidebar);
	const currentIndex = entries.findIndex((entry) => entry.isCurrent);
	const prev = applyPrevNextLinkConfig(entries[currentIndex - 1], paginationEnabled, config.prev);
	const next = applyPrevNextLinkConfig(
		currentIndex > -1 ? entries[currentIndex + 1] : undefined,
		paginationEnabled,
		config.next
	);
	return { prev, next };
}

/** Apply a prev/next link config to a navigation link. */
function applyPrevNextLinkConfig(
	link: Link | undefined,
	paginationEnabled: boolean,
	config: PrevNextLinkConfig | undefined
): Link | undefined {
	// Explicitly remove the link.
	if (config === false) return undefined;
	// Use the generated link if any.
	else if (config === true) return link;
	// If a link exists, update its label if needed.
	else if (typeof config === 'string' && link) {
		return { ...link, label: config };
	} else if (typeof config === 'object') {
		if (link) {
			// If a link exists, update both its label and href if needed.
			return {
				...link,
				label: config.label ?? link.label,
				href: config.link ?? link.href,
				// Explicitly remove sidebar link attributes for prev/next links.
				attrs: {},
			};
		} else if (config.link && config.label) {
			// If there is no link and the frontmatter contains both a URL and a label,
			// create a new link.
			return makeLink(config.link, config.label, '');
		}
	}
	// Otherwise, if the global config is enabled, return the generated link if any.
	return paginationEnabled ? link : undefined;
}

/** Remove the extension from a path. */
const stripExtension = (path: string) => path.replace(/\.\w+$/, '');
