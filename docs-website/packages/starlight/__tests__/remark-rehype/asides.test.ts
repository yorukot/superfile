import { createMarkdownProcessor } from '@astrojs/markdown-remark';
import { describe, expect, test } from 'vitest';
import { starlightAsides } from '../../integrations/asides';
import { createTranslationSystemFromFs } from '../../utils/translations-fs';
import { StarlightConfigSchema, type StarlightUserConfig } from '../../utils/user-config';

const starlightConfig = StarlightConfigSchema.parse({
	title: 'Asides Tests',
	locales: { en: { label: 'English' }, fr: { label: 'French' } },
	defaultLocale: 'en',
} satisfies StarlightUserConfig);

const useTranslations = createTranslationSystemFromFs(
	starlightConfig,
	// Using non-existent `_src/` to ignore custom files in this test fixture.
	{ srcDir: new URL('./_src/', import.meta.url) }
);

const processor = await createMarkdownProcessor({
	remarkPlugins: [
		...starlightAsides({
			starlightConfig,
			astroConfig: { root: new URL(import.meta.url), srcDir: new URL('./_src/', import.meta.url) },
			useTranslations,
		}),
	],
});

test('generates <aside>', async () => {
	const res = await processor.render(`
:::note
Some text
:::
`);
	expect(res.code).toMatchFileSnapshot('./snapshots/generates-aside.html');
});

describe('default labels', () => {
	test.each([
		['note', 'Note'],
		['tip', 'Tip'],
		['caution', 'Caution'],
		['danger', 'Danger'],
	])('%s has label %s', async (type, label) => {
		const res = await processor.render(`
:::${type}
Some text
:::
`);
		expect(res.code).includes(`aria-label="${label}"`);
		expect(res.code).includes(`</svg>${label}</p>`);
	});
});

describe('custom labels', () => {
	test.each(['note', 'tip', 'caution', 'danger'])('%s with custom label', async (type) => {
		const label = 'Custom Label';
		const res = await processor.render(`
:::${type}[${label}]
Some text
:::
  `);
		expect(res.code).includes(`aria-label="${label}"`);
		expect(res.code).includes(`</svg>${label}</p>`);
	});
});

test('ignores unknown directive variants', async () => {
	const res = await processor.render(`
:::unknown
Some text
:::
`);
	expect(res.code).toMatchInlineSnapshot('"<div><p>Some text</p></div>"');
});

test('handles complex children', async () => {
	const res = await processor.render(`
:::note
Paragraph [link](/href/).

![alt](/img.jpg)

<details>
<summary>See more</summary>

More.

</details>
:::
`);
	expect(res.code).toMatchFileSnapshot('./snapshots/handles-complex-children.html');
});

test('nested asides', async () => {
	const res = await processor.render(`
::::note
Note contents.

:::tip
Nested tip.
:::

::::
`);
	expect(res.code).toMatchFileSnapshot('./snapshots/nested-asides.html');
});

test('nested asides with custom titles', async () => {
	const res = await processor.render(`
:::::caution[Caution with a custom title]
Nested caution.

::::note
Nested note.

:::tip[Tip with a custom title]
Nested tip.
:::

::::

:::::
`);
	const labels = [...res.code.matchAll(/aria-label="(?<label>[^"]+)"/g)].map(
		(match) => match.groups?.label
	);
	expect(labels).toMatchInlineSnapshot(`
		[
		  "Caution with a custom title",
		  "Note",
		  "Tip with a custom title",
		]
	`);
	expect(res.code).toMatchFileSnapshot('./snapshots/nested-asides-custom-titles.html');
});

describe('translated labels in French', () => {
	test.each([
		['note', 'Note'],
		['tip', 'Astuce'],
		['caution', 'Attention'],
		['danger', 'Danger'],
	])('%s has label %s', async (type, label) => {
		const res = await processor.render(
			`
:::${type}
Some text
:::
`,
			// @ts-expect-error fileURL is part of MarkdownProcessor's options
			{ fileURL: new URL('./_src/content/docs/fr/index.md', import.meta.url) }
		);
		expect(res.code).includes(`aria-label="${label}"`);
		expect(res.code).includes(`</svg>${label}</p>`);
	});
});

test('runs without locales config', async () => {
	const processor = await createMarkdownProcessor({
		remarkPlugins: [
			...starlightAsides({
				starlightConfig: { locales: undefined },
				astroConfig: {
					root: new URL(import.meta.url),
					srcDir: new URL('./_src/', import.meta.url),
				},
				useTranslations,
			}),
		],
	});
	const res = await processor.render(':::note\nTest\n::');
	expect(res.code.includes('aria-label=Note"'));
});

test('tranforms back unhandled text directives', async () => {
	const res = await processor.render(
		`This is a:test of a sentence with a text:name[content]{key=val} directive.`
	);
	expect(res.code).toMatchInlineSnapshot(`
		"<p>This is a:test
		 of a sentence with a text:name[content]{key="val"}
		 directive.</p>"
	`);
});

test('tranforms back unhandled leaf directives', async () => {
	const res = await processor.render(`::video[Title]{v=xxxxxxxxxxx}`);
	expect(res.code).toMatchInlineSnapshot(`
		"<p>::video[Title]{v="xxxxxxxxxxx"}
		</p>"
	`);
});
