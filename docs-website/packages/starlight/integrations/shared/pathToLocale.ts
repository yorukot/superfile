import type { AstroConfig } from 'astro';
import type { StarlightConfig } from '../../types';

function slugToLocale(
	slug: string | undefined,
	localesConfig: StarlightConfig['locales']
): string | undefined {
	const locales = Object.keys(localesConfig || {});
	const baseSegment = slug?.split('/')[0];
	return baseSegment && locales.includes(baseSegment) ? baseSegment : undefined;
}

/** Get current locale from the full file path. */
export function pathToLocale(
	path: string | undefined,
	{
		starlightConfig,
		astroConfig,
	}: {
		starlightConfig: { locales: StarlightConfig['locales'] };
		astroConfig: { root: AstroConfig['root']; srcDir: AstroConfig['srcDir'] };
	}
): string | undefined {
	const srcDir = new URL(astroConfig.srcDir, astroConfig.root);
	const docsDir = new URL('content/docs/', srcDir);
	// Format path to unix style path.
	path = path?.replace(/\\/g, '/');
	// Ensure that the page path starts with a slash if the docs directory also does,
	// which makes stripping the docs path in the next step work on Windows, too.
	if (path && !path.startsWith('/') && docsDir.pathname.startsWith('/')) path = '/' + path;
	// Strip docs path leaving only content collection file ID.
	// Example: /Users/houston/repo/src/content/docs/en/guide.md => en/guide.md
	const slug = path?.replace(docsDir.pathname, '');
	return slugToLocale(slug, starlightConfig.locales);
}
