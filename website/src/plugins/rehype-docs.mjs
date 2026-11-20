import { SKIP, visit } from "unist-util-visit";
import { element, icon, text } from "./hast.mjs";

// Lucide "copy".
const copyIcon = [
	["rect", { width: 14, height: 14, x: 8, y: 8, rx: 2, ry: 2 }],
	["path", { d: "M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" }]
];

/**
 * Dresses Markdown output in the design system's components:
 * - a table becomes a Table (a box that scrolls sideways on small screens)
 * - a fenced code block becomes a code screen with its language and a copy button,
 *   the same markup as src/components/code.astro
 */
export function rehypeDocs() {
	return tree => {
		visit(tree, "element", (node, index, parent) => {
			if (node.tagName === "table") {
				node.properties.className = ["sf-table"];
				parent.children[index] = element("div", { className: ["sf-table-wrap"] }, [node]);
				return SKIP;
			}

			if (node.tagName === "pre") {
				const language = node.properties.dataLanguage;
				const label = language && language !== "plaintext" ? language : "text";
				// The same markup as src/components/CopyButton.astro.
				const copy = element("button", { className: ["sf-copy"], type: "button", ariaLabel: "Copy code", dataCopy: "" }, [icon(copyIcon), element("span", { dataCopyLabel: "" }, [text("copy")])]);
				const head = element("figcaption", { className: ["sf-code__head"], dataPagefindIgnore: "" }, [element("span", {}, [text(label)]), copy]);
				parent.children[index] = element("figure", { className: ["sf-code"] }, [head, node]);
				return SKIP;
			}
		});
	};
}
