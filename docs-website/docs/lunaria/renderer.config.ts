import { defineRendererConfig } from '@lunariajs/core';
import { TitleParagraph } from './components';

export default defineRendererConfig({
	slots: {
		afterTitle: TitleParagraph,
	},
});
