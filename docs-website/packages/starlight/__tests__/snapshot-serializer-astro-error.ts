import { AstroError } from 'astro/errors';
import type { SnapshotSerializer } from 'vitest';

export default {
	/** Check if a value should be handled by this serializer, i.e. if it is an `AstroError`. */
	test(val) {
		return !!val && AstroError.is(val);
	},
	/** Customize serialization of Astro errors to include the `hint`. Vitest only uses `message` by default. */
	serialize({ name, message, hint }: AstroError, config, indentation, depth, refs, printer) {
		const prettyError = `[${name}]:\n${indent(message)}\nHint:\n${indent(hint)}`;
		return printer(prettyError, config, indentation, depth, refs);
	},
} satisfies SnapshotSerializer;

/** Indent each line in `string` with a given character. */
function indent(string = '', indentation = '\t') {
	return string
		.split('\n')
		.map((line) => indentation + line)
		.join('\n');
}
