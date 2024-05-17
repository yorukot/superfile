import { pluginFramesTexts } from 'astro-expressive-code';
import type { StarlightConfig } from '../../types';
import type { createTranslationSystemFromFs } from '../../utils/translations-fs';
import { localeToLang } from '../shared/localeToLang';

export function addTranslations(
	config: StarlightConfig,
	useTranslations: ReturnType<typeof createTranslationSystemFromFs>
) {
	addTranslationsForLocale(config.defaultLocale.locale, config, useTranslations);
	if (config.isMultilingual) {
		for (const locale in config.locales) {
			if (locale === config.defaultLocale.locale || locale === 'root') continue;
			addTranslationsForLocale(locale, config, useTranslations);
		}
	}
}

function addTranslationsForLocale(
	locale: string | undefined,
	config: StarlightConfig,
	useTranslations: ReturnType<typeof createTranslationSystemFromFs>
) {
	const lang = localeToLang(config, locale);
	const t = useTranslations(locale);
	const translationKeys = [
		'expressiveCode.copyButtonCopied',
		'expressiveCode.copyButtonTooltip',
		'expressiveCode.terminalWindowFallbackTitle',
	] as const;
	translationKeys.forEach((key) => {
		const translation = t(key);
		if (!translation) return;
		const ecId = key.replace(/^expressiveCode\./, '');
		pluginFramesTexts.overrideTexts(lang, { [ecId]: translation });
	});
}
