import { glob } from "astro/loaders";
import { z } from "astro/zod";
import { defineCollection } from "astro:content";

// Pages carry more front matter than this (old releases especially); only these two fields are read.
const schema = z.object({
	title: z.string().optional(),
	description: z.string().optional()
});

export const collections = {
	// The current docs.
	docs: defineCollection({
		loader: glob({ pattern: "**/*.{md,mdx}", base: "./src/content/docs" }),
		schema
	}),
	// The docs of released versions, copied out of git tags by scripts/sync-versions.mjs.
	// An entry's id is "<tag>/<slug>", for example "v1.5.0/getting-started/tutorial".
	versions: defineCollection({
		loader: glob({
			pattern: "*/src/content/docs/**/*.{md,mdx}",
			base: "./.versions",
			generateId: ({ entry }) => entry.replace("/src/content/docs/", "/").replace(/\.mdx?$/, "")
		}),
		schema
	})
};
