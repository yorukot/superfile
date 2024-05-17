import { expect, test } from 'vitest';
import { processSteps } from '../../user-components/rehype-steps';

test('empty component throws an error', () => {
	expect(() => processSteps('')).toThrowErrorMatchingInlineSnapshot(
		`
		"[AstroUserError]:
			The \`<Steps>\` component expects its content to be a single ordered list (\`<ol>\`) but found no child elements.
		Hint:
			To learn more about the \`<Steps>\` component, see https://starlight.astro.build/guides/components/#steps"
	`
	);
});

test('component with non-element content throws an error', () => {
	expect(() => processSteps('<!-- comment -->Text node')).toThrowErrorMatchingInlineSnapshot(
		`
		"[AstroUserError]:
			The \`<Steps>\` component expects its content to be a single ordered list (\`<ol>\`) but found no child elements.
		Hint:
			To learn more about the \`<Steps>\` component, see https://starlight.astro.build/guides/components/#steps"
	`
	);
});

test('component with non-`<ol>` content throws an error', () => {
	expect(() => processSteps('<p>A paragraph is not an ordered list</p>'))
		.toThrowErrorMatchingInlineSnapshot(`
			"[AstroUserError]:
				The \`<Steps>\` component expects its content to be a single ordered list (\`<ol>\`) but found the following element: \`<p>\`.
			Hint:
				To learn more about the \`<Steps>\` component, see https://starlight.astro.build/guides/components/#steps
				
				Full HTML passed to \`<Steps>\`:
				
				<p>A paragraph is not an ordered list</p>
				"
		`);
});

test('component with multiple children throws an error', () => {
	expect(() =>
		processSteps(
			'<ol><li>List item</li></ol><p>I intended this to be part of the same list item</p><ol><li>Other list item</li></ol>'
		)
	).toThrowErrorMatchingInlineSnapshot(`
		"[AstroUserError]:
			The \`<Steps>\` component expects its content to be a single ordered list (\`<ol>\`) but found multiple child elements: \`<ol>\`, \`<p>\`, \`<ol>\`.
		Hint:
			To learn more about the \`<Steps>\` component, see https://starlight.astro.build/guides/components/#steps
			
			Full HTML passed to \`<Steps>\`:
			
			<ol>
			  <li>List item</li>
			</ol>
			<p>I intended this to be part of the same list item</p>
			<ol>
			  <li>Other list item</li>
			</ol>
			"
	`);
});

test('applies `role="list"` to child list', () => {
	const { html } = processSteps('<ol><li>Step one</li></ol>');
	expect(html).toMatchInlineSnapshot(`"<ol role="list" class="sl-steps"><li>Step one</li></ol>"`);
});

test('does not interfere with other attributes on the child list', () => {
	const { html } = processSteps('<ol start="5"><li>Step one</li></ol>');
	expect(html).toMatchInlineSnapshot(
		`"<ol start="5" role="list" class="sl-steps"><li>Step one</li></ol>"`
	);
});

test('applies `class="sl-list"` to child list', () => {
	const { html } = processSteps('<ol><li>Step one</li></ol>');
	expect(html).toContain('class="sl-steps"');
});

test('applies class name and preserves existing classes on a child list', () => {
	const testClass = 'test class-concat';
	const { html } = processSteps(`<ol class="${testClass}"><li>Step one</li></ol>`);
	expect(html).toContain(`class="${testClass} sl-steps"`);
	expect(html).toMatchInlineSnapshot(
		`"<ol class="test class-concat sl-steps" role="list"><li>Step one</li></ol>"`
	);
});
