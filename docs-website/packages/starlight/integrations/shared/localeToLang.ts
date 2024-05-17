import type { StarlightConfig } from '../../types';

/**
 * Get the BCP-47 language tag for the given locale.
 * @param locale Locale string or `undefined` for the root locale.
 */
export function localeToLang(config: StarlightConfig, locale: string | undefined): string {
	const lang = locale ? config.locales?.[locale]?.lang : config.locales?.root?.lang;
	const defaultLang = config.defaultLocale?.lang || config.defaultLocale?.locale;
	return lang || defaultLang || 'en';
}
