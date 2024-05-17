import type { AstroIntegrationLogger } from 'astro';
import { type StarlightPluginContext } from '../utils/plugins';

export function createTestPluginContext(): StarlightPluginContext {
	return {
		command: 'dev',
		// @ts-expect-error - we don't provide a full Astro config but only what is needed for the
		// plugins to run.
		config: { integrations: [] },
		isRestart: false,
		logger: new TestAstroIntegrationLogger(),
	};
}

class TestAstroIntegrationLogger {
	options = {} as AstroIntegrationLogger['options'];
	constructor(public label = 'test-integration-logger') {}
	fork = (label: string) => new TestAstroIntegrationLogger(label);
	info = () => undefined;
	warn = () => undefined;
	error = () => undefined;
	debug = () => undefined;
}
