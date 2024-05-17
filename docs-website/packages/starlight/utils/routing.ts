import type { GetStaticPathsItem } from 'astro';
import { type CollectionEntry, getCollection } from 'astro:content';
import config from 'virtual:starlight/user-config';
import {
	type LocaleData,
	localizedId,
	localizedSlug,
	slugToLocaleData,
	slugToParam,
} from './slugs';
import { validateLogoImports } from './validateLogoImports';

// Validate any user-provided logos imported correctly.
// We do this here so all pages trigger it and at the top level so it runs just once.
validateLogoImports();

export type StarlightDocsEntry = Omit<CollectionEntry<'docs'>, 'slug'> & {
	slug: string;
};

export interface Route extends LocaleData {
	/** Content collection entry for the current page. Includes frontmatter at `data`. */
	entry: StarlightDocsEntry;
	/** Locale metadata for the page content. Can be different from top-level locale values when a page is using fallback content. */
	entryMeta: LocaleData;
	/** The slug, a.k.a. permalink, for this page. */
	slug: string;
	/** The unique ID for this page. */
	id: string;
	/** True if this page is untranslated in the current language and using fallback content from the default locale. */
	isFallback?: true;
	[key: string]: unknown;
}

interface Path extends GetStaticPathsItem {
	params: { slug: string | undefined };
	props: Route;
}

/**
 * Astro is inconsistent in its `index.md` slug generation. In most cases,
 * `index` is stripped, but in the root of a collection, we get a slug of `index`.
 * We map that to an empty string for consistent behaviour.
 */
const normalizeIndexSlug = (slug: string) => (slug === 'index' ? '' : slug);

/** All entries in the docs content collection. */
const docs: StarlightDocsEntry[] = (
	(await getCollection('docs', ({ data }) => {
		// In production, filter out drafts.
		return import.meta.env.MODE !== 'production' || data.draft === false;
	})) ?? []
).map(({ slug, ...entry }) => ({
	...entry,
	slug: normalizeIndexSlug(slug),
}));

function getRoutes(): Route[] {
	const routes: Route[] = docs.map((entry) => ({
		entry,
		slug: entry.slug,
		id: entry.id,
		entryMeta: slugToLocaleData(entry.slug),
		...slugToLocaleData(entry.slug),
	}));

	// In multilingual sites, add required fallback routes.
	if (config.isMultilingual) {
		/** Entries in the docs content collection for the default locale. */
		const defaultLocaleDocs = getLocaleDocs(
			config.defaultLocale?.locale === 'root' ? undefined : config.defaultLocale?.locale
		);
		for (const key in config.locales) {
			if (key === config.defaultLocale.locale) continue;
			const localeConfig = config.locales[key];
			if (!localeConfig) continue;
			const locale = key === 'root' ? undefined : key;
			const localeDocs = getLocaleDocs(locale);
			for (const fallback of defaultLocaleDocs) {
				const slug = localizedSlug(fallback.slug, locale);
				const id = localizedId(fallback.id, locale);
				const doesNotNeedFallback = localeDocs.some((doc) => doc.slug === slug);
				if (doesNotNeedFallback) continue;
				routes.push({
					entry: fallback,
					slug,
					id,
					isFallback: true,
					lang: localeConfig.lang || 'en',
					locale,
					dir: localeConfig.dir,
					entryMeta: slugToLocaleData(fallback.slug),
				});
			}
		}
	}

	return routes;
}
export const routes = getRoutes();

function getPaths(): Path[] {
	return routes.map((route) => ({
		params: { slug: slugToParam(route.slug) },
		props: route,
	}));
}
export const paths = getPaths();

/**
 * Get all routes for a specific locale.
 * A locale of `undefined` is treated as the “root” locale, if configured.
 */
export function getLocaleRoutes(locale: string | undefined): Route[] {
	return filterByLocale(routes, locale);
}

/**
 * Get all entries in the docs content collection for a specific locale.
 * A locale of `undefined` is treated as the “root” locale, if configured.
 */
function getLocaleDocs(locale: string | undefined): StarlightDocsEntry[] {
	return filterByLocale(docs, locale);
}

/** Filter an array to find items whose slug matches the passed locale. */
function filterByLocale<T extends { slug: string }>(items: T[], locale: string | undefined): T[] {
	if (config.locales) {
		if (locale && locale in config.locales) {
			return items.filter((i) => i.slug === locale || i.slug.startsWith(locale + '/'));
		} else if (config.locales.root) {
			const langKeys = Object.keys(config.locales).filter((k) => k !== 'root');
			const isLangIndex = new RegExp(`^(${langKeys.join('|')})$`);
			const isLangDir = new RegExp(`^(${langKeys.join('|')})/`);
			return items.filter((i) => !isLangIndex.test(i.slug) && !isLangDir.test(i.slug));
		}
	}
	return items;
}
