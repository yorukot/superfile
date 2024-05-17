---
title: Frontmatter 参考
description: Starlight 支持的默认 frontmatter 字段的概述。
---

你可以通过设置 frontmatter 中的值来自定义 Starlight 中的单个 Markdown 和 MDX 页面。例如，一个常规页面可能会设置 `title` 和 `description` 字段：

```md {3-4}
---
# src/content/docs/example.md
title: 关于此项目
description: 了解更多关于此项目的信息。
---

欢迎来到关于页面！
```

## Frontmatter 字段

### `title` (必填)

**类型：** `string`

你必须为每个页面提供标题。它将显示在页面顶部、浏览器标签中和页面元数据中。

### `description`

**类型：** `string`

页面描述用于页面元数据，将被搜索引擎和社交媒体预览捕获。

### `slug`

**类型：** `string`

覆盖页面的slug。有关更多详细信息，请参阅 Astro文档中的 [ "定义自定义slugs"](https://docs.astro.build/zh-cn/guides/content-collections/#定义自定义-slugs) 部分。

### `editUrl`

**类型：** `string | boolean`

覆盖[全局 `editLink` 配置](/zh-cn/reference/configuration/#editlink)。设置为 false 可禁用特定页面的 “编辑页面” 链接，或提供此页面内容可编辑的备用 URL。

### `head`

**类型：** [`HeadConfig[]`](/zh-cn/reference/configuration/#headconfig)

你可以使用 `<head>` frontmatter 字段向页面的`<head>`添加其他标签。这意味着你可以将自定义样式、元数据或其他标签添加到单个页面。类似于[全局 `head` 选项](/zh-cn/reference/configuration/#head)。

```md
---
# src/content/docs/example.md
title: 关于我们
head:
  # 使用自定义 <title> 标签
  - tag: title
    content: 自定义关于我们页面标题
---
```

### `tableOfContents`

**类型：** `false | { minHeadingLevel?: number; maxHeadingLevel?: number; }`

覆盖[全局 `tableOfContents` 配置](/zh-cn/reference/configuration/#tableofcontents)。自定义要包含的标题级别，或设置为 `false` 以在此页面上隐藏目录。

```md
---
# src/content/docs/example.md
title: 目录中只有 H2 的页面
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 2
---
```

```md
---
# src/content/docs/example.md
title: 没有目录的页面
tableOfContents: false
---
```

### `template`

**类型：** `'doc' | 'splash'`  
**默认值：** `'doc'`

为页面选择布局模板。
页面默认使用 `'doc'` 布局。
设置为 `'splash'` 以使用没有任何侧边栏的更宽的布局，该布局专为落地页设计。

### `hero`

**类型：** [`HeroConfig`](#heroconfig)

添加一个 hero 组件到页面顶部。与 `template: splash` 配合使用效果更佳。

例如，此配置显示了一些常见选项，包括从你的仓库加载图像。

```md
---
# src/content/docs/example.md
title: 我的主页
template: splash
hero:
  title: '我的项目： Stellar Stuff Sooner'
  tagline: 把你的东西带到月球上，眨眼间又回来。
  image:
    alt: 一个闪闪发光、色彩鲜艳的标志
    file: ../../assets/logo.png
  actions:
    - text: 告诉我更多
      link: /getting-started/
      icon: right-arrow
      variant: primary
    - text: 在 GitHub 上查看
      link: https://github.com/astronaut/my-project
      icon: external
      attrs:
        rel: me
---
```

你可以在浅色和深色模式下显示不同版本的 hero 图像。

```md
---
# src/content/docs/example.md
hero:
  image:
    alt: 一个闪闪发光、色彩鲜艳的 logo
    dark: ../../assets/logo-dark.png
    light: ../../assets/logo-light.png
---
```

#### `HeroConfig`

```ts
interface HeroConfig {
  title?: string;
  tagline?: string;
  image?:
    | {
        // 你的仓库中的图像的相对路径。
        file: string;
        // 使图像对辅助技术可访问的 Alt 文本
        alt?: string;
      }
    | {
        // 使用深色模式的图像的相对路径。
        dark: string;
        // 使用浅色模式的图像的相对路径。
        light: string;
        // 使图像对辅助技术可访问的 Alt 文本
        alt?: string;
      }
    | {
        // 用于图像插槽的原始 HTML 。
        // 可以是自定义的 `<img>` 标签或内联的 `<svg>`。
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

**类型：** `{ content: string }`

在此页面顶部显示公告横幅。

`content` 的值可以包含链接或其他内容的 HTML。
例如，此页面显示了一个横幅，其中包含指向 `example.com` 的链接。

```md
---
# src/content/docs/example.md
title: 带有横幅的页面
banner:
  content: |
    我们刚刚发布了一下非常酷的东西！
    <a href="https://example.com">点击查看！</a>
---
```

### `lastUpdated`

**类型：** `Date | boolean`

覆盖[全局 `lastUpdated` 配置](/zh-cn/reference/configuration/#lastupdated)。如果指定了日期，它必须是有效的 [YAML 时间戳](https://yaml.org/type/timestamp.html)，并将覆盖存储在 Git 历史记录中的此页面的日期。

```md
---
# src/content/docs/example.md
title: 带有自定义更新日期的页面
lastUpdated: 2022-08-09
---
```

### `prev`

**类型：** `boolean | string | { link?: string; label?: string }`

覆盖[全局 `pagination` 配置](/zh-cn/reference/configuration/#pagination)。如果指定了字符串，则将替换生成的链接文本；如果指定了对象，则将同时覆盖链接和文本。

```md
---
# src/content/docs/example.md
# 隐藏上一页链接
prev: false
---
```

```md
---
# src/content/docs/example.md
# 将上一页链接更改为“继续教程”
prev: 继续教程
---
```

```md
---
# src/content/docs/example.md
# 同时覆盖上一页的链接和文本
prev:
  link: /unrelated-page/
  label: 一个不相关的页面
---
```

### `next`

**类型：** `boolean | string | { link?: string; label?: string }`

和 [`prev`](#prev) 一样，但是用于下一页链接。

```md
---
# src/content/docs/example.md
# 隐藏下一页链接
next: false
---
```

### `pagefind`

**类型：** `boolean`  
**默认值：** `true`

设置此页面是否应包含在 [Pagefind](https://pagefind.app/) 搜索索引中。设置为 `false` 以从搜索结果中排除页面：

```md
---
# src/content/docs/example.md
# 在搜索索引中隐藏此页面
pagefind: false
---
```

### `draft`

**类型：** `boolean`  
**默认值：** `false`

设置此页面是否应被视为草稿，并且不包含在 [生产版本](https://docs.astro.build/zh-cn/reference/cli-reference/#astro-build) 和 [自动生成的链接组](/zh-cn/guides/sidebar/#自动生成的分组) 中。设置为 `true` 可将页面标记为草稿，并使其仅在开发过程中可见。

```md
---
# src/content/docs/example.md
# 从生产版本中排除此页面
draft: true
---
```

### `sidebar`

**类型：** [`SidebarConfig`](#sidebarconfig)

在使用自动生成的链接分组时，控制如何在[侧边栏](/zh-cn/reference/configuration/#sidebar)中显示此页面。

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

**类型：** `string`  
**默认值：** 页面 [`title`](#title-必填)

在自动生成的链接组中显示时，设置侧边栏中此页面的标签。

```md
---
# src/content/docs/example.md
title: 关于此项目
sidebar:
  label: About
---
```

#### `order`

**类型：** `number`

当对链接组进行自动生成排序时，控制此页面的顺序。
数字越小，链接组中显示得越高。

```md
---
# src/content/docs/example.md
title: 要首先显示的页面
sidebar:
  order: 1
---
```

#### `hidden`

**类型：** `boolean`
**默认值：** `false`

防止此页面包含在自动生成的侧边栏组中。

```md
---
# src/content/docs/example.md
title: 从自动生成的侧边栏中隐藏的页面
sidebar:
  hidden: true
---
```

#### `badge`

**类型：** <code>string | <a href="/zh-cn/reference/configuration/#badgeconfig">BadgeConfig</a></code>

当在自动生成的链接组中显示时，在侧边栏中为页面添加徽章。

当使用字符串时，徽章将显示为默认的强调色。可选择的，传递一个 [`BadgeConfig` 对象](/zh-cn/reference/configuration/#badgeconfig) ，其中包含 `text` 和 `variant` 字段，可以自定义徽章。

```md
---
# src/content/docs/example.md
title: 带有徽章的页面
sidebar:
  # 使用与你的网站的强调色相匹配的默认类型
  badge: 新增
---
```

```md
---
# src/content/docs/example.md
title: 带有徽章的页面
sidebar:
  badge:
    text: 实验性
    variant: caution
---
```

#### `attrs`

**类型：** `Record<string, string | number | boolean | undefined>`

给自动生成的侧边栏分组的链接添加的 HTML 属性。

```md
---
# src/content/docs/example.md
title: 新标签页中打开页面
sidebar:
  # 在新标签页中打开页面
  attrs:
    target: _blank
---
```

## 自定义 frontmatter schema

Starlight 的 `docs` 内容集合的 frontmatter schema 在 `src/content/config.ts` 中使用 `docsSchema()` 辅助函数进行配置：

```ts {3,6}
// src/content/config.ts
import { defineCollection } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({ schema: docsSchema() }),
};
```

了解更多关于内容集合模式的信息，请参阅 Astro 文档中的 [“定义集合模式”](https://docs.astro.build/zh-cn/guides/content-collections/#定义集合模式) 部分。

`docsSchema()` 采用以下选项：

### `extend`

**类型：** Zod schema 或者返回 Zod schema 的函数  
**默认值：** `z.object({})`

通过在 `docsSchema()` 选项中设置 `extend` 来使用其他字段扩展 Starlight 的 schema。
值应该是一个 [Zod schema](https://docs.astro.build/zh-cn/guides/content-collections/#用-zod-定义数据类型)。

在下面的示例中，我们为 `description` 提供了一个更严格的类型，使其成为必填项，并添加了一个新的可选的 `category` 字段：

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: z.object({
        // 将内置字段设置为必填项。
        description: z.string(),
        // 将新字段添加到 schema 中。
        category: z.enum(['tutorial', 'guide', 'reference']).optional(),
      }),
    }),
  }),
};
```

要利用 [Astro `image()` 辅助函数](https://docs.astro.build/zh-cn/guides/images/#内容集合中的图像)，请使用返回 schema 扩展的函数：

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: ({ image }) => {
        return z.object({
          // 添加一个必须解析为本地图像的字段。
          cover: image(),
        });
      },
    }),
  }),
};
```
