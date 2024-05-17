---
title: Authoring Content in Markdown
description: An overview of the Markdown syntax Starlight supports.
---

Starlight supports the full range of [Markdown](https://daringfireball.net/projects/markdown/) syntax in `.md` files as well as frontmatter [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f) to define metadata such as a title and description.

Please be sure to check the [MDX docs](https://mdxjs.com/docs/what-is-mdx/#markdown) or [Markdoc docs](https://markdoc.dev/docs/syntax) if using those file formats, as Markdown support and usage can differ.

## Frontmatter

You can customize individual pages in Starlight by setting values in their frontmatter.
Frontmatter is set at the top of your files between `---` separators:

```md title="src/content/docs/example.md"
---
title: My page title
---

Page content follows the second `---`.
```

Every page must include at least a `title`.
See the [frontmatter reference](/reference/frontmatter/) for all available fields and how to add custom fields.

## Inline styles

Text can be **bold**, _italic_, or ~~strikethrough~~.

```md
Text can be **bold**, _italic_, or ~~strikethrough~~.
```

You can [link to another page](/getting-started/).

```md
You can [link to another page](/getting-started/).
```

You can highlight `inline code` with backticks.

```md
You can highlight `inline code` with backticks.
```

## Images

Images in Starlight use [Astro’s built-in optimized asset support](https://docs.astro.build/en/guides/assets/).

Markdown and MDX support the Markdown syntax for displaying images that includes alt-text for screen readers and assistive technology.

![An illustration of planets and stars featuring the word “astro”](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
![An illustration of planets and stars featuring the word “astro”](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

Relative image paths are also supported for images stored locally in your project.

```md
// src/content/docs/page-1.md

![A rocketship in space](../../assets/images/rocket.svg)
```

## Headings

You can structure content using a heading. Headings in Markdown are indicated by a number of `#` at the start of the line.

### How to structure page content in Starlight

Starlight is configured to automatically use your page title as a top-level heading and will include an "Overview" heading at top of each page's table of contents. We recommend starting each page with regular paragraph text content and using on-page headings from `<h2>` and down:

```md
---
title: Markdown Guide
description: How to use Markdown in Starlight
---

This page describes how to use Markdown in Starlight.

## Inline Styles

## Headings
```

### Automatic heading anchor links

Using headings in Markdown will automatically give you anchor links so you can link directly to certain sections of your page:

```md
---
title: My page of content
description: How to use Starlight's built-in anchor links
---

## Introduction

I can link to [my conclusion](#conclusion) lower on the same page.

## Conclusion

`https://my-site.com/page1/#introduction` navigates directly to my Introduction.
```

Level 2 (`<h2>`) and Level 3 (`<h3>`) headings will automatically appear in the page table of contents.

Learn more about how Astro processes heading `id`s in [the Astro Documentation](https://docs.astro.build/en/guides/markdown-content/#heading-ids)

## Asides

Asides (also known as “admonitions” or “callouts”) are useful for displaying secondary information alongside a page’s main content.

Starlight provides a custom Markdown syntax for rendering asides. Aside blocks are indicated using a pair of triple colons `:::` to wrap your content, and can be of type `note`, `tip`, `caution` or `danger`.

You can nest any other Markdown content types inside an aside, but asides are best suited to short and concise chunks of content.

### Note aside

:::note
Starlight is a documentation website toolkit built with [Astro](https://astro.build/). You can get started with this command:

```sh
npm create astro@latest -- --template starlight
```

:::

````md
:::note
Starlight is a documentation website toolkit built with [Astro](https://astro.build/). You can get started with this command:

```sh
npm create astro@latest -- --template starlight
```

:::
````

### Custom aside titles

You can specify a custom title for the aside in square brackets following the aside type, e.g. `:::tip[Did you know?]`.

:::tip[Did you know?]
Astro helps you build faster websites with [“Islands Architecture”](https://docs.astro.build/en/concepts/islands/).
:::

```md
:::tip[Did you know?]
Astro helps you build faster websites with [“Islands Architecture”](https://docs.astro.build/en/concepts/islands/).
:::
```

### More aside types

Caution and danger asides are helpful for drawing a user’s attention to details that may trip them up.
If you find yourself using these a lot, it may also be a sign that the thing you are documenting could benefit from being redesigned.

:::caution
If you are not sure you want an awesome docs site, think twice before using [Starlight](/).
:::

:::danger
Your users may be more productive and find your product easier to use thanks to helpful Starlight features.

- Clear navigation
- User-configurable colour theme
- [i18n support](/guides/i18n/)

:::

```md
:::caution
If you are not sure you want an awesome docs site, think twice before using [Starlight](/).
:::

:::danger
Your users may be more productive and find your product easier to use thanks to helpful Starlight features.

- Clear navigation
- User-configurable colour theme
- [i18n support](/guides/i18n/)

:::
```

## Blockquotes

> This is a blockquote, which is commonly used when quoting another person or document.
>
> Blockquotes are indicated by a `>` at the start of each line.

```md
> This is a blockquote, which is commonly used when quoting another person or document.
>
> Blockquotes are indicated by a `>` at the start of each line.
```

## Code blocks

A code block is indicated by a block with three backticks <code>```</code> at the start and end. You can indicate the programming language being used after the opening backticks.

```js
// Javascript code with syntax highlighting.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// Javascript code with syntax highlighting.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

### Expressive Code features

Starlight uses [Expressive Code](https://github.com/expressive-code/expressive-code/tree/main/packages/astro-expressive-code) to extend formatting possibilities for code blocks.
Expressive Code’s text markers and window frames plugins are enabled by default.
Code block rendering can be configured using Starlight’s [`expressiveCode` configuration option](/reference/configuration/#expressivecode).

#### Text markers

You can highlight specific lines or parts of your code blocks using [Expressive Code text markers](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#usage-in-markdown--mdx-documents) on the opening line of your code block.
Use curly braces (`{ }`) to highlight entire lines, and quotation marks to highlight strings of text.

There are three highlighting styles: neutral for calling attention to code, green for indicating inserted code, and red for indicating deleted code.
Both text and entire lines can be marked using the default marker, or in combination with `ins=` and `del=` to produce the desired highlighting.

Expressive Code provides several options for customizing the visual appearance of your code samples.
Many of these can be combined, for highly illustrative code samples.
Please explore the [Expressive Code documentation](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md) for the extensive options available.
Some of the most common examples are shown below:

- [Mark entire lines & line ranges using the `{ }` marker](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-entire-lines--line-ranges):

  ```js {2-3}
  function demo() {
    // This line (#2) and the next one are highlighted
    return 'This is line #3 of this snippet';
  }
  ```

  ````md
  ```js {2-3}
  function demo() {
    // This line (#2) and the next one are highlighted
    return 'This is line #3 of this snippet';
  }
  ```
  ````

- [Mark selections of text using the `" "` marker or regular expressions](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-individual-text-inside-lines):

  ```js "Individual terms" /Even.*supported/
  // Individual terms can be highlighted, too
  function demo() {
    return 'Even regular expressions are supported';
  }
  ```

  ````md
  ```js "Individual terms" /Even.*supported/
  // Individual terms can be highlighted, too
  function demo() {
    return 'Even regular expressions are supported';
  }
  ```
  ````

- [Mark text or lines as inserted or deleted with `ins` or `del`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#selecting-marker-types-mark-ins-del):

  ```js "return true;" ins="inserted" del="deleted"
  function demo() {
    console.log('These are inserted and deleted marker types');
    // The return statement uses the default marker type
    return true;
  }
  ```

  ````md
  ```js "return true;" ins="inserted" del="deleted"
  function demo() {
    console.log('These are inserted and deleted marker types');
    // The return statement uses the default marker type
    return true;
  }
  ```
  ````

- [Combine syntax highlighting with `diff`-like syntax](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#combining-syntax-highlighting-with-diff-like-syntax):

  ```diff lang="js"
    function thisIsJavaScript() {
      // This entire block gets highlighted as JavaScript,
      // and we can still add diff markers to it!
  -   console.log('Old code to be removed')
  +   console.log('New and shiny code!')
    }
  ```

  ````md
  ```diff lang="js"
    function thisIsJavaScript() {
      // This entire block gets highlighted as JavaScript,
      // and we can still add diff markers to it!
  -   console.log('Old code to be removed')
  +   console.log('New and shiny code!')
    }
  ```
  ````

#### Frames and titles

Code blocks can be rendered inside a window-like frame.
A frame that looks like a terminal window will be used for shell scripting languages (e.g. `bash` or `sh`).
Other languages display inside a code editor-style frame if they include a title.

A code block’s optional title can be set either with a `title="..."` attribute following the code block's opening backticks and language identifier, or with a file name comment in the first lines of the code.

- [Add a file name tab with a comment](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```js
  // my-test-file.js
  console.log('Hello World!');
  ```

  ````md
  ```js
  // my-test-file.js
  console.log('Hello World!');
  ```
  ````

- [Add a title to a Terminal window](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```bash title="Installing dependencies…"
  npm install
  ```

  ````md
  ```bash title="Installing dependencies…"
  npm install
  ```
  ````

- [Disable window frames with `frame="none"`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#overriding-frame-types)

  ```bash frame="none"
  echo "This is not rendered as a terminal despite using the bash language"
  ```

  ````md
  ```bash frame="none"
  echo "This is not rendered as a terminal despite using the bash language"
  ```
  ````

## Other common Markdown features

Starlight supports all other Markdown authoring syntax, such as lists and tables. See the [Markdown Cheat Sheet from The Markdown Guide](https://www.markdownguide.org/cheat-sheet/) for a quick overview of all the Markdown syntax elements.

## Advanced Markdown and MDX configuration

Starlight uses Astro’s Markdown and MDX renderer built on remark and rehype. You can add support for custom syntax and behavior by adding `remarkPlugins` or `rehypePlugins` in your Astro config file. See [“Configuring Markdown and MDX”](https://docs.astro.build/en/guides/markdown-content/#configuring-markdown-and-mdx) in the Astro docs to learn more.
