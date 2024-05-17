---
title: 在 Markdown 中创作内容
description: Starlight 支持的 Markdown 语法概述
---

Starlight 支持在 `.md` 文件中使用完整的 [Markdown](https://daringfireball.net/projects/markdown/) 语法，以及使用 [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f) 定义metadata 元数据，例如标题和描述。

如果使用这些文件格式，请务必检查 [MDX 文档](https://mdxjs.com/docs/what-is-mdx/#markdown) 或 [Markdoc 文档](https://markdoc.dev/docs/syntax)，因为 Markdown 的支持和用法可能会有所不同。

## Frontmatter

你可以通过设置 frontmatter 中的值来自定义 Starlight 里每个页面。
Frontmatter 是你的文件顶部在 `---` 中间的部分。

```md title="src/content/docs/example.md"
---
title: 我的页面标题
---

页面内容在第二个 `---` 后面。
```

每个页面都必须包含一个 `title`。
查看 [frontmatter 参考](/zh-cn/reference/frontmatter/) 了解所有可用字段以及如何添加自定义字段。

## 内联样式

文本可以是**粗体**，_斜体_，或~~删除线~~。

```md
文本可以是**粗体**，_斜体_，或~~删除线~~。
```

你可以 [链接到另一个页面](/zh-cn/getting-started/)。

```md
你可以 [链接到另一个页面](/zh-cn/getting-started/)。
```

你可以使用反引号高亮 `内联代码`。

```md
你可以使用反引号高亮 `内联代码`。
```

## 图片

Starlight 中的图片使用 [Astro 的内置优化资源支持](https://docs.astro.build/zh-cn/guides/assets/)。

Markdown 和 MDX 支持用于显示图片的 Markdown 语法，其中包括屏幕阅读器和辅助技术的 alt-text。

![一个星球和星星的插图，上面写着“astro”](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
![一个星球和星星的插图，上面写着“astro”](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

对于在项目中本地存储的图片，也支持图片的相对路径。

```md
// src/content/docs/page-1.md

![A rocketship in space](../../assets/images/rocket.svg)
```

## 标题

你可以使用标题来组织内容。Markdown 中的标题由行首的 `#` 数量来表示。

### 如何在 Starlight 中组织页面内容

Starlight 配置为自动使用页面标题作为一级标题，并将在每个页面的目录中包含一个“概述”标题。我们建议每个页面都从常规段落文本内容开始，并从 `<h2>` 开始使用页面标题：

```md
---
title: Markdown 指南
description: 如何在 Starlight 中使用 Markdown
---

本页面描述了如何在 Starlight 中使用 Markdown。

## 内联样式

## 标题
```

### 自动生成标题锚点链接

使用 Markdown 中的标题将自动为你提供锚点链接，以便你可以直接链接到页面的某些部分：

```md
---
title: 我的页面内容
description: 如何使用 Starlight 内置的锚点链接
---

## 介绍

我可以链接到同一页下面的[结论](#结论)。

## 结论

`https://my-site.com/page1/#introduction` 直接导航到我的介绍。
```

二级标题 (`<h2>`) 和 三级标题 (`<h3>`) 将自动出现在页面目录中。

在 [Astro 文档](https://docs.astro.build/zh-cn/guides/markdown-content/#标题-id)中了解 Astro 是如何处理标题 `id` 的。

## 旁白

旁白（也称为“警告”或“标注”）对于在页面的主要内容旁边显示辅助信息很有用。

Starlight 提供了一个自定义的 Markdown 语法来渲染旁白。旁白块使用一对三个冒号 `:::` 来包裹你的内容，并且可以是 `note`，`tip`，`caution` 或 `danger` 类型。

你可以在旁白中嵌套任何其他 Markdown 内容类型，但旁白最适合用于简短而简洁的内容块。

### Note 旁白

:::note
Starlight 是一个使用 [Astro](https://astro.build/) 构建的文档网站工具包。 你可以使用此命令开始：

```sh
npm create astro@latest -- --template starlight
```

:::

````md
:::note
Starlight 是一个使用 [Astro](https://astro.build/) 构建的文档网站工具包。 你可以使用此命令开始：

```sh
npm create astro@latest -- --template starlight
```

:::
````

### 自定义旁白标题

你可以在旁白类型后面的方括号中指定旁白的自定义标题，例如 `:::tip[你知道吗？]`。

:::tip[你知道吗？]
Astro 帮助你使用 [“群岛架构”](https://docs.astro.build/zh-cn/concepts/islands/) 构建更快的网站。
:::

```md
:::tip[你知道吗？]
Astro 帮助你使用 [“群岛架构”](https://docs.astro.build/zh-cn/concepts/islands/) 构建更快的网站。
:::
```

### 更多旁白类型

Caution 和 danger 旁白有助于吸引用户注意可能绊倒他们的细节。 如果你发现自己经常使用这些，这也可能表明你正在记录的内容可以从重新设计中受益。

:::caution
如果你不确定是否想要一个很棒的文档网站，请在使用 [Starlight](/zh-cn/) 之前三思。
:::

:::danger
借助有用的 Starlight 功能，你的用户可能会提高工作效率，并发现你的产品更易于使用。

- 清晰的导航
- 用户可配置的颜色主题
- [i18n 支持](/zh-cn/guides/i18n/)

:::

```md
:::caution
如果你不确定是否想要一个很棒的文档网站，请在使用 [Starlight](/zh-cn/) 之前三思。
:::

:::danger
借助有用的 Starlight 功能，你的用户可能会提高工作效率，并发现你的产品更易于使用。

- 清晰的导航
- 用户可配置的颜色主题
- [i18n 支持](/zh-cn/guides/i18n/)

:::
```

## 块引用

> 这是块引用，通常在引用其他人或文档时使用。
>
> 块引用以每行开头的 `>` 表示。

```md
> 这是块引用，通常在引用其他人或文档时使用。
>
> 块引用以每行开头的 `>` 表示。
```

## 代码块

代码块由三个反引号 <code>```</code> 开始和结束。你可以在开头的反引号后指定代码块的编程语言。

```js
// 带有语法高亮的 JavaScript 代码。
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// 带有语法高亮的 JavaScript 代码。
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

```md
长单行代码块不应换行。如果它们太长，它们应该水平滚动。这一行应该足够长长长长长长长长长长长长来证明这一点。
```

### Expressive Code 功能

Starlight 使用 [Expressive Code](https://github.com/expressive-code/expressive-code/tree/main/packages/astro-expressive-code) 来扩展代码块的格式化功能。
Expressive Code 的文本标记和窗口外框插件是默认启用的。
可以使用 Starlight 的 [`expressiveCode` 配置选项](/zh-cn/reference/configuration/#expressivecode) 来配置代码块的渲染。

#### 文本标记

你可以通过在代码块的起始行上使用 [Expressive Code 文本标记 (text markers)](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#usage-in-markdown--mdx-documents) 在代码块里突出显示特定行或代码块的一部分。
使用大括号 (`{ }`) 来突出显示整行，使用引号来突出显示文本字符串。

有三种突出显示样式：中性用于突出显示代码，绿色用于表示插入的代码，红色用于表示删除的代码。
字符串和整行都可以使用默认的标记，也可以与 `ins=` 和 `del=` 结合使用产生所需的突出显示效果。

Expressive Code 提供了几种自定义你的代码示例视觉外观的选项。
其中许多可以组合使用，以获得极具说明性的代码示例。
请探索 [Expressive Code 文档](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md) 以了解可用的众多设置项。
下面显示了一些最常见的示例：

- [使用 `{ }` 标记标出整行和行范围](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-entire-lines--line-ranges):

  ```js {2-3}
  function demo() {
    // 这一行 (#2) 以及下一行被高亮显示
    return '这是本代码段的第 3 行';
  }
  ```

  ````md
  ```js {2-3}
  function demo() {
    // 这一行 (#2) 以及下一行被高亮显示
    return '这是本代码段的第 3 行';
  }
  ```
  ````

- [使用 `" "` 标记或正则表达式标出文本字符串](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-individual-text-inside-lines):

  ```js "单个词组" /甚.*表达式/
  // 单个词组也能被高亮显示
  function demo() {
    return '甚至支持使用正则表达式';
  }
  ```

  ````md
  ```js "单个词组" /甚.*表达式/
  // 单个词组也能被高亮显示
  function demo() {
    return '甚至支持使用正则表达式';
  }
  ```
  ````

- [使用 `ins` 或 `del` 来标记行或文本为插入或删除](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#selecting-marker-types-mark-ins-del):

  ```js "return true;" ins="插入" del="删除"
  function demo() {
    console.log('这是插入以及删除类型的标记');
    // 返回语句使用默认标记类型
    return true;
  }
  ```

  ````md
  ```js "return true;" ins="插入" del="删除"
  function demo() {
    console.log('这是插入以及删除类型的标记');
    // 返回语句使用默认标记类型
    return true;
  }
  ```
  ````

- [混合语法高亮和类 `diff` 语法](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#combining-syntax-highlighting-with-diff-like-syntax):

  ```diff lang="js"
    function thisIsJavaScript() {
      // 这整个代码块都会被作为 JavaScript 高亮
      // 而且我们还可以给它添加 diff 标记！
  -   console.log('旧的不去')
  +   console.log('新的不来')
    }
  ```

  ````md
  ```diff lang="js"
    function thisIsJavaScript() {
      // 这整个代码块都会被作为 JavaScript 高亮
      // 而且我们还可以给它添加 diff 标记！
  -   console.log('旧的不去')
  +   console.log('新的不来')
    }
  ```
  ````

#### 边框和标题

代码块可以在类似窗口的框架中呈现。
默认情况下 shell 脚本语言（例如 `bash` 或 `sh`）会使用一个看起来像终端窗口的边框。
其他语言在提供了标题的情况下会在一个看起来像代码编辑器的边框中显示。

一个代码块的可选标题可以通过在代码块的开头反引号后面添加一个 `title="..."` 属性来设置，或者通过在代码的前几行中添加一个文件名注释来设置。

- [通过注释添加一个文件名标签](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

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

- [给终端窗口添加一个标题](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```bash title="安装依赖…"
  npm install
  ```

  ````md
  ```bash title="安装依赖…"
  npm install
  ```
  ````

- [使用 `frame="none"` 禁用边框](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#overriding-frame-types)

  ```bash frame="none"
  echo "这个代码块即使使用了 bash 语言也不会被显示成终端窗口"
  ```

  ````md
  ```bash frame="none"
  echo "这个代码块即使使用了 bash 语言也不会被显示成终端窗口"
  ```
  ````

## 其它通用 Markdown 语法

Starlight 支持所有其他 Markdown 语法，例如列表和表格。 请参阅 [Markdown 指南的 Markdown 速查表](https://www.markdownguide.org/cheat-sheet/) 以快速了解所有 Markdown 语法元素。

## 高级 Markdown 和 MDX 配置

Starlight 使用 Astro 的 Markdown 和 MDX 渲染器，该渲染器构建在 remark 和 rehype 之上。 你可以通过在 Astro 配置文件中添加 `remarkPlugins` 或 `rehypePlugins` 来添加对自定义语法和行为的支持。 请参阅 Astro 文档中的 [“配置 Markdown 和 MDX”](https://docs.astro.build/zh-cn/guides/markdown-content/#配置-markdown) 以了解更多信息。
