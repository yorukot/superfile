import { visit } from "unist-util-visit";
import { icon, text } from "./hast.mjs";

// Lucide lightbulb, triangle-alert, info and octagon-alert.
const icons = {
	tip: [
		["path", { d: "M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5" }],
		["path", { d: "M9 18h6" }],
		["path", { d: "M10 22h4" }]
	],
	caution: [
		["path", { d: "m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3" }],
		["path", { d: "M12 9v4" }],
		["path", { d: "M12 17h.01" }]
	],
	note: [
		["circle", { cx: "12", cy: "12", r: "10" }],
		["path", { d: "M12 16v-4" }],
		["path", { d: "M12 8h.01" }]
	],
	danger: [
		["path", { d: "M12 16h.01" }],
		["path", { d: "M12 8v4" }],
		[
			"path",
			{
				d: "M15.312 2a2 2 0 0 1 1.414.586l4.688 4.688A2 2 0 0 1 22 8.688v6.624a2 2 0 0 1-.586 1.414l-4.688 4.688a2 2 0 0 1-1.414.586H8.688a2 2 0 0 1-1.414-.586l-4.688-4.688A2 2 0 0 1 2 15.312V8.688a2 2 0 0 1 .586-1.414l4.688-4.688A2 2 0 0 1 8.688 2z"
			}
		]
	]
};

/**
 * Turns the :::tip, :::caution, :::note and :::danger blocks of remark-directive into callouts:
 * a Pane whose title is the kind of callout, with its icon. It writes plain HTML, so it works
 * the same in .md and .mdx files.
 */
export function remarkCallouts() {
	return tree => {
		visit(tree, ["containerDirective", "leafDirective"], node => {
			const kind = Object.hasOwn(icons, node.name) ? node.name : "note";
			const title = {
				type: "paragraph",
				children: [],
				data: { hName: "span", hProperties: { className: ["sf-pane__title", "sf-pane__title--icon"], dataPagefindIgnore: "" }, hChildren: [icon(icons[kind]), text(kind)] }
			};
			const body = { type: "blockquote", children: node.children, data: { hName: "div", hProperties: { className: ["sf-pane__body"] } } };
			node.data = { hName: "aside", hProperties: { className: ["sf-pane", ...(kind === "note" ? [] : [`sf-pane--${kind}`])] } };
			node.children = [title, body];
		});
	};
}
