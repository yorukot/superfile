import { describe, expect, test } from 'vitest';
import { processFileTree } from '../../user-components/rehype-file-tree';
import { Icons } from '../../components/Icons';

describe('validation', () => {
	test('throws an error with no content', () => {
		expect(() => processTestFileTree('')).toThrowErrorMatchingInlineSnapshot(
			`
			"[AstroUserError]:
				The \`<FileTree>\` component expects its content to be a single unordered list but found no child elements.
			Hint:
				To learn more about the \`<FileTree>\` component, see https://starlight.astro.build/guides/components/#file-tree"
		`
		);
	});

	test('throws an error with multiple root elements', () => {
		expect(() =>
			processTestFileTree('<p>test</p><ul><li>file</li></ul>')
		).toThrowErrorMatchingInlineSnapshot(
			`
			"[AstroUserError]:
				The \`<FileTree>\` component expects its content to be a single unordered list but found multiple child elements: \`<p>\` - \`<ul>\`.
			Hint:
				To learn more about the \`<FileTree>\` component, see https://starlight.astro.build/guides/components/#file-tree"
		`
		);
	});

	test('throws an error with no root ordered list', () => {
		expect(() => processTestFileTree('<ol><li>file</li></ol>')).toThrowErrorMatchingInlineSnapshot(
			`
			"[AstroUserError]:
				The \`<FileTree>\` component expects its content to be an unordered list but found the following element: \`<ol>\`.
			Hint:
				To learn more about the \`<FileTree>\` component, see https://starlight.astro.build/guides/components/#file-tree"
		`
		);
	});

	test('throws an error with no list item', () => {
		expect(() => processTestFileTree('<ul></ul>')).toThrowErrorMatchingInlineSnapshot(
			`
			"[AstroUserError]:
				The \`<FileTree>\` component expects its content to be an unordered list with at least one list item.
			Hint:
				To learn more about the \`<FileTree>\` component, see https://starlight.astro.build/guides/components/#file-tree"
		`
		);
	});
});

describe('processor', () => {
	test('processes a basic tree', () => {
		const html = processTestFileTree(`<ul>
  <li>root_file</li>
  <li>root_directory/
    <ul>
      <li>nested_file</li>
    </ul>
  <li>
</ul>`);

		expect(extractFileTree(html)).toMatchFileSnapshot('./snapshots/file-tree-basic.html');
	});

	test('does not add a comment node with no comments', () => {
		const html = processTestFileTree(`<ul><li>file</li></ul>`);

		expect(extractFileTree(html)).not.toContain('<span class="comment">');
	});

	test('processes text comments following the file name', () => {
		const html = processTestFileTree(`<ul><li>file this is a comment</li></ul>`);

		expect(extractFileTree(html)).toMatchFileSnapshot('./snapshots/file-tree-comment-text.html');
	});

	test('processes comment nodes', () => {
		const html = processTestFileTree(
			`<ul><li>file this is an <strong>important</strong> comment</li></ul>`
		);

		expect(extractFileTree(html)).toMatchFileSnapshot('./snapshots/file-tree-comment-nodes.html');
	});

	test('identifies directory with either a file name ending with a slash or a nested list', () => {
		const html = processTestFileTree(`<ul>
  <li>directory/</li>
  <li>another_directory
    <ul>
      <li>file</li>
    </ul>
  </li>
</ul>`);

		expect(extractFileTree(html).match(/<li class="directory">/g)).toHaveLength(2);
	});

	test('identifies placeholder with either 3 dots or an ellipsis', () => {
		const html = processTestFileTree(`<ul><li>...</li><li>…</li></ul>`);

		expect(extractFileTree(html).match(/<li class="[\w\s]*empty[\w\s]*">/g)).toHaveLength(2);
	});

	test('adds a placeholder to empty directories', () => {
		const html = processTestFileTree(`<ul><li>directory/</li></ul>`);

		expect(extractFileTree(html)).toContain('<li class="file empty">');
	});

	test('identifies highlighted with a strong tag', () => {
		const html = processTestFileTree(`<ul><li><strong>file</strong></li></ul>`);

		expect(extractFileTree(html)).toContain('<span class="highlight">');
	});
});

describe('icons', () => {
	test('adds a folder icon to directories with a screen-reader only label', () => {
		const html = processTestFileTree(`<ul><li>directory/</li></ul>`);

		expectHtmlToIncludeIcon(html, Icons['seti:folder']);
		expect(extractFileTree(html)).toContain('<span class="sr-only">Directory</span>');
	});

	test('adds a default file icon to unknown files', () => {
		const html = processTestFileTree(`<ul><li>test_file</li></ul>`);

		expectHtmlToIncludeIcon(html, Icons['seti:default']);
	});

	test('adds an icon to known files', () => {
		const html = processTestFileTree(`<ul><li>README.md</li></ul>`);

		expectHtmlToIncludeIcon(html, Icons['seti:info']);
	});

	test('adds an icon to known file extensions', () => {
		const html = processTestFileTree(`<ul><li>test.astro</li></ul>`);

		expectHtmlToIncludeIcon(html, Icons.astro);
	});

	test('adds an icon to known file partials', () => {
		const html = processTestFileTree(`<ul><li>TODO</li></ul>`);

		expectHtmlToIncludeIcon(html, Icons['seti:todo']);
	});

	test('does not add a special icon to file based on the last letter of the file name', () => {
		// The last letter of the file name is "c" and should not be matched to the icon for C files.
		const html = processTestFileTree(`<ul><li>testc</li></ul>`);

		expectHtmlToIncludeIcon(html, Icons['seti:default']);
	});
});

/** Calls the file tree processor with the given HTML and a default label. */
function processTestFileTree(html: string) {
	return processFileTree(html, 'Directory');
}

/** Extracts the file tree from the given HTML and optionally strips out the icon SVGs. */
function extractFileTree(html: string, stripIcons = true) {
	let tree = html.match(/<ul>.*<\/ul>/s)?.[0] ?? '';

	if (stripIcons) {
		tree = tree.replace(/<svg.*?<\/svg>/g, '<svg></svg>');
	}

	return tree;
}

function expectHtmlToIncludeIcon(html: string, icon: (typeof Icons)[keyof typeof Icons]) {
	return expect(extractFileTree(html, false)).toContain(icon.replace('/>', '>'));
}
