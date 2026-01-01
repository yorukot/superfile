import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import { pluginLineNumbers } from "@expressive-code/plugin-line-numbers";
import starlightGiscus from "starlight-giscus";
import sitemap from "@astrojs/sitemap";

const site = "https://superfile.dev/";

// https://astro.build/config
export default defineConfig({
  site: site,
  integrations: [
    sitemap(),
    starlight({
      title: "superfile",
      description: `superfile is a very fancy and modern terminal file manager that can complete the file operations you need!`,
      expressiveCode: {
        themes: ["dracula", "solarized-light"],
      },
      logo: {
        light: "/src/assets/superfile-day.svg",
        dark: "/src/assets/superfile-night.svg",
        replacesTitle: true,
      },
      components: {
        LastUpdated: "./src/components/LastUpdated.astro",
      },
      plugins: [
        starlightGiscus({
          repo: "yorukot/superfile",
          repoId: "R_kgDOLil1MA",
          category: "Docs Comments",
          categoryId: "DIC_kwDOLil1MM4CfbH7",
          mapping: "title",
          strict: false,
          reactionsEnabled: true,
          emitMetadata: false,
          inputPosition: "top",
          theme: "preferred_color_scheme",
          lang: "en",
          loading: "lazy",
        }),
      ],
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/yorukot/superfile",
        },
        {
          icon: "discord",
          label: "Discord",
          href: "https://discord.gg/YYtJ23Du7B",
        },
      ],
      head: [
        {
          tag: "meta",
          attrs: { property: "og:image", content: site + "og.jpg?v=1" },
        },
        {
          tag: "meta",
          attrs: { property: "twitter:image", content: site + "og.jpg?v=1" },
        },
        {
          tag: "link",
          attrs: { rel: "preconnect", href: "https://fonts.googleapis.com" },
        },
        {
          tag: "link",
          attrs: {
            rel: "preconnect",
            href: "https://fonts.gstatic.com",
            crossorigin: true,
          },
        },
        {
          tag: "link",
          attrs: {
            rel: "preload",
            href: "https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@500;600&display=swap",
            as: "style",
            onload: "this.onload=null;this.rel='stylesheet'",
          },
        },
        {
          tag: "noscript",
          content:
            '<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@500;600&display=swap">',
        },
        {
          tag: "script",
          attrs: {
            src: "https://cdn.jsdelivr.net/npm/@minimal-analytics/ga4/dist/index.js",
            async: true,
          },
        },
        {
          tag: "script",
          content: ` window.minimalAnalytics = {
            trackingId: 'G-WFLBCRZ7MC',
            autoTrack: true,
          };`,
        },
        {
          tag: "script",
          attrs: {
            defer: true,
            src: "https://umami.yorukot.me/script.js",
            "data-website-id": "9286f04f-4bcd-43f3-8ab7-b479c249f2a7",
          },
        },
      ],
      editLink: {
        baseUrl: "https://github.com/yorukot/superfile/edit/main/website/",
      },
      sidebar: [
        {
          label: "Overview",
          link: "/overview",
        },
        {
          label: "Start Here",
          items: [
            {
              label: "Installation",
              link: "/getting-started/installation/",
            },
            {
              label: "Tutorial",
              link: "/getting-started/tutorial/",
            },
            {
              label: "Image Preview",
              link: "/getting-started/image-preview/",
            },
          ],
        },
        {
          label: "Configure",
          items: [
            {
              label: "All config file path",
              link: "/configure/config-file-path",
            },
            {
              label: "superfile config",
              link: "/configure/superfile-config/",
            },
            {
              label: "Custom hotkeys",
              link: "/configure/custom-hotkeys/",
            },
            {
              label: "Custom theme",
              link: "/configure/custom-theme",
            },
            {
              label: "Enable plugin",
              link: "/configure/enable-plugin",
            },
          ],
        },
        {
          label: "List",
          items: [
            {
              label: "Hotkey list",
              link: "/list/hotkey-list/",
            },
            {
              label: "Theme list",
              link: "/list/theme-list/",
            },
            {
              label: "Plugin list",
              link: "/list/plugin-list/",
            },
          ],
        },
        {
          label: "Contribute",
          items: [
            {
              label: "How to contribute",
              link: "/contribute/how-to-contribute",
            },
            {
              label: "File structure",
              link: "/contribute/file-struct",
            },
            {
              label: "Implementation Info",
              link: "/contribute/implementation-info",
            },
          ],
        },
        {
          label: "Troubleshooting",
          link: "/troubleshooting",
        },
        {
          label: "How to contribute",
          link: "/contribute/how-to-contribute",
        },
        {
          label: "Changelog",
          link: "/changelog",
        },
      ],
      customCss: ["./src/styles/custom.css"],
      lastUpdated: true,
    }),
  ],
  // Process images with sharp: https://docs.astro.build/en/guides/assets/#using-sharp
  image: {
    service: {
      entrypoint: "astro/assets/services/sharp",
      config: {
        limitInputPixels: false,
      },
    },
  },
});
