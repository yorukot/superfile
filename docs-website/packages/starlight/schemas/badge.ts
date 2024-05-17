import { z } from 'astro/zod';

const badgeSchema = () =>
	z.object({
		variant: z.enum(['note', 'danger', 'success', 'caution', 'tip', 'default']).default('default'),
		text: z.string(),
	});

export const BadgeConfigSchema = () =>
	z
		.union([z.string(), badgeSchema()])
		.transform((badge) => {
			if (typeof badge === 'string') {
				return { variant: 'default' as const, text: badge };
			}
			return badge;
		})
		.optional();

export type Badge = z.output<ReturnType<typeof badgeSchema>>;
