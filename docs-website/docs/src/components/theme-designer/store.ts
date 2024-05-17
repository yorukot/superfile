import { map } from './atom';

export const presets = {
	ocean: {
		accent: { hue: 240, chroma: 0.27 },
		gray: { hue: 220, chroma: 0.025 },
	},
	forest: {
		accent: { hue: 140, chroma: 0.27 },
		gray: { hue: 140, chroma: 0.03 },
	},
	oxide: {
		accent: { hue: 30, chroma: 0.27 },
		gray: { hue: 30, chroma: 0.02 },
	},
	nebula: {
		accent: { hue: 320, chroma: 0.27 },
		gray: { hue: 305, chroma: 0.07 },
	},
	default: {
		accent: { hue: 269, chroma: 0.27 },
		gray: { hue: 270, chroma: 0.016 },
	},
};

export const store = {
	accent: map(presets.default.accent),
	gray: map(presets.default.gray),
};

export const usePreset = (name: string) => {
	if (name in presets) {
		const { accent, gray } = presets[name as keyof typeof presets];
		store.accent.set(accent);
		store.gray.set(gray);
	}
};

const MAX_CHROMA = 0.27;

export const useRandom = () => {
	store.accent.set({ hue: randomHue(), chroma: MAX_CHROMA - randomChroma() });
	store.gray.set({ hue: randomHue(), chroma: randomChroma() });
};

const randomHue = () => Math.round(Math.random() * 360);
const randomChroma = () => Math.pow(Math.random(), 3) * MAX_CHROMA;
