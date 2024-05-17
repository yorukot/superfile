---
title: プラグイン
description: StarlightのプラグインAPIの概要。
tableOfContents:
  maxHeadingLevel: 4
---

Starlightのプラグインにより、Starlightの設定、UI、および動作をカスタマイズできます。このリファレンスページでは、プラグインがアクセス可能なAPIについて説明します。

Starlightのプラグインを使用する方法について、詳しくは[設定方法のリファレンス](/ja/reference/configuration/#plugins)を参照してください。また、利用可能なプラグインの一覧については、[プラグインショーケース](/ja/resources/plugins/)を確認してください。

## API早見表

Starlightのプラグインは次の構造をもちます。各プロパティとフックパラメータの詳細については、以下を参照してください。

```ts
interface StarlightPlugin {
  name: string;
  hooks: {
    setup: (options: {
      config: StarlightUserConfig;
      updateConfig: (newConfig: StarlightUserConfig) => void;
      addIntegration: (integration: AstroIntegration) => void;
      astroConfig: AstroConfig;
      command: 'dev' | 'build' | 'preview';
      isRestart: boolean;
      logger: AstroIntegrationLogger;
    }) => void | Promise<void>;
  };
}
```

## `name`

**type:** `string`

プラグインには、自身を説明する一意の名前を指定する必要があります。名前は、このプラグインに関連する[ログメッセージ](#logger)を出力するときに使用されます。また、あるプラグインの存在を検出するために使用される場合もあります。

## `hooks`

フックは、Starlightが特定のタイミングでプラグインコードを実行するために呼び出す関数です。現在、Starlightは`setup`フックのみサポートしています。

### `hooks.setup`

プラグインのセットアップ関数は、Starlightが（[`astro:config:setup`](https://docs.astro.build/ja/reference/integrations-reference/#astroconfigsetup)インテグレーションフックにおいて）初期化される際に呼び出されます。`setup`フックは、Starlightの設定を更新したり、Astroのインテグレーションを追加したりするために使用できます。

このフックは、以下のオプションとともに呼び出されます。

#### `config`

**type:** `StarlightUserConfig`

ユーザーが提供した[Starlightの設定](/ja/reference/configuration/)の、読み取り専用の複製です。この設定は、現在のプラグインより前に置かれた他のプラグインによって更新されている可能性があります。

#### `updateConfig`

**type:** `(newConfig: StarlightUserConfig) => void`

ユーザーが提供した[Starlightの設定](/ja/reference/configuration/)を更新するためのコールバック関数です。上書きしたいルートレベルの設定キーを指定します。ネストされた設定値を更新するには、ネストされたオブジェクトの全体を指定する必要があります。

既存の設定オプションをオーバーライドせず拡張するには、既存の値を新しい値へと展開します。以下の例では、`config.social`を新しい`social`オブジェクトに展開し、既存の設定に新しい[`social`](/ja/reference/configuration/#social)メディアアカウントを追加しています。

```ts {6-11}
// plugin.ts
export default {
  name: 'add-twitter-plugin',
  hooks: {
    setup({ config, updateConfig }) {
      updateConfig({
        social: {
          ...config.social,
          twitter: 'https://twitter.com/astrodotbuild',
        },
      });
    },
  },
};
```

#### `addIntegration`

**type:** `(integration: AstroIntegration) => void`

プラグインが必要とする[Astroのインテグレーション](https://docs.astro.build/ja/reference/integrations-reference/)を追加するためのコールバック関数です。

以下の例では、プラグインはまず[AstroのReactインテグレーション](https://docs.astro.build/ja/guides/integrations-guide/react/)が設定されているかどうかを確認し、設定されていない場合は`addIntegration()`を使用して追加します。

```ts {14} "addIntegration,"
// plugin.ts
import react from '@astrojs/react';

export default {
  name: 'plugin-using-react',
  hooks: {
    setup({ addIntegration, astroConfig }) {
      const isReactLoaded = astroConfig.integrations.find(
        ({ name }) => name === '@astrojs/react'
      );

      // まだロードされていない場合のみ、Reactインテグレーションを追加します。
      if (!isReactLoaded) {
        addIntegration(react());
      }
    },
  },
};
```

#### `astroConfig`

**type:** `AstroConfig`

ユーザーが提供した[Astroの設定](https://docs.astro.build/ja/reference/configuration-reference/)の、読み取り専用の複製です。

#### `command`

**type:** `'dev' | 'build' | 'preview'`

Starlightを実行するために使用されたコマンドです。

- `dev` - プロジェクトは`astro dev`により実行されています
- `build` - プロジェクトは`astro build`により実行されています
- `preview` - プロジェクトは`astro preview`により実行されています

#### `isRestart`

**type:** `boolean`

開発サーバーが起動したときは`false`、リロードがトリガーされたときは`true`となります。再起動が発生するよくある理由としては、開発サーバーが実行されている間にユーザーが`astro.config.mjs`を編集した場合などがあります。

#### `logger`

**type:** `AstroIntegrationLogger`

ログを書き込むために使用する[Astroインテグレーションロガー](https://docs.astro.build/ja/reference/integrations-reference/#astrointegrationlogger)のインスタンスです。すべてのログメッセージは、プラグイン名が接頭辞として付加されます。

```ts {6}
// plugin.ts
export default {
  name: 'long-process-plugin',
  hooks: {
    setup({ logger }) {
      logger.info('時間が掛かる処理を開始します…');
      // 何らかの時間が掛かる処理…
    },
  },
};
```

上記の例では、指定されたinfoメッセージを含むメッセージがログに出力されます。

```shell
[long-process-plugin] 時間が掛かる処理を開始します…
```
