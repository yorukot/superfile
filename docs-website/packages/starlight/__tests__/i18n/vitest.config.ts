import { defineVitestConfig } from '../test-config';

export default defineVitestConfig({
	title: 'i18n with no root locale',
	defaultLocale: 'en',
	locales: {
		fr: { label: 'French' },
		en: { label: 'English', lang: 'en-US' },
		ar: { label: 'Arabic', dir: 'rtl' },
		'pt-br': { label: 'Brazilian Portuguese', lang: 'pt-BR' },
	},
});
