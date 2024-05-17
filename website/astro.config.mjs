import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import { pluginLineNumbers } from '@expressive-code/plugin-line-numbers';

const site = 'https://starter.obytes.com/';

// https://astro.build/config
export default defineConfig({
  site: 'https://superfile.github.io',
  integrations: [
    starlight({
      title: 'Superfile',
      description: `Superfile is a very fancy and modern terminal file manager that can complete the file operations you need!`,
      expressiveCode: {
        themes: ['dracula', 'solarized-light'],
      },
      logo: {
        light: '/src/assets/superfile-day.svg',
        dark: '/src/assets/superfile-night.svg',
        replacesTitle: true,
      },
      components: {
        LastUpdated: './src/components/LastUpdated.astro',
      },
      social: {
        github: 'https://github.com/mhnightcat/superfile',
      },
      head: [
        {
          tag: 'meta',
          attrs: { property: 'og:image', content: site + 'og.jpg?v=1' },
        },
        {
          tag: 'meta',
          attrs: { property: 'twitter:image', content: site + 'og.jpg?v=1' },
        },
        {
          tag: 'link',
          attrs: { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        },
        {
          tag: 'link',
          attrs: {
            rel: 'preconnect',
            href: 'https://fonts.gstatic.com',
            crossorigin: true,
          },
        },
        {
          tag: 'link',
          attrs: {
            rel: 'stylesheet',
            href: 'https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@500;600&display=swap',
          },
        },
        {
          tag: 'script',
          attrs: {
            src: 'https://cdn.jsdelivr.net/npm/@minimal-analytics/ga4/dist/index.js',
            async: true,
          },
        },
        {
          tag: 'script',
          content: ` window.minimalAnalytics = {
            trackingId: 'G-GQ45JJD1JC',
            autoTrack: true,
          };`,
        },
      ],
      editLink: {
				baseUrl: 'https://github.com/mhnightcat/superfile/edit/main/website/',
			},
      sidebar: [
        {
          label: 'Overview',
          link: '/overview',
        },
        {
          label: 'Start Here',
          items: [
            // Each item here is one entry in the navigation menu.
            {
              label: 'Installation',
              link: '/getting-started/installation/',
            },
            {
              label: 'Tutorial',
              link: '/getting-started/tutorial/',
            },
            {
              label: 'Hotkey list',
              link: '/getting-started/hotkey-list/',
            }
          ],
        },
        {
          label: 'Changelog',
          link: '/changelog',
        },
      ],
      customCss: ['./src/styles/custom.css'],
      lastUpdated: true,
    }),
  ],
  // Process images with sharp: https://docs.astro.build/en/guides/assets/#using-sharp
  image: {
    service: {
      entrypoint: 'astro/assets/services/sharp',
    },
  },
});
