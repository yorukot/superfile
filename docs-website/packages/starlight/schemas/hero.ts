import { z } from 'astro/zod';
import type { SchemaContext } from 'astro:content';
import { Icons } from '../components/Icons';

type IconName = keyof typeof Icons;
const iconNames = Object.keys(Icons) as [IconName, ...IconName[]];

export const HeroSchema = ({ image }: SchemaContext) =>
	z.object({
		/**
		 * The large title text to show. If not provided, will default to the top-level `title`.
		 * Can include HTML.
		 */
		title: z.string().optional(),
		/**
		 * A short bit of text about your project.
		 * Will be displayed in a smaller size below the title.
		 */
		tagline: z.string().optional(),
		/** The image to use in the hero. You can provide either a relative `file` path or raw `html`. */
		image: z
			.union([
				z.object({
					/** Alt text for screenreaders and other assistive technologies describing your hero image. */
					alt: z.string().default(''),
					/** Relative path to an image file in your repo, e.g. `../../assets/hero.png`. */
					file: image(),
				}),
				z.object({
					/** Alt text for screenreaders and other assistive technologies describing your hero image. */
					alt: z.string().default(''),
					/** Relative path to an image file in your repo to use in dark mode, e.g. `../../assets/hero-dark.png`. */
					dark: image(),
					/** Relative path to an image file in your repo to use in light mode, e.g. `../../assets/hero-light.png`. */
					light: image(),
				}),
				z
					.object({
						/** Raw HTML string instead of an image file. Useful for inline SVGs or more complex hero content. */
						html: z.string(),
					})
					.transform(({ html }) => ({ html, alt: '' })),
			])
			.optional(),
		/** An array of call-to-action links displayed at the bottom of the hero. */
		actions: z
			.object({
				/** Text label displayed in the link. */
				text: z.string(),
				/** Value for the link’s `href` attribute, e.g. `/page` or `https://mysite.com`. */
				link: z.string(),
				/** Button style to use. One of `primary`, `secondary`, or `minimal` (the default). */
				variant: z.enum(['primary', 'secondary', 'minimal']).default('minimal'),
				/**
				 * An optional icon to display alongside the link text.
				 * Can be an inline `<svg>` or the name of one of Starlight’s built-in icons.
				 */
				icon: z
					.union([z.enum(iconNames), z.string().startsWith('<svg')])
					.transform((icon) => {
						const parsedIcon = z.enum(iconNames).safeParse(icon);
						return parsedIcon.success
							? ({ type: 'icon', name: parsedIcon.data } as const)
							: ({ type: 'raw', html: icon } as const);
					})
					.optional(),
				/** HTML attributes to add to the link */
				attrs: z.record(z.union([z.string(), z.number(), z.boolean()])).optional(),
			})
			.array()
			.default([]),
	});
