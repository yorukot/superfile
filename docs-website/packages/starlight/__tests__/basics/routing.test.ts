import { getCollection } from 'astro:content';
import config from 'virtual:starlight/user-config';
import { expect, test, vi } from 'vitest';
import { routes } from '../../utils/routing';

vi.mock('astro:content', async () =>
	(await import('../test-utils')).mockedAstroContent({
		docs: [
			['404.md', { title: 'Not found' }],
			['index.mdx', { title: 'Home page' }],
			['guides/authoring-content.md', { title: 'Authoring content', draft: true }],
		],
	})
);

test('test suite is using correct env', () => {
	expect(config.title).toMatchObject({ en: 'Basics' });
});

test('route slugs are normalized', () => {
	const indexRoute = routes.find((route) => route.id.startsWith('index.md'));
	expect(indexRoute?.slug).toBe('');
});

test('routes contain copy of original doc as entry', async () => {
	const docs = await getCollection('docs');
	for (const route of routes) {
		const doc = docs.find((doc) => doc.id === route.id);
		if (!doc) throw new Error('Expected to find doc for route ' + route.id);
		// Compare without slug as slugs can be normalized.
		const { slug: _, ...entry } = route.entry;
		const { slug: __, ...input } = doc;
		expect(entry).toEqual(input);
	}
});

test('routes have locale data added', () => {
	for (const route of routes) {
		expect(route.lang).toBe('en');
		expect(route.dir).toBe('ltr');
		expect(route.locale).toBeUndefined();
	}
});

test('routes includes drafts except in production', async () => {
	expect(routes.find((route) => route.id === 'guides/authoring-content.md')).toBeTruthy();

	// Reset the modules registry so that re-importing `utils/routing.ts` re-evaluates the module and
	// re-computes the routes. Re-importing the module is necessary because top-level imports cannot
	// be re-evaluated.
	vi.resetModules();
	// Set the mode to production.
	vi.stubEnv('MODE', 'production');
	// Re-import the module to re-evaluate it.
	const { routes: prodRoutes } = await import('../../utils/routing');

	expect(prodRoutes.find((route) => route.id === 'guides/authoring-content.md')).toBeFalsy();

	vi.unstubAllEnvs();
	vi.resetModules();
});
