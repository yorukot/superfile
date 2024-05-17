import { defineConfig } from 'vitest/config';

export default defineConfig({
	test: {
		coverage: {
			reportsDirectory: './__coverage__',
			thresholds: {
				autoUpdate: true,
				lines: 94,
				functions: 100,
				branches: 85,
				statements: 94,
			},
		},
	},
});
