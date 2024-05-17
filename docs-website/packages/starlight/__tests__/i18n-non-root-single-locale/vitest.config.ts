import { defineVitestConfig } from '../test-config';

export default defineVitestConfig({
	title: 'i18n with a non-root single locale',
	defaultLocale: 'fr',
	locales: {
		fr: { label: 'Français', lang: 'fr-CA' },
	},
});
