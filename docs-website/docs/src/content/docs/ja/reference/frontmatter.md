---
title: フロントマター
description: Starlightがデフォルトでサポートするフロントマターのフィールドについて。
---

Starlightでは、フロントマターに値を設定することで、MarkdownとMDXのページを個別にカスタマイズできます。たとえば通常のページでは、`title`と`description`フィールドを設定します。

```md {3-4}
---
# src/content/docs/example.md
title: このプロジェクトについて
description: 私が取り組んでいるプロジェクトについてもっと知る。
---

概要ページへようこそ！
```

## フロントマターのフィールド

### `title`（必須）

**type:** `string`

すべてのページにタイトルを指定する必要があります。これは、ページの上部、ブラウザのタブ、およびページのメタデータとして表示されます。

### `description`

**type:** `string`

ページに関する説明文はページのメタデータとして使用され、また検索エンジンやソーシャルメディアのプレビューでも使用されます。

### `slug`

**type**: `string`

ページのスラグを上書きします。詳しくは、Astroドキュメントの[「カスタムスラグの定義」](https://docs.astro.build/ja/guides/content-collections/#カスタムスラグの定義)を参照してください。

### `editUrl`

**type:** `string | boolean`

[グローバルの `editLink` 設定](/ja/reference/configuration/#editlink)を上書きします。`false`を設定して特定のページの「ページを編集」リンクを無効にするか、あるいはこのページのコンテンツを編集する代替URLを指定します。

### `head`

**type:** [`HeadConfig[]`](/ja/reference/configuration/#headconfig)

フロントマターの`head`フィールドを使用して、ページの`<head>`にタグを追加できます。これにより、カスタムスタイル、メタデータ、またはその他のタグを単一のページに追加できます。[グローバルの`head`オプション](/ja/reference/configuration/#head)と同様です。

```md
---
# src/content/docs/example.md
title: 私たちについて
head:
  # カスタム<title>タグを使う
  - tag: title
    content: カスタムのタイトル
---
```

### `tableOfContents`

**type:** `false | { minHeadingLevel?: number; maxHeadingLevel?: number; }`

[グローバルの`tableOfContents`設定](/ja/reference/configuration/#tableofcontents)を上書きします。表示したい見出しのレベルをカスタマイズするか、あるいは`false`に設定して目次を非表示とします。

```md
---
# src/content/docs/example.md
title: 目次にH2のみを表示するページ
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 2
---
```

```md
---
# src/content/docs/example.md
title: 目次のないページ
tableOfContents: false
---
```

### `template`

**type:** `'doc' | 'splash'`  
**default:** `'doc'`

ページのレイアウトテンプレートを設定します。ページはデフォルトで`'doc'`レイアウトを使用します。ランディングページ向けにサイドバーのない幅広のレイアウトを使用するには、`'splash'`を設定します。

### `hero`

**type:** [`HeroConfig`](#heroconfig)

ヒーローコンポーネントをページの上部に追加します。`template: splash`との相性が良いでしょう。

リポジトリからの画像の読み込みなど、よく使われるオプションの設定例は以下となります。

```md
---
# src/content/docs/example.md
title: 私のホームページ
template: splash
hero:
  title: '私のプロジェクト: 最高品質を、最速で'
  tagline: 瞬く間に月まで一往復。
  image:
    alt: キラリと光る、鮮やかなロゴ
    file: ../../assets/logo.png
  actions:
    - text: もっと知りたい
      link: /getting-started/
      icon: right-arrow
      variant: primary
    - text: GitHubで見る
      link: https://github.com/astronaut/my-project
      icon: external
      attrs:
        rel: me
---
```

ライトモードとダークモードで、異なるバージョンのヒーロー画像を表示できます。

```md
---
# src/content/docs/example.md
hero:
  image:
    alt: キラリと光る、鮮やかなロゴ
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
        // リポジトリ内の画像への相対パス。
        file: string;
        // この画像を支援技術からアクセス可能にするための代替テキスト。
        alt?: string;
      }
    | {
        // ダークモードで使用する、リポジトリ内の画像への相対パス。
        dark: string;
        // ライトモードで使用する、リポジトリ内の画像への相対パス。
        light: string;
        // この画像を支援技術からアクセス可能にするための代替テキスト。
        alt?: string;
      }
    | {
        // 画像のスロットに使用する生のHTML。
        // カスタムの`<img>`タグやインラインの`<svg>`などが使えます。
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

ページの上部にお知らせ用のバナーを表示します。

`content`の値には、リンクやその他のコンテンツ用のHTMLを含められます。たとえば以下のページでは、`example.com`へのリンクを含むバナーを表示しています。

```md
---
# src/content/docs/example.md
title: バナーを含むページ
banner:
  content: |
    素晴らしいサイトをリリースしました！
    <a href="https://example.com">確認してみてください</a>
---
```

### `lastUpdated`

**type:** `Date | boolean`

[グローバルの`lastUpdated`オプション](/ja/reference/configuration/#lastupdated)を上書きします。日付を指定する場合は有効な[YAMLタイムスタンプ](https://yaml.org/type/timestamp.html)である必要があり、ページのGit履歴に保存されている日付を上書きします。

```md
---
# src/content/docs/example.md
title: 最終更新日をカスタマイズしたページ
lastUpdated: 2022-08-09
---
```

### `prev`

**type:** `boolean | string | { link?: string; label?: string }`

[グローバルの`pagination`オプション](/ja/reference/configuration/#pagination)を上書きします。文字列を指定すると生成されるリンクテキストが置き換えられ、オブジェクトを指定するとリンクとテキストの両方を上書きできます。

```md
---
# src/content/docs/example.md
# 前のページへのリンクを非表示にする
prev: false
---
```

```md
---
# src/content/docs/example.md
# 前のページへのリンクテキストを上書きする
prev: チュートリアルを続ける
---
```

```md
---
# src/content/docs/example.md
# 前のページへのリンクとテキストを上書きする
prev:
  link: /unrelated-page/
  label: その他のページをチェックする
---
```

### `next`

**type:** `boolean | string | { link?: string; label?: string }`

次のページへのリンクに対して、[`prev`](#prev)と同様の設定ができます。

```md
---
# src/content/docs/example.md
# 次のページへのリンクを非表示にする
next: false
---
```

### `pagefind`

**type:** `boolean`  
**default:** `true`

ページを[Pagefind](https://pagefind.app/)の検索インデックスに含めるかどうかを設定します。ページを検索結果から除外するには、`false`に設定します。

```md
---
# src/content/docs/example.md
# このページを検索インデックスから外す
pagefind: false
---
```

### `sidebar`

**type:** [`SidebarConfig`](#sidebarconfig)

自動生成されるリンクのグループを使用している際に、[サイドバー](/ja/reference/configuration/#sidebar)にページをどのように表示するかを設定します。

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
**default:** ページの[`title`](#title必須)

自動生成されるリンクのグループ内に表示される、ページのサイドバー上でのラベルを設定します。

```md
---
# src/content/docs/example.md
title: このプロジェクトについて
sidebar:
  label: 概要
---
```

#### `order`

**type:** `number`

自動生成されるリンクのグループをソートする際の、このページの順番を設定します。小さな数値ほどリンクグループの上部に表示されます。

```md
---
# src/content/docs/example.md
title: 最初に表示するページ
sidebar:
  order: 1
---
```

#### `hidden`

**type:** `boolean`
**default:** `false`

自動生成されるサイドバーのグループにこのページを含めないようにします。

```md
---
# src/content/docs/example.md
title: 自動生成されるサイドバーで非表示にするページ
sidebar:
  hidden: true
---
```

#### `badge`

**type:** <code>string | <a href="/ja/reference/configuration/#badgeconfig">BadgeConfig</a></code>

自動生成されるリンクのグループに表示されたとき、サイドバーのページにバッジを追加します。文字列を使用すると、バッジはデフォルトのアクセントカラーで表示されます。オプションで、`text`と`variant`フィールドをもつ[`BadgeConfig`オブジェクト](/ja/reference/configuration/#badgeconfig)を渡してバッジをカスタマイズできます。

```md
---
# src/content/docs/example.md
title: バッジを含むページ
sidebar:
  # サイトのアクセントカラーに合わせたデフォルトのバリアントを使用します
  badge: New
---
```

```md
---
# src/content/docs/example.md
title: バッジを含むページ
sidebar:
  badge:
    text: 実験的
    variant: caution
---
```

#### `attrs`

**type:** `Record<string, string | number | boolean | undefined>`

自動生成されるリンクのグループ内に表示されるサイドバーのページリンクに追加するHTML属性。

```md
---
# src/content/docs/example.md
title: 新しいタブで開くページ
sidebar:
  # 新しいタブでページを開きます
  attrs:
    target: _blank
---
```

## フロントマタースキーマをカスタマイズする

Starlightの`docs`コンテンツコレクションのフロントマタースキーマは、`docsSchema()`ヘルパーを使用して`src/content/config.ts`で設定されています。

```ts {3,6}
// src/content/config.ts
import { defineCollection } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({ schema: docsSchema() }),
};
```

コンテンツコレクションのスキーマについて詳しくは、Astroドキュメントの[「コレクションスキーマの定義」](https://docs.astro.build/ja/guides/content-collections/#コレクションスキーマの定義)を参照してください。

`docsSchema()`は以下のオプションを受け取ります。

### `extend`

**type:** ZodスキーマまたはZodスキーマを返す関数  
**default:** `z.object({})`

`docsSchema()`のオプションで`extend`を設定すると、Starlightのスキーマを追加のフィールドで拡張できます。値は[Zodスキーマ](https://docs.astro.build/ja/guides/content-collections/#zodによるデータ型の定義)である必要があります。

次の例では、`description`を必須にするために厳し目の型を指定し、さらにオプションの`category`フィールドを新規追加しています。

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: z.object({
        // 組み込みのフィールドをオプションから必須に変更します。
        description: z.string(),
        // 新しいフィールドをスキーマに追加します。
        category: z.enum(['tutorial', 'guide', 'reference']).optional(),
      }),
    }),
  }),
};
```

[Astroの`image()`ヘルパー](https://docs.astro.build/ja/guides/images/#コンテンツコレクションと画像)を利用するには、拡張したスキーマを返す関数を使用します。

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: ({ image }) => {
        return z.object({
          // ローカルの画像へと解決されるフィールドを追加します。
          cover: image(),
        });
      },
    }),
  }),
};
```
