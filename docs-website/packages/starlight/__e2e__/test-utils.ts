import { fileURLToPath } from 'node:url';
import { test as baseTest, type Page } from '@playwright/test';
import { build, preview } from 'astro';

export { expect, type Locator } from '@playwright/test';

// Setup a test environment that will build and start a preview server for a given fixture path and
// provide a Starlight Playwright fixture accessible from within all tests.
export async function testFactory(fixturePath: string) {
	let previewServer: PreviewServer | undefined;

	const test = baseTest.extend<{ starlight: StarlightPage }>({
		starlight: async ({ page }, use) => {
			if (!previewServer) {
				throw new Error('Could not find a preview server to run tests against.');
			}

			await use(new StarlightPage(previewServer, page));
		},
	});

	test.beforeAll(async () => {
		const root = fileURLToPath(new URL(fixturePath, import.meta.url));
		await build({ logLevel: 'error', root });
		previewServer = await preview({ logLevel: 'error', root });
	});

	test.afterAll(async () => {
		await previewServer?.stop();
	});

	return test;
}

// A Playwright test fixture accessible from within all tests.
class StarlightPage {
	constructor(
		private readonly previewServer: PreviewServer,
		private readonly page: Page
	) {}

	// Navigate to a URL relative to the server used during a test run and return the resource response.
	goto(url: string) {
		return this.page.goto(this.resolveUrl(url));
	}

	// Resolve a URL relative to the server used during a test run.
	resolveUrl(url: string) {
		return `http://localhost:${this.previewServer.port}${url.replace(/^\/?/, '/')}`;
	}
}

type PreviewServer = Awaited<ReturnType<typeof preview>>;
