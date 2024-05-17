import { expect, test, vi } from 'vitest';
import {
	generateStarlightPageRouteData,
	type StarlightPageProps,
} from '../../utils/starlight-page';

vi.mock('virtual:starlight/collection-config', async () => {
	const { z } = await vi.importActual<typeof import('astro:content')>('astro:content');
	return (await import('../test-utils')).mockedCollectionConfig({
		extend: z.object({
			// Make the built-in description field required.
			description: z.string(),
			// Add a new optional field.
			category: z.string().optional(),
		}),
	});
});

const starlightPageProps: StarlightPageProps = {
	frontmatter: { title: 'This is a test title' },
};

test('throws a validation error if a built-in field required by the user schema is not passed down', async () => {
	// The first line should be a user-friendly error message describing the exact issue and the second line should be
	// the missing description field.
	expect(() =>
		generateStarlightPageRouteData({
			props: starlightPageProps,
			url: new URL('https://example.com/test-slug'),
		})
	).rejects.toThrowErrorMatchingInlineSnapshot(`
		"[AstroUserError]:
			Invalid frontmatter props passed to the \`<StarlightPage/>\` component.
		Hint:
			**description**: Required"
	`);
});

test('returns new field defined in the user schema', async () => {
	const category = 'test category';
	const data = await generateStarlightPageRouteData({
		props: {
			...starlightPageProps,
			frontmatter: {
				...starlightPageProps.frontmatter,
				description: 'test description',
				// @ts-expect-error - Custom field defined in the user schema.
				category,
			},
		},
		url: new URL('https://example.com/test-slug'),
	});
	// @ts-expect-error - Custom field defined in the user schema.
	expect(data.entry.data.category).toBe(category);
});
