import { describe, expect, test, vi } from 'vitest';
import { fileWithBase, pathWithBase } from '../../utils/base';

describe('fileWithBase()', () => {
	describe('with no base', () => {
		test('does not prepend anything', () => {
			expect(fileWithBase('/img.svg')).toBe('/img.svg');
		});
		test('adds leading slash if needed', () => {
			expect(fileWithBase('img.svg')).toBe('/img.svg');
		});
	});

	describe('with base', () => {
		test('prepends base', async () => {
			// Reset the modules registry so that re-importing `../../utils/base` re-evaluates the module
			// and re-computes the base. Re-importing the module is necessary because top-level imports
			// cannot be re-evaluated.
			vi.resetModules();
			// Set the base URL.
			vi.stubEnv('BASE_URL', '/base/');
			// Re-import the module to re-evaluate it.
			const { fileWithBase } = await import('../../utils/base');

			expect(fileWithBase('/img.svg')).toBe('/base/img.svg');

			vi.unstubAllEnvs();
			vi.resetModules();
		});
	});
});

describe('pathWithBase()', () => {
	describe('with no base', () => {
		test('does not prepend anything', () => {
			expect(pathWithBase('/path/')).toBe('/path/');
		});
		test('adds leading slash if needed', () => {
			expect(pathWithBase('path')).toBe('/path');
		});
	});

	describe('with base', () => {
		test('prepends base', async () => {
			// See the first test with a base in this file for an explanation of the environment stubbing.
			vi.resetModules();
			vi.stubEnv('BASE_URL', '/base/');
			const { pathWithBase } = await import('../../utils/base');

			expect(pathWithBase('/path/')).toBe('/base/path/');

			vi.unstubAllEnvs();
			vi.resetModules();
		});
	});
});
