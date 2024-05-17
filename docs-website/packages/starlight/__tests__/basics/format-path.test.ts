import { describe, expect, test } from 'vitest';
import { createPathFormatter } from '../../utils/createPathFormatter';

type FormatPathOptions = Parameters<typeof createPathFormatter>[0];
const formatPath = (href: string, opts: FormatPathOptions) => createPathFormatter(opts)(href);

describe.each<{ options: FormatPathOptions; tests: Array<{ path: string; expected: string }> }>([
	{
		options: { format: 'file', trailingSlash: 'ignore' },
		tests: [
			// index page
			{ path: '/', expected: '/index.html' },
			// with trailing slash
			{ path: '/reference/configuration/', expected: '/reference/configuration.html' },
			// without trailing slash
			{ path: '/api/v1/users', expected: '/api/v1/users.html' },
			// with file extension
			{ path: '/guides/components.html', expected: '/guides/components.html' },
			// with file extension and trailing slash
			{ path: '/guides/components.html/', expected: '/guides/components.html' },
		],
	},
	{
		options: { format: 'file', trailingSlash: 'always' },
		tests: [
			// index page
			{ path: '/', expected: '/index.html' },
			// with trailing slash
			{ path: '/reference/configuration/', expected: '/reference/configuration.html' },
			// without trailing slash
			{ path: '/api/v1/users', expected: '/api/v1/users.html' },
			// with file extension
			{ path: '/guides/components.html', expected: '/guides/components.html' },
			// with file extension and trailing slash
			{ path: '/guides/components.html/', expected: '/guides/components.html' },
		],
	},
	{
		options: { format: 'file', trailingSlash: 'never' },
		tests: [
			// index page
			{ path: '/', expected: '/index.html' },
			// with trailing slash
			{ path: '/reference/configuration/', expected: '/reference/configuration.html' },
			// without trailing slash
			{ path: '/api/v1/users', expected: '/api/v1/users.html' },
			// with file extension
			{ path: '/guides/components.html', expected: '/guides/components.html' },
			// with file extension and trailing slash
			{ path: '/guides/components.html/', expected: '/guides/components.html' },
		],
	},
	{
		options: { format: 'directory', trailingSlash: 'always' },
		tests: [
			// index page
			{ path: '/', expected: '/' },
			// with trailing slash
			{ path: '/reference/configuration/', expected: '/reference/configuration/' },
			// without trailing slash
			{ path: '/api/v1/users', expected: '/api/v1/users/' },
			// with file extension
			{ path: '/guides/components.html', expected: '/guides/components/' },
			// with file extension and trailing slash
			{ path: '/guides/components.html/', expected: '/guides/components/' },
		],
	},
	{
		options: { format: 'directory', trailingSlash: 'never' },
		tests: [
			// index page
			{ path: '/', expected: '/' },
			// with trailing slash
			{ path: '/reference/configuration/', expected: '/reference/configuration' },
			// without trailing slash
			{ path: '/api/v1/users', expected: '/api/v1/users' },
			// with file extension
			{ path: '/guides/components.html', expected: '/guides/components' },
			// with file extension and trailing slash
			{ path: '/guides/components.html/', expected: '/guides/components' },
		],
	},
	{
		options: { format: 'directory', trailingSlash: 'ignore' },
		tests: [
			// index page
			{ path: '/', expected: '/' },
			// with trailing slash
			{ path: '/reference/configuration/', expected: '/reference/configuration/' },
			// without trailing slash
			{ path: '/api/v1/users', expected: '/api/v1/users' },
			// with file extension
			{ path: '/guides/components.html', expected: '/guides/components' },
			// with file extension and trailing slash
			{ path: '/guides/components.html/', expected: '/guides/components' },
		],
	},
])(
	'formatPath() with { format: $options.format, trailingSlash: $options.trailingSlash }',
	({ options, tests }) => {
		test.each(tests)('returns $expected for $path', ({ path, expected }) => {
			expect(formatPath(path, options)).toBe(expected);
		});
	}
);
