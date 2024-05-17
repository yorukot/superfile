import { rehype } from 'rehype';
import { expect, test } from 'vitest';
import { rehypeRtlCodeSupport } from '../../integrations/code-rtl-support';

const processor = rehype().data('settings', { fragment: true }).use(rehypeRtlCodeSupport());

test('applies `dir="auto"` to inline code', async () => {
	const input = `<p>Some text with <code>inline code</code>.</p>`;
	const output = String(await processor.process(input));
	expect(output).not.toEqual(input);
	expect(output).includes('dir="auto"');
	expect(output).toMatchInlineSnapshot(
		`"<p>Some text with <code dir="auto">inline code</code>.</p>"`
	);
});

test('applies `dir="ltr"` to code blocks', async () => {
	const input = `<p>Some text in a paragraph:</p><pre><code>console.log('test')</code></pre>`;
	const output = String(await processor.process(input));
	expect(output).not.toEqual(input);
	expect(output).includes('dir="ltr"');
	expect(output).toMatchInlineSnapshot(
		`"<p>Some text in a paragraph:</p><pre dir="ltr"><code>console.log('test')</code></pre>"`
	);
});
