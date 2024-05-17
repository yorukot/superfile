# Contributor Manual

We welcome contributions of any size and contributors of any skill level.
As an open source project, we believe in giving back to our contributors.
We are happy to help with guidance on PRs, technical writing, and turning any feature idea into a reality.

> **Tip for new contributors:**
> Take a look at [GitHub's Docs](https://docs.github.com/en/get-started/quickstart/hello-world) for helpful information on working with GitHub.

This document is an active work in progress — like Starlight itself! Feel free to join us in [the Astro Discord server][discord] to join the discussion. Look for the `#starlight` channel and say “Hi!” when you arrive.

## Types of contributions

There are lots of ways to contribute to Starlight.

Maintaining Starlight requires writing Astro code, as well as addressing accessibility, styling, and UX concerns.
This repository also contains the code for the Starlight docs website.
Help writing docs, catching typos and errors, as well as translating docs into other languages is always welcome.

You can also get involved by leaving feedback on [issues][issues] or reviewing [pull requests][pulls] by other contributors.

We encourage you to:

- [**Open an issue**][new-issue] to let us know of bugs in Starlight, documentation you found unclear, or other issues you run into.

- [**Look at existing issues**][issues] (especially those labelled [“good first issue”][gfi]) to find ways to contribute.

- **Make a pull request (PR)** to address an open issue or to fix obvious problems.
  Read more about [making a PR in GitHub’s docs][pr-docs]

- [**Review existing PRs**][pulls] to help us merge contributions sooner.

- [**Add or update translations**](#translations). We need help translating both Starlight’s UI and documentation.

## About this repo

This repo is a “monorepo,” meaning it contains several projects in one. It contains the Starlight docs site in [`docs/`](./docs/) and the packages that make up Starlight in [`packages/`](./packages/).

### Setting up a development environment

You can [develop locally](#developing-locally) or use an online coding development environment like [GitHub Codespaces](#developing-using-github-codespaces) or [Gitpod](#developing-using-gitpod) to get started quickly.

#### Developing locally

**Prerequisites:** Developing Starlight requires [Node.js](https://nodejs.org/en) (v16 or higher) and [pnpm](https://pnpm.io/) (v8.2 or higher). Make sure you have these installed before following these steps.

1. **Fork Starlight** to your personal GitHub account by clicking <kbd>Fork</kbd> on the [main Starlight repo page][sl].

2. **Clone your fork** of Starlight to your computer. Replace `YOUR-USERNAME` in the command below with your GitHub username to clone in a Terminal:

   ```sh
   git clone https://github.com/YOUR-USERNAME/starlight.git
   ```

3. **Change directory** to the cloned repo:

   ```sh
   cd starlight
   ```

4. **Install dependencies** with `pnpm`:

   ```sh
   pnpm i
   ```

#### Developing using Gitpod

**Prerequisites:** Developing Starlight using Gitpod requires a free [Gitpod account](https://gitpod.io).

1. **Open the Gitpod URL** [https://gitpod.io/#https://github.com/withastro/starlight](https://gitpod.io/#https://github.com/withastro/starlight). You can alternatively install a [Gitpod browser extension](https://www.gitpod.io/docs/configure/user-settings/browser-extension) which will add a "Gitpod" button when viewing [Starlight's repo on GitHub](https://github.com/withastro/starlight).

2. **Install dependencies** with `pnpm`:

   ```sh
   pnpm i
   ```

#### Developing using GitHub Codespaces

1. **Create a new codespace** via https://codespaces.new/withastro/starlight

2. If running the docs site, pass the `--host` flag to avoid “502 Bad Gateway” errors:

   ```sh
   cd docs
   pnpm dev --host
   ```

The dev container used for GitHub Codespaces can also be used with [other supporting tools](https://containers.dev/supporting), including VS Code.

### Making a Pull Request

When making a pull request containing changes impacting users to Starlight or any related packages (`packages/*`), be sure to [add a changeset](https://github.com/changesets/changesets/blob/main/docs/adding-a-changeset.md#i-am-in-a-multi-package-repository-a-mono-repo) that will describe the changes to users.
Documentation-only (`docs/*`) and non-package (`examples/*`) changes do not need changesets.

```sh
pnpm exec changeset
```

## Testing

### Testing visual changes while you work

Run the Astro dev server on the docs site to see how changes you make impact a project using Starlight.

To do this, move into the `docs/` directory from the root of the repo and then run `pnpm dev`:

```sh
cd docs
pnpm dev
```

You should then be able to open <http://localhost:4321> and see your changes.

> **Note**
> Changes to the Starlight integration will require you to quit and restart the dev server to take effect.

### Check for broken links in the docs site

When adding or translating content in the Starlight docs site, you can check all internal links are valid.
All GitHub PRs are checked this way automatically, but testing locally can help if you want to confirm changes are correct before committing them.

To do this, move into the `docs/` directory from the root of the repo and then build the site with the `CHECK_LINKS` environment variable:

```sh
cd docs
CHECK_LINKS=true pnpm build
```

If there are any broken links, the build will fail and log which pages need to be fixed.

### Unit tests

The Starlight package includes unit tests in [`packages/starlight/__tests__/`](./packages/starlight/__tests__/), which are run using [Vitest][vitest].

To run tests, move into the Starlight package and then run `pnpm test`:

```sh
cd packages/starlight
pnpm test
```

This will run tests and then listen for changes, re-running tests when files change.

#### Test environments

A lot of Starlight code relies on Vite virtual modules provided either by Astro or by Starlight itself. Each subdirectory of `packages/starlight/__tests__/` should contain a `vitest.config.ts` file that uses the `defineVitestConfig()` helper to define a valid test environment for tests in that directory. This helper takes a single argument, which provides a Starlight user config object:

```ts
// packages/starlight/__tests__/basics/vitest.config.ts
import { defineVitestConfig } from '../test-config';

export default defineVitestConfig({
  title: 'Basics',
});
```

This allows you to run tests of Starlight code against different combinations of Starlight configuration options.

#### Mocking content collections

Starlight relies on a user’s `docs` and (optional) `i18n` content collections, which aren’t available during testing. You can use a top-level `vi.mock()` call and the `mockedAstroContent` helper to set up fake collection entries for the current test file:

```js
import { describe, expect, test, vi } from 'vitest';

vi.mock('astro:content', async () =>
  (await import('../test-utils')).mockedAstroContent({
    docs: [
      ['index.mdx', { title: 'Home Page' }],
      ['environmental-impact.md', { title: 'Eco-friendly docs' }],
    ],
    i18n: [['en', { 'page.editLink': 'Modify this doc!' }]],
  })
);
```

#### Test coverage

To see how much of Starlight’s code is currently being tested, run `pnpm test:coverage` from the Starlight package:

```sh
cd packages/starlight
pnpm test:coverage
```

This will print a table to your terminal and also generate an HTML report you can load in a web browser by opening [`packages/starlight/__coverage__/index.html`](./packages/starlight/__coverage__/index.html).

### End-to-end (E2E) tests

Starlight also includes E2E tests in [`packages/starlight/__e2e__/`](./packages/starlight/__e2e__/), which are run using [Playwright][playwright].

To run these tests, move into the Starlight package and then run `pnpm test:e2e`:

```sh
cd packages/starlight
pnpm test:e2e
```

#### Test fixtures

Each subdirectory of `packages/starlight/__e2e__/fixtures` should contain the basic files needed to run Starlight (`package.json`, `astro.config.mjs`, a content collection configuration in `src/content/config.ts` and some content to render in `src/content/docs/`).

The `testFactory()` helper can be used in a test file to define the fixture which will be built and loaded in a preview server during a set of tests.

```ts
// packages/starlight/__e2e__/feature.test.ts
import { testFactory } from './test-utils';

const test = await testFactory('./fixtures/basics/');
```

This allows you to run tests against different combinations of Astro and Starlight configuration options for various content.

#### When to add E2E tests?

E2E are most useful for testing what happens on a page after it has been loaded by a browser. They run slower than unit tests so they should be used sparingly when unit tests aren’t sufficient.

## Translations

Translations help make Starlight accessible to more people.

### Translating Starlight’s UI

Starlight’s UI comes with some built-in text elements. For example, the table of contents on a Starlight page has a heading of “On this page” and the theme picker shows “Light”, “Dark”, and “Auto” labels. Starlight aims to provide these in as many languages as possible.

Help out by adding or updating translation files in [`packages/starlight/translations`](./packages/starlight/translations/).
Each language’s JSON file follows the [translation structure described in Starlight’s docs](https://starlight.astro.build/guides/i18n/#translate-starlights-ui).

📺 **Prefer a visual walkthrough?** [Watch an introduction to Starlight’s translation files.](https://scrimba.com/scrim/cpb44bt3)

### Translating Starlight’s docs

Starlight’s documentation is also translated into multiple languages. You can find the source code for the site in [the `docs/` directory](./docs/) of this repository.

Help out by:

- Reviewing [open translation PRs][pulls]
- Updating out-of-date translated pages
- Adding an untranslated page

Visit **<https://i18n.starlight.astro.build>** to track translation progress for the currently supported languages.

#### Adding a new language to Starlight’s docs

To add a language, you will need its BCP-47 tag and a label. See [“Adding a new language”](https://contribute.docs.astro.build/guides/i18n/#adding-a-new-language) in the Astro docs repo for some helpful tips around choosing these.

- Add your language to the `locales` config in `docs/astro.config.mjs`
- Add your language to the `locales` config in `docs/lunaria.config.json`
- Add your language’s subtag to the i18n label config in `.github/labeler.yml`
- Add your language to the `pa11y` script’s `--sitemap-exclude` flag in `package.json`
- Create the first translated page for your language.
  This must be the Starlight landing page: `docs/src/content/docs/{language}/index.mdx`.
- Open a pull request on GitHub to add your changes to Starlight!

## Understanding Starlight

- Starlight is built as an Astro integration.
  Read the [Astro Integration API docs][api-docs] to learn more about how integrations work.

  The Starlight integration is exported from [`packages/starlight/index.ts`](./packages/starlight/index.ts).
  It sets up Starlight’s routing logic, parses user config, and adds configuration to a Starlight user’s Astro project.

- Most pages in a Starlight project are built using a single [`packages/starlight/index.astro`](./packages/starlight/index.astro) route.
  If you’ve worked on an Astro site before, much of this should look familiar: it’s an Astro component and uses a number of other components to build a page based on user content.

- Starlight consumes a user’s content from the `'docs'` [content collection](https://docs.astro.build/en/guides/content-collections/).
  This allows us to specify the permissible frontmatter via [a Starlight-specific schema](./packages/starlight/schema.ts) and get predictable data while providing clear error messages if a user sets invalid frontmatter in a page.

- Components that require JavaScript for their functionality are all written without a UI framework, most often as custom elements.
  This helps keep Starlight lightweight and makes it easier for a user to choose to add components from a framework of their choice to their project.

- Components that require client-side JavaScript or CSS should use JavaScript/CSS features that are well-supported by browsers.

  You can find a list of supported browsers and their versions using this [browserslist query](https://browsersl.ist/#q=%3E+0.5%25%2C+not+dead%2C+Chrome+%3E%3D+88%2C+Edge+%3E%3D+88%2C+Firefox+%3E%3D+98%2C+Safari+%3E%3D+15.4%2C+iOS+%3E%3D+15.4%2C+not+op_mini+all). To check whether or not a feature is supported, you can visit the [Can I use](https://caniuse.com) website and search for the feature.

[discord]: https://astro.build/chat
[issues]: https://github.com/withastro/starlight/issues
[sl]: https://github.com/withastro/starlight/pulls
[pulls]: https://github.com/withastro/starlight/pulls
[new-issue]: https://github.com/withastro/starlight/issues/new/choose
[pr-docs]: https://docs.github.com/en/get-started/quickstart/contributing-to-projects#making-a-pull-request
[gfi]: https://github.com/withastro/starlight/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22+
[api-docs]: https://docs.astro.build/en/reference/integrations-reference/
[vitest]: https://vitest.dev/
[playwright]: https://playwright.dev/

## Showcase

We love to see websites built with Starlight and share them with the community on our [showcase](https://starlight.astro.build/resources/showcase/) page.
If you’ve built a documentation site with Starlight, adding it to the showcase is just a pull request away!

1. Set up a development environment by following the [“Setting up a development environment”](#setting-up-a-development-environment) instructions.
2. Add a screenshot of your site to the `docs/src/assets/showcase/` directory. The image file must:
   - Be a `.png` file and named after your site’s domain, e.g. `example.com.png`.
   - Have the dimensions of 800 × 450 pixels.
3. Add a new entry for your website in `docs/src/components/showcase-sites.astro`.

   - The new entry must be appended at the end of the existing list of sites.
   - The `title` attribute must be the name of your site with no extra details.
   - The `href` attribute must be the URL of your Starlight site. If your documentation is hosted on a subdomain or subdirectory, include that in the URL.
   - The `thumbnail` attribute must be the filename of the screenshot you added in step 2.

   ```diff
     <Card title="Example" href="https://example.net" thumbnail="example.net.png" />
     <Card title="Last Example" href="https://example.org" thumbnail="example.org.png" />
   + <Card title="New Example" href="https://example.com" thumbnail="example.com.png" />
   </FluidGrid>
   ```

4. Open a pull request on GitHub to add your changes.
