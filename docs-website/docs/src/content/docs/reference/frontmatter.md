---
title: Frontmatter Reference
description: An overview of the default frontmatter fields Starlight supports.
---

You can customize individual Markdown and MDX pages in Starlight by setting values in their frontmatter. For example, a regular page might set `title` and `description` fields:

```md {3-4}
---
# src/content/docs/example.md
title: About this project
description: Learn more about the project I’m working on.
---

Welcome to the about page!
```

## Frontmatter fields

### `title` (required)

**type:** `string`

You must provide a title for every page. This will be displayed at the top of the page, in browser tabs, and in page metadata.

### `description`

**type:** `string`

The page description is used for page metadata and will be picked up by search engines and in social media previews.

### `slug`

**type**: `string`

Override the slug of the page. See [“Defining custom slugs”](https://docs.astro.build/en/guides/content-collections/#defining-custom-slugs) in the Astro docs for more details.

### `editUrl`

**type:** `string | boolean`

Overrides the [global `editLink` config](/reference/configuration/#editlink). Set to `false` to disable the “Edit page” link for a specific page or provide an alternative URL where the content of this page is editable.

### `head`

**type:** [`HeadConfig[]`](/reference/configuration/#headconfig)

You can add additional tags to your page’s `<head>` using the `head` frontmatter field. This means you can add custom styles, metadata or other tags to a single page. Similar to the [global `head` option](/reference/configuration/#head).

```md
---
# src/content/docs/example.md
title: About us
head:
  # Use a custom <title> tag
  - tag: title
    content: Custom about title
---
```

### `tableOfContents`

**type:** `false | { minHeadingLevel?: number; maxHeadingLevel?: number; }`

Overrides the [global `tableOfContents` config](/reference/configuration/#tableofcontents).
Customize the heading levels to be included or set to `false` to hide the table of contents on this page.

```md
---
# src/content/docs/example.md
title: Page with only H2s in the table of contents
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 2
---
```

```md
---
# src/content/docs/example.md
title: Page with no table of contents
tableOfContents: false
---
```

### `template`

**type:** `'doc' | 'splash'`  
**default:** `'doc'`

Set the layout template for this page.
Pages use the `'doc'` layout by default.
Set to `'splash'` to use a wider layout without any sidebars designed for landing pages.

### `hero`

**type:** [`HeroConfig`](#heroconfig)

Add a hero component to the top of this page. Works well with `template: splash`.

For example, this config shows some common options, including loading an image from your repository.

```md
---
# src/content/docs/example.md
title: My Home Page
template: splash
hero:
  title: 'My Project: Stellar Stuff Sooner'
  tagline: Take your stuff to the moon and back in the blink of an eye.
  image:
    alt: A glittering, brightly colored logo
    file: ~/assets/logo.png
  actions:
    - text: Tell me more
      link: /getting-started/
      icon: right-arrow
      variant: primary
    - text: View on GitHub
      link: https://github.com/astronaut/my-project
      icon: external
      attrs:
        rel: me
---
```

You can display different versions of the hero image in light and dark modes.

```md
---
# src/content/docs/example.md
hero:
  image:
    alt: A glittering, brightly colored logo
    dark: ~/assets/logo-dark.png
    light: ~/assets/logo-light.png
---
```

#### `HeroConfig`

```ts
interface HeroConfig {
  title?: string;
  tagline?: string;
  image?:
    | {
        // Relative path to an image in your repository.
        file: string;
        // Alt text to make the image accessible to assistive technology
        alt?: string;
      }
    | {
        // Relative path to an image in your repository to be used for dark mode.
        dark: string;
        // Relative path to an image in your repository to be used for light mode.
        light: string;
        // Alt text to make the image accessible to assistive technology
        alt?: string;
      }
    | {
        // Raw HTML to use in the image slot.
        // Could be a custom `<img>` tag or inline `<svg>`.
        html: string;
      };
  actions?: Array<{
    text: string;
    link: string;
    variant: 'primary' | 'secondary' | 'minimal';
    icon: string;
    attrs?: Record<string, string | number | boolean>;
  }>;
}
```

### `banner`

**type:** `{ content: string }`

Displays an announcement banner at the top of this page.

The `content` value can include HTML for links or other content.
For example, this page displays a banner including a link to `example.com`.

```md
---
# src/content/docs/example.md
title: Page with a banner
banner:
  content: |
    We just launched something cool!
    <a href="https://example.com">Check it out</a>
---
```

### `lastUpdated`

**type:** `Date | boolean`

Overrides the [global `lastUpdated` option](/reference/configuration/#lastupdated). If a date is specified, it must be a valid [YAML timestamp](https://yaml.org/type/timestamp.html) and will override the date stored in Git history for this page.

```md
---
# src/content/docs/example.md
title: Page with a custom last update date
lastUpdated: 2022-08-09
---
```

### `prev`

**type:** `boolean | string | { link?: string; label?: string }`

Overrides the [global `pagination` option](/reference/configuration/#pagination). If a string is specified, the generated link text will be replaced and if an object is specified, both the link and the text will be overridden.

```md
---
# src/content/docs/example.md
# Hide the previous page link
prev: false
---
```

```md
---
# src/content/docs/example.md
# Override the previous page link text
prev: Continue the tutorial
---
```

```md
---
# src/content/docs/example.md
# Override both the previous page link and text
prev:
  link: /unrelated-page/
  label: Check out this other page
---
```

### `next`

**type:** `boolean | string | { link?: string; label?: string }`

Same as [`prev`](#prev) but for the next page link.

```md
---
# src/content/docs/example.md
# Hide the next page link
next: false
---
```

### `pagefind`

**type:** `boolean`  
**default:** `true`

Set whether this page should be included in the [Pagefind](https://pagefind.app/) search index. Set to `false` to exclude a page from search results:

```md
---
# src/content/docs/example.md
# Hide this page from the search index
pagefind: false
---
```

### `draft`

**type:** `boolean`  
**default:** `false`

Set whether this page should be considered a draft and not be included in [production builds](https://docs.astro.build/en/reference/cli-reference/#astro-build) and [autogenerated link groups](/guides/sidebar/#autogenerated-groups). Set to `true` to mark a page as a draft and make it only visible during development.

```md
---
# src/content/docs/example.md
# Exclude this page from production builds
draft: true
---
```

### `sidebar`

**type:** [`SidebarConfig`](#sidebarconfig)

Control how this page is displayed in the [sidebar](/reference/configuration/#sidebar), when using an autogenerated link group.

#### `SidebarConfig`

```ts
interface SidebarConfig {
  label?: string;
  order?: number;
  hidden?: boolean;
  badge?: string | BadgeConfig;
  attrs?: Record<string, string | number | boolean | undefined>;
}
```

#### `label`

**type:** `string`  
**default:** the page [`title`](#title-required)

Set the label for this page in the sidebar when displayed in an autogenerated group of links.

```md
---
# src/content/docs/example.md
title: About this project
sidebar:
  label: About
---
```

#### `order`

**type:** `number`

Control the order of this page when sorting an autogenerated group of links.
Lower numbers are displayed higher up in the link group.

```md
---
# src/content/docs/example.md
title: Page to display first
sidebar:
  order: 1
---
```

#### `hidden`

**type:** `boolean`  
**default:** `false`

Prevents this page from being included in an autogenerated sidebar group.

```md
---
# src/content/docs/example.md
title: Page to hide from autogenerated sidebar
sidebar:
  hidden: true
---
```

#### `badge`

**type:** <code>string | <a href="/reference/configuration/#badgeconfig">BadgeConfig</a></code>

Add a badge to the page in the sidebar when displayed in an autogenerated group of links.
When using a string, the badge will be displayed with a default accent color.
Optionally, pass a [`BadgeConfig` object](/reference/configuration/#badgeconfig) with `text` and `variant` fields to customize the badge.

```md
---
# src/content/docs/example.md
title: Page with a badge
sidebar:
  # Uses the default variant matching your site’s accent color
  badge: New
---
```

```md
---
# src/content/docs/example.md
title: Page with a badge
sidebar:
  badge:
    text: Experimental
    variant: caution
---
```

#### `attrs`

**type:** `Record<string, string | number | boolean | undefined>`

HTML attributes to add to the page link in the sidebar when displayed in an autogenerated group of links.

```md
---
# src/content/docs/example.md
title: Page opening in a new tab
sidebar:
  # Opens the page in a new tab
  attrs:
    target: _blank
---
```

## Customize frontmatter schema

The frontmatter schema for Starlight’s `docs` content collection is configured in `src/content/config.ts` using the `docsSchema()` helper:

```ts {3,6}
// src/content/config.ts
import { defineCollection } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({ schema: docsSchema() }),
};
```

Learn more about content collection schemas in [“Defining a collection schema”](https://docs.astro.build/en/guides/content-collections/#defining-a-collection-schema) in the Astro docs.

`docsSchema()` takes the following options:

### `extend`

**type:** Zod schema or function that returns a Zod schema  
**default:** `z.object({})`

Extend Starlight’s schema with additional fields by setting `extend` in the `docsSchema()` options.
The value should be a [Zod schema](https://docs.astro.build/en/guides/content-collections/#defining-datatypes-with-zod).

In the following example, we provide a stricter type for `description` to make it required and add a new optional `category` field:

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: z.object({
        // Make a built-in field required instead of optional.
        description: z.string(),
        // Add a new field to the schema.
        category: z.enum(['tutorial', 'guide', 'reference']).optional(),
      }),
    }),
  }),
};
```

To take advantage of the [Astro `image()` helper](https://docs.astro.build/en/guides/images/#images-in-content-collections), use a function that returns your schema extension:

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: ({ image }) => {
        return z.object({
          // Add a field that must resolve to a local image.
          cover: image(),
        });
      },
    }),
  }),
};
```
