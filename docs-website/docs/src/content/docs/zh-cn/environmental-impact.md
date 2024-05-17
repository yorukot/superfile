---
title: 环保文档
description: 了解 Starlight 如何帮助你构建更环保的文档站点并减少碳足迹。
---

网络行业对气候影响的估计范围从[2%][sf]到[全球碳排放量的4%][bbc]，大致相当于航空业的排放量。
计算网站生态影响的因素很多，但本指南包含了一些减少文档站点环境足迹的技巧。

好消息是，选择 Starlight 是一个很好的开始。
根据网站碳计算器的数据，本站点[比 99% 的网页更环保][sl-carbon]，每个页面访问产生 0.01g 的 CO₂。

## 页面大小

网页传输的数据越多，它所需的能源资源就越多。
根据[来自 HTTP 存档的数据][http]，2023 年 4 月，中位数网页需要用户下载超过 2,000 KB。

Starlight 构建的页面尽可能轻量。
例如，在第一次访问时，用户将下载少于 50 KB 的压缩数据，仅为 HTTP 存档中位数的 2.5%。
通过良好的缓存策略，后续导航可以下载少至 10 KB。

### 图片

虽然 Starlight 提供了一个很好的基线，但是你添加到文档页面的图片可能会快速增加页面的大小。
Starlight 使用 Astro 的 [优化的资源支持][assets] 来优化 Markdown 和 MDX 文件中的本地图片。

### UI 组件

使用 React 或 Vue 等 UI 框架构建的组件可以轻松地向页面添加大量 JavaScript。
因为 Starlight 是基于 Astro 构建的，所以这样的组件默认情况下不会加载**任何客户端 JavaScript**，这要归功于 [Astro 岛屿][islands]。

### 缓存

缓存用于控制浏览器存储和重用已下载数据的时间。
良好的缓存策略确保用户在内容更改时尽快获得新内容，但也避免了在内容未更改时反复下载相同的内容。

最常见的配置缓存的方法是使用 [`Cache-Control` HTTP headers][cache]。
使用 Starlight 时，你可以为 `/_astro/` 目录中的所有内容设置长时间缓存。
该目录包含可以安全永久缓存的 CSS、JavaScript 和其他捆绑的资源，从而减少不必要的下载：

```
Cache-Control: public, max-age=604800, immutable
```

如何配置缓存取决于你的 Web 主机。例如，Vercel 会自动为你应用这种缓存策略，而 Netlify 可以通过在项目中添加 `public/_headers` 文件来[为 Netlify 设置自定义 headers][ntl-headers]：

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## 功耗

网页的制作方式会影响其在用户设备上的运行功率。
通过使用最少的 JavaScript，Starlight 减少了用户手机、平板电脑或电脑加载和呈现网页所需的处理能力。

当添加诸如分析跟踪脚本或 JavaScript 重型内容（如视频嵌入）之类的功能时，请注意这些功能会增加页面的功耗。
如果你需要分析，请考虑选择像 [Cabin][cabin]、[Fathom][fathom] 或 [Plausible][plausible] 这样的轻量级选项。
等待[用户交互时加载视频][lazy-video]可以改善 YouTube 和 Vimeo 视频等嵌入。
[`astro-embed`][embed] 等包可以帮助常见的服务。

:::tip[你知道吗？]
解析和编译 JavaScript 是浏览器必须执行的最昂贵的任务之一。
与渲染相同大小的 JPEG 图像相比，[JavaScript 处理所需的时间可能超过 JPEG 的 30 倍][cost-of-js]。
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## 托管

网页托管在哪里会对你的文档站点的环保程度产生很大的影响。
数据中心和服务器农场可能会对环境产生很大的影响，包括高耗电量和大量使用水资源。

选择使用可再生能源的主机将意味着你的站点的碳排放量更低。[绿色网络目录][gwb]是一个可以帮助你找到主机公司的工具。

[gwb]: https://www.thegreenwebfoundation.org/directory/

## 比较

好奇和其他文档框架相比如何？
下面使用 [Website Carbon Calculator][wcc] 的测试比较了使用不同工具构建的类似页面。

| 框架                        | 每页访问量产生 CO₂ |
| --------------------------- | ------------------ |
| [Starlight][sl-carbon]      | 0.01g              |
| [VitePress][vp-carbon]      | 0.05g              |
| [Docus][dc-carbon]          | 0.05g              |
| [Sphinx][sx-carbon]         | 0.07g              |
| [MkDocs][mk-carbon]         | 0.10g              |
| [Nextra][nx-carbon]         | 0.11g              |
| [docsify][dy-carbon]        | 0.11g              |
| [Docusaurus][ds-carbon]     | 0.24g              |
| [Read the Docs][rtd-carbon] | 0.24g              |
| [GitBook][gb-carbon]        | 0.71g              |

<small>数据收集于 2023 年 5 月 14 日。点击链接查看最新数据。</small>

[sl-carbon]: https://www.websitecarbon.com/website/starlight-astro-build-getting-started/
[vp-carbon]: https://www.websitecarbon.com/website/vitepress-dev-guide-what-is-vitepress/
[dc-carbon]: https://www.websitecarbon.com/website/docus-dev-introduction-getting-started/
[sx-carbon]: https://www.websitecarbon.com/website/sphinx-doc-org-en-master-usage-quickstart-html/
[mk-carbon]: https://www.websitecarbon.com/website/mkdocs-org-getting-started/
[nx-carbon]: https://www.websitecarbon.com/website/nextra-site-docs-docs-theme-start/
[dy-carbon]: https://www.websitecarbon.com/website/docsify-js-org/
[ds-carbon]: https://www.websitecarbon.com/website/docusaurus-io-docs/
[rtd-carbon]: https://www.websitecarbon.com/website/docs-readthedocs-io-en-stable-index-html/
[gb-carbon]: https://www.websitecarbon.com/website/docs-gitbook-com/

## 更多资源

### 工具

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Carbon Control](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### 文章和演讲

- [“Building a greener web”](https://youtu.be/EfPoOt7T5lg)，Michelle Barker 的演讲
- [“Sustainable Web Development Strategies Within An Organization”](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/)，Michelle Barker 的文章
- [“A sustainable web for everyone”](https://2021.stateofthebrowser.com/speakers/tom-greenwood/)，Tom Greenwood 的演讲
- [“How Web Content Can Affect Power Usage”](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/)，Benjamin Poulain 和 Simon Fraser 的文章

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/zh-cn/guides/assets/
[islands]: https://docs.astro.build/zh-cn/concepts/islands/
[wcc]: https://www.websitecarbon.com/
