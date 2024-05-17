import { expect, test } from 'vitest';
import { processPanels, TabItemTagname } from '../../user-components/rehype-tabs';

const TabItem = ({ label, slot, icon }: { label: string; slot: string; icon?: string }) => {
	const iconAttr = icon ? ` data-icon="${icon}"` : '';
	return `<${TabItemTagname} data-label="${label}"${iconAttr}>${slot}</${TabItemTagname}>`;
};

/** Get an array of HTML strings, one for each `<section>` created by rehype-tabs for each tab item. */
const extractSections = (html: string) =>
	[...html.matchAll(/<section.*?<\/section>/g)].map(([section]) => section);

test('empty component returns no html or panels', () => {
	const { panels, html } = processPanels('');
	expect(html).toEqual('');
	expect(panels).toEqual([]);
});

test('non-tab-item content is passed unchanged', () => {
	const input = '<p>Random paragraph</p>';
	const { panels, html } = processPanels(input);
	expect(html).toEqual(input);
	expect(panels).toEqual([]);
});

test('tab items are processed', () => {
	const label = 'Test';
	const slot = '<p>Random paragraph</p>';
	const input = TabItem({ label, slot });
	const { panels, html } = processPanels(input);

	expect(html).toMatchInlineSnapshot(
		`"<section id="tab-panel-0" aria-labelledby="tab-0" role="tabpanel" tabindex="0"><p>Random paragraph</p></section>"`
	);
	expect(panels).toHaveLength(1);
	expect(panels?.[0]?.label).toBe(label);
	expect(panels?.[0]?.panelId).toMatchInlineSnapshot('"tab-panel-0"');
	expect(panels?.[0]?.tabId).toMatchInlineSnapshot('"tab-0"');
	expect(panels?.[0]?.icon).not.toBeDefined();
});

test('only first item is not hidden', () => {
	const labels = ['One', 'Two', 'Three'];
	const input = labels.map((label) => TabItem({ label, slot: `<div>${label}</div>` })).join('');
	const { panels, html } = processPanels(input);

	expect(panels).toHaveLength(3);
	expect(html).toMatchInlineSnapshot(
		`"<section id="tab-panel-1" aria-labelledby="tab-1" role="tabpanel" tabindex="0"><div>One</div></section><section id="tab-panel-2" aria-labelledby="tab-2" role="tabpanel" tabindex="0" hidden><div>Two</div></section><section id="tab-panel-3" aria-labelledby="tab-3" role="tabpanel" tabindex="0" hidden><div>Three</div></section>"`
	);
	const sections = extractSections(html);
	expect(sections).toMatchInlineSnapshot(`
		[
		  "<section id="tab-panel-1" aria-labelledby="tab-1" role="tabpanel" tabindex="0"><div>One</div></section>",
		  "<section id="tab-panel-2" aria-labelledby="tab-2" role="tabpanel" tabindex="0" hidden><div>Two</div></section>",
		  "<section id="tab-panel-3" aria-labelledby="tab-3" role="tabpanel" tabindex="0" hidden><div>Three</div></section>",
		]
	`);
	expect(sections.map((section) => section.includes('hidden'))).toEqual([false, true, true]);
});

test('applies incrementing ID and aria-labelledby to each tab item', () => {
	const labels = ['One', 'Two', 'Three'];
	const input = labels.map((label) => TabItem({ label, slot: `<div>${label}</div>` })).join('');
	const { panels, html } = processPanels(input);

	// IDs are incremented globally to ensure they are unique, so we need to extract from the panel data.
	const firstTabIdMatches = panels?.[0]?.tabId.match(/^tab-(\d)+$/);
	const firstTabId = parseInt(firstTabIdMatches![1]!, 10);

	extractSections(html).forEach((section, index) => {
		expect(section).includes(`id="tab-panel-${firstTabId + index}"`);
		expect(section).includes(`aria-labelledby="tab-${firstTabId + index}"`);
	});
});

test('applies tabindex="0" to tab items without focusable content', () => {
	const input = [
		TabItem({ label: 'Focusable', slot: `<div><a href="/home/">Home</a></div>` }),
		TabItem({ label: 'Not Focusable', slot: `<div>Plain text</div>` }),
		TabItem({
			label: 'Focusable Nested',
			slot: `<div><p><span><input type="text"></span></p></div>`,
		}),
	].join('');
	const { html } = processPanels(input);
	expect(html).toMatchInlineSnapshot(
		`"<section id="tab-panel-7" aria-labelledby="tab-7" role="tabpanel"><div><a href="/home/">Home</a></div></section><section id="tab-panel-8" aria-labelledby="tab-8" role="tabpanel" tabindex="0" hidden><div>Plain text</div></section><section id="tab-panel-9" aria-labelledby="tab-9" role="tabpanel" hidden><div><p><span><input type="text"></span></p></div></section>"`
	);
	const sections = extractSections(html);
	expect(sections[0]).not.includes('tabindex="0"');
	expect(sections[1]).includes('tabindex="0"');
	expect(sections[2]).not.includes('tabindex="0"');
});

test('processes a tab item icon', () => {
	const icon = 'star';
	const input = TabItem({ label: 'Test', slot: '<p>Random paragraph</p>', icon });
	const { panels, html } = processPanels(input);

	expect(html).toMatchInlineSnapshot(
		`"<section id="tab-panel-10" aria-labelledby="tab-10" role="tabpanel" tabindex="0"><p>Random paragraph</p></section>"`
	);
	expect(panels).toHaveLength(1);
	expect(panels?.[0]?.icon).toBe(icon);
});
