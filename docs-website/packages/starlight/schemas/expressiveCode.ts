import { z } from 'astro/zod';
import type { StarlightExpressiveCodeOptions } from '../integrations/expressive-code';

export const ExpressiveCodeSchema = () =>
	z
		.union([
			z.custom<StarlightExpressiveCodeOptions>((value) => typeof value === 'object' && value),
			z.boolean(),
		])
		.describe(
			'Define how code blocks are rendered by passing options to Expressive Code, or disable the integration by passing `false`.'
		)
		.optional();
