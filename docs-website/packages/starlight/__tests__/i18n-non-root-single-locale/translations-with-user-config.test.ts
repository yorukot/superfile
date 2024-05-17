import { describe, expect, test, vi } from 'vitest';
import translations from '../../translations';
import { useTranslations } from '../../utils/translations';

vi.mock('astro:content', async () =>
	(await import('../test-utils')).mockedAstroContent({
		i18n: [['fr-CA', { 'page.editLink': 'Modifier cette doc!' }]],
	})
);

describe('useTranslations()', () => {
	test('uses user-defined translations', () => {
		const t = useTranslations('fr');
		expect(t('page.editLink')).toBe('Modifier cette doc!');
		expect(t('page.editLink')).not.toBe(translations.fr?.['page.editLink']);
	});
});
