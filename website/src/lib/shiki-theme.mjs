// The Shiki theme for every code block. Code always sits on a "screen"
// (the charcoal `term-bg`), in both the light and the dark site theme, so these
// are the `term-*` tokens from src/styles/tokens.css.
const ink = "#EFE8DD";
const muted = "#A3ABB2";
const signal = "#FF7A4F";
const teal = "#8FC2C5";
const amber = "#F7C078";

export const superfileShikiTheme = {
	name: "superfile",
	type: "dark",
	colors: {
		"editor.background": "#303841",
		"editor.foreground": ink
	},
	tokenColors: [
		{ scope: ["comment", "punctuation.definition.comment"], settings: { foreground: muted, fontStyle: "italic" } },
		{
			scope: [
				"keyword",
				"storage",
				"entity.name.tag",
				"support.type.property-name",
				"variable.other.key",
				"entity.name.section",
				"entity.other.attribute-name",
				"meta.object-literal.key",
				"variable.parameter.option"
			],
			settings: { foreground: signal }
		},
		{ scope: ["string", "string.quoted", "markup.inline.raw", "markup.fenced_code"], settings: { foreground: teal } },
		{ scope: ["constant", "constant.numeric", "constant.language", "variable.language", "support.constant"], settings: { foreground: amber } },
		{ scope: ["entity.name.function", "support.function", "variable", "punctuation", "meta.brace"], settings: { foreground: ink } },
		{ scope: ["markup.heading", "markup.bold"], settings: { foreground: amber, fontStyle: "bold" } },
		{ scope: ["markup.inserted"], settings: { foreground: teal } },
		{ scope: ["markup.deleted"], settings: { foreground: signal } }
	]
};
