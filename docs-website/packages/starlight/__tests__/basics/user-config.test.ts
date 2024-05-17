import { expect, test } from 'vitest';
import { StarlightConfigSchema } from '../../utils/user-config';

test('preserve social config order', () => {
	const config = StarlightConfigSchema.parse({
		title: 'Test',
		social: {
			twitch: 'https://www.twitch.tv/bholmesdev',
			github: 'https://github.com/withastro/starlight',
			discord: 'https://astro.build/chat',
		},
	});
	expect(Object.keys(config.social || {})).toEqual(['twitch', 'github', 'discord']);
});
