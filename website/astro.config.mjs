import { unified } from "@astrojs/markdown-remark";
import mdx from "@astrojs/mdx";
import sitemap from "@astrojs/sitemap";
import { defineConfig } from "astro/config";
import remarkDirective from "remark-directive";
import { remarkDirectivesHandler } from "./src/plugins/remarkDirectives.mjs";

const site = "https://superfile.dev/";

export default defineConfig({
	site,
	compressHTML: true,
	integrations: [mdx(), sitemap()],

	markdown: {
		processor: unified({
			gfm: true,
			remarkPlugins: [remarkDirective, remarkDirectivesHandler]
		}),
		shikiConfig: {
			theme: "catppuccin-mocha",
			wrap: false
		}
	},

	image: {
		service: {
			entrypoint: "astro/assets/services/sharp",
			config: { limitInputPixels: false }
		}
	}
});
