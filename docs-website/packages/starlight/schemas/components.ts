import { z } from 'astro/zod';

export function ComponentConfigSchema() {
	return z
		.object({
			/*
			HEAD ----------------------------------------------------------------------------------------
			*/

			/**
			 * Component rendered inside each page’s `<head>`.
			 * Includes important tags including `<title>`, and `<meta charset="utf-8">`.
			 *
			 * Override this component as a last resort. Prefer the `head` option Starlight config if possible.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/Head.astro `Head` default implementation}
			 */
			Head: z.string().default('@astrojs/starlight/components/Head.astro'),

			/**
			 * Component rendered inside `<head>` that sets up dark/light theme support.
			 * The default implementation includes an inline script and a `<template>` used by the
			 * script in `ThemeSelect.astro`.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeProvider.astro `ThemeProvider` default implementation}
			 */
			ThemeProvider: z.string().default('@astrojs/starlight/components/ThemeProvider.astro'),

			/*
			BODY ----------------------------------------------------------------------------------------
			*/

			/**
			 * Component rendered as the first element inside `<body>` which links to the main page
			 * content for accessibility. The default implementation is hidden until a user focuses it
			 * by tabbing with their keyboard.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/SkipLink.astro `SkipLink` default implementation}
			 */
			SkipLink: z.string().default('@astrojs/starlight/components/SkipLink.astro'),

			/*
			LAYOUT --------------------------------------------------------------------------------------
			*/

			/**
			 * Layout component wrapped around most of the page content.
			 * The default implementation sets up the header–sidebar–main layout and includes
			 * `header` and `sidebar` named slots along with a default slot for the main content.
			 * It also renders `<MobileMenuToggle />` to support toggling the sidebar navigation
			 * on small (mobile) viewports.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageFrame.astro `PageFrame` default implementation}
			 */
			PageFrame: z.string().default('@astrojs/starlight/components/PageFrame.astro'),
			/**
			 * Component rendered inside `<PageFrame>` that is responsible for toggling the
			 * sidebar navigation on small (mobile) viewports.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuToggle.astro `MobileMenuToggle` default implementation}
			 */
			MobileMenuToggle: z.string().default('@astrojs/starlight/components/MobileMenuToggle.astro'),

			/**
			 * Layout component wrapped around the main content column and right sidebar (table of contents).
			 * The default implementation handles the switch between a single-column, small-viewport layout
			 * and a two-column, larger-viewport layout.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/TwoColumnContent.astro `TwoColumnContent` default implementation}
			 */
			TwoColumnContent: z.string().default('@astrojs/starlight/components/TwoColumnContent.astro'),

			/*
			HEADER --------------------------------------------------------------------------------------
			*/

			/**
			 * Header component displayed at the top of every page.
			 * The default implementation displays `<SiteTitle />`, `<Search />`, `<SocialIcons />`,
			 * `<ThemeSelect />`, and `<LanguageSelect />`.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/Header.astro `Header` default implementation}
			 */
			Header: z.string().default('@astrojs/starlight/components/Header.astro'),
			/**
			 * Component rendered at the start of the site header to render the site title.
			 * The default implementation includes logic for rendering logos defined in Starlight config.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/SiteTitle.astro `SiteTitle` default implementation}
			 */
			SiteTitle: z.string().default('@astrojs/starlight/components/SiteTitle.astro'),
			/**
			 * Component used to render Starlight’s search UI. The default implementation includes the
			 * button in the header and the code for displaying a search modal when it is clicked.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/Search.astro `Search` default implementation}
			 */
			Search: z.string().default('@astrojs/starlight/components/Search.astro'),
			/**
			 * Component rendered in the site header including social icon links. The default
			 * implementation uses the `social` option in Starlight config to render icons and links.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/SocialIcons.astro `SocialIcons` default implementation}
			 */
			SocialIcons: z.string().default('@astrojs/starlight/components/SocialIcons.astro'),
			/**
			 * Component rendered in the site header that allows users to select their preferred color scheme.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeSelect.astro `ThemeSelect` default implementation}
			 */
			ThemeSelect: z.string().default('@astrojs/starlight/components/ThemeSelect.astro'),
			/**
			 * Component rendered in the site header that allows users to switch to a different language.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/LanguageSelect.astro `LanguageSelect` default implementation}
			 */
			LanguageSelect: z.string().default('@astrojs/starlight/components/LanguageSelect.astro'),

			/*
			SIDEBAR -------------------------------------------------------------------------------------
			*/

			/**
			 * Component rendered before page content that contains global navigation.
			 * The default implementation displays as a sidebar on wide enough viewports and inside a
			 * drop-down menu on small (mobile) viewports. It also renders `<MobileMenuFooter />` to
			 * show additional items inside the mobile menu.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/Sidebar.astro `Sidebar` default implementation}
			 */
			Sidebar: z.string().default('@astrojs/starlight/components/Sidebar.astro'),
			/**
			 * Component rendered at the bottom of the mobile drop-down menu.
			 * The default implementation renders `<ThemeSelect />` and `<LanguageSelect />`.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuFooter.astro `MobileMenuFooter` default implementation}
			 */
			MobileMenuFooter: z.string().default('@astrojs/starlight/components/MobileMenuFooter.astro'),

			/*
			TOC -----------------------------------------------------------------------------------------
			*/

			/**
			 * Component rendered before the main page’s content to display a table of contents.
			 * The default implementation renders `<TableOfContents />` and `<MobileTableOfContents />`.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageSidebar.astro `PageSidebar` default implementation}
			 */
			PageSidebar: z.string().default('@astrojs/starlight/components/PageSidebar.astro'),
			/**
			 * Component that renders the current page’s table of contents on wider viewports.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/TableOfContents.astro `TableOfContents` default implementation}
			 */
			TableOfContents: z.string().default('@astrojs/starlight/components/TableOfContents.astro'),
			/**
			 * Component that renders the current page’s table of contents on small (mobile) viewports.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileTableOfContents.astro `MobileTableOfContents` default implementation}
			 */
			MobileTableOfContents: z
				.string()
				.default('@astrojs/starlight/components/MobileTableOfContents.astro'),

			/*
			CONTENT HEADER ------------------------------------------------------------------------------
			*/

			/**
			 * Banner component rendered at the top of each page. The default implementation uses the
			 * page’s `banner` frontmatter value to decide whether or not to render.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/Banner.astro `Banner` default implementation}
			 */
			Banner: z.string().default('@astrojs/starlight/components/Banner.astro'),

			/**
			 * Layout component used to wrap sections of the main content column.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/ContentPanel.astro `ContentPanel` default implementation}
			 */
			ContentPanel: z.string().default('@astrojs/starlight/components/ContentPanel.astro'),

			/**
			 * Component containing the `<h1>` element for the current page.
			 *
			 * Implementations should ensure they set `id="_top"` on the `<h1>` element as in the default
			 * implementation.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageTitle.astro `PageTitle` default implementation}
			 */
			PageTitle: z.string().default('@astrojs/starlight/components/PageTitle.astro'),

			/**
			 * Notice displayed to users on pages where a translation for the current language is not
			 * available. Only used on multilingual sites.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/FallbackContentNotice.astro `FallbackContentNotice` default implementation}
			 */
			FallbackContentNotice: z
				.string()
				.default('@astrojs/starlight/components/FallbackContentNotice.astro'),

			/**
			 * Notice displayed to users on draft pages. Only used in development mode.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/DraftContentNotice.astro `DraftContentNotice` default implementation}
			 */
			DraftContentNotice: z
				.string()
				.default('@astrojs/starlight/components/DraftContentNotice.astro'),

			/**
			 * Component rendered at the top of the page when `hero` is set in frontmatter. The default
			 * implementation shows a large title, tagline, and call-to-action links alongside an optional image.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/Hero.astro `Hero` default implementation}
			 */
			Hero: z.string().default('@astrojs/starlight/components/Hero.astro'),

			/*
			CONTENT -------------------------------------------------------------------------------------
			*/

			/**
			 * Component rendered around each page’s main content.
			 * The default implementation sets up basic styles to apply to Markdown content.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/MarkdownContent.astro `MarkdownContent` default implementation}
			 */
			MarkdownContent: z.string().default('@astrojs/starlight/components/MarkdownContent.astro'),

			/*
			CONTENT FOOTER ------------------------------------------------------------------------------
			*/

			/**
			 * Footer component displayed at the bottom of each documentation page.
			 * The default implementation displays `<LastUpdated />`, `<Pagination />`, and `<EditLink />`.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/Footer.astro `Footer` default implementation}
			 */
			Footer: z.string().default('@astrojs/starlight/components/Footer.astro'),
			/**
			 * Component rendered in the page footer to display the last-updated date.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/LastUpdated.astro `LastUpdated` default implementation}
			 */
			LastUpdated: z.string().default('@astrojs/starlight/components/LastUpdated.astro'),
			/**
			 * Component rendered in the page footer to display navigation arrows between previous/next pages.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/Pagination.astro `Pagination` default implementation}
			 */
			Pagination: z.string().default('@astrojs/starlight/components/Pagination.astro'),
			/**
			 * Component rendered in the page footer to display a link to where the page can be edited.
			 *
			 * @see {@link https://github.com/withastro/starlight/blob/main/packages/starlight/components/EditLink.astro `EditLink` default implementation}
			 */
			EditLink: z.string().default('@astrojs/starlight/components/EditLink.astro'),
		})
		.default({});
}
