import { describe, expect, test } from 'vitest';
import translations from '../../translations';
import { useTranslations } from '../../utils/translations';

describe('built-in translations', () => {
	test('includes English', () => {
		expect(translations).toHaveProperty('en');
	});
});

describe('useTranslations()', () => {
	test('works when no i18n collection is available', () => {
		const t = useTranslations(undefined);
		expect(t).toBeTypeOf('function');
		expect(t('page.editLink')).toBe(translations.en?.['page.editLink']);
	});

	test('returns default locale for unknown language', () => {
		const locale = 'xx';
		expect(translations).not.toHaveProperty(locale);
		const t = useTranslations(locale);
		expect(t('page.editLink')).toBe(translations.en?.['page.editLink']);
	});

	test('uses built-in translations for regional variants', () => {
		const t = useTranslations('pt-br');
		expect(t('page.nextLink')).toBe(translations.pt?.['page.nextLink']);
		expect(t('page.nextLink')).not.toBe(translations.en?.['page.nextLink']);
	});
});
