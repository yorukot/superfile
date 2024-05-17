import { useMode, modeOklch, modeRgb, formatHex, clampChroma } from 'culori/fn';

const rgb = useMode(modeRgb);
export const oklch = useMode(modeOklch);

/** Convert an OKLCH color to an RGB hex code. */
export const oklchToHex = (l: number, c: number, h: number) => {
	const okLchColor = oklch(`oklch(${l}% ${c} ${h})`)!;
	const rgbColor = rgb(clampChroma(okLchColor, 'oklch'));
	return formatHex(rgbColor);
};

/** Generate dark and light palettes based on user-selected hue and chroma values. */
export function getPalettes(config: {
	accent: { hue: number; chroma: number };
	gray: { hue: number; chroma: number };
}) {
	const {
		accent: { hue: ah, chroma: ac },
		gray: { hue: gh, chroma: gc },
	} = config;
	return {
		dark: {
			// Accents
			'accent-low': oklchToHex(25.94, ac / 3, ah),
			accent: oklchToHex(52.28, ac, ah),
			'accent-high': oklchToHex(83.38, ac / 3, ah),
			// Grays
			white: oklchToHex(100, 0, 0),
			'gray-1': oklchToHex(94.77, gc / 2.5, gh),
			'gray-2': oklchToHex(81.34, gc / 2, gh),
			'gray-3': oklchToHex(63.78, gc, gh),
			'gray-4': oklchToHex(46.01, gc, gh),
			'gray-5': oklchToHex(34.09, gc, gh),
			'gray-6': oklchToHex(27.14, gc, gh),
			black: oklchToHex(20.94, gc / 2, gh),
		},
		light: {
			// Accents
			'accent-low': oklchToHex(87.81, ac / 4, ah),
			accent: oklchToHex(52.95, ac, ah),
			'accent-high': oklchToHex(31.77, ac / 2, ah),
			// Grays
			white: oklchToHex(20.94, gc / 2, gh),
			'gray-1': oklchToHex(27.14, gc, gh),
			'gray-2': oklchToHex(34.09, gc, gh),
			'gray-3': oklchToHex(46.01, gc, gh),
			'gray-4': oklchToHex(63.78, gc, gh),
			'gray-5': oklchToHex(81.34, gc / 2, gh),
			'gray-6': oklchToHex(94.77, gc / 2.5, gh),
			'gray-7': oklchToHex(97.35, gc / 5, gh),
			black: oklchToHex(100, 0, 0),
		},
	};
}

/*
This is the default Starlight color palette in OKLCH.
Used as a reference when designing the algorithm for mapping
user hue and chroma pairs to a similar palette.

Dark mode:
  accent lo   oklch(25.94% 0.09 273.5)
  accent      oklch(52.28% 0.266 268.7)
  accent hi   oklch(83.38% 0.084 279.5)

  white       oklch(100% 0 0)
  gray-1      oklch(94.77% 0.008 278.19)
  gray-2      oklch(81.34% 0.011 274.87)
  gray-3      oklch(63.78% 0.019 265.84)
  gray-4      oklch(46.01% 0.021 270.93)
  gray-5      oklch(34.09% 0.017 267.07)
  gray-6      oklch(27.14% 0.015 267.03)
  black       oklch(20.94% 0.01 268.4)
  
Light mode:
  accent lo   oklch(87.81% 0.056 280.2)
  accent      oklch(52.95% 0.243 270.2)
  accent hi   oklch(31.77% 0.177 267.2)
  
  white       oklch(20.94% 0.01 268.4)
  gray-1      oklch(27.14% 0.015 267.03)
  gray-2      oklch(34.09% 0.017 267.07)
  gray-3      oklch(46.01% 0.021 270.93)
  gray-4      oklch(63.78% 0.019 265.84)
  gray-5      oklch(81.34% 0.011 274.87)
  gray-6      oklch(94.77% 0.008 278.19)
  gray-7      oklch(97.35% 0.004 286.32)
  black       oklch(100% 0 0)
*/
