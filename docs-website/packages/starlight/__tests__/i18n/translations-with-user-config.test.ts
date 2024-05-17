import { describe, expect, test, vi } from 'vitest';
import translations from '../../translations';
import { useTranslations } from '../../utils/translations';

vi.mock('astro:content', async () =>
	(await import('../test-utils')).mockedAstroContent({
		i18n: [
			['en-US', { 'page.editLink': 'Modify this doc!' }],
			['pt-BR', { 'page.editLink': 'Modifique esse doc!' }],
		],
	})
);

describe('useTranslations()', () => {
	test('uses user-defined translations', () => {
		const t = useTranslations(undefined);
		expect(t('page.editLink')).toBe('Modify this doc!');
		expect(t('page.editLink')).not.toBe(translations.en?.['page.editLink']);
	});

	test('uses user-defined regional translations when available', () => {
		const t = useTranslations('pt-br');
		expect(t('page.editLink')).toBe('Modifique esse doc!');
		expect(t('page.editLink')).not.toBe(translations.pt?.['page.editLink']);
	});
});
