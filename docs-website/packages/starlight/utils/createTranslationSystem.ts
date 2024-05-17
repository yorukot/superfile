import type { i18nSchemaOutput } from '../schemas/i18n';
import builtinTranslations from '../translations/index';
import type { StarlightConfig } from './user-config';

export function createTranslationSystem<T extends i18nSchemaOutput>(
	userTranslations: Record<string, T>,
	config: Pick<StarlightConfig, 'defaultLocale' | 'locales'>
) {
	/** User-configured default locale. */
	const defaultLocale = config.defaultLocale?.locale || 'root';

	/** Default map of UI strings based on Starlight and user-configured defaults. */
	const defaults = buildDictionary(
		builtinTranslations.en!,
		userTranslations.en,
		builtinTranslations[defaultLocale] || builtinTranslations[stripLangRegion(defaultLocale)],
		userTranslations[defaultLocale]
	);

	/**
	 * Generate a utility function that returns UI strings for the given `locale`.
	 *
	 * Also includes an `all()` method for getting the entire dictionary.
	 *
	 * @param {string | undefined} [locale]
	 * @example
	 * const t = useTranslations('en');
	 * const label = t('search.label');
	 * // => 'Search'
	 * const dictionary = t.all();
	 * // => { 'skipLink.label': 'Skip to content', 'search.label': 'Search', ... }
	 */
	return function useTranslations(locale: string | undefined) {
		const lang = localeToLang(locale, config.locales, config.defaultLocale);
		const dictionary = buildDictionary(
			defaults,
			builtinTranslations[lang] || builtinTranslations[stripLangRegion(lang)],
			userTranslations[lang]
		);
		const t = <K extends keyof typeof dictionary>(key: K) => dictionary[key];
		t.all = () => dictionary;
		return t;
	};
}

/**
 * Strips the region subtag from a BCP-47 lang string.
 * @param {string} [lang]
 * @example
 * const lang = stripLangRegion('en-GB'); // => 'en'
 */
function stripLangRegion(lang: string) {
	return lang.replace(/-[a-zA-Z]{2}/, '');
}

/**
 * Get the BCP-47 language tag for the given locale.
 * @param locale Locale string or `undefined` for the root locale.
 */
function localeToLang(
	locale: string | undefined,
	locales: StarlightConfig['locales'],
	defaultLocale: StarlightConfig['defaultLocale']
): string {
	const lang = locale ? locales?.[locale]?.lang : locales?.root?.lang;
	const defaultLang = defaultLocale?.lang || defaultLocale?.locale;
	return lang || defaultLang || 'en';
}

type BuiltInStrings = (typeof builtinTranslations)['en'];

/** Build a dictionary by layering preferred translation sources. */
function buildDictionary<T extends Record<string, string | undefined>>(
	base: BuiltInStrings,
	...dictionaries: (T | BuiltInStrings | undefined)[]
): BuiltInStrings & T {
	const dictionary = { ...base };
	// Iterate over alternate dictionaries to avoid overwriting preceding values with `undefined`.
	for (const dict of dictionaries) {
		for (const key in dict) {
			const value = dict[key as keyof typeof dict];
			if (value) dictionary[key as keyof typeof dictionary] = value;
		}
	}
	return dictionary as BuiltInStrings & T;
}
