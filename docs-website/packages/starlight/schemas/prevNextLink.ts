import { z } from 'astro/zod';

export const PrevNextLinkConfigSchema = () =>
	z
		.union([
			z.boolean(),
			z.string(),
			z
				.object({
					/** The navigation link URL. */
					link: z.string().optional(),
					/** The navigation link text. */
					label: z.string().optional(),
				})
				.strict(),
		])
		.optional();

export type PrevNextLinkUserConfig = z.input<ReturnType<typeof PrevNextLinkConfigSchema>>;
export type PrevNextLinkConfig = z.output<ReturnType<typeof PrevNextLinkConfigSchema>>;
