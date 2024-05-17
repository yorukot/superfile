import { expect, testFactory, type Locator } from './test-utils';

const test = await testFactory('./fixtures/basics/');

test('syncs tabs with a click event', async ({ page, starlight }) => {
	await starlight.goto('/tabs');

	const tabs = page.locator('starlight-tabs');
	const pkgTabsA = tabs.nth(0);
	const pkgTabsB = tabs.nth(2);

	// Select the pnpm tab in the first set of synced tabs.
	await pkgTabsA.getByRole('tab').filter({ hasText: 'pnpm' }).click();

	await expectSelectedTab(pkgTabsA, 'pnpm', 'pnpm command');
	await expectSelectedTab(pkgTabsB, 'pnpm', 'another pnpm command');

	// Select the yarn tab in the second set of synced tabs.
	await pkgTabsB.getByRole('tab').filter({ hasText: 'yarn' }).click();

	await expectSelectedTab(pkgTabsB, 'yarn', 'another yarn command');
	await expectSelectedTab(pkgTabsA, 'yarn', 'yarn command');
});

test('syncs tabs with a keyboard event', async ({ page, starlight }) => {
	await starlight.goto('/tabs');

	const tabs = page.locator('starlight-tabs');
	const pkgTabsA = tabs.nth(0);
	const pkgTabsB = tabs.nth(2);

	// Select the pnpm tab in the first set of synced tabs with the keyboard.
	await pkgTabsA.getByRole('tab', { selected: true }).press('ArrowRight');

	await expectSelectedTab(pkgTabsA, 'pnpm', 'pnpm command');
	await expectSelectedTab(pkgTabsB, 'pnpm', 'another pnpm command');

	// Select back the npm tab in the second set of synced tabs with the keyboard.
	const selectedTabB = pkgTabsB.getByRole('tab', { selected: true });
	await selectedTabB.press('ArrowRight');
	await selectedTabB.press('ArrowLeft');
	await selectedTabB.press('ArrowLeft');

	await expectSelectedTab(pkgTabsA, 'npm', 'npm command');
	await expectSelectedTab(pkgTabsB, 'npm', 'another npm command');
});

test('syncs only tabs using the same sync key', async ({ page, starlight }) => {
	await starlight.goto('/tabs');

	const tabs = page.locator('starlight-tabs');
	const pkgTabsA = tabs.nth(0);
	const unsyncedTabs = tabs.nth(1);
	const styleTabs = tabs.nth(3);

	// Select the pnpm tab in the set of tabs synced with the 'pkg' key.
	await pkgTabsA.getByRole('tab').filter({ hasText: 'pnpm' }).click();

	await expectSelectedTab(unsyncedTabs, 'one', 'tab 1');
	await expectSelectedTab(styleTabs, 'css', 'css code');
});

test('supports synced tabs with different tab items', async ({ page, starlight }) => {
	await starlight.goto('/tabs');

	const tabs = page.locator('starlight-tabs');
	const pkgTabsA = tabs.nth(0);
	const pkgTabsB = tabs.nth(2); // This set contains an extra tab item.

	// Select the bun tab in the second set of synced tabs.
	await pkgTabsB.getByRole('tab').filter({ hasText: 'bun' }).click();

	await expectSelectedTab(pkgTabsA, 'npm', 'npm command');
	await expectSelectedTab(pkgTabsB, 'bun', 'another bun command');
});

test('persists the focus when syncing tabs', async ({ page, starlight }) => {
	await starlight.goto('/tabs');

	const pkgTabsA = page.locator('starlight-tabs').nth(0);

	// Focus the selected tab in the set of tabs synced with the 'pkg' key.
	await pkgTabsA.getByRole('tab', { selected: true }).focus();
	// Select the pnpm tab in the set of tabs synced with the 'pkg' key using the keyboard.
	await page.keyboard.press('ArrowRight');

	expect(
		await pkgTabsA
			.getByRole('tab', { selected: true })
			.evaluate((node) => document.activeElement === node)
	).toBe(true);
});

test('preserves tabs position when alternating between tabs with different content heights', async ({
	page,
	starlight,
}) => {
	await starlight.goto('/tabs-variable-height');

	const tabs = page.locator('starlight-tabs').nth(1);
	const selectedTab = tabs.getByRole('tab', { selected: true });

	// Scroll to the second set of synced tabs and focus the selected tab.
	await tabs.scrollIntoViewIfNeeded();
	await selectedTab.focus();

	// Get the bounding box of the tabs.
	const initialBoundingBox = await tabs.boundingBox();

	// Select the second tab which has a different height.
	await selectedTab.press('ArrowRight');

	// Ensure the tabs vertical position is exactly the same after selecting the second tab.
	// Note that a small difference could be the result of the base line-height having a fractional part which can cause a
	// sub-pixel difference in some browsers like Chrome or Firefox.
	expect((await tabs.boundingBox())?.y).toBe(initialBoundingBox?.y);
});

test('syncs tabs with the same sync key if they do not consistenly use icons', async ({
	page,
	starlight,
}) => {
	await starlight.goto('/tabs');

	const tabs = page.locator('starlight-tabs');
	const pkgTabsA = tabs.nth(0); // This set does not use icons for tab items.
	const pkgTabsB = tabs.nth(4); // This set uses icons for tab items.

	// Select the pnpm tab in the first set of synced tabs.
	await pkgTabsA.getByRole('tab').filter({ hasText: 'pnpm' }).click();

	await expectSelectedTab(pkgTabsA, 'pnpm', 'pnpm command');
	await expectSelectedTab(pkgTabsB, 'pnpm', 'another pnpm command');

	// Select the yarn tab in the second set of synced tabs.
	await pkgTabsB.getByRole('tab').filter({ hasText: 'yarn' }).click();

	await expectSelectedTab(pkgTabsB, 'yarn', 'another yarn command');
	await expectSelectedTab(pkgTabsA, 'yarn', 'yarn command');
});

async function expectSelectedTab(tabs: Locator, label: string, panel: string) {
	expect((await tabs.getByRole('tab', { selected: true }).textContent())?.trim()).toBe(label);
	expect((await tabs.getByRole('tabpanel').textContent())?.trim()).toBe(panel);
}
