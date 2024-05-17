---
title: 環境に優しいドキュメント
description: Starlightがどのように環境に優しいドキュメントサイトを作成し、カーボンフットプリントを減らすのに役立つかについて学びます。
---

ウェブ産業が気候に与える影響は、世界の炭素排出量の[2％][sf]から[4％][bbc]であり、航空業界の排出量にほぼ匹敵すると推定されています。ウェブサイトの生態学的影響の計算には多くの複雑な要因がありますが、このガイドではドキュメントサイトの環境への負荷を減らすためのいくつかのヒントを紹介します。

幸いなことに、Starlightを選ぶことは素晴らしいスタートです。Website Carbon Calculatorによると、このサイトは[テスト対象のウェブページの99％よりもクリーン][sl-carbon]であり、ページ訪問あたり0.01gのCO₂を生成します。

## ページの重さ

ウェブページが転送するデータが多いほど、より多くのエネルギー資源が必要になります。2023年4月現在、[HTTPアーカイブのデータ][http]によると、中央値のウェブページでは、ユーザーは2,000KB以上のデータをダウンロードする必要がありました。

Starlightは、可能な限り軽量なページを作成します。たとえば、初回訪問時にユーザーがダウンロードする圧縮データは50KB未満であり、HTTPアーカイブの中央値のわずか2.5％にすぎません。優れたキャッシュ戦略により、後続のナビゲーションのダウンロードは10KB程度に抑えられます。

### 画像

Starlightは優れたベースラインを提供しますが、ドキュメントページに追加した画像はページの重さを急速に増加させます。Starlightは、Astroの[最適化されたアセットサポート][assets]を使用して、MarkdownとMDXファイル内のローカル画像を最適化します。

### UIコンポーネント

ReactやVueなどのUIフレームワークで作成されたコンポーネントを使うと、大量のJavaScriptがページに追加される可能性があります。StarlightはAstro上に構築されており、[Astroアイランド][islands]のおかげで、このようなコンポーネントがデフォルトで**ロードするクライアントサイドJavaScriptはゼロ**となります。

### キャッシュ

キャッシュは、ブラウザがすでにダウンロードしたデータをどのくらいの期間保存して再利用するかを制御するために使用されます。優れたキャッシュ戦略は、コンテンツが変更されたときにユーザーができるだけ早く新しいコンテンツを取得することを保証し、また変更されていないときに同じコンテンツを何度も無駄にダウンロードするのを避けます。

キャッシュを設定する最も一般的な方法は、[`Cache-Control` HTTPヘッダー][cache]を使用するものです。Starlightを使用する場合、`/_astro/`ディレクトリ内のすべてに対し長めのキャッシュ時間を設定できます。このディレクトリには、CSS、JavaScript、その他のバンドルされたアセットが含まれており、安全に永続的なキャッシュを設定できるため、不要なダウンロードを減らすことができます。

```
Cache-Control: public, max-age=604800, immutable
```

キャッシュの設定方法は、ウェブホストによって異なります。たとえば、Vercelは設定不要でこのキャッシュ戦略を適用しますが、Netlifyではプロジェクトに`public/_headers`ファイルを追加することで[Netlify用カスタムヘッダー][ntl-headers]を設定できます。

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## 電力消費

ウェブページをどのようにビルドするかは、ユーザーのデバイス上で実行するために必要な電力に影響します。Starlightは、最小限のJavaScriptを使用することで、ユーザーの携帯電話、タブレット、そしてコンピューターがページをロードしてレンダリングするために必要な処理量を削減します。

アナリティクスのトラッキングスクリプトのような機能や、ビデオ埋め込みのようなJavaScriptを多用するコンテンツを追加すると、ページの電力消費量が増加する可能性があるため注意してください。アナリティクスが必要な場合は、[Cabin][cabin]、[Fathom][fathom]、あるいは[Plausible][plausible]のような軽量なオプションを選択することを検討してください。YouTubeやVimeoのような埋め込みは、[ユーザーの操作に応じて動画をロードする][lazy-video]ことで改善できます。[`astro-embed`][embed]のようなパッケージは、一般的なサービスに対し有効です。

:::tip[ご存知でしたか？]
JavaScriptの解析とコンパイルは、ブラウザが実行する最も高コストなタスクの1つです。同じサイズのJPEG画像をレンダリングするのに比べ、[JavaScriptの処理には30倍以上の時間がかかる][cost-of-js]ことがあります。
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## ホスティング

ウェブページをホスティングする場所は、あなたのドキュメントサイトがどれだけ環境に優しいかどうかに大きく影響します。データセンターやサーバーファームは、高い電力消費や水の集中的な使用など、生態系に大きな影響を与える可能性があります。

再生可能エネルギーを使用するホストを選択すれば、サイトの炭素排出量を少なくできます。[Green Web Directory][gwb]は、ホスティング会社を見つけるのに役立つツールの1つです。

[gwb]: https://www.thegreenwebfoundation.org/directory/

## 比較

他のドキュメントフレームワークとの比較に興味がありますか？[Website Carbon Calculator][wcc]を用いた以下のテストでは、異なるツールで作成された類似のページを比較しています。

| フレームー枠                | ページ訪問ごとのCO₂ |
| --------------------------- | ------------------- |
| [Starlight][sl-carbon]      | 0.01g               |
| [VitePress][vp-carbon]      | 0.05g               |
| [Docus][dc-carbon]          | 0.05g               |
| [Sphinx][sx-carbon]         | 0.07g               |
| [MkDocs][mk-carbon]         | 0.10g               |
| [Nextra][nx-carbon]         | 0.11g               |
| [docsify][dy-carbon]        | 0.11g               |
| [Docusaurus][ds-carbon]     | 0.24g               |
| [Read the Docs][rtd-carbon] | 0.24g               |
| [GitBook][gb-carbon]        | 0.71g               |

<small>データは2023年5月14日に収集されたものです。リンクをクリックすると、最新の数値が表示されます。</small>

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

## その他のリソース

### ツール

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Carbon Control](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### 記事と講演

- [“Building a greener web”](https://youtu.be/EfPoOt7T5lg)、Michelle Barkerによる講演
- [“Sustainable Web Development Strategies Within An Organization”](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/)、Michelle Barkerによる記事
- [“A sustainable web for everyone”](https://2021.stateofthebrowser.com/speakers/tom-greenwood/)、Tom Greenwoodによる講演
- [“How Web Content Can Affect Power Usage”](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/)、Benjamin PoulainとSimon Fraserによる記事

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/ja/guides/assets/
[islands]: https://docs.astro.build/ja/concepts/islands/
[wcc]: https://www.websitecarbon.com/
