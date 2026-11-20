export const element = (tagName, properties = {}, children = []) => ({ type: "element", tagName, properties, children });
export const text = value => ({ type: "text", value });

/** A Lucide icon, from the shapes inside its 24×24 svg: [["path", { d: "…" }], …]. */
export const icon = shapes =>
	element(
		"svg",
		{ viewBox: "0 0 24 24", fill: "none", stroke: "currentColor", strokeWidth: 2, strokeLinecap: "round", strokeLinejoin: "round", ariaHidden: "true" },
		shapes.map(([tagName, properties]) => element(tagName, properties))
	);
