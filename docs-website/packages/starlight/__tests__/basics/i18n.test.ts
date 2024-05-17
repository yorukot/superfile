import { describe, expect, test } from 'vitest';
import { pickLang } from '../../utils/i18n';

describe('pickLang', () => {
	const dictionary = { en: 'Hello', fr: 'Bonjour' };

	test('returns the requested language string', () => {
		expect(pickLang(dictionary, 'en')).toBe('Hello');
		expect(pickLang(dictionary, 'fr')).toBe('Bonjour');
	});

	test('returns undefined for unknown languages', () => {
		expect(pickLang(dictionary, 'ar' as any)).toBeUndefined();
	});
});
