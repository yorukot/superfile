import type { AstroConfig } from 'astro';
import { fileWithBase, pathWithBase } from './base';
import {
	ensureHtmlExtension,
	ensureTrailingSlash,
	stripHtmlExtension,
	stripTrailingSlash,
} from './path';

interface FormatPathOptions {
	format?: AstroConfig['build']['format'];
	trailingSlash?: AstroConfig['trailingSlash'];
}

const formatStrategies = {
	file: {
		addBase: fileWithBase,
		handleExtension: (href: string) => ensureHtmlExtension(href),
	},
	directory: {
		addBase: pathWithBase,
		handleExtension: (href: string) => stripHtmlExtension(href),
	},
};

const trailingSlashStrategies = {
	always: ensureTrailingSlash,
	never: stripTrailingSlash,
	ignore: (href: string) => href,
};

/** Format a path based on the project config. */
function formatPath(
	href: string,
	{ format = 'directory', trailingSlash = 'ignore' }: FormatPathOptions
) {
	// @ts-expect-error — TODO: add support for `preserve` (https://github.com/withastro/starlight/issues/1781)
	const formatStrategy = formatStrategies[format];
	const trailingSlashStrategy = trailingSlashStrategies[trailingSlash];

	// Add base
	href = formatStrategy.addBase(href);

	// Handle extension
	href = formatStrategy.handleExtension(href);

	// Skip trailing slash handling for `build.format: 'file'`
	if (format === 'file') return href;

	// Handle trailing slash
	href = href === '/' ? href : trailingSlashStrategy(href);

	return href;
}

export function createPathFormatter(opts: FormatPathOptions) {
	return (href: string) => formatPath(href, opts);
}
