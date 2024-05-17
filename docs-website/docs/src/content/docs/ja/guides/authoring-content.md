---
title: Markdownでのコンテンツ作成
description: StarlightがサポートするMarkdown構文の概要。
---

Starlightでは、`.md`ファイルにおいて[Markdown](https://daringfireball.net/projects/markdown/)構文のすべての機能を利用できます。また、タイトルや説明文（description）などのメタデータを定義するためのフロントマター[YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f)もサポートしています。

MDXやMarkdocを使用する場合、サポートされるMarkdownの機能や使用方法が異なることがあるため、[MDXドキュメント](https://mdxjs.com/docs/what-is-mdx/#markdown)や[Markdocドキュメント](https://markdoc.dev/docs/syntax)を必ず確認してください。

## フロントマター

フロントマターに値を設定して、Starlightの個々のページをカスタマイズできます。フロントマターは、ファイル先頭の`---`によって区切られた区間に設定します。

```md title="src/content/docs/example.md"
---
title: ページのタイトル
---

ページのコンテンツは、2つ目の`---`の後に続きます。
```

すべてのページには、少なくとも`title`が必要です。利用可能なすべてのフィールドと、カスタムフィールドの追加方法については、[フロントマターのリファレンス](/ja/reference/frontmatter/)を参照してください。

## インラインスタイル

テキストは**太字**、_斜体_、または~~取り消し線~~にできます。

```md
テキストは**太字**、_斜体_、または~~取り消し線~~にできます。
```

[別のページにリンク](/ja/getting-started/)できます。

```md
[別のページにリンク](/ja/getting-started/)できます。
```

バックティックで`インラインコード`を強調できます。

```md
バックティックで`インラインコード`を強調できます。
```

## 画像

Starlightは、[Astro組み込みのアセット最適化機能](https://docs.astro.build/ja/guides/assets/)を使用して画像を表示します。

MarkdownとMDXは、スクリーンリーダーや支援技術のための代替テキストを含む画像を表示するためのMarkdown構文をサポートしています。

![「astro」という単語を中心に据えた惑星と恒星のイラスト](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
![「astro」という単語を中心に据えた惑星と恒星のイラスト](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

プロジェクトにローカルに保存されている画像についても、相対パスを使用して表示できます。

```md
// src/content/docs/page-1.md

![宇宙空間に浮かぶロケット](../../assets/images/rocket.svg)
```

## 見出し

見出しによりコンテンツを構造化できます。Markdownの見出しは、行の先頭にある`#`の数で示されます。

### ページコンテンツを構造化する方法

Starlightは、ページタイトルをトップレベルの見出しとして自動的に使用し、また各ページの目次の先頭に「概要」という見出しを含めます。各ページを通常のテキストコンテンツで開始し、ページ上の見出しは`<h2>`以下を使用することをおすすめします。

```md
---
title: Markdownガイド
description: StarlightでのMarkdownの使い方
---

このページでは、StarlightでMarkdownを使用する方法について説明します。

## インラインスタイル

## 見出し
```

### 見出しの自動アンカーリンク

Markdownで見出しを使用するとアンカーリンクが自動的に付与されるため、ページの特定のセクションに直接リンクできます。

```md
---
title: 私のページコンテンツ
description: Starlightの組み込みアンカーリンクの使い方
---

## はじめに

同じページの下部にある[結論](#結論)にリンクできます。

## 結論

`https://my-site.com/page1/#はじめに`は、「はじめに」に直接移動します。
```

レベル2（`<h2>`）とレベル3（`<h3>`）の見出しは、ページの目次に自動的に表示されます。

Astroが見出しの`id`をどのように処理するかについて、詳しくは[Astroドキュメント](https://docs.astro.build/ja/guides/markdown-content/#見出しid)を参照してください。

## 補足情報

補足情報（「警告」や「吹き出し」とも呼ばれます）は、ページのメインコンテンツと並べて補助的な情報を表示するのに便利です。

Starlightは、補足情報をレンダリングするためのカスタムMarkdown構文を提供しています。補足情報のブロックは、コンテンツを囲む3つのコロン`:::`によって示し、`note`（注釈）、`tip`（ヒント）、`caution`（注意）、`danger`（危険）というタイプに設定できます。

他のMarkdownコンテンツを補足情報の中にネストできますが、補足情報は短く簡潔なコンテンツに最も適しています。

### 注釈

:::note
Starlightは、[Astro](https://astro.build/)製のドキュメントサイト用ツールキットです。次のコマンドではじめられます。

```sh
npm create astro@latest -- --template starlight
```

:::

````md
:::note
Starlightは、[Astro](https://astro.build/)製のドキュメントサイト用ツールキットです。次のコマンドではじめられます。

```sh
npm create astro@latest -- --template starlight
```

:::
````

### カスタムのタイトル

補足情報のカスタムタイトルを、補足情報のタイプの後ろに角括弧で囲んで指定できます。たとえば`:::tip[ご存知でしたか？]`のようにできます。

:::tip[ご存知でしたか？]
Astroでは[「アイランドアーキテクチャ」](https://docs.astro.build/ja/concepts/islands/)を使用して、より高速なWebサイトを構築できます。
:::

```md
:::tip[ご存知でしたか？]
Astroでは[「アイランドアーキテクチャ」](https://docs.astro.build/ja/concepts/islands/)を使用して、より高速なWebサイトを構築できます。
:::
```

### その他のタイプ

注意（Caution）と危険（Danger）の補足は、ユーザーがつまずく可能性のある細かい点に注意を向けさせるのに役立ちます。もしこれらを多用しているとすれば、それはあなたがドキュメントを書いている対象の設計を見直す余地があることのサインかもしれません。

:::caution
もしあなたが素晴らしいドキュメントサイトを望んでいないのであれば、[Starlight](/ja/)は不要かもしれません。
:::

:::danger
Starlightの便利な機能のおかげで、ユーザーはより生産的になり、プロダクトはより使いやすくなるかもしれません。

- わかりやすいナビゲーション
- ユーザーが設定可能なカラーテーマ
- [国際化機能](/ja/guides/i18n/)

:::

```md
:::caution
もしあなたが素晴らしいドキュメントサイトを望んでいないのであれば、[Starlight](/ja/)は不要かもしれません。
:::

:::danger
Starlightの便利な機能のおかげで、ユーザーはより生産的になり、プロダクトはより使いやすくなるかもしれません。

- わかりやすいナビゲーション
- ユーザーが設定可能なカラーテーマ
- [国際化機能](/ja/guides/i18n/)

:::
```

## 引用

> これは引用です。他の人の言葉やドキュメントを引用するときによく使われます。
>
> 引用は、各行の先頭に`>`を付けることで示されます。

```md
> これは引用です。他の人の言葉やドキュメントを引用するときによく使われます。
>
> 引用は、各行の先頭に`>`を付けることで示されます。
```

## コード

コードのブロックは、先頭と末尾に3つのバックティック<code>```</code>を持つブロックで示されます。コードブロックを開始するバックティックの後ろに、使用されているプログラミング言語を指定できます。

```js
// シンタックスハイライトされたJavascriptコード。
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// シンタックスハイライトされたJavascriptコード。
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

### Expressive Code機能

Starlightは、コードブロックのフォーマットを拡張するために[Expressive Code](https://github.com/expressive-code/expressive-code/tree/main/packages/astro-expressive-code)を使用しています。Expressive Codeのテキストマーカーとウィンドウフレームプラグインはデフォルトで有効になっています。コードブロックのレンダリングは、Starlightの[`expressiveCode`設定オプション](/ja/reference/configuration/#expressivecode)により設定できます。

#### テキストマーカー

[Expressive Codeのテキストマーカー](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#usage-in-markdown--mdx-documents)をコードブロックの先頭で使うことで、コードブロックの特定の行や部分をハイライトできます。波括弧（`{ }`）を使って行全体をハイライトし、引用符を使ってテキストの文字列をハイライトします。

ハイライトのスタイルは3つあります。コードに注意を向けるための中立的なスタイル、挿入されたコードを示す緑色のスタイル、削除されたコードを示す赤色のスタイルです。テキストと行全体の両方を、デフォルトのマーカー、または`ins=`と`del=`を組み合わせてマークし、目的のハイライトを生成できます。

Expressive Codeには、コードサンプルの外観をカスタマイズするためのさまざまなオプションが用意されています。これらの多くは組み合わせることができ、非常に明快なコードサンプルを作成できます。利用可能な多くのオプションについては、[Expressive Codeのドキュメント](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md)を確認してください。最も一般的な例をいくつか以下に示します。

- [行全体と行の範囲を`{ }`マーカーを使ってマークする](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-entire-lines--line-ranges):

  ```js {2-3}
  function demo() {
    // この行（2行目）と次の行はハイライトされます
    return 'このスニペットの3行目です';
  }
  ```

  ````md
  ```js {2-3}
  function demo() {
    // この行（2行目）と次の行はハイライトされます
    return 'このスニペットの3行目です';
  }
  ```
  ````

- [`" "`マーカーまたは正規表現を使って選択されたテキストをマークする](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-individual-text-inside-lines):

  ```js "個別の用語" /正規表現.*います/
  // 個別の用語もハイライトできます
  function demo() {
    return '正規表現もサポートされています';
  }
  ```

  ````md
  ```js "個別の用語" /正規表現.*います/
  // 個別の用語もハイライトできます
  function demo() {
    return '正規表現もサポートされています';
  }
  ```
  ````

- [追加、削除されたテキストや行を、`ins`と`del`でマークする](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#selecting-marker-types-mark-ins-del):

  ```js "return true;" ins="挿入" del="削除"
  function demo() {
    console.log('これらは挿入と削除のマーカーです');
    // return文はデフォルトのマーカータイプを使用します
    return true;
  }
  ```

  ````md
  ```js "return true;" ins="挿入" del="削除"
  function demo() {
    console.log('これらは挿入と削除のマーカーです');
    // return文はデフォルトのマーカータイプを使用します
    return true;
  }
  ```
  ````

- [構文ハイライトと`diff`風の構文を組み合わせる](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#combining-syntax-highlighting-with-diff-like-syntax):

  ```diff lang="js"
    function thisIsJavaScript() {
      // このブロック全体はJavaScriptとしてハイライトされますが、
      // diffマーカーの追加も可能です！
  -   console.log('削除される古いコード')
  +   console.log('新しいキラキラコード！')
    }
  ```

  ````md
  ```diff lang="js"
    function thisIsJavaScript() {
      // このブロック全体はJavaScriptとしてハイライトされますが、
      // diffマーカーの追加も可能です！
  -   console.log('削除される古いコード')
  +   console.log('新しいキラキラコード！')
    }
  ```
  ````

#### フレームとタイトル

コードブロックをウィンドウのようなフレームの中にレンダリングできます。シェルスクリプト言語（`bash`や`sh`など）には、ターミナルウィンドウのようなフレームが使用されます。その他の言語は、タイトルを含んでいる場合、コードエディタスタイルのフレーム内に表示されます。

`title="..."`属性を、コードブロックの開始を表わすバックティックと言語識別子の後ろに続けて記述するか、コードの最初の行にファイル名コメントを記述することで、コードブロックにオプションでタイトルを設定できます。

- [コメントによりファイル名タブを追加する](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

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

- [Terminalウィンドウにタイトルを追加する](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```bash title="依存関係のインストール中…"
  npm install
  ```

  ````md
  ```bash title="依存関係のインストール中…"
  npm install
  ```
  ````

- [`frame="none"`によりウィンドウフレームを無効化する](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#overriding-frame-types)

  ```bash frame="none"
  echo "bash言語を使用していますが、これはターミナルとしてレンダリングされません"
  ```

  ````md
  ```bash frame="none"
  echo "bash言語を使用していますが、これはターミナルとしてレンダリングされません"
  ```
  ````

## その他のMarkdown機能

Starlightは、リストやテーブルなど、その他のMarkdown記法をすべてサポートしています。Markdownのすべての構文要素の概要については、[The Markdown GuideのMarkdownチートシート](https://www.markdownguide.org/cheat-sheet/)を参照してください。

## 高度なMarkdownとMDXの設定

Starlightは、remarkとrehypeをベースとした、AstroのMarkdown・MDXレンダラーを使用しています。Astroの設定ファイルに`remarkPlugins`または`rehypePlugins`を追加することで、カスタム構文や動作をサポートできます。詳しくは、Astroドキュメントの[「MarkdownとMDXの設定」](https://docs.astro.build/ja/guides/markdown-content/#markdownとmdxの設定)を参照してください。
