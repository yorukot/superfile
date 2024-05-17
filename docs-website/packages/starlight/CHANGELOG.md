# @astrojs/starlight

## 0.22.4

### Patch Changes

- [#1871](https://github.com/withastro/starlight/pull/1871) [`03bb126b`](https://github.com/withastro/starlight/commit/03bb126b74d9adaba1be2f0df3f34566815dd77f) Thanks [@delucis](https://github.com/delucis)! - Adds a `blueSky` icon and social link option

- [#1873](https://github.com/withastro/starlight/pull/1873) [`13f33b81`](https://github.com/withastro/starlight/commit/13f33b81fd51d18165eb52f2a0c02b890084e4bd) Thanks [@ekfuhrmann](https://github.com/ekfuhrmann)! - Adds 1 new icon: `alpine`

- [#1857](https://github.com/withastro/starlight/pull/1857) [`32cdfaf0`](https://github.com/withastro/starlight/commit/32cdfaf0155e65ff6fbe9c0cfacd6969ab0015d9) Thanks [@tarikcoskun](https://github.com/tarikcoskun)! - Updates Turkish UI translations

- [#1736](https://github.com/withastro/starlight/pull/1736) [`cfa94a34`](https://github.com/withastro/starlight/commit/cfa94a346ef10804b90db28d217be175e1c1d5ed) Thanks [@julien-deramond](https://github.com/julien-deramond)! - Prevent list items from overflowing Markdown content

## 0.22.3

### Patch Changes

- [#1838](https://github.com/withastro/starlight/pull/1838) [`9fe84754`](https://github.com/withastro/starlight/commit/9fe847544f1edb85bf5b25cd81db39227814335e) Thanks [@delucis](https://github.com/delucis)! - Adds extra information to the errors thrown by the `<Steps>` component to help locate misformatted code

- [#1863](https://github.com/withastro/starlight/pull/1863) [`50be60bb`](https://github.com/withastro/starlight/commit/50be60bbc5cbc42db42e868b9e8f128b4dcbd6a5) Thanks [@torn4dom4n](https://github.com/torn4dom4n)! - Update Vietnamese translation

- [#1837](https://github.com/withastro/starlight/pull/1837) [`a33a1223`](https://github.com/withastro/starlight/commit/a33a12231772c1dc4b7cc2db3477a6802f3ef53e) Thanks [@delucis](https://github.com/delucis)! - Adds three new icons: `comment`, `comment-alt`, `heart`

- [#1842](https://github.com/withastro/starlight/pull/1842) [`c7838636`](https://github.com/withastro/starlight/commit/c7838636edb8d60a2422ce76a2db511b9cebbb70) Thanks [@delucis](https://github.com/delucis)! - Moves the `href` used in the site title link to Starlight’s route data object. This makes it possible for overrides to change the title link while reusing Starlight’s default component implemenation.

- [#1840](https://github.com/withastro/starlight/pull/1840) [`cb85563c`](https://github.com/withastro/starlight/commit/cb85563c9a3d4eb2925ad884e6a4e8698a15381b) Thanks [@MiahaCybersec](https://github.com/MiahaCybersec)! - Adds 1 new icon: `hackerone`

## 0.22.2

### Patch Changes

- [#1811](https://github.com/withastro/starlight/pull/1811) [`fe06aa13`](https://github.com/withastro/starlight/commit/fe06aa1307208ef9f5b249181ec29837f96940c2) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fixes a `<Tabs>` sync issue when inconsistently using the `icon` prop or not on `<TabItem>` components.

- [#1826](https://github.com/withastro/starlight/pull/1826) [`52ea7381`](https://github.com/withastro/starlight/commit/52ea7381e131338a03cffb3499ba1699951cea1e) Thanks [@dragomano](https://github.com/dragomano)! - Updates Russian UI translations

## 0.22.1

### Patch Changes

- [`1c0fc384`](https://github.com/withastro/starlight/commit/1c0fc3849771713d5a3e7a572bdbf1483ae5551b) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fixes an issue where the `siteTitle` property would not be set when using the `<StarlightPage />` component.

## 0.22.0

### Minor Changes

- [#640](https://github.com/withastro/starlight/pull/640) [`7dc503ea`](https://github.com/withastro/starlight/commit/7dc503ea7993123a4aeff453d08de41cac887353) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds support for syncing multiple sets of tabs on the same page.

- [#1620](https://github.com/withastro/starlight/pull/1620) [`ca0678ca`](https://github.com/withastro/starlight/commit/ca0678ca556d739bda9648edc1b79c764fdea851) Thanks [@emjio](https://github.com/emjio)! - Adds support for translating the site title

  ⚠️ **Potentially breaking change:** The shape of the `title` field on Starlight’s internal config object has changed. This used to be a string, but is now an object.

  If you are relying on `config.title` (for example in a custom `<SiteTitle>` or `<Head>` component), you will need to update your code. We recommend using the new [`siteTitle` prop](https://starlight.astro.build/reference/overrides/#sitetitle) available to component overrides:

  ```astro
  ---
  import type { Props } from '@astrojs/starlight/props';

  // The site title for this page’s language:
  const { siteTitle } = Astro.props;
  ---
  ```

- [#1613](https://github.com/withastro/starlight/pull/1613) [`61493e55`](https://github.com/withastro/starlight/commit/61493e55f1a80362af13f98d665018376e987439) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds new `draft` frontmatter option to exclude a page from production builds.

- [#640](https://github.com/withastro/starlight/pull/640) [`7dc503ea`](https://github.com/withastro/starlight/commit/7dc503ea7993123a4aeff453d08de41cac887353) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Updates the default `line-height` from `1.8` to `1.75`. This change avoids having a line height with a fractional part which can cause scripts accessing dimensions involving the line height to get an inconsistent rounded value in various browsers.

  If you want to preserve the previous `line-height`, you can add the following custom CSS to your site:

  ```css
  :root {
  	--sl-line-height: 1.8;
  }
  ```

- [#1720](https://github.com/withastro/starlight/pull/1720) [`749ddf85`](https://github.com/withastro/starlight/commit/749ddf85a21d8ed1bfedbe60dee676cdd8784e96) Thanks [@jacobdalamb](https://github.com/jacobdalamb)! - Updates `astro-expressive-code` dependency to the latest minor release (0.35) and exposes a new `@astrojs/starlight/expressive-code/hast` module for users who need to use Expressive Code’s version of `hast`.

  This includes a potentially breaking change if you use custom Expressive Code plugins. See the [Expressive Code release notes](https://expressive-code.com/releases/#0340) for full details.

- [#1769](https://github.com/withastro/starlight/pull/1769) [`bd5f1cbd`](https://github.com/withastro/starlight/commit/bd5f1cbd5aef9e2d78e7f7187eb07deee87399d0) Thanks [@ncjones](https://github.com/ncjones)! - Adds support for [accessing frontmatter data as a variable](https://docs.astro.build/en/guides/integrations-guide/markdoc/#access-frontmatter-from-your-markdoc-content) when using Markdoc

### Patch Changes

- [#1788](https://github.com/withastro/starlight/pull/1788) [`681a4273`](https://github.com/withastro/starlight/commit/681a427366755fec71ba65d45e36f7d1267cf387) Thanks [@dragomano](https://github.com/dragomano)! - Adds Russian translations for Expressive Code labels

- [#1780](https://github.com/withastro/starlight/pull/1780) [`4db6025a`](https://github.com/withastro/starlight/commit/4db6025a1c5c56cac2e3a98bd2e13124402445c7) Thanks [@MiahaCybersec](https://github.com/MiahaCybersec)! - Adds 1 new icon: `signal`

- [#1785](https://github.com/withastro/starlight/pull/1785) [`65009c9c`](https://github.com/withastro/starlight/commit/65009c9cf8b0570303ecb87713e1c2968a704437) Thanks [@dreyfus92](https://github.com/dreyfus92)! - Adds 5 new icons: `node`, `cloudflare`, `vercel`, `netlify` and `deno`

- [#1786](https://github.com/withastro/starlight/pull/1786) [`d05d693a`](https://github.com/withastro/starlight/commit/d05d693afcf1771b8269dfe2ccc94f8952c643e8) Thanks [@delucis](https://github.com/delucis)! - Fixes type inference for i18n strings added by extending the default schema

- [#1777](https://github.com/withastro/starlight/pull/1777) [`6949404b`](https://github.com/withastro/starlight/commit/6949404b24a1c8254fd32d75122fdfbaf896fe29) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fixes an issue where TypeScript could fail to serialize the frontmatter schema when configured to emit declaration files

- [#1734](https://github.com/withastro/starlight/pull/1734) [`4493dcfa`](https://github.com/withastro/starlight/commit/4493dcfac5171f839b1b0e39444a15ce696adee4) Thanks [@delucis](https://github.com/delucis)! - Refactors `<ThemeSelect>` custom element logic to improve performance

- [#1731](https://github.com/withastro/starlight/pull/1731) [`f08b0dff`](https://github.com/withastro/starlight/commit/f08b0dff9638bbe7704ac2ba2e855d8d1464ba76) Thanks [@techfg](https://github.com/techfg)! - Fixes responding to system color scheme changes when theme is `auto`

- [#1793](https://github.com/withastro/starlight/pull/1793) [`2616f0c7`](https://github.com/withastro/starlight/commit/2616f0c7acf39a99c9f92b3db4108cae81120034) Thanks [@Mrahmani71](https://github.com/Mrahmani71)! - Updates the Farsi UI translations

## 0.21.5

### Patch Changes

- [#1728](https://github.com/withastro/starlight/pull/1728) [`0a75680d`](https://github.com/withastro/starlight/commit/0a75680ddd2f3325ab9ad7ac910f7c884b89a9ed) Thanks [@delucis](https://github.com/delucis)! - Adds 1 new icon: `pkl`

- [#1709](https://github.com/withastro/starlight/pull/1709) [`c5cd1811`](https://github.com/withastro/starlight/commit/c5cd181186b42422f3e47052bf8182cb490bda6b) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fixes a UI strings translation issue for sites configured with a single non-root language different from English.

- [#1723](https://github.com/withastro/starlight/pull/1723) [`3b29b3ab`](https://github.com/withastro/starlight/commit/3b29b3ab824f538f27e20310cb08786a92c7bd65) Thanks [@OliverSpeir](https://github.com/OliverSpeir)! - Fixes accessibility by using `aria-selected="false"` for inactive tabs instead of removing `aria-selected="true"` in the tablist of Starlight’s `<Tabs>` component

- [#1706](https://github.com/withastro/starlight/pull/1706) [`f171ac4d`](https://github.com/withastro/starlight/commit/f171ac4d6396eb2538598d85957670df50938b6a) Thanks [@jorenbroekema](https://github.com/jorenbroekema)! - Fixes some minor type errors

## 0.21.4

### Patch Changes

- [#1703](https://github.com/withastro/starlight/pull/1703) [`b26238f2`](https://github.com/withastro/starlight/commit/b26238f22990dcf8ba002bea6a50c66f20ad5786) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fixes aside custom titles rendering for nested asides.

- [#1708](https://github.com/withastro/starlight/pull/1708) [`a72cb966`](https://github.com/withastro/starlight/commit/a72cb96600798c1fbc7558f8fd24556ca442d312) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fixes translation issues with Expressive Code when using a default language other than English

## 0.21.3

### Patch Changes

- [#1622](https://github.com/withastro/starlight/pull/1622) [`3a074bad`](https://github.com/withastro/starlight/commit/3a074bad6c139bb9d6169d95ec79bc0fc1ecbdfe) Thanks [@SamuelLHuber](https://github.com/SamuelLHuber)! - Adds 1 new icon: `farcaster`

- [#1616](https://github.com/withastro/starlight/pull/1616) [`a86f9b71`](https://github.com/withastro/starlight/commit/a86f9b71b795fb6dcd0409ca568e43d25525b964) Thanks [@dragomano](https://github.com/dragomano)! - Updates Russian UI strings

- [#1698](https://github.com/withastro/starlight/pull/1698) [`67b892fd`](https://github.com/withastro/starlight/commit/67b892fd5290dfd0eeb95f4e60b6427bdc82110f) Thanks [@liruifengv](https://github.com/liruifengv)! - Adds 1 new icon: `starlight`

- [#1687](https://github.com/withastro/starlight/pull/1687) [`6fa9ea7e`](https://github.com/withastro/starlight/commit/6fa9ea7e8d4d601cf8f49b61dafb1ebb557d1718) Thanks [@mingjunlu](https://github.com/mingjunlu)! - Translates `fileTree.directory` UI string into Traditional Chinese.

## 0.21.2

### Patch Changes

- [#1628](https://github.com/withastro/starlight/pull/1628) [`24c0823c`](https://github.com/withastro/starlight/commit/24c0823c61b1e9850575766876f2e1035541cfd1) Thanks [@o-az](https://github.com/o-az)! - Adds 1 new icon: `nix`

- [#1614](https://github.com/withastro/starlight/pull/1614) [`78fc9042`](https://github.com/withastro/starlight/commit/78fc90426d58d6c36dcb8215e3181476d0702f50) Thanks [@kpodurgiel](https://github.com/kpodurgiel)! - Adds Polish UI translations

- [#1596](https://github.com/withastro/starlight/pull/1596) [`13ed30cd`](https://github.com/withastro/starlight/commit/13ed30cd335798177dfe24a27851d2c14d2fe80a) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds support for toggling the built-in search modal using the `Ctrl+k` keyboard shortcut.

- [#1608](https://github.com/withastro/starlight/pull/1608) [`4096e1b7`](https://github.com/withastro/starlight/commit/4096e1b77b3464338e5489d00cec4c29a1cd3c32) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Removes nested CSS from the `<FileTree>` component to prevent a potential warning when using Tailwind CSS.

- [#1626](https://github.com/withastro/starlight/pull/1626) [`67459cb4`](https://github.com/withastro/starlight/commit/67459cb4021859f4a45d50a5f993d2c849f340a3) Thanks [@hippotastic](https://github.com/hippotastic)! - Fixes a bundling issue that caused imports from `@astrojs/starlight/components` to fail when using the config setting `expressiveCode: false`.

## 0.21.1

### Patch Changes

- [#1584](https://github.com/withastro/starlight/pull/1584) [`8851d5cd`](https://github.com/withastro/starlight/commit/8851d5cd0d8f8439320ef729ca57a59418db52b9) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds 2 new icons: `apple` and `linux`.

- [#1577](https://github.com/withastro/starlight/pull/1577) [`0ba77890`](https://github.com/withastro/starlight/commit/0ba77890e0dcf54a849c735efd870327c10972aa) Thanks [@morinokami](https://github.com/morinokami)! - Translates `fileTree.directory` UI string into Japanese.

- [#1593](https://github.com/withastro/starlight/pull/1593) [`fa7ed245`](https://github.com/withastro/starlight/commit/fa7ed2458caf6261d16c5f43365cedbcb8572a48) Thanks [@liruifengv](https://github.com/liruifengv)! - Translates `fileTree.directory` UI string into simplified Chinese.

- [#1585](https://github.com/withastro/starlight/pull/1585) [`bd4e278f`](https://github.com/withastro/starlight/commit/bd4e278f7fe7d7335494602db29a63002fd45059) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Translates `fileTree.directory` UI string into French.

- [#1587](https://github.com/withastro/starlight/pull/1587) [`c5794260`](https://github.com/withastro/starlight/commit/c5794260251ed414a396089782a1788539c92dd3) Thanks [@Eveeifyeve](https://github.com/Eveeifyeve)! - Adds 1 new icon: `homebrew`.

## 0.21.0

### Minor Changes

- [#1568](https://github.com/withastro/starlight/pull/1568) [`5f99a71d`](https://github.com/withastro/starlight/commit/5f99a71ddfe92568b1cd3c0bfe5ebfd139797c1a) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds support for optionally setting an icon on a `<TabItem>` component to make it easier to visually distinguish between tabs.

- [#1308](https://github.com/withastro/starlight/pull/1308) [`9a918a5b`](https://github.com/withastro/starlight/commit/9a918a5b4902f43729f4d023257772710af3a12b) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds `<FileTree>` component to display the structure of a directory.

- [#1308](https://github.com/withastro/starlight/pull/1308) [`9a918a5b`](https://github.com/withastro/starlight/commit/9a918a5b4902f43729f4d023257772710af3a12b) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds 144 new file-type icons from the [Seti UI icon set](https://github.com/jesseweed/seti-ui#current-icons), available with the `seti:` prefix, e.g. `seti:javascript`.

- [#1564](https://github.com/withastro/starlight/pull/1564) [`d880065e`](https://github.com/withastro/starlight/commit/d880065e29a632823a08adcb6158a59fd9557270) Thanks [@delucis](https://github.com/delucis)! - Adds a `<Steps>` component for styling more complex guided tasks.

- [#1308](https://github.com/withastro/starlight/pull/1308) [`9a918a5b`](https://github.com/withastro/starlight/commit/9a918a5b4902f43729f4d023257772710af3a12b) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds 5 new icons: `astro`, `biome`, `bun`, `mdx`, and `pnpm`.

## 0.20.1

### Patch Changes

- [#1553](https://github.com/withastro/starlight/pull/1553) [`8e091147`](https://github.com/withastro/starlight/commit/8e09114755d37322d6e97b0dc90a5dfd781de8cc) Thanks [@hippotastic](https://github.com/hippotastic)! - Updates Expressive Code to v0.33.4 to fix potential race condition bug in Shiki.

## 0.20.0

### Minor Changes

- [#1541](https://github.com/withastro/starlight/pull/1541) [`1043052f`](https://github.com/withastro/starlight/commit/1043052f3890a577a73276472f3773924909406b) Thanks [@hippotastic](https://github.com/hippotastic)! - Updates `astro-expressive-code` dependency to the latest minor release (0.33).

  This unlocks support for [word wrap](https://expressive-code.com/key-features/word-wrap/) and [line numbers](https://expressive-code.com/plugins/line-numbers/), as well as updating the syntax highlighter to the latest Shiki release, which includes new and updated language grammars.

  See the [Expressive Code release notes](https://expressive-code.com/releases/) for more information including details of potentially breaking changes.

### Patch Changes

- [#1542](https://github.com/withastro/starlight/pull/1542) [`b3b7a606`](https://github.com/withastro/starlight/commit/b3b7a6069952d5f27a49b2fd097aa4db065e1718) Thanks [@delucis](https://github.com/delucis)! - Improves error messages shown by Starlight for configuration errors.

- [#1544](https://github.com/withastro/starlight/pull/1544) [`65dc6586`](https://github.com/withastro/starlight/commit/65dc6586ef7c1754875db1d48c49e709051a0b13) Thanks [@torn4dom4n](https://github.com/torn4dom4n)! - Update Vietnamese UI translations

## 0.19.1

### Patch Changes

- [#1527](https://github.com/withastro/starlight/pull/1527) [`163bc84`](https://github.com/withastro/starlight/commit/163bc848e173eecca92d1cb034045fdb42aa4ff1) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Exports the `StarlightPageProps` TypeScript type representing the props expected by the `<StarlightPage />` component.

- [#1504](https://github.com/withastro/starlight/pull/1504) [`fc83a05`](https://github.com/withastro/starlight/commit/fc83a05235b74be2bfe6ba8e7f95a8a5a618ead3) Thanks [@mingjunlu](https://github.com/mingjunlu)! - Adds Traditional Chinese UI translations

- [#1534](https://github.com/withastro/starlight/pull/1534) [`aada680`](https://github.com/withastro/starlight/commit/aada6805abc0068f07393585b86978ef5200439c) Thanks [@delucis](https://github.com/delucis)! - Improves DX of the `sidebar` prop used by the new `<StarlightPage>` component.

## 0.19.0

### Minor Changes

- [#1485](https://github.com/withastro/starlight/pull/1485) [`2cb3578`](https://github.com/withastro/starlight/commit/2cb35782dace67c7c418a31005419fa95493b3d3) Thanks [@timokoessler](https://github.com/timokoessler)! - Add support for setting html attributes of hero action links

- [#1175](https://github.com/withastro/starlight/pull/1175) [`dd11b95`](https://github.com/withastro/starlight/commit/dd11b9538abdf4b5ba2ef70e07c0edda03e95add) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds a new `<StarlightPage>` component to use the Starlight layout in custom pages.

  To learn more about this new feature, check out the new [“Using Starlight’s design in custom pages” guide](https://starlight.astro.build/guides/pages/#using-starlights-design-in-custom-pages).

- [#1499](https://github.com/withastro/starlight/pull/1499) [`97bf523`](https://github.com/withastro/starlight/commit/97bf523923fb9678c12f58fcdbe36757f0e56ceb) Thanks [@delucis](https://github.com/delucis)! - Adds a new `<Aside>` component

  The new component is in addition to the existing custom Markdown syntax.

## 0.18.1

### Patch Changes

- [#1487](https://github.com/withastro/starlight/pull/1487) [`6a72bda`](https://github.com/withastro/starlight/commit/6a72bda8c5569e2eda68fdf258ae9b1dc8b320d6) Thanks [@NavyStack](https://github.com/NavyStack)! - Improves Korean UI translations

- [#1489](https://github.com/withastro/starlight/pull/1489) [`b0d36de`](https://github.com/withastro/starlight/commit/b0d36de3398d4895603a787b612b1f0747defbdc) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fixes a potential text rendering issue with text containing colons.

## 0.18.0

### Minor Changes

- [#1454](https://github.com/withastro/starlight/pull/1454) [`1d9ef56`](https://github.com/withastro/starlight/commit/1d9ef567907bed7210e75ab3460f536c0768a87f) Thanks [@Fryuni](https://github.com/Fryuni)! - Makes Starlight compatible with [on-demand server rendering](https://docs.astro.build/en/guides/server-side-rendering/) (sometimes referred to as server-side rendering or SSR).

  Starlight pages are always prerendered, even when using `output: 'server'`.

- [#1454](https://github.com/withastro/starlight/pull/1454) [`1d9ef56`](https://github.com/withastro/starlight/commit/1d9ef567907bed7210e75ab3460f536c0768a87f) Thanks [@Fryuni](https://github.com/Fryuni)! - Enables Astro’s [`experimental.globalRoutePriority`](https://docs.astro.build/en/reference/configuration-reference/#experimentalglobalroutepriority) option and bumps the minimum required Astro version.

  ⚠️ **BREAKING CHANGE** The minimum supported Astro version is now 4.2.7. Upgrade Astro and Starlight together:

  ```sh
  npx @astrojs/upgrade
  ```

## 0.17.4

### Patch Changes

- [#1473](https://github.com/withastro/starlight/pull/1473) [`29da505`](https://github.com/withastro/starlight/commit/29da505474174fefaec4e27a2c2c3e90e3f68a31) Thanks [@delucis](https://github.com/delucis)! - Fixes a CSS bug for users with JavaScript disabled

- [#1465](https://github.com/withastro/starlight/pull/1465) [`ce3108c`](https://github.com/withastro/starlight/commit/ce3108cf6ecb77d12db973485d21e0fc7fd63ca6) Thanks [@delucis](https://github.com/delucis)! - Updates internal MDX, sitemap, and Expressive Code dependencies to the latest versions

## 0.17.3

### Patch Changes

- [#1461](https://github.com/withastro/starlight/pull/1461) [`2e17880`](https://github.com/withastro/starlight/commit/2e17880957d1aae2a84c77500afa9b66e5292a6a) Thanks [@liruifengv](https://github.com/liruifengv)! - Improves the table of contents title translation in Simplified Chinese

- [#1462](https://github.com/withastro/starlight/pull/1462) [`4741ccc`](https://github.com/withastro/starlight/commit/4741cccc8adbef500bcaf95416a1c61a90761c06) Thanks [@delucis](https://github.com/delucis)! - Fixes overflow of very long site titles on narrow viewports

- [#1459](https://github.com/withastro/starlight/pull/1459) [`9a8e0ec`](https://github.com/withastro/starlight/commit/9a8e0ec59cba0e088512ea9b6d17224085f3a178) Thanks [@delucis](https://github.com/delucis)! - Fixes a bug where table of contents highlighting could break given very specific combinations of content and viewport size

- [#1458](https://github.com/withastro/starlight/pull/1458) [`8c88642`](https://github.com/withastro/starlight/commit/8c88642875e8344396074a780e28fb0860b249f8) Thanks [@delucis](https://github.com/delucis)! - Silences i18n content collection warnings for projects without custom translations.

## 0.17.2

### Patch Changes

- [#1442](https://github.com/withastro/starlight/pull/1442) [`1a642e4`](https://github.com/withastro/starlight/commit/1a642e4d74ee4c30e85bce37b41888b1eae0544a) Thanks [@delucis](https://github.com/delucis)! - Fixes URLs in language picker for sites with `build.format: 'file'`

- [#1440](https://github.com/withastro/starlight/pull/1440) [`2ea1e88`](https://github.com/withastro/starlight/commit/2ea1e883186660b48f0ea8c4da7fead5fb74e313) Thanks [@hippotastic](https://github.com/hippotastic)! - Adds JS support to the `@astrojs/starlight/expressive-code` export to allow importing from non-TS environments.

## 0.17.1

### Patch Changes

- [#1437](https://github.com/withastro/starlight/pull/1437) [`655aed4`](https://github.com/withastro/starlight/commit/655aed4840cae59e9abd64b4b585e60f1cfab209) Thanks [@hippotastic](https://github.com/hippotastic)! - Adds Starlight-specific types to `defineEcConfig` function and exports `StarlightExpressiveCodeOptions`.

  This provides Starlight types and IntelliSense support for your Expressive Code configuration options inside an `ec.config.mjs` file. See the [Expressive Code documentation](https://expressive-code.com/key-features/code-component/#using-an-ecconfigmjs-file) for more information.

- [#1420](https://github.com/withastro/starlight/pull/1420) [`275f87f`](https://github.com/withastro/starlight/commit/275f87fd7fc676b9ab323354078c06894e0832c7) Thanks [@abdelhalimjean](https://github.com/abdelhalimjean)! - Fix rare `font-family` issue if users have a font installed with a name of `""`

- [#1365](https://github.com/withastro/starlight/pull/1365) [`a0af7cc`](https://github.com/withastro/starlight/commit/a0af7cc696da987a76edab96cdd2329779e87724) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Correctly format Pagefind search result links when `trailingSlash: 'never'` is used

## 0.17.0

### Minor Changes

- [#1389](https://github.com/withastro/starlight/pull/1389) [`21b3620`](https://github.com/withastro/starlight/commit/21b36201aa1e01c8395d0f24b2fa4e32b90550bb) Thanks [@connor-baer](https://github.com/connor-baer)! - Adds new `disable404Route` config option to disable injection of Astro’s default 404 route

- [#1395](https://github.com/withastro/starlight/pull/1395) [`ce05dfb`](https://github.com/withastro/starlight/commit/ce05dfb4b1e9b90fad057d5d4328e4445f986b3b) Thanks [@hippotastic](https://github.com/hippotastic)! - Adds a new [`<Code>` component](https://starlight.astro.build/guides/components/#code) to render dynamic code strings with Expressive Code

## 0.16.0

### Minor Changes

- [#1383](https://github.com/withastro/starlight/pull/1383) [`490c6ef`](https://github.com/withastro/starlight/commit/490c6eff34ab408c4f55777b7b0caa16787dd3d4) Thanks [@delucis](https://github.com/delucis)! - Refactors Starlight’s internal virtual module system for components to avoid circular references

  This is a change to an internal API.
  If you were importing the internal `virtual:starlight/components` module, this no longer exists.
  Update your imports to use the individual virtual modules now available for each component, for example `virtual:starlight/components/EditLink`.

- [#1151](https://github.com/withastro/starlight/pull/1151) [`134292d`](https://github.com/withastro/starlight/commit/134292ddd89683007d7de25545d39738a82c626c) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Fixes sidebar auto-generation issue when a file and a directory, located at the same level, have identical names.

  For example, `src/content/docs/guides.md` and `src/content/docs/guides/example.md` will now both be included and `src/content/docs/guides.md` is treated in the same way a `src/content/docs/guides/index.md` file would be.

- [#1386](https://github.com/withastro/starlight/pull/1386) [`0163634`](https://github.com/withastro/starlight/commit/0163634abb8578ce7a3d7ceea36432e98ea70e78) Thanks [@delucis](https://github.com/delucis)! - Tightens `line-height` on `<LinkCard>` titles to fix regression from original design

  If you want to preserve the previous `line-height`, you can add the following custom CSS to your site:

  ```css
  .sl-link-card a {
  	line-height: 1.6;
  }
  ```

- [#1376](https://github.com/withastro/starlight/pull/1376) [`8398432`](https://github.com/withastro/starlight/commit/8398432aa4a0f38e2dd4452dfcdf7033c5713334) Thanks [@delucis](https://github.com/delucis)! - Tweaks vertical spacing in Markdown content styles.

  This is a subtle change to Starlight’s default content styling that should improve most sites:

  - Default vertical spacing between content items is reduced from `1.5rem` to `1rem`.
  - Spacing before headings is now relative to font size, meaning higher-level headings have slightly more spacing and lower-level headings slightly less.

  The overall impact is to tighten up content that belongs together and improve the visual hierarchy of headings to break up sections.

  Although this is a subtle change, we recommend visually inspecting your site in case this impacts layout of any custom CSS or components.

  If you want to preserve the previous spacing, you can add the following custom CSS to your site:

  ```css
  /* Restore vertical spacing to match Starlight v0.15 and below. */
  .sl-markdown-content
  	:not(a, strong, em, del, span, input, code)
  	+ :not(a, strong, em, del, span, input, code, :where(.not-content *)) {
  	margin-top: 1.5rem;
  }
  .sl-markdown-content
  	:not(h1, h2, h3, h4, h5, h6)
  	+ :is(h1, h2, h3, h4, h5, h6):not(:where(.not-content *)) {
  	margin-top: 2.5rem;
  }
  ```

- [#1372](https://github.com/withastro/starlight/pull/1372) [`773880d`](https://github.com/withastro/starlight/commit/773880de87b79bf3107dbc32df29a86dd11e4e6f) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Updates the table of contents highlighting styles to prevent UI shifts when scrolling through a page.

  If you want to preserve the previous, buggy styling, you can add the following custom CSS to your site:

  ```css
  starlight-toc a[aria-current='true'],
  starlight-toc a[aria-current='true']:hover,
  starlight-toc a[aria-current='true']:focus {
  	font-weight: 600;
  	color: var(--sl-color-text-invert);
  	background-color: var(--sl-color-text-accent);
  }
  ```

## 0.15.4

### Patch Changes

- [#1378](https://github.com/withastro/starlight/pull/1378) [`0f4a31d`](https://github.com/withastro/starlight/commit/0f4a31da4b6d384c569e8556dcc559dc8bfbfebd) Thanks [@delucis](https://github.com/delucis)! - Updates dependencies: `@astrojs/mdx`, `@astrojs/sitemap`, and `astro-expressive-code`

## 0.15.3

### Patch Changes

- [#1303](https://github.com/withastro/starlight/pull/1303) [`3eefd21`](https://github.com/withastro/starlight/commit/3eefd21f2267648b17bc2d6874350fd5dd8bbcb2) Thanks [@lilnasy](https://github.com/lilnasy)! - chore: fix type errors in Starlight internals

- [#1351](https://github.com/withastro/starlight/pull/1351) [`932c022`](https://github.com/withastro/starlight/commit/932c0229d7d8d55f30161ccc36c908140c1f252a) Thanks [@roberto-butti](https://github.com/roberto-butti)! - Adds Italian translation for `search.devWarning` UI

- [#1298](https://github.com/withastro/starlight/pull/1298) [`c7e995c`](https://github.com/withastro/starlight/commit/c7e995cb018179789b5ee45bae5fdd9c20309945) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Fixes incorrect sorting behavior for some autogenerated sidebars

- [#1347](https://github.com/withastro/starlight/pull/1347) [`8994d00`](https://github.com/withastro/starlight/commit/8994d007266e0bd8e6116b306ccd9e24c9710411) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Refactor `getLastUpdated` to use `node:child_process` instead of `execa`.

- [#1353](https://github.com/withastro/starlight/pull/1353) [`90fe8da`](https://github.com/withastro/starlight/commit/90fe8da15c8eb227817c2232345ac359aef6bab5) Thanks [@delucis](https://github.com/delucis)! - Fixes sidebar scrollbar hiding behind navbar

## 0.15.2

### Patch Changes

- [#1254](https://github.com/withastro/starlight/pull/1254) [`e9659e8`](https://github.com/withastro/starlight/commit/e9659e869cd0c9ad0b7388397b0fff8e2a9db27a) Thanks [@Pukimaa](https://github.com/Pukimaa)! - Adds Open Collective social link icon

- [#1295](https://github.com/withastro/starlight/pull/1295) [`c3732a9`](https://github.com/withastro/starlight/commit/c3732a9bb5cb7907f00a3ed5e65534f48a5ff6b9) Thanks [@juchym](https://github.com/juchym)! - Improve Ukrainian UI translations

## 0.15.1

### Patch Changes

- [#1273](https://github.com/withastro/starlight/pull/1273) [`ae53155`](https://github.com/withastro/starlight/commit/ae531557aa4d42bd27c15f8f08bb3ca8242c9beb) Thanks [@natemoo-re](https://github.com/natemoo-re)! - Updates `<SocialIcon />` styling for improved accessibility. Specifically, the component now meets the [Target Size (Minimum)](https://www.w3.org/WAI/WCAG22/Understanding/target-size-minimum.html) success criteria defined by [Web Content Accessibility Guidelines (WCAG) 2.2](https://www.w3.org/TR/WCAG22/).

- [#1289](https://github.com/withastro/starlight/pull/1289) [`9bd343f`](https://github.com/withastro/starlight/commit/9bd343fb1efab90a0aa03a95b1928a53c1674000) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds French translations for Expressive Code UI

- [#1280](https://github.com/withastro/starlight/pull/1280) [`6b1693d`](https://github.com/withastro/starlight/commit/6b1693d55552a48316a31d986e1cbaf695f10a61) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Adds Spanish translations for Expressive Code UI

- [#1276](https://github.com/withastro/starlight/pull/1276) [`667f23d`](https://github.com/withastro/starlight/commit/667f23d615742b44bb18ace39d981f8797b8ac55) Thanks [@hippotastic](https://github.com/hippotastic)! - Updates `astro-expressive-code` dependency to the latest version

- [#1266](https://github.com/withastro/starlight/pull/1266) [`c9edf30`](https://github.com/withastro/starlight/commit/c9edf30b16f66757797dcaa5161b4afc18027476) Thanks [@alex-way](https://github.com/alex-way)! - Removes redundant subprocess calls in git last-updated time utility to improve performance

- [#1278](https://github.com/withastro/starlight/pull/1278) [`e88abb0`](https://github.com/withastro/starlight/commit/e88abb0cc8b329500c15bc77aaed3907ec7dc507) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Exports the `StarlightUserConfig` TypeScript type representing the user's Starlight configuration received by plugins.

## 0.15.0

### Minor Changes

- [#1238](https://github.com/withastro/starlight/pull/1238) [`02a808e`](https://github.com/withastro/starlight/commit/02a808e4a0b9ac2383576e3495f6a766b663d773) Thanks [@delucis](https://github.com/delucis)! - Add support for Astro v4, drop support for Astro v3

  ⚠️ **BREAKING CHANGE** Astro v3 is no longer supported. Make sure you [update Astro](https://docs.astro.build/en/guides/upgrade-to/v4/) and any other integrations at the same time as updating Starlight.

  Use the new `@astrojs/upgrade` command to upgrade Astro and Starlight together:

  ```sh
  npx @astrojs/upgrade
  ```

- [#1242](https://github.com/withastro/starlight/pull/1242) [`d8fc9e1`](https://github.com/withastro/starlight/commit/d8fc9e15bd2ae4c945b5a3856a6ce3b5629e8b29) Thanks [@delucis](https://github.com/delucis)! - Enables link prefetching on hover by default

  Astro v4’s [prefetch](https://docs.astro.build/en/guides/prefetch) support is now enabled by default. If `prefetch` is not set in `astro.config.mjs`, Starlight will use `prefetch: { prefetchAll: true, defaultStrategy: 'hover' }` by default.

  If you want to preserve previous behaviour, disable link prefetching in `astro.config.mjs`:

  ```js
  import { defineConfig } from 'astro/config';
  import starlight from '@astrojs/starlight';

  export default defineConfig({
  	// Disable link prefetching:
  	prefetch: false,

  	integrations: [
  		starlight({
  			// ...
  		}),
  	],
  });
  ```

### Patch Changes

- [#1226](https://github.com/withastro/starlight/pull/1226) [`909afa2`](https://github.com/withastro/starlight/commit/909afa2d468099e237bfbd25eda56270b7b00082) Thanks [@tlandmangh](https://github.com/tlandmangh)! - Add Dutch translations of default aside labels

- [#1243](https://github.com/withastro/starlight/pull/1243) [`ee234eb`](https://github.com/withastro/starlight/commit/ee234ebddcba8d07e2c879f33e38631c8955ffcf) Thanks [@khajimatov](https://github.com/khajimatov)! - Fix typo in Russian untranslated content notice

- [#1170](https://github.com/withastro/starlight/pull/1170) [`bcc2301`](https://github.com/withastro/starlight/commit/bcc2301c06796edec3923c666078e82eaf5a1990) Thanks [@tmcw](https://github.com/tmcw)! - Fix timezone-reliance in LastUpdated

- [#1203](https://github.com/withastro/starlight/pull/1203) [`4601449`](https://github.com/withastro/starlight/commit/4601449894bbbd619e4149788113090b67697fe1) Thanks [@orhun](https://github.com/orhun)! - Adds Matrix social link icon

## 0.14.0

### Minor Changes

- [#1144](https://github.com/withastro/starlight/pull/1144) [`7c0b8cb`](https://github.com/withastro/starlight/commit/7c0b8cb334c501678f7ab87cce372cddfdde34ed) Thanks [@delucis](https://github.com/delucis)! - Adds a configuration option to disable site indexing with Pagefind and the default search UI

- [#942](https://github.com/withastro/starlight/pull/942) [`efd7fdc`](https://github.com/withastro/starlight/commit/efd7fdcb55b39988f157c1a4b2c368c86a39520f) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Adds plugin API

  See the [plugins reference](https://starlight.astro.build/reference/plugins/) to learn more about creating plugins for Starlight using this new API.

- [#1135](https://github.com/withastro/starlight/pull/1135) [`e5a863a`](https://github.com/withastro/starlight/commit/e5a863a98b2e5335e122ca440dcb84e9426939b4) Thanks [@delucis](https://github.com/delucis)! - Exposes localized UI strings in route data

  Component overrides can now access a `labels` object in their props which includes all the localized UI strings for the current page.

- [#1162](https://github.com/withastro/starlight/pull/1162) [`00d101b`](https://github.com/withastro/starlight/commit/00d101b159bfa4bb307a66ccae53dd417d9564e0) Thanks [@delucis](https://github.com/delucis)! - Adds support for extending Starlight’s content collection schemas

## 0.13.1

### Patch Changes

- [#1111](https://github.com/withastro/starlight/pull/1111) [`cb19d07`](https://github.com/withastro/starlight/commit/cb19d07d6192ffb732ac6fcf9df04d4f098bfc1f) Thanks [@at-the-vr](https://github.com/at-the-vr)! - Fix minor punctuation typo in Hindi UI string

- [#1156](https://github.com/withastro/starlight/pull/1156) [`631c5ae`](https://github.com/withastro/starlight/commit/631c5aeccba60254ff649712f93ba30495775edf) Thanks [@votemike](https://github.com/votemike)! - Updates `@astrojs/sitemap` dependency to the latest version

- [#1109](https://github.com/withastro/starlight/pull/1109) [`0c25c1f`](https://github.com/withastro/starlight/commit/0c25c1f33bbfe311724784530c30ada44eb5de19) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Internal: fix import issue with expressive-code

## 0.13.0

### Minor Changes

- [#1023](https://github.com/withastro/starlight/pull/1023) [`a3b80f7`](https://github.com/withastro/starlight/commit/a3b80f71037504f2b8d7f1a641924215091122bb) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Respect the `trailingSlash` and `build.format` Astro options when creating Starlight navigation links.

  ⚠️ **Potentially breaking change:**
  This change will cause small changes in link formatting for most sites.
  These are unlikely to break anything, but if you care about link formatting, you may want to change some Astro settings.

  If you want to preserve Starlight’s previous behavior, set `trailingSlash: 'always'` in your `astro.config.mjs`:

  ```js
  import { defineConfig } from 'astro/config';
  import starlight from '@astrojs/starlight';

  export default defineConfig({
  	trailingSlash: 'always',
  	integrations: [
  		starlight({
  			// ...
  		}),
  	],
  });
  ```

- [#742](https://github.com/withastro/starlight/pull/742) [`c6a4bcb`](https://github.com/withastro/starlight/commit/c6a4bcb7982c54c513f20c96a9b2aaf9ac09094b) Thanks [@hippotastic](https://github.com/hippotastic)! - Adds Expressive Code as Starlight’s default code block renderer

  ⚠️ **Potentially breaking change:**
  This addition changes how Markdown code blocks are rendered. By default, Starlight will now use [Expressive Code](https://github.com/expressive-code/expressive-code/tree/main/packages/astro-expressive-code).
  If you were already customizing how code blocks are rendered and don't want to use the [features provided by Expressive Code](https://starlight.astro.build/guides/authoring-content/#expressive-code-features), you can preserve the previous behavior by setting the new config option `expressiveCode` to `false`.

  If you had previously added Expressive Code manually to your Starlight project, you can now remove the manual set-up in `astro.config.mjs`:

  - Move your configuration to Starlight’s new `expressiveCode` option.
  - Remove the `astro-expressive-code` integration.

  For example:

  ```diff
  import starlight from '@astrojs/starlight';
  import { defineConfig } from 'astro/config';
  - import expressiveCode from 'astro-expressive-code';

  export default defineConfig({
    integrations: [
  -   expressiveCode({
  -     themes: ['rose-pine'],
  -   }),
      starlight({
        title: 'My docs',
  +     expressiveCode: {
  +       themes: ['rose-pine'],
  +     },
      }),
    ],
  });
  ```

  Note that the built-in Starlight version of Expressive Code sets some opinionated defaults that are different from the `astro-expressive-code` defaults. You may need to set some `styleOverrides` if you wish to keep styles exactly the same.

- [#517](https://github.com/withastro/starlight/pull/517) [`5b549cb`](https://github.com/withastro/starlight/commit/5b549cb634f51d28bf9a7f92ad0d82c1671e788a) Thanks [@liruifengv](https://github.com/liruifengv)! - Add i18n support for default aside labels

### Patch Changes

- [#1088](https://github.com/withastro/starlight/pull/1088) [`4fe5537`](https://github.com/withastro/starlight/commit/4fe553749a6708fdb119b12a2dbc6b10a980bde1) Thanks [@Lootjs](https://github.com/Lootjs)! - i18n(ru): added Russian aside labels translation

- [#1083](https://github.com/withastro/starlight/pull/1083) [`e03a653`](https://github.com/withastro/starlight/commit/e03a65313365b7dbe6095727b28b4e639c446f68) Thanks [@at-the-vr](https://github.com/at-the-vr)! - i18n(hi): Add Hindi language support

- [#1075](https://github.com/withastro/starlight/pull/1075) [`2f2adf2`](https://github.com/withastro/starlight/commit/2f2adf29f2a13d5ff0f1577207210745a5ae7405) Thanks [@russbiggs](https://github.com/russbiggs)! - Add Slack social link icon

- [#1065](https://github.com/withastro/starlight/pull/1065) [`2d72ed6`](https://github.com/withastro/starlight/commit/2d72ed67c666b26eae44649e70aecef3db815d19) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Ignore search keyboard shortcuts for elements with contents that are editable

- [#1081](https://github.com/withastro/starlight/pull/1081) [`f27f781`](https://github.com/withastro/starlight/commit/f27f781556d37e73d0b1d902de745b67f8e4f24d) Thanks [@farisphp](https://github.com/farisphp)! - i18n(id): Add Indonesian aside labels translation

- [#1082](https://github.com/withastro/starlight/pull/1082) [`ce27486`](https://github.com/withastro/starlight/commit/ce27486fabd3884ed4bca9372ebd72a0597ab765) Thanks [@bogdaaamn](https://github.com/bogdaaamn)! - i18n(ro): Add Romanian UI translations

## 0.12.1

### Patch Changes

- [#1069](https://github.com/withastro/starlight/pull/1069) [`b86f360`](https://github.com/withastro/starlight/commit/b86f3608f03be9455ec1d5ba11820c9bf601ad1e) Thanks [@Genteure](https://github.com/Genteure)! - Fix sidebar highlighting and navigation buttons for pages with path containing non-ASCII characters

- [#1025](https://github.com/withastro/starlight/pull/1025) [`0d1e75e`](https://github.com/withastro/starlight/commit/0d1e75e17269ddac3eb15b7dfb4480da1bb01c6c) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Internal: fix import issue in translation string loading mechanism

- [#1044](https://github.com/withastro/starlight/pull/1044) [`a5a9754`](https://github.com/withastro/starlight/commit/a5a9754f111b97abfd277d99759e9857aa0fb22b) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fix last updated dates for pages displaying fallback content

- [#1049](https://github.com/withastro/starlight/pull/1049) [`c27495d`](https://github.com/withastro/starlight/commit/c27495da61f9376236519ed3f08a169f245a189c) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Expose Markdown content styles in `@astrojs/starlight/style/markdown.css`

## 0.12.0

### Minor Changes

- [#995](https://github.com/withastro/starlight/pull/995) [`5bf4457`](https://github.com/withastro/starlight/commit/5bf44577634935b9fa6d50b040abcd680035075f) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Adds support for adding sidebar badges to group headings

- [#988](https://github.com/withastro/starlight/pull/988) [`977fe13`](https://github.com/withastro/starlight/commit/977fe135a74661300589898abe98aec73cad9ed3) Thanks [@magicDGS](https://github.com/magicDGS)! - Include social icon links in mobile menu

- [#280](https://github.com/withastro/starlight/pull/280) [`72cca2d`](https://github.com/withastro/starlight/commit/72cca2d07644f00595da6ebf7d603adb282f359d) Thanks [@cbontems](https://github.com/cbontems)! - Support light & dark variants of the hero image.

  ⚠️ **Potentially breaking change:** The `hero.image` schema is now slightly stricter than previously.

  The `hero.image.html` property can no longer be used alongside the `hero.image.alt` or `hero.image.file` properties.
  Previously, `html` was ignored when used with `file` and `alt` was ignored when used with `html`.
  Now, those combinations will throw errors.
  If you encounter errors, remove the `image.hero` property that is not in use.

### Patch Changes

- [#1004](https://github.com/withastro/starlight/pull/1004) [`7f92213`](https://github.com/withastro/starlight/commit/7f92213a0b93de5a844816841a6bc9cdd371de0c) Thanks [@nunhes](https://github.com/nunhes)! - Add Galician language support

- [#1003](https://github.com/withastro/starlight/pull/1003) [`f1fdb50`](https://github.com/withastro/starlight/commit/f1fdb50daebe79548c7789d3f7dd968b261d2da7) Thanks [@delucis](https://github.com/delucis)! - Internal: refactor translation string loading to make translations available to Starlight integration code

## 0.11.2

### Patch Changes

- [#944](https://github.com/withastro/starlight/pull/944) [`7a6446e`](https://github.com/withastro/starlight/commit/7a6446ebc61ba9cc36d1dcfd13db15c1533751ab) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fix issue with sidebar autogenerated groups configured with a directory containing leading or trailing slash

- [#985](https://github.com/withastro/starlight/pull/985) [`92b3b57`](https://github.com/withastro/starlight/commit/92b3b575404d0dc34f720c0ba29d8ed50be98f58) Thanks [@delucis](https://github.com/delucis)! - Fix edit URLs for pages displaying fallback content

- [#986](https://github.com/withastro/starlight/pull/986) [`0470734`](https://github.com/withastro/starlight/commit/0470734e0dc323f9945e06bee4338c2f777ba0d6) Thanks [@dreyfus92](https://github.com/dreyfus92)! - Prevent overscrolling on mobile table of contents by setting 'overscroll-behavior: contain'.

- [#924](https://github.com/withastro/starlight/pull/924) [`39d6302`](https://github.com/withastro/starlight/commit/39d6302db42ee1105d690bbd3a66053e6b21e15a) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Remove extra margin from markdown lists that uses inline code

- [#814](https://github.com/withastro/starlight/pull/814) [`1e517d9`](https://github.com/withastro/starlight/commit/1e517d92a4cf146d2dcae58f7c4299d6f25ea73e) Thanks [@julien-deramond](https://github.com/julien-deramond)! - Prevent text from overflowing pagination items

## 0.11.1

### Patch Changes

- [#892](https://github.com/withastro/starlight/pull/892) [`2b30321`](https://github.com/withastro/starlight/commit/2b30321bde801bb9945d73dc954e25b40f4324fa) Thanks [@delucis](https://github.com/delucis)! - Add Patreon social link icon

- [#854](https://github.com/withastro/starlight/pull/854) [`71a52a1`](https://github.com/withastro/starlight/commit/71a52a16c44e3568128c83070541235133c44436) Thanks [@mehalter](https://github.com/mehalter)! - Add Reddit icon

- [#852](https://github.com/withastro/starlight/pull/852) [`344c92e`](https://github.com/withastro/starlight/commit/344c92e1b8bca5f92ec087df6cccf5c611eefdff) Thanks [@Lootjs](https://github.com/Lootjs)! - Improve Russian language support

- [#891](https://github.com/withastro/starlight/pull/891) [`395920c`](https://github.com/withastro/starlight/commit/395920c46e7b24cfff31800b3426ab375078e5c1) Thanks [@Frikadellios](https://github.com/Frikadellios)! - Add Ukrainian language support

- [#890](https://github.com/withastro/starlight/pull/890) [`63ea8e8`](https://github.com/withastro/starlight/commit/63ea8e86643b050c6be6f9a6167f6642b039c709) Thanks [@delucis](https://github.com/delucis)! - Update `execa` dependency to v8

- [#859](https://github.com/withastro/starlight/pull/859) [`eaa7a90`](https://github.com/withastro/starlight/commit/eaa7a902c7b7638b326709fd5203d932b20ed3fa) Thanks [@oggnimodd](https://github.com/oggnimodd)! - Improve Indonesian language support

- [#864](https://github.com/withastro/starlight/pull/864) [`b84aff2`](https://github.com/withastro/starlight/commit/b84aff2b9cccbc35c8619763c7f36841abe6344b) Thanks [@mehalter](https://github.com/mehalter)! - Optimize UI icon SVG paths

## 0.11.0

### Minor Changes

- [#774](https://github.com/withastro/starlight/pull/774) [`903a579`](https://github.com/withastro/starlight/commit/903a57942ceb99b68672c3fa54622b39cc5d76f8) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Support adding HTML attributes to sidebar links from config and frontmatter

- [#796](https://github.com/withastro/starlight/pull/796) [`372ec96`](https://github.com/withastro/starlight/commit/372ec96d31d0c1a9aa8bc1605de2b424bf9bd5af) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Add the `@astrojs/sitemap` and `@astrojs/mdx` integrations only if they are not detected in the Astro configuration.

  ⚠️ **BREAKING CHANGE** The minimum supported version of Astro is now v3.2.0. Make sure you update Astro at the same time as updating Starlight:

  ```sh
  npm install astro@latest
  ```

- [#447](https://github.com/withastro/starlight/pull/447) [`b45719b`](https://github.com/withastro/starlight/commit/b45719b581353f8d8f0ce0a9b5c89132e902377b) Thanks [@andremralves](https://github.com/andremralves)! - Add `titleDelimiter` configuration option and include site title in page `<title>` tags

  ⚠️ **BREAKING CHANGE** — Previously, every page’s `<title>` only included its individual frontmatter title.
  Now, `<title>` tags include the page title, a delimiter character (`|` by default), and the site title.
  For example, in the Startlight docs, `<title>Configuration Reference</title>` is now `<title>Configuration Reference | Starlight</title>`.

  If you have a page where you need to override this new behaviour, set a custom title using the `head` frontmatter property:

  ```md
  ---
  title: My Page
  head:
    - tag: title
      content: Custom Title
  ---
  ```

- [#709](https://github.com/withastro/starlight/pull/709) [`140e729`](https://github.com/withastro/starlight/commit/140e729a8bf12f805ae0b7e2b5ad959cf68d8e22) Thanks [@delucis](https://github.com/delucis)! - Add support for overriding Starlight’s built-in components

  ⚠️ **BREAKING CHANGE** — The page footer is now included on pages with `template: splash` in their frontmatter. Previously, this was not the case. If you are using `template: splash` and want to continue to hide footer elements, disable them in your frontmatter:

  ```md
  ---
  title: Landing page
  template: splash
  # Disable unwanted footer elements as needed
  editUrl: false
  lastUpdated: false
  prev: false
  next: false
  ---
  ```

  ⚠️ **BREAKING CHANGE** — This change involved refactoring the structure of some of Starlight’s built-in components slightly. If you were previously overriding these using other techniques, you may need to adjust your code.

### Patch Changes

- [#815](https://github.com/withastro/starlight/pull/815) [`b7b23a2`](https://github.com/withastro/starlight/commit/b7b23a2c90a25fe8ea08338379b83d19c74d9037) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Add Facebook and email icons

- [#810](https://github.com/withastro/starlight/pull/810) [`dbe977b`](https://github.com/withastro/starlight/commit/dbe977b6ce3efcffefab850eca08bef316b41e53) Thanks [@hasham-qaiser](https://github.com/hasham-qaiser)! - Use `<span>` instead of `<h2>` in sidebar group headings

- [#807](https://github.com/withastro/starlight/pull/807) [`7c73dd1`](https://github.com/withastro/starlight/commit/7c73dd146ee294f9092346a0b0041990cc648a13) Thanks [@torn4dom4n](https://github.com/torn4dom4n)! - Add Vietnamese translations for Starlight UI

- [#756](https://github.com/withastro/starlight/pull/756) [`f55a8f0`](https://github.com/withastro/starlight/commit/f55a8f014a7addc46e971dd6b7148f4545acd16c) Thanks [@julien-deramond](https://github.com/julien-deramond)! - Prevent text from overflowing in several cases

## 0.10.4

### Patch Changes

- [#752](https://github.com/withastro/starlight/pull/752) [`6833ee1`](https://github.com/withastro/starlight/commit/6833ee12159ea5be23885c41da80569a89aafa33) Thanks [@apinet](https://github.com/apinet)! - Add X social link logo

- [#789](https://github.com/withastro/starlight/pull/789) [`2528fb0`](https://github.com/withastro/starlight/commit/2528fb011a388f3920af9994012bd7db2b1654c3) Thanks [@delucis](https://github.com/delucis)! - Update bundled version of `@astrojs/mdx` to v1.1.0

- [#794](https://github.com/withastro/starlight/pull/794) [`a0de12d`](https://github.com/withastro/starlight/commit/a0de12d596ab2b1c0d79d6d63d2d6cb9fe6d2644) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Add Telegram icon

- [#792](https://github.com/withastro/starlight/pull/792) [`a8358df`](https://github.com/withastro/starlight/commit/a8358df7849b342de342a4a2e88e019c39d5bbe8) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Add RSS icon

- [#778](https://github.com/withastro/starlight/pull/778) [`957d2c3`](https://github.com/withastro/starlight/commit/957d2c32123ab0c223fcef03f91ce3ad7be7aa21) Thanks [@jermanuts](https://github.com/jermanuts)! - Improve Arabic UI translations

## 0.10.3

### Patch Changes

- [#783](https://github.com/withastro/starlight/pull/783) [`f94727e`](https://github.com/withastro/starlight/commit/f94727e7d286a6910f913a572b27eb17c42f1729) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Fix GitHub edit link to include src path from project config

- [#781](https://github.com/withastro/starlight/pull/781) [`a293ef9`](https://github.com/withastro/starlight/commit/a293ef9ebb10a07db456156d8bdacc4ff6a2ca38) Thanks [@dreyfus92](https://github.com/dreyfus92)! - Removed role from Banner component to avoid duplication in header.

- [#745](https://github.com/withastro/starlight/pull/745) [`006d606`](https://github.com/withastro/starlight/commit/006d60695761ec10e5c4e715ed2212cd1fbedda0) Thanks [@TheOtterlord](https://github.com/TheOtterlord)! - Prevent Starlight crashing when the content folder doesn't exist, or is empty

- [#775](https://github.com/withastro/starlight/pull/775) [`2ef3036`](https://github.com/withastro/starlight/commit/2ef303649a0b66a6ec6a216815e05d41bf22b594) Thanks [@delucis](https://github.com/delucis)! - Fix content collection schema compatibility with Astro 3.1 and higher

- [#773](https://github.com/withastro/starlight/pull/773) [`423d575`](https://github.com/withastro/starlight/commit/423d575cc8227e4db86a85c70c45c0f3f7a184d2) Thanks [@tlandmangh](https://github.com/tlandmangh)! - Fix Dutch UI translation for “Previous page” links

## 0.10.2

### Patch Changes

- [#735](https://github.com/withastro/starlight/pull/735) [`2da8692`](https://github.com/withastro/starlight/commit/2da86929c8041f6585790c3baf1cba42220650cc) Thanks [@delucis](https://github.com/delucis)! - Use Starlight font custom property in Pagefind modal

- [#735](https://github.com/withastro/starlight/pull/735) [`2da8692`](https://github.com/withastro/starlight/commit/2da86929c8041f6585790c3baf1cba42220650cc) Thanks [@delucis](https://github.com/delucis)! - Fix RTL styling in Pagefind modal

- [#739](https://github.com/withastro/starlight/pull/739) [`a9de4a7`](https://github.com/withastro/starlight/commit/a9de4a7dcf9ec8c5c801e8a6cbb0d7faf2c34db7) Thanks [@radenpioneer](https://github.com/radenpioneer)! - Add Indonesian UI translation

- [#747](https://github.com/withastro/starlight/pull/747) [`7589515`](https://github.com/withastro/starlight/commit/75895154b11cf9368d4d6b45647b156ce32a88f0) Thanks [@nirtamir2](https://github.com/nirtamir2)! - Add Hebrew UI translations

## 0.10.1

### Patch Changes

- [#726](https://github.com/withastro/starlight/pull/726) [`f3157c6`](https://github.com/withastro/starlight/commit/f3157c6065943af39995b6dbae5f63cf424bd9a3) Thanks [@delucis](https://github.com/delucis)! - Fix a rare bug in table of contents when handling headings that increase by more than one level on a page.

- [#729](https://github.com/withastro/starlight/pull/729) [`80c6ab1`](https://github.com/withastro/starlight/commit/80c6ab1c1ec48805e74c53b615a78d65127eeacb) Thanks [@delucis](https://github.com/delucis)! - Upgrade Pagefind to v1.0.3

- [#715](https://github.com/withastro/starlight/pull/715) [`e726155`](https://github.com/withastro/starlight/commit/e7261559f2539a0ceefd36a28e4fbbc17f5970b8) Thanks [@itsmatteomanf](https://github.com/itsmatteomanf)! - feat: prevent scroll on body when search is open

## 0.10.0

### Minor Changes

- [#692](https://github.com/withastro/starlight/pull/692) [`2a58e1a`](https://github.com/withastro/starlight/commit/2a58e1aa068d01833a0ab9e74e4b46cccaee1775) Thanks [@delucis](https://github.com/delucis)! - Upgrade Pagefind to v1 and display page headings in search results

### Patch Changes

- [#708](https://github.com/withastro/starlight/pull/708) [`136cfb1`](https://github.com/withastro/starlight/commit/136cfb180f22db116cfdb62fd93d21daff596946) Thanks [@julien-deramond](https://github.com/julien-deramond)! - Fix main content column width for pages without a table of contents

- [#682](https://github.com/withastro/starlight/pull/682) [`660a5f5`](https://github.com/withastro/starlight/commit/660a5f57adf0340de21df3e364aada38255bb06c) Thanks [@vedmalex](https://github.com/vedmalex)! - Add Russian language support

## 0.9.1

### Patch Changes

- [#647](https://github.com/withastro/starlight/pull/647) [`ea57726`](https://github.com/withastro/starlight/commit/ea5772655274a3900310cb700836fdd2f6dba7cd) Thanks [@bgmort](https://github.com/bgmort)! - Fix translated 404 pages not being excluded from search results

- [#667](https://github.com/withastro/starlight/pull/667) [`9828f73`](https://github.com/withastro/starlight/commit/9828f739b73e2f377c1450b9e11f0914722ee440) Thanks [@delucis](https://github.com/delucis)! - Break inline `<code>` across lines to avoid overflow

- [#642](https://github.com/withastro/starlight/pull/642) [`e623d92`](https://github.com/withastro/starlight/commit/e623d92c2fddc0ff5fe83d2554266885d683a906) Thanks [@fk](https://github.com/fk)! - Don't hard-code nav height in table of contents highlighting script

- [#676](https://github.com/withastro/starlight/pull/676) [`6419006`](https://github.com/withastro/starlight/commit/641900615aa9a9a128d6934e65a57ba89e503cfd) Thanks [@vedmalex](https://github.com/vedmalex)! - Upgrade and pin Pagefind to latest beta release.

- [#647](https://github.com/withastro/starlight/pull/647) [`ea57726`](https://github.com/withastro/starlight/commit/ea5772655274a3900310cb700836fdd2f6dba7cd) Thanks [@bgmort](https://github.com/bgmort)! - Add frontmatter option to exclude a page from Pagefind search results

## 0.9.0

### Minor Changes

- [#626](https://github.com/withastro/starlight/pull/626) [`5dd22b8`](https://github.com/withastro/starlight/commit/5dd22b875dc19a32c48692082fbd934e2b70da63) Thanks [@delucis](https://github.com/delucis)! - Throw an error for duplicate MDX or sitemap integrations

- [#615](https://github.com/withastro/starlight/pull/615) [`7b75b3e`](https://github.com/withastro/starlight/commit/7b75b3eb7e6f7870a0adef2d6534ff48309fdb0e) Thanks [@delucis](https://github.com/delucis)! - Bump minimum required Astro version to 3.0

  ⚠️ **BREAKING CHANGE** Astro v2 is no longer supported. Make sure you [update Astro](https://docs.astro.build/en/guides/upgrade-to/v3/) and any other integrations at the same time as updating Starlight.

## 0.8.1

### Patch Changes

- [#612](https://github.com/withastro/starlight/pull/612) [`1b367e3`](https://github.com/withastro/starlight/commit/1b367e3f65e3736b5f91c9853a487f7f5d174a6f) Thanks [@KubaJastrz](https://github.com/KubaJastrz)! - Avoid applying hovered `<select>` text color to its `<options>`

## 0.8.0

### Minor Changes

- [#529](https://github.com/withastro/starlight/pull/529) [`c2d0e7f`](https://github.com/withastro/starlight/commit/c2d0e7f2699e60a48a3a9074eee6439dee8624a1) Thanks [@delucis](https://github.com/delucis)! - For improved compatibility with Tailwind, some Starlight built-in class names are now prefixed with `"sl-"`.

  While not likely, if you were relying on one of these internal class names in your own components or custom CSS, you will need to update to use the prefixed version.

  - **Before:** `flex`, `md:flex`, `lg:flex`, `block`, `md:block`, `lg:block`, `hidden`, `md:hidden`, `lg:hidden`.
  - **After:** `sl-flex`, `md:sl-flex`, `lg:sl-flex`, `sl-block`, `md:sl-block`, `lg:sl-block`, `sl-hidden`, `md:sl-hidden`, `lg:sl-hidden`.

- [#593](https://github.com/withastro/starlight/pull/593) [`5b8af95`](https://github.com/withastro/starlight/commit/5b8af95049781954eabc3895027218b3de8ff054) Thanks [@delucis](https://github.com/delucis)! - Add announcement banner feature

- [#516](https://github.com/withastro/starlight/pull/516) [`70a32a1`](https://github.com/withastro/starlight/commit/70a32a1736c776febb34cf0ca3014f375ff9fec8) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Support adding badges to sidebar links from config file and frontmatter

### Patch Changes

- [#569](https://github.com/withastro/starlight/pull/569) [`a7691f8`](https://github.com/withastro/starlight/commit/a7691f82fdabb1c4e6b14bcfa8289aaceb929997) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Use locale when sorting autogenerated sidebar entries

- [#559](https://github.com/withastro/starlight/pull/559) [`5726353`](https://github.com/withastro/starlight/commit/5726353f1a4815df9ab5a7acd7aea6af53383adc) Thanks [@delucis](https://github.com/delucis)! - Add Stack Overflow icon

## 0.7.3

### Patch Changes

- [#525](https://github.com/withastro/starlight/pull/525) [`87caf21`](https://github.com/withastro/starlight/commit/87caf21adaac98cf8342dac4db97ace327849616) Thanks [@delucis](https://github.com/delucis)! - Improve inline code and code block support in RTL languages

- [#537](https://github.com/withastro/starlight/pull/537) [`56c19bc`](https://github.com/withastro/starlight/commit/56c19bc871f2a4f205d4b0bb833fd81a3ed2e0f0) Thanks [@carlgleisner](https://github.com/carlgleisner)! - Add Swedish UI translations.

- [#528](https://github.com/withastro/starlight/pull/528) [`f5e5503`](https://github.com/withastro/starlight/commit/f5e55036987db98dbd0be7a84eb7819a48234a2f) Thanks [@jsparkdev](https://github.com/jsparkdev)! - add Korean language support

## 0.7.2

### Patch Changes

- [#506](https://github.com/withastro/starlight/pull/506) [`5e3133c`](https://github.com/withastro/starlight/commit/5e3133c42232b201b981cf4b3bc1c3dd56b09fa5) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Improve table of content current item highlight behavior

- [#499](https://github.com/withastro/starlight/pull/499) [`fcff49e`](https://github.com/withastro/starlight/commit/fcff49ee4260ad68e80833712e161cbb978a2562) Thanks [@D3vil0p3r](https://github.com/D3vil0p3r)! - Add icons for Instagram

- [#502](https://github.com/withastro/starlight/pull/502) [`3c87a16`](https://github.com/withastro/starlight/commit/3c87a16de3c867ad89294a0ea84d63eca2e74d7a) Thanks [@Mrahmani71](https://github.com/Mrahmani71)! - Add Farsi UI translations

- [#496](https://github.com/withastro/starlight/pull/496) [`cd28392`](https://github.com/withastro/starlight/commit/cd28392ac73ac0ba1a441328fcd1d65d7d441366) Thanks [@lorenzolewis](https://github.com/lorenzolewis)! - Fix `lastUpdated` date position to be consistent

- [#402](https://github.com/withastro/starlight/pull/402) [`d8669b8`](https://github.com/withastro/starlight/commit/d8669b869761ac15d1d611eda7dd94a62ce0fd7a) Thanks [@chopfitzroy](https://github.com/chopfitzroy)! - Fix content sometimes appearing above the mobile table of contents.

## 0.7.1

### Patch Changes

- [#488](https://github.com/withastro/starlight/pull/488) [`da35556`](https://github.com/withastro/starlight/commit/da35556eb95f2d397dfce03cc4acfacb0dcf1e89) Thanks [@mayank99](https://github.com/mayank99)! - Improved accessibility of LinkCard by only including the title as part of the link text, and using a pseudo-element to keep the card clickable.

- [#489](https://github.com/withastro/starlight/pull/489) [`35cd82e`](https://github.com/withastro/starlight/commit/35cd82e7f8622772a5155add99ad8baf61ae08a1) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Respect `hidden` sidebar frontmatter property when no sidebar configuration is provided

## 0.7.0

### Minor Changes

- [#441](https://github.com/withastro/starlight/pull/441) [`0119a49`](https://github.com/withastro/starlight/commit/0119a49b9a5f7844e7689df5577e8132bf871535) Thanks [@lorenzolewis](https://github.com/lorenzolewis)! - Add support for hiding entries from an autogenerated sidebar:

  ```md
  ---
  title: About this project
  sidebar:
    hidden: true
  ---
  ```

- [#470](https://github.com/withastro/starlight/pull/470) [`d076aec`](https://github.com/withastro/starlight/commit/d076aec856921c2fe8a5204a0c31580a846af180) Thanks [@delucis](https://github.com/delucis)! - Drop support for the `--sl-hue-accent` CSS custom property.

  ⚠️ **BREAKING CHANGE** — In previous Starlight versions you could control the accent color by setting the `--sl-hue-accent` custom property. This could result in inaccessible color contrast and unpredictable results.

  You must now set accent colors directly. If you relied on setting `--sl-hue-accent`, migrate by setting light and dark mode colors in your custom CSS:

  ```css
  :root {
  	--sl-hue-accent: 234;
  	--sl-color-accent-low: hsl(var(--sl-hue-accent), 54%, 20%);
  	--sl-color-accent: hsl(var(--sl-hue-accent), 100%, 60%);
  	--sl-color-accent-high: hsl(var(--sl-hue-accent), 100%, 87%);
  }

  :root[data-theme='light'] {
  	--sl-color-accent-high: hsl(var(--sl-hue-accent), 80%, 30%);
  	--sl-color-accent: hsl(var(--sl-hue-accent), 90%, 60%);
  	--sl-color-accent-low: hsl(var(--sl-hue-accent), 88%, 90%);
  }
  ```

  The [new color theme editor](https://starlight.astro.build/guides/css-and-tailwind/#color-theme-editor) might help if you’d prefer to set a new color scheme.

- [#397](https://github.com/withastro/starlight/pull/397) [`73eb5e6`](https://github.com/withastro/starlight/commit/73eb5e6ac6511dc4a6f5c4ca6c0c60d521f1db3c) Thanks [@lorenzolewis](https://github.com/lorenzolewis)! - Add `LinkCard` component

### Patch Changes

- [#460](https://github.com/withastro/starlight/pull/460) [`2e0fb90`](https://github.com/withastro/starlight/commit/2e0fb9053e96839287071e8a9c523796570cb0f6) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fix current page highlight in sidebar for URLs with no trailing slash

- [#467](https://github.com/withastro/starlight/pull/467) [`461a5d5`](https://github.com/withastro/starlight/commit/461a5d5c0424b03fb95b7ff7b27c944d04430244) Thanks [@delucis](https://github.com/delucis)! - Fix type error for downstream `tsc` users

- [#475](https://github.com/withastro/starlight/pull/475) [`06a205e`](https://github.com/withastro/starlight/commit/06a205e0e673f505bbb87dfcfcb0f35b051677e9) Thanks [@Yan-Thomas](https://github.com/Yan-Thomas)! - Locales whose language tag includes a regional subtag now use built-in UI translations for their base language. For example, a locale with a language of `pt-BR` will use our `pt` UI translations.

- [#473](https://github.com/withastro/starlight/pull/473) [`6a7692a`](https://github.com/withastro/starlight/commit/6a7692ae3178f9f9f727cc17b8ae860604afd78f) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fix issue with nested `<Tabs>` components

## 0.6.1

### Patch Changes

- [#442](https://github.com/withastro/starlight/pull/442) [`42c0abd`](https://github.com/withastro/starlight/commit/42c0abdd245f2f6595d67e203965f463829ef870) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Increase Markdown table border contrast

- [#443](https://github.com/withastro/starlight/pull/443) [`cb8bcec`](https://github.com/withastro/starlight/commit/cb8bcec533c9a7849eda01a4a4157b4726c9902c) Thanks [@delucis](https://github.com/delucis)! - Add icons for Bitbucket, Gitter, CodePen, and Microsoft Teams

- [#445](https://github.com/withastro/starlight/pull/445) [`a80e180`](https://github.com/withastro/starlight/commit/a80e180ca5abb85aa0c9db111ef5ae8e0c1bb539) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Prevent repeated table of contents mark on mobile

## 0.6.0

### Minor Changes

- [#424](https://github.com/withastro/starlight/pull/424) [`4485d90`](https://github.com/withastro/starlight/commit/4485d90fbddf7c9458b43f9d9b7560b41ec9e98f) Thanks [@delucis](https://github.com/delucis)! - Add support for customising autogenerated sidebar link labels from page frontmatter, overriding the page title:

  ```md
  ---
  title: About this project
  sidebar:
    label: About
  ---
  ```

- [#359](https://github.com/withastro/starlight/pull/359) [`e733311`](https://github.com/withastro/starlight/commit/e73331133b0e2574a139409ba76d97cc1bd52a82) Thanks [@IDurward](https://github.com/IDurward)! - Add support for defining the order of auto-generated link groups in the sidebar using a frontmatter value:

  ```md
  ---
  title: Page to display first
  sidebar:
    order: 1
  ---
  ```

### Patch Changes

- [#413](https://github.com/withastro/starlight/pull/413) [`5a9d8f1`](https://github.com/withastro/starlight/commit/5a9d8f11d59bd48322a1f2ff90e68333c3207ee1) Thanks [@delucis](https://github.com/delucis)! - Fix site title overflow bug for longer titles on narrow screens

- [#381](https://github.com/withastro/starlight/pull/381) [`6e62909`](https://github.com/withastro/starlight/commit/6e629095e78da4bfd422cd0a9cd9beb0d85d9a1a) Thanks [@lorenzolewis](https://github.com/lorenzolewis)! - Preserve order of `social` config in navbar

- [#419](https://github.com/withastro/starlight/pull/419) [`38ff53c`](https://github.com/withastro/starlight/commit/38ff53c216898efaa8c07394500e82da1d68ee8a) Thanks [@lorenzolewis](https://github.com/lorenzolewis)! - Improve styling of sidebar entries that wrap onto multiple lines

- [#418](https://github.com/withastro/starlight/pull/418) [`c7b2a4e`](https://github.com/withastro/starlight/commit/c7b2a4e9c8c55564be75f0c0901e38577ac764ec) Thanks [@delucis](https://github.com/delucis)! - Set `tab-size: 2` on content code blocks to override default browser value of `8`

- [#399](https://github.com/withastro/starlight/pull/399) [`31b8a5a`](https://github.com/withastro/starlight/commit/31b8a5aed2bca363c1b05c683b020e596b70bf4a) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Add new global `favicon` option defaulting to `'/favicon.svg'` to set the path of the default favicon for your website. Additional icons can be specified using the `head` option.

- [#414](https://github.com/withastro/starlight/pull/414) [`e951671`](https://github.com/withastro/starlight/commit/e95167174e3eab3790328b8e42517abcbca04ff3) Thanks [@delucis](https://github.com/delucis)! - Add GitLab to social link icons

## 0.5.6

### Patch Changes

- [#383](https://github.com/withastro/starlight/pull/383) [`0ebc47e`](https://github.com/withastro/starlight/commit/0ebc47e52dc420240c8cb724c01f98dc22bdfc60) Thanks [@delucis](https://github.com/delucis)! - Fix edge case where index files in an index directory would end up with the wrong slug

- [#373](https://github.com/withastro/starlight/pull/373) [`308b3aa`](https://github.com/withastro/starlight/commit/308b3aaeb0122af81a12514a81e160910f93d7a7) Thanks [@lorenzolewis](https://github.com/lorenzolewis)! - Fix visual overflow for wide logos

- [#385](https://github.com/withastro/starlight/pull/385) [`fb35397`](https://github.com/withastro/starlight/commit/fb35397f107f7bbb2cb4929b7837f105f565a659) Thanks [@lorenzolewis](https://github.com/lorenzolewis)! - Fix nested elements in markdown content

- [#386](https://github.com/withastro/starlight/pull/386) [`e6f6f30`](https://github.com/withastro/starlight/commit/e6f6f304437203d5ee6770092ac79063448b821f) Thanks [@huijing](https://github.com/huijing)! - Prevent search keyboard shortcuts from triggering when input elements are focused

## 0.5.5

### Patch Changes

- [`a161c05`](https://github.com/withastro/starlight/commit/a161c05b74d2300c1fe49bfd8e111cc45c9a5bff) Thanks [@delucis](https://github.com/delucis)! - Fix missing metadata required for `astro add` support

## 0.5.4

### Patch Changes

- [#360](https://github.com/withastro/starlight/pull/360) [`8415df6`](https://github.com/withastro/starlight/commit/8415df63e502d517b68d7665d9257726e3dde246) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Fix build warnings when using the TypeScript [`verbatimModuleSyntax`](https://www.typescriptlang.org/tsconfig#verbatimModuleSyntax) compiler option

## 0.5.3

### Patch Changes

- [#352](https://github.com/withastro/starlight/pull/352) [`a2e23be`](https://github.com/withastro/starlight/commit/a2e23be71d9f3592a9ac615981233bf4e9f3af6b) Thanks [@TheOtterlord](https://github.com/TheOtterlord)! - Fix page scrolling when the window resizes, while the mobile nav is open

- [#353](https://github.com/withastro/starlight/pull/353) [`65b2b75`](https://github.com/withastro/starlight/commit/65b2b7561a185be29ff7f773bf9432dc3c4da2e4) Thanks [@liruifengv](https://github.com/liruifengv)! - Add Simplified Chinese language support

## 0.5.2

### Patch Changes

- [#343](https://github.com/withastro/starlight/pull/343) [`d618678`](https://github.com/withastro/starlight/commit/d618678b1901c621e1c8d2dc1a34ee299582b14e) Thanks [@delucis](https://github.com/delucis)! - Fix escaping of non-relative user config file paths for custom CSS and logos

## 0.5.1

### Patch Changes

- [#336](https://github.com/withastro/starlight/pull/336) [`2b3302b`](https://github.com/withastro/starlight/commit/2b3302b80451f318fb05a5e8a7284feb28999e66) Thanks [@delucis](https://github.com/delucis)! - Add support for LinkedIn, Threads, and Twitch social icon links

- [#335](https://github.com/withastro/starlight/pull/335) [`757c65f`](https://github.com/withastro/starlight/commit/757c65ffc468fd2c782312b476fa7659d0cfd198) Thanks [@delucis](https://github.com/delucis)! - Fix relative path resolution on Windows

- [#332](https://github.com/withastro/starlight/pull/332) [`0600c1a`](https://github.com/withastro/starlight/commit/0600c1a917bf86efa6b2d053aa47e3a4b17e8049) Thanks [@sasoria](https://github.com/sasoria)! - Add Norwegian UI translations

- [#328](https://github.com/withastro/starlight/pull/328) [`e478848`](https://github.com/withastro/starlight/commit/e478848de1c41a46f58d0ac0d62d7b7272cf1241) Thanks [@astridx](https://github.com/astridx)! - Add missing accessible labels for Codeberg and YouTube social links

## 0.5.0

### Minor Changes

- [#313](https://github.com/withastro/starlight/pull/313) [`dc42569`](https://github.com/withastro/starlight/commit/dc42569bddfae2c48ea60c0dd5cc70643a129a68) Thanks [@delucis](https://github.com/delucis)! - Add a `not-content` CSS class that allows users to opt out of Starlight’s default content styling

- [#297](https://github.com/withastro/starlight/pull/297) [`fb15a9b`](https://github.com/withastro/starlight/commit/fb15a9b65252ac5fa32304096fbdb49ecdd6009b) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Improve `<Tabs>` component keyboard interactions

- [#303](https://github.com/withastro/starlight/pull/303) [`69b7d4c`](https://github.com/withastro/starlight/commit/69b7d4c23761a45dc2b9ea75c6c9c904a885ba5d) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Add new global `pagination` option defaulting to `true` to define whether or not the previous and next page links are shown in the footer. A page can override this setting or the link text and/or URL using the new `prev` and `next` frontmatter fields.

### Patch Changes

- [#318](https://github.com/withastro/starlight/pull/318) [`5db3e6e`](https://github.com/withastro/starlight/commit/5db3e6ea2e5cb7d9552fc54567811358851fb533) Thanks [@delucis](https://github.com/delucis)! - Support relative paths in Starlight config for `customCSS` and `logo` paths

## 0.4.2

### Patch Changes

- [#308](https://github.com/withastro/starlight/pull/308) [`c3aa4c6`](https://github.com/withastro/starlight/commit/c3aa4c6aa18f7f6859ad1c0acc28f0da59a84760) Thanks [@delucis](https://github.com/delucis)! - Fix use of default monospace font stack

- [#286](https://github.com/withastro/starlight/pull/286) [`a2aedfc`](https://github.com/withastro/starlight/commit/a2aedfc7f9555b44f5f33aad7f4a98b207a11b47) Thanks [@mzaien](https://github.com/mzaien)! - Add Arabic UI translations

## 0.4.1

### Patch Changes

- [#300](https://github.com/withastro/starlight/pull/300) [`377a25d`](https://github.com/withastro/starlight/commit/377a25dc4c51c060e751aeba4d3f946a41de907a) Thanks [@cbontems](https://github.com/cbontems)! - Fix broken link on 404 page when `defaultLocale: 'root'` is set in `astro.config.mjs`

- [#289](https://github.com/withastro/starlight/pull/289) [`dffca46`](https://github.com/withastro/starlight/commit/dffca461633940847e9177913053885c5e8b5f29) Thanks [@RyanRBrown](https://github.com/RyanRBrown)! - Fix saturation of purple text in light theme

- [#301](https://github.com/withastro/starlight/pull/301) [`d47639d`](https://github.com/withastro/starlight/commit/d47639d50b53fa691c1d9b0f30f82ebf7f6ddf7e) Thanks [@delucis](https://github.com/delucis)! - Enable inline stylesheets for Astro versions ≥2.6.0

## 0.4.0

### Minor Changes

- [#259](https://github.com/withastro/starlight/pull/259) [`8102389`](https://github.com/withastro/starlight/commit/810238934ae1a95c53042ca2875bb4033aad0114) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Add support for collapsed sidebar groups

- [#254](https://github.com/withastro/starlight/pull/254) [`faa70de`](https://github.com/withastro/starlight/commit/faa70de584bf596fdd7184c4a8622d67d1410ecf) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Expose `<Icon>` component

- [#256](https://github.com/withastro/starlight/pull/256) [`048e948`](https://github.com/withastro/starlight/commit/048e948bce650d559517850c73d827733b8164c4) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Add new global `lastUpdated` option defaulting to `false` to define whether or not the last updated date is shown in the footer. A page can override this setting or the generated date using the new `lastUpdated` frontmatter field.

  ⚠️ Breaking change. Starlight will no longer show this date by default. To keep the previous behavior, you must explicitly set `lastUpdated` to `true` in your configuration.

  ```diff
  starlight({
  + lastUpdated: true,
  }),
  ```

### Patch Changes

- [#264](https://github.com/withastro/starlight/pull/264) [`ed1e46b`](https://github.com/withastro/starlight/commit/ed1e46beb1bc054ecdba36ecfe566ecaeaf8799b) Thanks [@astridx](https://github.com/astridx)! - Add new icon for displaying codeberg.org in social links.

- [#260](https://github.com/withastro/starlight/pull/260) [`01b65b1`](https://github.com/withastro/starlight/commit/01b65b1adf012474daf5678b4a709e3a7a484814) Thanks [@ElianCodes](https://github.com/ElianCodes)! - Add Dutch UI translations

- [#269](https://github.com/withastro/starlight/pull/269) [`fdc18b5`](https://github.com/withastro/starlight/commit/fdc18b5476957f8017a0fa1489c6fed89d5a9480) Thanks [@baspinarenes](https://github.com/baspinarenes)! - Add Turkish UI translations

- [#270](https://github.com/withastro/starlight/pull/270) [`1d3e705`](https://github.com/withastro/starlight/commit/1d3e705256fa0668db73b01c898e7e3b3b505c49) Thanks [@cbontems](https://github.com/cbontems)! - Improve French UI translations

- [#272](https://github.com/withastro/starlight/pull/272) [`6b23ebc`](https://github.com/withastro/starlight/commit/6b23ebc9974828837a2de9175297664e5d28a999) Thanks [@cbontems](https://github.com/cbontems)! - Add YouTube social link support

- [#267](https://github.com/withastro/starlight/pull/267) [`af2e43c`](https://github.com/withastro/starlight/commit/af2e43c7325a8f7fa6c9f867a3ec864daae39e96) Thanks [@nikcio](https://github.com/nikcio)! - Add Danish UI translations

- [#273](https://github.com/withastro/starlight/pull/273) [`d4f5134`](https://github.com/withastro/starlight/commit/d4f5134c91b393ac448efc3f44849fc886a05551) Thanks [@Waxer59](https://github.com/Waxer59)! - Fix typo in Spanish UI translations

## 0.3.1

### Patch Changes

- [#257](https://github.com/withastro/starlight/pull/257) [`0502327`](https://github.com/withastro/starlight/commit/050232770c8a2d22f317def8e734e4abb36387c6) Thanks [@JosefJezek](https://github.com/JosefJezek)! - Add Czech language support

- [#261](https://github.com/withastro/starlight/pull/261) [`2062b9e`](https://github.com/withastro/starlight/commit/2062b9e21695b5dc3f116731dbfdf609e495018c) Thanks [@delucis](https://github.com/delucis)! - Fix autogenerated navigation for pages using fallback content

## 0.3.0

### Minor Changes

- [#237](https://github.com/withastro/starlight/pull/237) [`4279d75`](https://github.com/withastro/starlight/commit/4279d7512a8261b576056471f5aa1ede1e6aae4a) Thanks [@HiDeoo](https://github.com/HiDeoo)! - Use path instead of slugified path for auto-generated sidebar item configuration

  ⚠️ Potentially breaking change. If your docs directory names don’t match their URLs, for example they contain whitespace like `docs/my docs/`, and you were referencing these in an `autogenerate` sidebar group as `my-docs`, update your config to reference these with the directory name instead of the slugified version:

  ```diff
  autogenerate: {
  - directory: 'my-docs',
  + directory: 'my docs',
  }
  ```

- [#226](https://github.com/withastro/starlight/pull/226) [`1aa2187`](https://github.com/withastro/starlight/commit/1aa2187944dde4419e523f0087139f5a21efd826) Thanks [@delucis](https://github.com/delucis)! - Add support for custom 404 pages.

### Patch Changes

- [#234](https://github.com/withastro/starlight/pull/234) [`91309ae`](https://github.com/withastro/starlight/commit/91309ae13250c5fd9f91a8e1843f16430773ff15) Thanks [@morinokami](https://github.com/morinokami)! - Add Japanese translation for `search.devWarning`

- [#227](https://github.com/withastro/starlight/pull/227) [`fbdecfa`](https://github.com/withastro/starlight/commit/fbdecfab47effb0cba7cbc9233a7b6bffdded320) Thanks [@Yan-Thomas](https://github.com/Yan-Thomas)! - Add missing i18n support to the Search component's dev warning.

- [#244](https://github.com/withastro/starlight/pull/244) [`f1bcbeb`](https://github.com/withastro/starlight/commit/f1bcbebeb441b6bb9ed6a1ab2414791e9d5de6ef) Thanks [@Waxer59](https://github.com/Waxer59)! - Add Spanish translation for `search.devWarning`

## 0.2.0

### Minor Changes

- [#171](https://github.com/withastro/starlight/pull/171) [`198c3f0`](https://github.com/withastro/starlight/commit/198c3f001410f259dab7d085136a37afe863cfa4) Thanks [@delucis](https://github.com/delucis)! - Add Starlight generator tag to HTML output

- [#217](https://github.com/withastro/starlight/pull/217) [`490fd98`](https://github.com/withastro/starlight/commit/490fd98d4e7b38ec01c568eee0ab00844e59c53d) Thanks [@delucis](https://github.com/delucis)! - Updated sidebar styles. Sidebars now support top-level links and groups are styled with a subtle border and indentation to improve comprehension of nesting.

- [#178](https://github.com/withastro/starlight/pull/178) [`d046c55`](https://github.com/withastro/starlight/commit/d046c55a62290c15f2e09faf4359f02df9492f6d) Thanks [@delucis](https://github.com/delucis)! - Add support for translating the Pagefind search modal

- [#210](https://github.com/withastro/starlight/pull/210) [`cb5b121`](https://github.com/withastro/starlight/commit/cb5b1210e23548e2983865a4b38308b0f54dc7ce) Thanks [@delucis](https://github.com/delucis)! - Change page title ID to `_top` for cleaner hash URLs

  ⚠️ Potentially breaking change if you were linking manually to `#starlight__overview` anywhere. If you were, update these links to use `#_top` instead.

### Patch Changes

- [#208](https://github.com/withastro/starlight/pull/208) [`09fc565`](https://github.com/withastro/starlight/commit/09fc565d44bd3abb4508541b458531de8624036f) Thanks [@delucis](https://github.com/delucis)! - Update `@astrojs/mdx` and `@astrojs/sitemap` to latest

- [#216](https://github.com/withastro/starlight/pull/216) [`54905c5`](https://github.com/withastro/starlight/commit/54905c502c5e6de5516e36ddcd4969893572baa5) Thanks [@morinokami](https://github.com/morinokami)! - Encode heading id when finding current link

## 0.1.4

### Patch Changes

- [#190](https://github.com/withastro/starlight/pull/190) [`a3809e4`](https://github.com/withastro/starlight/commit/a3809e4f1e14f3949e9e25f7ffbdea2920408edb) Thanks [@gabrielemercolino](https://github.com/gabrielemercolino)! - Added Italian language support

- [#193](https://github.com/withastro/starlight/pull/193) [`c9ca4eb`](https://github.com/withastro/starlight/commit/c9ca4ebe10f4776999e3fff4ac4c19ac0a714bac) Thanks [@BryceRussell](https://github.com/BryceRussell)! - Fix bottom padding for sidebar on larger screen sizes

## 0.1.3

### Patch Changes

- [#183](https://github.com/withastro/starlight/pull/183) [`89e0a04`](https://github.com/withastro/starlight/commit/89e0a04c26639246f550957acda2285e50417729) Thanks [@delucis](https://github.com/delucis)! - Fix disclosure caret rotation in sidebar sub-groups

- [#177](https://github.com/withastro/starlight/pull/177) [`bdafdb0`](https://github.com/withastro/starlight/commit/bdafdb050d9d4e2501485dff37b71a3175a3b0c8) Thanks [@rviscomi](https://github.com/rviscomi)! - Fix Markdown table overflow

- [#185](https://github.com/withastro/starlight/pull/185) [`4844915`](https://github.com/withastro/starlight/commit/4844915c25c9cdfb852e41584d93abfd85b82d08) Thanks [@delucis](https://github.com/delucis)! - Support setting an SVG as the hero image file

## 0.1.2

### Patch Changes

- [#174](https://github.com/withastro/starlight/pull/174) [`6ab31b4`](https://github.com/withastro/starlight/commit/6ab31b4900166f952c1ca5ec4e4a1ef66f31be97) Thanks [@rviscomi](https://github.com/rviscomi)! - Split `withBase` URL helper to fix use with files.

- [#168](https://github.com/withastro/starlight/pull/168) [`cb18eef`](https://github.com/withastro/starlight/commit/cb18eef4fda8227a6c5ec73589526dd7fbb8f4a6) Thanks [@BryceRussell](https://github.com/BryceRussell)! - Fix bottom padding on left sidebar

- [#167](https://github.com/withastro/starlight/pull/167) [`990ec53`](https://github.com/withastro/starlight/commit/990ec53dee099fdb6d113a3be5ef375c73e6945a) Thanks [@BryceRussell](https://github.com/BryceRussell)! - Add `bundlePath` option to Pagefind configuration

- [`4f666ba`](https://github.com/withastro/starlight/commit/4f666ba4fad7118a31bae819eb6be068da9e4d94) Thanks [@delucis](https://github.com/delucis)! - Fix focus outline positioning in tabs

## 0.1.1

### Patch Changes

- [#155](https://github.com/withastro/starlight/pull/155) Thanks [@thomasbnt](https://github.com/thomasbnt)! - Add French language support

- [#158](https://github.com/withastro/starlight/pull/158) [`92d82f5`](https://github.com/withastro/starlight/commit/92d82f534c6ff8513b01f9f26748d9980a6d4c79) Thanks [@kevinzunigacuellar](https://github.com/kevinzunigacuellar)! - Fix word wrapping in search modal on narrow screens

## 0.1.0

### Minor Changes

- [`43f3a02`](https://github.com/withastro/starlight/commit/43f3a024d6903072780f158dd95fb04b9f678535) Thanks [@delucis](https://github.com/delucis)! - Release v0.1.0

## 0.0.19

### Patch Changes

- [`fab453c`](https://github.com/withastro/starlight/commit/fab453c27a26a3928c0b355306d45b313a5fc531) Thanks [@delucis](https://github.com/delucis)! - Design tweak: larger sidebar text with more spacing

- [#134](https://github.com/withastro/starlight/pull/134) [`5f4acdf`](https://github.com/withastro/starlight/commit/5f4acdf75102f4431f5c60f65912db8d690b098c) Thanks [@Yan-Thomas](https://github.com/Yan-Thomas)! - Add Portuguese language support

- [`8805fbf`](https://github.com/withastro/starlight/commit/8805fbf30a2c26208aaf6d29ee53586f2dbf6cce) Thanks [@delucis](https://github.com/delucis)! - Add box-shadow to prev/next page links as per designs

- [`81ef58e`](https://github.com/withastro/starlight/commit/81ef58eac2a53672773d7d564068539190960127) Thanks [@delucis](https://github.com/delucis)! - Design tweak: slightly less horizontal padding in header component on narrower viewports

- [`8c103b3`](https://github.com/withastro/starlight/commit/8c103b3b44e4b91159f8225fffdf9ba843f9c395) Thanks [@delucis](https://github.com/delucis)! - Design tweak: pad bottom of page content slightly

- [#129](https://github.com/withastro/starlight/pull/129) [`bbcb277`](https://github.com/withastro/starlight/commit/bbcb277591514705fcc39665068aa331cfa2a653) Thanks [@delucis](https://github.com/delucis)! - Fix bug setting writing direction from a single root locale

## 0.0.18

### Patch Changes

- [`a76ae4d`](https://github.com/withastro/starlight/commit/a76ae4d4c459eae1690ed7fe6d4f4debb0137975) Thanks [@delucis](https://github.com/delucis)! - Add new icons for use in starter project

## 0.0.17

### Patch Changes

- [#107](https://github.com/withastro/starlight/pull/107) [`2f2d3ee`](https://github.com/withastro/starlight/commit/2f2d3eed1e7ed48d75205cfc3169719da7fdae1a) Thanks [@delucis](https://github.com/delucis)! - Small CSS size optimisation

- [#105](https://github.com/withastro/starlight/pull/105) [`55fec5d`](https://github.com/withastro/starlight/commit/55fec5d7e15da0e7365cee196d091bf5d15129c9) Thanks [@delucis](https://github.com/delucis)! - Add `<Card>` and `<CardGrid>` components for landing pages and other uses

## 0.0.16

### Patch Changes

- [#103](https://github.com/withastro/starlight/pull/103) [`ccb919d`](https://github.com/withastro/starlight/commit/ccb919d6580955e3428430a704f4a33fbc55a78d) Thanks [@delucis](https://github.com/delucis)! - Support adding a hero section to pages

- [#101](https://github.com/withastro/starlight/pull/101) [`6a2c0df`](https://github.com/withastro/starlight/commit/6a2c0df4d5586b70b46c854061df67b028e73630) Thanks [@TheOtterlord](https://github.com/TheOtterlord)! - Add better error messages for starlight config

## 0.0.15

### Patch Changes

- [`ded79af`](https://github.com/withastro/starlight/commit/ded79af43fad5ae0ec35739f655bf9e0c141a559) Thanks [@delucis](https://github.com/delucis)! - Add missing skip link to 404 page

- [#99](https://github.com/withastro/starlight/pull/99) [`d162b2f`](https://github.com/withastro/starlight/commit/d162b2fc0795248fa89d45f2e5d4207126a59256) Thanks [@delucis](https://github.com/delucis)! - Fix “next page” arrow showing on pages not in sidebar

- [#99](https://github.com/withastro/starlight/pull/99) [`d162b2f`](https://github.com/withastro/starlight/commit/d162b2fc0795248fa89d45f2e5d4207126a59256) Thanks [@delucis](https://github.com/delucis)! - Add support for a “splash” layout

- [#99](https://github.com/withastro/starlight/pull/99) [`d162b2f`](https://github.com/withastro/starlight/commit/d162b2fc0795248fa89d45f2e5d4207126a59256) Thanks [@delucis](https://github.com/delucis)! - Support hiding right sidebar table of contents

- [#99](https://github.com/withastro/starlight/pull/99) [`d162b2f`](https://github.com/withastro/starlight/commit/d162b2fc0795248fa89d45f2e5d4207126a59256) Thanks [@delucis](https://github.com/delucis)! - Move edit page link to page footer so it is accessible on mobile

## 0.0.14

### Patch Changes

- [#95](https://github.com/withastro/starlight/pull/95) [`de24b54`](https://github.com/withastro/starlight/commit/de24b54971577912979a3fb67570f4c95efe27a6) Thanks [@delucis](https://github.com/delucis)! - Support translations in sidebar config

- [#97](https://github.com/withastro/starlight/pull/97) [`2d51762`](https://github.com/withastro/starlight/commit/2d517623fb8670d4e7f2656eacff0d5beb27d95a) Thanks [@morinokami](https://github.com/morinokami)! - Add Japanese language support

## 0.0.13

### Patch Changes

- [`8688778`](https://github.com/withastro/starlight/commit/86887786d158e4cdb9e4bd021b2232eb6dba284c) Thanks [@delucis](https://github.com/delucis)! - Fix small CSS compatibility issue

- [#93](https://github.com/withastro/starlight/pull/93) [`c6d7960`](https://github.com/withastro/starlight/commit/c6d7960c8673886eb2b17843e78a897133c05fe2) Thanks [@delucis](https://github.com/delucis)! - Fix default locale routing bug when not using root locale

- [`d8a171b`](https://github.com/withastro/starlight/commit/d8a171b6b45c73151485fe8f08630fb6a1cc12a6) Thanks [@delucis](https://github.com/delucis)! - Fix autogenerated sidebar bug with index routes in subdirectories

- [`d8b9f32`](https://github.com/withastro/starlight/commit/d8b9f3260daaceed8a31eedcc44bd00733f96254) Thanks [@delucis](https://github.com/delucis)! - Fix false positive in sidebar autogeneration logic

- [#92](https://github.com/withastro/starlight/pull/92) [`02821d2`](https://github.com/withastro/starlight/commit/02821d2c8a58c485697cb8f0770c6ba63e709b2a) Thanks [@delucis](https://github.com/delucis)! - Update Pagefind to latest v1 alpha

- [`51fe914`](https://github.com/withastro/starlight/commit/51fe91468fea125ec33cb6d6b1b66f147302fdc0) Thanks [@delucis](https://github.com/delucis)! - Guarantee route and autogenerated sidebar sort order

- [`116c4f5`](https://github.com/withastro/starlight/commit/116c4f5eb0ddf4dddbd10005bc72a7e6cb880a67) Thanks [@delucis](https://github.com/delucis)! - Fix minor dev layout bug in Search modal for RTL languages

## 0.0.12

### Patch Changes

- [#85](https://github.com/withastro/starlight/pull/85) [`c86c1d6`](https://github.com/withastro/starlight/commit/c86c1d6e93d978d13e42bbc449e0225a06793ba3) Thanks [@BryceRussell](https://github.com/BryceRussell)! - Improve outside click detection on the search modal

## 0.0.11

### Patch Changes

- [`4da9cbd`](https://github.com/withastro/starlight/commit/4da9cbd7e643c97acbf4c2016aefb8f712bc9869) Thanks [@delucis](https://github.com/delucis)! - Fix typo in English UI strings

## 0.0.10

### Patch Changes

- [#78](https://github.com/withastro/starlight/pull/78) [`d3ee6fc`](https://github.com/withastro/starlight/commit/d3ee6fc643de7a320a6bb83432cdcfbb0a4e4289) Thanks [@delucis](https://github.com/delucis)! - Add support for customising and translating Starlight’s UI.

  Users can provide translations in JSON files in `src/content/i18n/` which is a data collection. For example, a `src/content/i18n/de.json` might translate the search UI:

  ```json
  {
  	"search.label": "Suchen",
  	"search.shortcutLabel": "(Drücke / zum Suchen)"
  }
  ```

  This change also allows Starlight to provide built-in support for more languages than just English and adds German & Spanish support.

- [#76](https://github.com/withastro/starlight/pull/76) [`5e82073`](https://github.com/withastro/starlight/commit/5e8207350dba0fce92fa101d311db627e2157654) Thanks [@lloydjatkinson](https://github.com/lloydjatkinson)! - Scale down code block font size to match Figma design

- [#78](https://github.com/withastro/starlight/pull/78) [`d3ee6fc`](https://github.com/withastro/starlight/commit/d3ee6fc643de7a320a6bb83432cdcfbb0a4e4289) Thanks [@delucis](https://github.com/delucis)! - Require a minimum Astro version of 2.5.0

## 0.0.9

### Patch Changes

- [#72](https://github.com/withastro/starlight/pull/72) [`3dc1d0c`](https://github.com/withastro/starlight/commit/3dc1d0c342c6db4e30b016035fa446101b1805a2) Thanks [@delucis](https://github.com/delucis)! - Fix vertical alignment of social icons in site header

- [#63](https://github.com/withastro/starlight/pull/63) [`823e351`](https://github.com/withastro/starlight/commit/823e351e1c1fc68ca3c20ab35c9a7d9b13760a70) Thanks [@liruifengv](https://github.com/liruifengv)! - Make sidebar groups to collapsible

- [#72](https://github.com/withastro/starlight/pull/72) [`3dc1d0c`](https://github.com/withastro/starlight/commit/3dc1d0c342c6db4e30b016035fa446101b1805a2) Thanks [@delucis](https://github.com/delucis)! - Fix image aspect ratio in Markdown content

## 0.0.8

### Patch Changes

- [#62](https://github.com/withastro/starlight/pull/62) [`a91191e`](https://github.com/withastro/starlight/commit/a91191e8308ffa746a3eadeea61e39412f32f926) Thanks [@delucis](https://github.com/delucis)! - Make `base` support consistent, including when `trailingSlash: 'never'` is set.

- [#61](https://github.com/withastro/starlight/pull/61) [`608f34c`](https://github.com/withastro/starlight/commit/608f34cbbe485c39730f33828971397f9c8a3534) Thanks [@liruifengv](https://github.com/liruifengv)! - Fix toc headingsObserver rootMargin

- [#66](https://github.com/withastro/starlight/pull/66) [`9ca67d8`](https://github.com/withastro/starlight/commit/9ca67d8984f76c22e5411d7352aa8e0bd4514f42) Thanks [@Yan-Thomas](https://github.com/Yan-Thomas)! - Make site title width fit the content

- [#64](https://github.com/withastro/starlight/pull/64) [`4460e55`](https://github.com/withastro/starlight/commit/4460e55de210fd9a23a762fe76f5c32297d68d76) Thanks [@delucis](https://github.com/delucis)! - Fix table of contents intersection observer for all possible viewport sizes.

- [#67](https://github.com/withastro/starlight/pull/67) [`38c2c1f`](https://github.com/withastro/starlight/commit/38c2c1f1ed25d6efe6ab2637ca2d9fbcdafcd240) Thanks [@TheOtterlord](https://github.com/TheOtterlord)! - Fix background color on select component

- [#57](https://github.com/withastro/starlight/pull/57) [`5b6cccb`](https://github.com/withastro/starlight/commit/5b6cccb7f9ee5810c75fbbb45496e2a1d022f7dd) Thanks [@BryceRussell](https://github.com/BryceRussell)! - Update site title link to include locale

## 0.0.7

### Patch Changes

- [#55](https://github.com/withastro/starlight/pull/55) [`8597b9c`](https://github.com/withastro/starlight/commit/8597b9c1002f8c5073d25ae5cacd4060ded2f8c8) Thanks [@delucis](https://github.com/delucis)! - Fix routing logic to handle `index.md` slug differences between docs collection root and nested directories.

- [#54](https://github.com/withastro/starlight/pull/54) [`db728d6`](https://github.com/withastro/starlight/commit/db728d61afa5cea060c66f746a4cc4ab3e1c3bcd) Thanks [@TheOtterlord](https://github.com/TheOtterlord)! - Add padding to scroll preventing headings being obscured by nav

- [#51](https://github.com/withastro/starlight/pull/51) [`3adbdbb`](https://github.com/withastro/starlight/commit/3adbdbbb71a4b3648984fa1028fa116d0aff9a7d) Thanks [@delucis](https://github.com/delucis)! - Support displaying a custom logo in the nav bar.

- [#51](https://github.com/withastro/starlight/pull/51) [`3adbdbb`](https://github.com/withastro/starlight/commit/3adbdbbb71a4b3648984fa1028fa116d0aff9a7d) Thanks [@delucis](https://github.com/delucis)! - All Starlight projects now use Astro’s experimental optimized asset support.

## 0.0.6

### Patch Changes

- [#47](https://github.com/withastro/starlight/pull/47) [`e96d9a7`](https://github.com/withastro/starlight/commit/e96d9a7628c5c04fe34dbc65ddd6fabdc0667a6d) Thanks [@delucis](https://github.com/delucis)! - Fix CSS ordering issue caused by imports in 404 route.

- [#47](https://github.com/withastro/starlight/pull/47) [`e96d9a7`](https://github.com/withastro/starlight/commit/e96d9a7628c5c04fe34dbc65ddd6fabdc0667a6d) Thanks [@delucis](https://github.com/delucis)! - Highlight current page section in table of contents.

- [`1028119`](https://github.com/withastro/starlight/commit/10281196aba65075e4ac202dc0f23927c44403ee) Thanks [@delucis](https://github.com/delucis)! - Use default locale in `404.astro`.

- [`05f8fd4`](https://github.com/withastro/starlight/commit/05f8fd4c3114e4c25075b35086c5b3e7d0ff49d7) Thanks [@delucis](https://github.com/delucis)! - Include `initial-scale=1` in viewport meta tag.

- [#47](https://github.com/withastro/starlight/pull/47) [`e96d9a7`](https://github.com/withastro/starlight/commit/e96d9a7628c5c04fe34dbc65ddd6fabdc0667a6d) Thanks [@delucis](https://github.com/delucis)! - Fix usage of `aria-current` in navigation sidebar to use `page` value.

- [#48](https://github.com/withastro/starlight/pull/48) [`a49485d`](https://github.com/withastro/starlight/commit/a49485def3fe4f505e90bf934eedcb135b3d3f51) Thanks [@delucis](https://github.com/delucis)! - Improve right sidebar layout.

## 0.0.5

### Patch Changes

- [#42](https://github.com/withastro/starlight/pull/42) [`c6c1b67`](https://github.com/withastro/starlight/commit/c6c1b6727140a76c42c661f406000cc6e9b175de) Thanks [@delucis](https://github.com/delucis)! - Support setting custom `<head>` tags in config or frontmatter.

## 0.0.4

### Patch Changes

- [#40](https://github.com/withastro/starlight/pull/40) [`e22dd76`](https://github.com/withastro/starlight/commit/e22dd76136f9749bb5d43f96241385faccfc90a1) Thanks [@delucis](https://github.com/delucis)! - Generate sitemaps for Starlight sites

- [#38](https://github.com/withastro/starlight/pull/38) [`623b577`](https://github.com/withastro/starlight/commit/623b577319b1dea2d6c42f1b680139fb858d85d6) Thanks [@delucis](https://github.com/delucis)! - Add tab components for use in MDX.

## 0.0.3

### Patch Changes

- [`1f62c75`](https://github.com/withastro/starlight/commit/1f62c75297ed44da721350840744823876a64bc5) Thanks [@delucis](https://github.com/delucis)! - Add repo, issues & homepage links to package.json

- [#36](https://github.com/withastro/starlight/pull/36) [`62265e4`](https://github.com/withastro/starlight/commit/62265e4a653d161483e3844b568ab150334e9238) Thanks [@delucis](https://github.com/delucis)! - Collapse table of contents in dropdown on narrower viewports

## 0.0.2

### Patch Changes

- [#35](https://github.com/withastro/starlight/pull/35) [`ea999c4`](https://github.com/withastro/starlight/commit/ea999c4afa0e72da5a5510d30a7bc13380e10a4a) Thanks [@delucis](https://github.com/delucis)! - Add support for social media links in the site header.

- [#33](https://github.com/withastro/starlight/pull/33) [`833393a`](https://github.com/withastro/starlight/commit/833393a1cd8fc83a3b9c6055cfb8d160ab40a40c) Thanks [@MoustaphaDev](https://github.com/MoustaphaDev)! - Import zod from `astro/zod` to fix an issue related to importing the `astro:content` virtual module

## 0.0.1

### Patch Changes

- [#30](https://github.com/withastro/starlight/pull/30) [`dc6e04d`](https://github.com/withastro/starlight/commit/dc6e04d2cf09339713743bc113a57d1880fec85d) Thanks [@delucis](https://github.com/delucis)! - Set up npm publishing with changesets
