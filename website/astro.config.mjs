import mdx from "@astrojs/mdx";
import sitemap from "@astrojs/sitemap";
import { defineConfig } from "astro/config";
import remarkDirective from "remark-directive";
import { remarkDirectivesHandler } from "./src/plugins/remarkDirectives.mjs";

const site = "https://superfile.dev/";

export default defineConfig({
	site,
	integrations: [mdx(), sitemap()],

	markdown: {
		// Required for GFM tables in MDX content; undefined is treated as false by @astrojs/mdx.
		gfm: true,
		shikiConfig: {
			theme: "catppuccin-mocha",
			wrap: false
		},
		remarkPlugins: [remarkDirective, remarkDirectivesHandler]
	},

	image: {
		service: {
			entrypoint: "astro/assets/services/sharp",
			config: { limitInputPixels: false }
		}
	}
});
