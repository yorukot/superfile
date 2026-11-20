import { unified } from "@astrojs/markdown-remark";
import mdx from "@astrojs/mdx";
import sitemap from "@astrojs/sitemap";
import { defineConfig } from "astro/config";
import remarkDirective from "remark-directive";
import { superfileShikiTheme } from "./src/lib/shiki-theme.mjs";
import { rehypeDocs } from "./src/plugins/rehype-docs.mjs";
import { remarkCallouts } from "./src/plugins/remark-callouts.mjs";

export default defineConfig({
	site: "https://superfile.dev",
	// Astro 7 strips the whitespace between inline elements by default ("jsx"); the pages are written for the HTML rule.
	compressHTML: true,
	prefetch: { prefetchAll: true },
	integrations: [mdx(), sitemap()],

	markdown: {
		// remark and rehype plugins need the unified processor; Astro's default (Sätteri) does not run them.
		processor: unified({
			remarkPlugins: [remarkDirective, remarkCallouts],
			rehypePlugins: [rehypeDocs]
		}),
		shikiConfig: { theme: superfileShikiTheme }
	},

	// The demo GIFs are larger than sharp's default pixel limit.
	image: { service: { config: { limitInputPixels: false } } }
});
