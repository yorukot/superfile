import opentype, { type Font, Glyph } from 'opentype.js';
import { seti, starlight } from '../config';
import type { Definitions } from '../../starlight/user-components/rehype-file-tree';
import { getSetiIconName } from './seti';

// This matches the default precision used by the SVGO default preset.
const pathDecimalPrecision = 3;

/** Extract SVG paths from the Seti UI icon font from a list of icon names matching font glyphs. */
export function getIconSvgPaths(
	icons: string[],
	definitions: Definitions,
	fontBuffer: ArrayBuffer
) {
	const iconSvgs: Record<string, string> = {};

	let font: Font;

	try {
		font = opentype.parse(fontBuffer);
	} catch (error) {
		throw new Error('Failed to parse icon SVGs.', { cause: error });
	}

	for (const icon of icons) {
		let glyph: Glyph;

		try {
			// Find the glyph matching the icon name.
			glyph = font.nameToGlyph(icon);
		} catch (error) {
			// If the glyph is not found, this means that multiple icons share the same glyph and we have
			// a mapping for such case.
			const alias = getFontGlyphAlias(icon);

			// When an alias is found, we update the definitions to use the alias instead of the original
			// icon name and continue to the next icon as there is no need to extract an SVG.
			updateDefinitionsWithAlias(definitions, icon, alias);
			continue;
		}

		// We need to compute various metrics to ensure the icon properly fits the viewBox size.
		const { fontSize, offsetX, offsetY } = getComputedFontSizeToFit(
			glyph,
			starlight.iconViewBoxSize
		);
		const path = glyph.getPath(offsetX, offsetY, fontSize);
		const iconName = getSetiIconName(icon);
		iconSvgs[iconName] = path.toSVG(pathDecimalPrecision);
	}

	return iconSvgs;
}

/**
 * Compute the font size and offsets to fit a glyph in a viewBox of a given size.
 * We first try to fit the glyph in the viewBox and increase the font size until it no longer fits
 * to get the best fit possible.
 */
function getComputedFontSizeToFit(glyph: opentype.Glyph, size: number, fontSize = size) {
	const { width, height } = getGlyphPathDimensions(glyph, fontSize);

	// If one of the glyph path dimensions is greater than the viewBox size, we can stop here and return
	// the previous font size and offsets.
	if (width > size || height > size) {
		const fittingFontSize = fontSize - 1;
		const { x1, y1, width, height } = getGlyphPathDimensions(glyph, fittingFontSize);

		return {
			fontSize: fittingFontSize,
			// We need to ensure that the glyph is centered in the viewBox.
			offsetX: (x1 + (width - size) / 2) * -1,
			offsetY: (y1 + (height - size) / 2) * -1,
		};
	}

	// If the glyph path dimensions are smaller than the viewBox size, we can increase the font size
	// and try again.
	return getComputedFontSizeToFit(glyph, size, fontSize + 1);
}

/** Returns the bounding box, width, and height of a glyph path. */
function getGlyphPathDimensions(glyph: opentype.Glyph, fontSize: number) {
	const path = glyph.getPath(0, 0, fontSize);
	const boundingBox = path.getBoundingBox();

	return {
		...boundingBox,
		width: boundingBox.x2 - boundingBox.x1,
		height: boundingBox.y2 - boundingBox.y1,
	};
}

/**
 * Return the alias of a font glyph when multiple icons share the same glyph or throw an error if
 * no mapping is found.
 */
function getFontGlyphAlias(icon: string): string {
	const alias = seti.aliases[icon as keyof typeof seti.aliases];

	if (!alias) {
		throw new Error(
			`Failed to find a glyph for the icon '${icon}'. This usually means that the icon is sharing a glyph with another icon and such association must be defined in the 'seti.aliases' configuration.`
		);
	}

	return alias;
}

/** Update the definitions to use an alias instead of a specific icon name. */
function updateDefinitionsWithAlias(definitions: Definitions, icon: string, alias: string) {
	const prefixedIcon = getSetiIconName(icon);
	const prefixedAlias = getSetiIconName(alias);

	updateDefinitionsRecordWithAlias(definitions.files, prefixedIcon, prefixedAlias);
	updateDefinitionsRecordWithAlias(definitions.extensions, prefixedIcon, prefixedAlias);
	updateDefinitionsRecordWithAlias(definitions.partials, prefixedIcon, prefixedAlias);
}

/** Update a definitions record to use an alias instead of a specific icon name. */
function updateDefinitionsRecordWithAlias(
	record: Record<string, string>,
	icon: string,
	alias: string
) {
	for (const [key, value] of Object.entries(record)) {
		if (value === icon) {
			record[key] = alias;
		}
	}
}
