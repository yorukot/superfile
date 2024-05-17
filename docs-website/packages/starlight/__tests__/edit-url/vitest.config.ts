import { defineVitestConfig } from '../test-config';

export default defineVitestConfig({
	title: 'Docs With Edit Links',
	editLink: {
		baseUrl: 'https://github.com/withastro/starlight/edit/main/docs/',
	},
});
