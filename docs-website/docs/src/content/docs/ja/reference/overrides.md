---
title: オーバーライド
description: Starlightのオーバーライドでサポートされているコンポーネントとコンポーネントプロパティの概要。
tableOfContents:
  maxHeadingLevel: 4
---

Starlightの[`components`](/ja/reference/configuration/#components)設定オプションに置き換え対象のコンポーネントへのパスを指定することで、Starlightの組み込みコンポーネントをオーバーライドできます。このページでは、オーバーライド可能なすべてのコンポーネントと、GitHub上にあるコンポーネントのデフォルト実装へのリンクの一覧を記載しています。

[コンポーネントのオーバーライドガイド](/ja/guides/overriding-components/)も参照してください。

## コンポーネントprops

すべてのコンポーネントは、現在のページに関する情報を含んでいる、標準の`Astro.props`オブジェクトにアクセスできます。

カスタムコンポーネントに型を付けるには、Starlightから`Props`型をインポートします。

```astro
---
// src/components/Custom.astro
import type { Props } from '@astrojs/starlight/props';

const { hasSidebar } = Astro.props;
//      ^ type: boolean
---
```

これにより、`Astro.props`にアクセスする際、オートコンプリートと型が有効になります。

### Props

Starlightは、以下のpropsをカスタムコンポーネントに渡します。

#### `dir`

**Type:** `'ltr' | 'rtl'`

ページの書字方向。

#### `lang`

**Type:** `string`

このページのロケールのBCP-47言語タグ。たとえば`en`、`zh-CN`、`pt-BR`など。

#### `locale`

**Type:** `string | undefined`

言語が配信されるベースパス。ルートロケールスラグの場合は`undefined`となります。

#### `slug`

**Type:** `string`

コンテンツファイル名から生成されたページのスラグ。

#### `id`

**Type:** `string`

コンテンツファイル名に基づくページの一意のID。

#### `isFallback`

**Type:** `true | undefined`

このページが現在の言語で未翻訳であり、デフォルトロケールのフォールバックコンテンツを使用している場合は`true`となります。多言語サイトでのみ使用されます。

#### `entryMeta`

**Type:** `{ dir: 'ltr' | 'rtl'; lang: string }`

ページコンテンツのロケールメタデータ。ページがフォールバックコンテンツを使用している場合、トップレベルのロケール値とは異なる場合があります。

#### `entry`

現在のページのAstroコンテンツコレクションのエントリー。`entry.data`には、現在のページのフロントマターの値が含まれます。

```ts
entry: {
  data: {
    title: string;
    description: string | undefined;
    // その他の値
  }
}
```

このオブジェクトの構造については、[Astroのコレクションエントリー型](https://docs.astro.build/ja/reference/api-reference/#collection-entry-type)リファレンスを参照してください。

#### `sidebar`

**Type:** `SidebarEntry[]`

ページのサイトナビゲーション用サイドバーのエントリー。

#### `hasSidebar`

**Type:** `boolean`

ページにサイドバーを表示するかどうか。

#### `pagination`

**Type:** `{ prev?: Link; next?: Link }`

ページネーションの設定が有効な場合にサイドバーに表示される、前のページと次のページへのリンク。

#### `toc`

**Type:** `{ minHeadingLevel: number; maxHeadingLevel: number; items: TocItem[] } | undefined`

目次の設定が有効な場合、このページの目次。

#### `headings`

**Type:** `{ depth: number; slug: string; text: string }[]`

現在のページから抽出されたすべてのMarkdown見出しの配列。Starlightの設定オプションをもとに目次コンポーネントを作成したい場合は、[`toc`](#toc)を使用してください。

#### `lastUpdated`

**Type:** `Date | undefined`

最終更新日の設定が有効な場合、このページが最後に更新された日時を表わすJavaScriptの`Date`オブジェクト。

#### `editUrl`

**Type:** `URL | undefined`

ページの編集設定が有効な場合、このページを編集可能なアドレスの`URL`オブジェクト。

#### `labels`

**Type:** `Record<string, string>`

現在のページのローカライズされたUI文字列を含んだオブジェクト。利用可能なすべてのキーの一覧については、[「StarlightのUIを翻訳する」](/ja/guides/i18n/#starlightのuiを翻訳する)ガイドを参照してください。

---

## コンポーネント

### ヘッド

以下のコンポーネントは、各ページの`<head>`要素内にレンダリングされます。[`<head>`内に配置可能な要素](https://developer.mozilla.org/ja/docs/Web/HTML/Element/head#関連情報)のみを含めるようにしてください。

#### `Head`

**デフォルトコンポーネント:** [`Head.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Head.astro)

各ページの`<head>`内にレンダリングされるコンポーネント。`<title>`や`<meta charset="utf-8">`などの重要なタグが含まれます。

このコンポーネントをオーバーライドするのは最後の手段としてください。可能な限り、Starlightの設定オプション[`head`](/ja/reference/configuration/#head)を使用してください。

#### `ThemeProvider`

**デフォルトコンポーネント:** [`ThemeProvider.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeProvider.astro)

ダーク/ライトテーマのサポートを設定するための、`<head>`内にレンダリングされるコンポーネント。デフォルトの実装では、インラインスクリプトと[`<ThemeSelect />`](#themeselect)で使用される`<template>`が含まれています。

---

### アクセシビリティ

#### `SkipLink`

**デフォルトコンポーネント:** [`SkipLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SkipLink.astro)

アクセシビリティのために`<body>`内の最初の要素としてレンダリングされる、メインページコンテンツへのリンクのコンポーネント。デフォルトの実装では、ユーザーがキーボードでタブを押してフォーカスするまで非表示となります。

---

### レイアウト

以下のコンポーネントは、Starlightのコンポーネントのレイアウトと、異なるブレークポイント間のビューの管理を担当します。これらをオーバーライドすると著しく複雑になるため、可能な限り、より低レベルのコンポーネントをオーバーライドすることをおすすめします。

#### `PageFrame`

**デフォルトコンポーネント:** [`PageFrame.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageFrame.astro)

ページコンテンツの大部分をラップするレイアウトコンポーネント。デフォルトの実装では、ヘッダー・サイドバー・メインのレイアウトをセットし、`header`と`sidebar`の名前付きスロットと、メインコンテンツのデフォルトスロットを含みます。また、小さな（モバイル）ビューポートでのサイドバーナビゲーションの切り替えをサポートするために、[`<MobileMenuToggle />`](#mobilemenutoggle)をレンダリングします。

#### `MobileMenuToggle`

**デフォルトコンポーネント:** [`MobileMenuToggle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuToggle.astro)

小さな（モバイル）ビューポートでのサイドバーナビゲーションの切り替えを担当する、[`<PageFrame>`](#pageframe)内にレンダリングされるコンポーネント。

#### `TwoColumnContent`

**デフォルトコンポーネント:** [`TwoColumnContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TwoColumnContent.astro)

メインコンテンツのカラムと右サイドバー（目次）をラップするレイアウトコンポーネント。デフォルトの実装では、1カラムの小さなビューポート向けレイアウトと、2カラムの大きなビューポート向けレイアウトの切り替えをおこないます。

---

### ヘッダー

以下のコンポーネントは、Starlightのトップナビゲーションバーをレンダリングします。

#### `Header`

**デフォルトコンポーネント:** [`Header.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Header.astro)

すべてのページの上部に表示されるヘッダーコンポーネント。デフォルトの実装では、[`<SiteTitle />`](#sitetitle)、[`<Search />`](#search)、[`<SocialIcons />`](#socialicons)、[`<ThemeSelect />`](#themeselect)、[`<LanguageSelect />`](#languageselect)を表示します。

#### `SiteTitle`

**デフォルトコンポーネント:** [`SiteTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SiteTitle.astro)

サイトタイトルをレンダリングするためにサイトヘッダーの先頭にレンダリングされるコンポーネント。デフォルトの実装では、Starlightの設定で定義されたロゴをレンダリングするためのロジックが含まれています。

#### `Search`

**デフォルトコンポーネント:** [`Search.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Search.astro)

Starlightの検索UIをレンダリングするために使用されるコンポーネント。デフォルトの実装では、ヘッダー内のボタンと、クリックされたときに検索モーダルを表示し、[PagefindのUI](https://pagefind.app/)をロードするためのコードが含まれています。

[`pagefind`](/ja/reference/configuration/#pagefind)が無効になっている場合、デフォルトの検索コンポーネントはレンダリングされません。ただし、`Search`をオーバーライドすると、`pagefind`設定オプションが`false`であっても常にカスタムコンポーネントがレンダリングされます。これにより、Pagefindを無効にしたときに、代替となる検索プロバイダーのUIを追加できます。

#### `SocialIcons`

**デフォルトコンポーネント:** [`SocialIcons.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SocialIcons.astro)

ソーシャルアイコンへのリンクを含む、サイトヘッダーにレンダリングされるコンポーネント。デフォルトの実装では、Starlightの設定の[`social`](/ja/reference/configuration/#social)オプションを使用して、アイコンとリンクをレンダリングします。

#### `ThemeSelect`

**デフォルトコンポーネント:** [`ThemeSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeSelect.astro)

ユーザーが好みのカラースキームを選択できるようにするための、サイトヘッダーにレンダリングされるコンポーネント。

#### `LanguageSelect`

**デフォルトコンポーネント:** [`LanguageSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LanguageSelect.astro)

ユーザーが別の言語に切り替えられるようにするための、サイトヘッダーにレンダリングされるコンポーネント。

---

### グローバルサイドバー

Starlightのグローバルサイドバーには、メインのサイトナビゲーションが含まれます。小さなビューポートでは、これはドロップダウンメニューの背後に隠されます。

#### `Sidebar`

**デフォルトコンポーネント:** [`Sidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Sidebar.astro)

グローバルナビゲーションを含んだ、ページコンテンツの前にレンダリングされるコンポーネント。デフォルトの実装では、十分に広いビューポートではサイドバーとして、小さな（モバイル）ビューポートではドロップダウンメニューの内側に表示されます。また、モバイルメニュー内に追加のアイテムを表示するために、[`<MobileMenuFooter />`](#mobilemenufooter)をレンダリングします。

#### `MobileMenuFooter`

**デフォルトコンポーネント:** [`MobileMenuFooter.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuFooter.astro)

モバイルドロップダウンメニューの下部にレンダリングされるコンポーネント。デフォルトの実装では、[`<ThemeSelect />`](#themeselect)と[`<LanguageSelect />`](#languageselect)をレンダリングします。

---

### ページサイドバー

Starlightのページサイドバーは、現在のページの見出しを列挙する目次の表示を担当しています。小さなビューポートでは、これは固定されたドロップダウンメニューへと折りたたまれます。

#### `PageSidebar`

**デフォルトコンポーネント:** [`PageSidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageSidebar.astro)

目次を表示するために、メインページのコンテンツの前にレンダリングされるコンポーネント。デフォルトの実装では、[`<TableOfContents />`](#tableofcontents)と[`<MobileTableOfContents />`](#mobiletableofcontents)をレンダリングします。

#### `TableOfContents`

**デフォルトコンポーネント:** [`TableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TableOfContents.astro)

現在のページの目次を、大きめのビューポートにおいてレンダリングするコンポーネント。

#### `MobileTableOfContents`

**デフォルトコンポーネント:** [`MobileTableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileTableOfContents.astro)

現在のページの目次を、小さな（モバイル）ビューポートにおいてレンダリングするコンポーネント。

---

### コンテンツ

以下のコンポーネントは、ページコンテンツのメインカラムにレンダリングされます。

#### `Banner`

**デフォルトコンポーネント:** [`Banner.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Banner.astro)

各ページの上部にレンダリングされるバナーコンポーネント。デフォルトの実装では、ページの[`banner`](/ja/reference/frontmatter/#banner)フロントマターの値を使用して、レンダリングするかどうかを決定します。

#### `ContentPanel`

**デフォルトコンポーネント:** [`ContentPanel.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ContentPanel.astro)

メインコンテンツカラムのセクションをラップするために使用されるレイアウトコンポーネント。

#### `PageTitle`

**デフォルトコンポーネント:** [`PageTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageTitle.astro)

現在のページの`<h1>`要素を含むコンポーネント。

デフォルトの実装と同様に、`<h1>`要素に`id="_top"`を設定する必要があります。

#### `FallbackContentNotice`

**デフォルトコンポーネント:** [`FallbackContentNotice.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/FallbackContentNotice.astro)

現在の言語の翻訳が利用できないページにおいて、ユーザーに表示される通知。多言語サイトでのみ使用されます。

#### `Hero`

**デフォルトコンポーネント:** [`Hero.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Hero.astro)

フロントマターで[`hero`](/ja/reference/frontmatter/#hero)が設定されている場合に、ページの上部にレンダリングされるコンポーネント。デフォルトの実装では、大きなタイトル、タグライン、コールトゥアクション（call-to-action）リンク、オプションの画像を表示します。

#### `MarkdownContent`

**デフォルトコンポーネント:** [`MarkdownContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MarkdownContent.astro)

各ページのメインコンテンツの周囲にレンダリングされるコンポーネント。デフォルトの実装では、Markdownコンテンツに適用する基本的なスタイルをセットします。

Markdownコンテンツのスタイルは`@astrojs/starlight/style/markdown.css`にも公開されており、`.sl-markdown-content`CSSクラスにスコープされています。

---

### フッター

以下のコンポーネントは、ページコンテンツのメインカラムの下部にレンダリングされます。

#### `Footer`

**デフォルトコンポーネント:** [`Footer.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Footer.astro)

各ページの下部に表示されるフッターコンポーネント。デフォルトの実装では、[`<LastUpdated />`](#lastupdated)、[`<Pagination />`](#pagination)、[`<EditLink />`](#editlink)を表示します。

#### `LastUpdated`

**デフォルトコンポーネント:** [`LastUpdated.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LastUpdated.astro)

最終更新日を表示するために、ページフッターにレンダリングされるコンポーネント。

#### `EditLink`

**デフォルトコンポーネント:** [`EditLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/EditLink.astro)

ページを編集できる場所へのリンクを表示するために、ページフッターにレンダリングされるコンポーネント。

#### `Pagination`

**デフォルトコンポーネント:** [`Pagination.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Pagination.astro)

前のページと次のページとの間にナビゲーション用矢印を表示するために、ページフッターにレンダリングされるコンポーネント。
