/**
 * This is a modified version of Astro's error map.
 * source: https://github.com/withastro/astro/blob/main/packages/astro/src/content/error-map.ts
 */

import { AstroError } from 'astro/errors';
import type { z } from 'astro:content';

type TypeOrLiteralErrByPathEntry = {
	code: 'invalid_type' | 'invalid_literal';
	received: unknown;
	expected: unknown[];
};

/**
 * Parse data with a Zod schema and throw a nicely formatted error if it is invalid.
 *
 * @param schema The Zod schema to use to parse the input.
 * @param input Input data that should match the schema.
 * @param message Error message preamble to use if the input fails to parse.
 * @returns Validated data parsed by Zod.
 */
export function parseWithFriendlyErrors<T extends z.Schema>(
	schema: T,
	input: z.input<T>,
	message: string
): z.output<T> {
	const parsedConfig = schema.safeParse(input, { errorMap });
	if (!parsedConfig.success) {
		throw new AstroError(message, parsedConfig.error.issues.map((i) => i.message).join('\n'));
	}
	return parsedConfig.data;
}

const errorMap: z.ZodErrorMap = (baseError, ctx) => {
	const baseErrorPath = flattenErrorPath(baseError.path);
	if (baseError.code === 'invalid_union') {
		// Optimization: Combine type and literal errors for keys that are common across ALL union types
		// Ex. a union between `{ key: z.literal('tutorial') }` and `{ key: z.literal('blog') }` will
		// raise a single error when `key` does not match:
		// > Did not match union.
		// > key: Expected `'tutorial' | 'blog'`, received 'foo'
		let typeOrLiteralErrByPath: Map<string, TypeOrLiteralErrByPathEntry> = new Map();
		for (const unionError of baseError.unionErrors.map((e) => e.errors).flat()) {
			if (unionError.code === 'invalid_type' || unionError.code === 'invalid_literal') {
				const flattenedErrorPath = flattenErrorPath(unionError.path);
				if (typeOrLiteralErrByPath.has(flattenedErrorPath)) {
					typeOrLiteralErrByPath.get(flattenedErrorPath)!.expected.push(unionError.expected);
				} else {
					typeOrLiteralErrByPath.set(flattenedErrorPath, {
						code: unionError.code,
						received: (unionError as any).received,
						expected: [unionError.expected],
					});
				}
			}
		}
		const messages: string[] = [prefix(baseErrorPath, 'Did not match union.')];
		const details: string[] = [...typeOrLiteralErrByPath.entries()]
			// If type or literal error isn't common to ALL union types,
			// filter it out. Can lead to confusing noise.
			.filter(([, error]) => error.expected.length === baseError.unionErrors.length)
			.map(([key, error]) =>
				key === baseErrorPath
					? // Avoid printing the key again if it's a base error
					  `> ${getTypeOrLiteralMsg(error)}`
					: `> ${prefix(key, getTypeOrLiteralMsg(error))}`
			);

		if (details.length === 0) {
			const expectedShapes: string[] = [];
			for (const unionError of baseError.unionErrors) {
				const expectedShape: string[] = [];
				for (const issue of unionError.issues) {
					// If the issue is a nested union error, show the associated error message instead of the
					// base error message.
					if (issue.code === 'invalid_union') {
						return errorMap(issue, ctx);
					}
					const relativePath = flattenErrorPath(issue.path)
						.replace(baseErrorPath, '')
						.replace(leadingPeriod, '');
					if ('expected' in issue && typeof issue.expected === 'string') {
						expectedShape.push(
							relativePath ? `${relativePath}: ${issue.expected}` : issue.expected
						);
					} else {
						expectedShape.push(relativePath);
					}
				}
				expectedShapes.push(`{ ${expectedShape.join('; ')} }`);
			}
			if (expectedShapes.length) {
				details.push('> Expected type `' + expectedShapes.join(' | ') + '`');
				details.push('> Received `' + stringify(ctx.data) + '`');
			}
		}

		return {
			message: messages.concat(details).join('\n'),
		};
	} else if (baseError.code === 'invalid_literal' || baseError.code === 'invalid_type') {
		return {
			message: prefix(
				baseErrorPath,
				getTypeOrLiteralMsg({
					code: baseError.code,
					received: (baseError as any).received,
					expected: [baseError.expected],
				})
			),
		};
	} else if (baseError.message) {
		return { message: prefix(baseErrorPath, baseError.message) };
	} else {
		return { message: prefix(baseErrorPath, ctx.defaultError) };
	}
};

const getTypeOrLiteralMsg = (error: TypeOrLiteralErrByPathEntry): string => {
	if (error.received === 'undefined') return 'Required';
	const expectedDeduped = new Set(error.expected);
	switch (error.code) {
		case 'invalid_type':
			return `Expected type \`${unionExpectedVals(expectedDeduped)}\`, received \`${stringify(
				error.received
			)}\``;
		case 'invalid_literal':
			return `Expected \`${unionExpectedVals(expectedDeduped)}\`, received \`${stringify(
				error.received
			)}\``;
	}
};

const prefix = (key: string, msg: string) => (key.length ? `**${key}**: ${msg}` : msg);

const unionExpectedVals = (expectedVals: Set<unknown>) =>
	[...expectedVals].map((expectedVal) => stringify(expectedVal)).join(' | ');

const flattenErrorPath = (errorPath: (string | number)[]) => errorPath.join('.');

/** `JSON.stringify()` a value with spaces around object/array entries. */
const stringify = (val: unknown) =>
	JSON.stringify(val, null, 1).split(newlinePlusWhitespace).join(' ');
const newlinePlusWhitespace = /\n\s*/;
const leadingPeriod = /^\./;
