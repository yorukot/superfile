# superfile website

The source of [superfile.dev](https://superfile.dev): the home page and the docs. It is an [Astro](https://astro.build) site.

## Commands

Run these from `website/`.

| Command        | What it does                                                      |
| :------------- | :---------------------------------------------------------------- |
| `pnpm install` | Installs the dependencies.                                        |
| `pnpm dev`     | Starts a local server at `localhost:4321`.                        |
| `pnpm build`   | Builds the site to `dist/` and creates the Pagefind search index. |
| `pnpm preview` | Serves the build. Search only works here, not in `pnpm dev`.      |
| `pnpm format`  | Formats every file with Prettier. CI checks this.                 |

## Where things are

| Path                        | What it holds                                                                                            |
| :-------------------------- | :------------------------------------------------------------------------------------------------------- |
| `src/pages/index.astro`     | The home page. Each of its sections is a component in `src/components/home/`.                            |
| `src/content/docs/`         | The docs, in Markdown and MDX. A file's path is its URL.                                                 |
| `src/layouts/Docs.astro`    | The docs layout. Its parts (sidebar, version switcher, table of contents) are in `src/components/docs/`. |
| `src/lib/navigation.ts`     | The docs sidebar. Add a new docs page here.                                                              |
| `docs.config.mjs`           | How the docs are versioned, and which version opens at the site root.                                    |
| `scripts/sync-versions.mjs` | Copies the docs of every release out of its git tag into `.versions/`. It runs before a build.           |
| `src/components/`           | What every page shares: nav, footer, search, theme toggle, command and copy button.                      |
| `src/lib/`                  | Data read at build time: releases from the changelog, theme colors from the theme files, and so on.      |
| `src/styles/tokens.css`     | The design tokens: colors for both themes, spacing, radii and type.                                      |
| `src/styles/components.css` | The design system's shared components. Classes start with `sf-`.                                         |
| `src/styles/prose.css`      | How rendered Markdown looks.                                                                             |

## How the code is organized

- A component holds its own markup, its own scoped `<style>` and, if it does something, its own `<script>`. Only what several components share is global: the tokens and the `sf-` classes.
- Pages are ordinary documents, without a client-side router. A script runs once, when its page loads.
- Colors come from the tokens. Each is written once as `light-dark(light, dark)`, and `<html data-theme>` picks the side.
- A component's script is rendered right after the component. Do not count on `:first-child` or `:last-child` between sibling components: use `:first-of-type` or a class.

## Design rules

The site follows the superfile design system. The rules that matter most when you add something:

- There are two kinds of container. A box has a 1px border and no fill. A screen has the charcoal `term-bg` fill and no border, and holds terminal content. Do not mix them, and do not add shadows or gradients.
- A fill on the page means "active": the selected row, the active tab, the key to press next.
- No text is smaller than 14px.
- Write the name in lowercase: superfile. UI text is lowercase mono, and headlines are plain sentences.
