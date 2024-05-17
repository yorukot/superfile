import type { MarkdownHeading } from 'astro';
import { fileURLToPath } from 'node:url';
import project from 'virtual:starlight/project-context';
import config from 'virtual:starlight/user-config';
import { generateToC, type TocItem } from './generateToC';
import { getNewestCommitDate } from './git';
import { getPrevNextLinks, getSidebar, type SidebarEntry } from './navigation';
import { ensureTrailingSlash } from './path';
import type { Route } from './routing';
import { localizedId } from './slugs';
import { useTranslations } from './translations';
import { formatPath } from './format-path';

export interface PageProps extends Route {
	headings: MarkdownHeading[];
}

export interface StarlightRouteData extends Route {
	/** Title of the site. */
	siteTitle: string;
	/** URL or path used as the link when clicking on the site title. */
	siteTitleHref: string;
	/** Array of Markdown headings extracted from the current page. */
	headings: MarkdownHeading[];
	/** Site navigation sidebar entries for this page. */
	sidebar: SidebarEntry[];
	/** Whether or not the sidebar should be displayed on this page. */
	hasSidebar: boolean;
	/** Links to the previous and next page in the sidebar if enabled. */
	pagination: ReturnType<typeof getPrevNextLinks>;
	/** Table of contents for this page if enabled. */
	toc: { minHeadingLevel: number; maxHeadingLevel: number; items: TocItem[] } | undefined;
	/** JS Date object representing when this page was last updated if enabled. */
	lastUpdated: Date | undefined;
	/** URL object for the address where this page can be edited if enabled. */
	editUrl: URL | undefined;
	/** Record of UI strings localized for the current page. */
	labels: ReturnType<ReturnType<typeof useTranslations>['all']>;
}

export function generateRouteData({
	props,
	url,
}: {
	props: PageProps;
	url: URL;
}): StarlightRouteData {
	const { entry, locale, lang } = props;
	const sidebar = getSidebar(url.pathname, locale);
	const siteTitle = getSiteTitle(lang);
	return {
		...props,
		siteTitle,
		siteTitleHref: getSiteTitleHref(locale),
		sidebar,
		hasSidebar: entry.data.template !== 'splash',
		pagination: getPrevNextLinks(sidebar, config.pagination, entry.data),
		toc: getToC(props),
		lastUpdated: getLastUpdated(props),
		editUrl: getEditUrl(props),
		labels: useTranslations(locale).all(),
	};
}

export function getToC({ entry, locale, headings }: PageProps) {
	const tocConfig =
		entry.data.template === 'splash'
			? false
			: entry.data.tableOfContents !== undefined
			? entry.data.tableOfContents
			: config.tableOfContents;
	if (!tocConfig) return;
	const t = useTranslations(locale);
	return {
		...tocConfig,
		items: generateToC(headings, { ...tocConfig, title: t('tableOfContents.overview') }),
	};
}

function getLastUpdated({ entry }: PageProps): Date | undefined {
	const { lastUpdated: frontmatterLastUpdated } = entry.data;
	const { lastUpdated: configLastUpdated } = config;

	if (frontmatterLastUpdated ?? configLastUpdated) {
		const currentFilePath = fileURLToPath(new URL('src/content/docs/' + entry.id, project.root));
		try {
			return frontmatterLastUpdated instanceof Date
				? frontmatterLastUpdated
				: getNewestCommitDate(currentFilePath);
		} catch {
			// If the git command fails, ignore the error.
			return undefined;
		}
	}

	return undefined;
}

function getEditUrl({ entry, id, isFallback }: PageProps): URL | undefined {
	const { editUrl } = entry.data;
	// If frontmatter value is false, editing is disabled for this page.
	if (editUrl === false) return;

	let url: string | undefined;
	if (typeof editUrl === 'string') {
		// If a URL was provided in frontmatter, use that.
		url = editUrl;
	} else if (config.editLink.baseUrl) {
		const srcPath = project.srcDir.replace(project.root, '');
		const filePath = isFallback ? localizedId(id, config.defaultLocale.locale) : id;
		// If a base URL was added in Starlight config, synthesize the edit URL from it.
		url = ensureTrailingSlash(config.editLink.baseUrl) + srcPath + 'content/docs/' + filePath;
	}
	return url ? new URL(url) : undefined;
}

/** Get the site title for a given language. **/
export function getSiteTitle(lang: string): string {
	const defaultLang = config.defaultLocale.lang as string;
	if (lang && config.title[lang]) {
		return config.title[lang] as string;
	}
	return config.title[defaultLang] as string;
}

export function getSiteTitleHref(locale: string | undefined): string {
	return formatPath(locale || '/');
}
