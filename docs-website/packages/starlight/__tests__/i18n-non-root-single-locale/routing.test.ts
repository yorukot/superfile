import { expect, test, vi } from 'vitest';
import { routes } from '../../utils/routing';

vi.mock('astro:content', async () =>
	(await import('../test-utils')).mockedAstroContent({
		docs: [
			['fr/index.mdx', { title: 'Accueil' }],
			// @ts-expect-error — Using a slug not present in Starlight docs site
			['en/index.mdx', { title: 'Home page' }],
		],
	})
);

test('route slugs are normalized', () => {
	const indexRoute = routes.find((route) => route.id.startsWith('fr/index.md'));
	expect(indexRoute?.slug).toBe('fr');
});

test('routes for the configured locale have locale data added', () => {
	for (const route of routes) {
		if (route.id.startsWith('fr')) {
			expect(route.lang).toBe('fr-CA');
			expect(route.dir).toBe('ltr');
			expect(route.locale).toBe('fr');
		} else {
			expect(route.lang).toBe('fr-CA');
			expect(route.dir).toBe('ltr');
			expect(route.locale).toBeUndefined();
		}
	}
});

test('does not mark any route as fallback routes', () => {
	const fallbacks = routes.filter((route) => route.isFallback);
	expect(fallbacks.length).toBe(0);
});
