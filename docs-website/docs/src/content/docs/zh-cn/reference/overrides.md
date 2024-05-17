---
title: 重写参考
description: Starlight 支持重写的组件以及组件参数的概述。
tableOfContents:
  maxHeadingLevel: 4
---

你可以通过在 Starlight 的 [`components`](/zh-cn/reference/configuration/#components) 配置选项中提供替代组件的路径来替换掉 (即重写) Starlight 的内置组件。
本页面列出了所有可被重写的组件和它们默认实现的 GitHub 链接。

在[重写组件指南](/zh-cn/guides/overriding-components/)中了解更多。

## 组件参数

所有组件都可以使用标准的 `Astro.props` 对象，该对象包含有关当前页面的信息。

从 Starlight 导入 `Props` 类型来为你的自定义组件定义类型：

```astro
---
// src/components/Custom.astro
import type { Props } from '@astrojs/starlight/props';

const { hasSidebar } = Astro.props;
//      ^ type: boolean
---
```

这样当使用 `Astro.props` 时就会有自动补全和类型提示。

### 参数

Starlight 会将以下参数传递给你的自定义组件。

#### `dir`

**类型：** `'ltr' | 'rtl'`

当前页面的文本方向。

#### `lang`

**类型：** `string`

当前页面的 BCP-47 语言标签，例如 `en`、`zh-CN` 或 `pt-BR`。

#### `locale`

**类型：** `string | undefined`

当前语言的根路径。对于默认语言来说是 `undefined`。

#### `siteTitle`

**类型：** `string`

根据页面语言设置的网站标题。

#### `siteTitleHref`

**类型：** `string`

网站标题的 `href` 属性值，链接回首页，例如 `/`。
对于多语言站点，将包含当前的语言环境，例如 `/en/` 或 `/zh-cn/`。

#### `slug`

**类型：** `string`

从内容文件名生成的页面 slug。

#### `id`

**类型：** `string`

基于内容文件名的页面的唯一 ID。

#### `isFallback`

**类型：** `true | undefined`

如果此页面在当前语言中未翻译，回退到使用默认语言的内容，则为 `true`。
仅在多语言站点中使用。

#### `entryMeta`

**类型：** `{ dir: 'ltr' | 'rtl'; lang: string }`

页面内容的语言环境元数据 (locale metadata)。当页面使用回退内容时可以与顶级语言环境设置值不同。

#### `entry`

当前页面所对应的 Astro 内容集合条目。
在 `entry.data` 中包含当前页面的 frontmatter 内容。

```ts
entry: {
  data: {
    title: string;
    description: string | undefined;
    // 等
  }
}
```

在 [Astro 的集合条目类型](https://docs.astro.build/zh-cn/reference/api-reference/#集合条目类型)参考中了解更多关于此对象的信息。

#### `sidebar`

**类型：** `SidebarEntry[]`

当前页面的侧边栏条目。

#### `hasSidebar`

**类型：** `boolean`

当前页面是否应该显示侧边栏。

#### `pagination`

**类型：** `{ prev?: Link; next?: Link }`

如果启用了，当前页面在侧边栏中的上一页和下一页的链接。

#### `toc`

**类型：** `{ minHeadingLevel: number; maxHeadingLevel: number; items: TocItem[] } | undefined`

如果启用了，当前页面的目录 (table of contents)。

#### `headings`

**类型：** `{ depth: number; slug: string; text: string }[]`

从当前页面提取的所有 Markdown 标题的数组。
如果你想要构建一个遵循 Starlight 配置选项的目录组件，请使用 [`toc`](#toc)。

#### `lastUpdated`

**类型：** `Date | undefined`

如果启用了，表示此页面最后更新时间的 JavaScript `Date` 对象。

#### `editUrl`

**类型：** `URL | undefined`

如果启用了，表示可以编辑此页面的地址的 JavaScript `URL` 对象。

#### `labels`

**类型：** `Record<string, string>`

一个包含为当前页面本地化的 UI 字符串的对象。请参阅 [“翻译 Starlight UI”](/zh-cn/guides/i18n/#翻译-starlight-的-ui) 指南以获取所有可用键的列表。

---

## 组件

### 头部

这些组件在每个页面的 `<head>` 元素内渲染。
它们应只包含[允许在 `<head>` 中使用的元素](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/head#相关链接)。

#### `Head`

**默认组件：** [`Head.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Head.astro)

在每个页面的 `<head>` 元素内渲染的组件。
包含 `<title>` 和 `<meta charset="utf-8">` 等重要标签。

重写此组件应作为最后手段。
如果可能，请优先使用 Starlight 配置中的 [`head`](/zh-cn/reference/configuration/#head) 选项。

#### `ThemeProvider`

**默认组件：** [`ThemeProvider.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeProvider.astro)

在 `<head>` 内渲染的用于提供暗色/亮色主题支持的组件。
默认实现包括一个内联脚本和一个 `<template>`，该脚本在 [`<ThemeSelect />`](#themeselect) 中使用该模板。

---

### 无障碍

#### `SkipLink`

**默认组件：** [`SkipLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SkipLink.astro)

在 `<body>` 内渲染的第一个元素，它链接到主页面内容以实现无障碍访问。
默认实现默认为隐藏状态，只有在用户使用键盘通过 tab 键聚焦到它时才会显示。

---

### 布局

这些组件负责在不同的断点 (breakpoints) 上布局 Starlight 的组件、管理视图。
重写这些组件会有很大的复杂性。
如果可能，请优先重写较低级别的组件。

#### `PageFrame`

**默认组件：** [`PageFrame.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageFrame.astro)

包在绝大部分页面内容外的布局组件。
默认实现提供了头部—侧边栏—主内容的布局，并包含 `header` 和 `sidebar` 命名插槽以及主内容的默认插槽。
它还渲染了 [`<MobileMenuToggle />`](#mobilemenutoggle) 以支持在小 (移动) 视口 (viewports) 上切换侧边栏导航。

#### `MobileMenuToggle`

**默认组件：** [`MobileMenuToggle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuToggle.astro)

在 [`<PageFrame>`](#pageframe) 内渲染的负责在小 (移动) 视口上切换侧边栏导航的组件。

#### `TwoColumnContent`

**默认组件：** [`TwoColumnContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TwoColumnContent.astro)

包在主内容列和右侧栏 (目录) 外的布局组件。
默认实现实现了在单列、小视口布局和两列、较大视口布局之间的切换。

---

### 导航

这些组件渲染 Starlight 的顶部导航栏。

#### `Header`

**默认组件：** [`Header.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Header.astro)

在每个页面顶部显示的导航栏组件。
默认实现显示了 [`<SiteTitle />`](#sitetitle-1)、[`<Search />`](#search)、[`<SocialIcons />`](#socialicons)、[`<ThemeSelect />`](#themeselect) 和 [`<LanguageSelect />`](#languageselect)。

#### `SiteTitle`

**默认组件：** [`SiteTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SiteTitle.astro)

在导航栏开头渲染的组件，用于渲染站点标题。
默认实现包含在 Starlight 配置中定义的 logo 的渲染逻辑。

#### `Search`

**默认组件：** [`Search.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Search.astro)

用于渲染 Starlight 搜索 UI 的组件。
默认实现包含在导航栏中的按钮和在点击时显示搜索模态框以及加载 [Pagefind UI](https://pagefind.app/) 的代码。

当 [`pagefind`](/zh-cn/reference/configuration/#pagefind) 被禁用时，默认的搜索组件不会被渲染。
然而，如果你重写了 `Search`，你的自定义组件将总是被渲染，即使 `pagefind` 配置选项是 `false`。
这允许你在禁用 Pagefind 时为其他搜索提供商添加 UI。

#### `SocialIcons`

**默认组件：** [`SocialIcons.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SocialIcons.astro)

在导航栏中渲染的组件，用于渲染社交图标链接。
默认实现使用 Starlight 配置中的 [`social`](/zh-cn/reference/configuration/#social) 选项来渲染图标和链接。

#### `ThemeSelect`

**默认组件：** [`ThemeSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeSelect.astro)

在导航栏中渲染的组件，用于允许用户选择深浅主题偏好。

#### `LanguageSelect`

**默认组件：** [`LanguageSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LanguageSelect.astro)

在导航栏中渲染的组件，用于允许用户切换到不同的语言。

---

### 全局侧边栏

Starlight 的全局侧边栏包含了主站点导航。
在较窄的视口上，它会隐藏在下拉菜单中。

#### `Sidebar`

**默认组件：** [`Sidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Sidebar.astro)

在页面内容之前渲染的包含全局导航的组件。
默认实现在足够宽的视口上显示为侧边栏，在小 (移动) 视口上显示为下拉菜单。
它还渲染了 [`<MobileMenuFooter />`](#mobilemenufooter) 以在移动菜单中显示额外的项目。

#### `MobileMenuFooter`

**默认组件：** [`MobileMenuFooter.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuFooter.astro)

在移动下拉菜单最底部中渲染的组件。
默认实现渲染了 [`<ThemeSelect />`](#themeselect) 和 [`<LanguageSelect />`](#languageselect)。

---

### 页面侧边栏

Starlight 的页面侧边栏负责显示当前页面的子标题的目录。
在较窄的视口上，它会缩为一个固定的下拉菜单。

#### `PageSidebar`

**默认组件：** [`PageSidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageSidebar.astro)

在页面内容之前渲染的包含目录的组件。
默认实现渲染了 [`<TableOfContents />`](#tableofcontents) 和 [`<MobileTableOfContents />`](#mobiletableofcontents)。

#### `TableOfContents`

**默认组件：** [`TableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TableOfContents.astro)

在较宽的视口上渲染当前页面的目录的组件。

#### `MobileTableOfContents`

**默认组件：** [`MobileTableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileTableOfContents.astro)

在小 (移动) 视口上渲染当前页面的目录的组件。

---

### 内容

这些组件在页面主内容列中渲染。

#### `Banner`

**默认组件：** [`Banner.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Banner.astro)

横幅 (Banner) 组件在每个页面的顶部渲染。
默认实现使用页面的 [`banner`](/zh-cn/reference/frontmatter/#banner) frontmatter 值来决定是否渲染。

#### `ContentPanel`

**默认组件：** [`ContentPanel.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ContentPanel.astro)

包在页面主内容列中的段落外的布局组件。

#### `PageTitle`

**默认组件：** [`PageTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageTitle.astro)

包含当前页面的 `<h1>` 元素的组件。

自定义实现应确保在 `<h1>` 元素上设置 `id="_top"`，就像默认实现中一样。

#### `DraftContentNotice`

**默认组件：** [`DraftContentNotice.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/DraftContentNotice.astro)

在开发过程中，当当前页面被标记为草稿时，向用户显示的通知。

#### `FallbackContentNotice`

**默认组件：** [`FallbackContentNotice.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/FallbackContentNotice.astro)

在本页面当前语言没有翻译时显示给用户的通知。
仅在多语言站点中使用。

#### `Hero`

**默认组件：** [`Hero.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Hero.astro)

当设置了 frontmatter 中的 [`hero`](/zh-cn/reference/frontmatter/#hero) 时在页面顶部渲染的组件。
默认实现显示了一个大标题、标语、动作链接 (call-to-action links) 和可选的图片。

#### `MarkdownContent`

**默认组件：** [`MarkdownContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MarkdownContent.astro)

在页面主内容列中渲染 Markdown 内容的组件。
默认实现为 Markdown 内容提供了基本的样式。

Markdown 内容样式也暴露在 `@astrojs/starlight/style/markdown.css` 中，并且作用域限制在 `.sl-markdown-content` CSS 类中。

---

### 页脚

这些组件在页面主内容列的底部渲染。

#### `Footer`

**默认组件：** [`Footer.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Footer.astro)

页脚 (Footer) 组件在每个页面的底部渲染。
默认实现显示了 [`<LastUpdated />`](#lastupdated)、[`<Pagination />`](#pagination) 和 [`<EditLink />`](#editlink)。

#### `LastUpdated`

**默认组件：** [`LastUpdated.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LastUpdated.astro)

在页脚中渲染的组件，用于显示最后更新日期。

#### `EditLink`

**默认组件：** [`EditLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/EditLink.astro)

在页脚中渲染的组件，用于显示指向页面编辑地址的链接。

#### `Pagination`

**默认组件：** [`Pagination.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Pagination.astro)

在页脚中渲染的组件，用于显示上一页和下一页的导航箭头。
