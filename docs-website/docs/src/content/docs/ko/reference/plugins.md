---
title: 플러그인 참조
description: Starlight 플러그인 API의 개요입니다.
tableOfContents:
  maxHeadingLevel: 4
---

Starlight 플러그인은 Starlight 구성, UI, 동작을 사용자 정의할 수 있을 뿐만 아니라 공유 및 재사용하기도 쉽습니다.
이 참조 페이지에는 플러그인이 액세스할 수 있는 API가 문서화되어 있습니다.

[구성 참조](/ko/reference/configuration/#plugins)에서 Starlight 플러그인 사용에 대해 자세히 알아보거나 [플러그인 쇼케이스](/ko/resources/plugins/)를 방문하여 사용 가능한 플러그인 목록을 확인하세요.

## 빠른 API 참조

Starlight 플러그인은 다음과 같은 형태를 가지고 있습니다.
다양한 속성과 후크 매개변수에 대한 자세한 내용은 아래를 참조하세요.

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

**타입:** `string`

플러그인은 자신을 설명하는 고유한 이름을 제공해야 합니다. 이름은 이 플러그인과 관련된 [메시지를 로깅](#logger)할 때 사용되며 다른 플러그인에서 이 플러그인의 존재를 감지하는 데 사용될 수도 있습니다.

## `hooks`

후크는 Starlight가 특정 시간에 플러그인 코드를 실행하기 위해 호출하는 함수입니다. 현재 Starlight는 단일 `setup` 후크를 지원합니다.

### `hooks.setup`

Starlight가 초기화될 때 호출되는 플러그인 설정 함수입니다 ([`astro:config:setup`](https://docs.astro.build/ko/reference/integrations-reference/#astroconfigsetup) 통합 후크 실행 중 호출).
`setup` 후크를 사용하여 Starlight 구성을 업데이트하거나 Astro 통합을 추가할 수 있습니다.

이 후크는 다음 옵션들과 함께 호출됩니다.

#### `config`

**타입:** `StarlightUserConfig`

사용자 제공 [Starlight 구성](/ko/reference/configuration/)의 읽기 전용 복사본입니다.
이 구성은 현재 플러그인 이전에 구성된 다른 플러그인에 의해 업데이트되었을 수 있습니다.

#### `updateConfig`

**타입:** `(newConfig: StarlightUserConfig) => void`

사용자가 제공한 [Starlight 구성](/ko/reference/configuration/)을 업데이트하는 콜백 함수입니다.
재정의하려는 루트 수준 구성 키를 제공합니다.
중첩된 구성 값을 업데이트하려면 전체 중첩 객체를 제공해야 합니다.

기존 구성 옵션을 재정의하지 않고 확장하려면 기존 값을 새 값에 전개하여 확장하세요.
다음 예시에서는 `config.social`에 전개 연산자를 사용하여 새로운 `social` 객체를 확장합니다. 이를 통해, 새로운 [`social`](/ko/reference/configuration/#social) 미디어 계정을 기존 구성에 추가합니다.

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

**타입:** `(integration: AstroIntegration) => void`

플러그인에 필요한 [Astro 통합](https://docs.astro.build/ko/reference/integrations-reference/)을 추가하기 위한 콜백 함수입니다.

다음 예시에서 플러그인은 먼저 [Astro의 React 통합](https://docs.astro.build/ko/guides/integrations-guide/react/)이 구성되어 있는지 확인하고, 구성되어 있지 않으면 `addIntegration()`을 사용하여 이를 추가합니다.

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

      // 아직 불러오지 않은 경우에만 React 통합을 추가합니다.
      if (!isReactLoaded) {
        addIntegration(react());
      }
    },
  },
};
```

#### `astroConfig`

**타입:** `AstroConfig`

사용자 제공 [Astro 구성](https://docs.astro.build/ko/reference/configuration-reference/)의 읽기 전용 복사본입니다.

#### `command`

**타입:** `'dev' | 'build' | 'preview'`

Starlight를 실행하는 데 사용되는 명령:

- `dev` - 프로젝트는 `astro dev`로 실행
- `build` - 프로젝트는 `astro build`로 실행
- `preview` - 프로젝트는 `astro preview`로 실행

#### `isRestart`

**타입:** `boolean`

개발 서버가 시작되면 `false`, 서버가 다시 시작되면 `true`입니다.
개발 서버 재시작의 일반적인 이유에는 개발 서버가 실행되는 동안 사용자가 `astro.config.mjs`를 편집하는 경우가 포함됩니다.

#### `logger`

**타입:** `AstroIntegrationLogger`

로그를 작성하는 데 사용할 수 있는 [Astro 통합 로거](https://docs.astro.build/ko/reference/integrations-reference/#astrointegrationlogger)의 인스턴스입니다.
기록된 모든 메시지에는 플러그인 이름이 앞에 추가됩니다.

```ts {6}
// plugin.ts
export default {
  name: 'long-process-plugin',
  hooks: {
    setup({ logger }) {
      logger.info('시간이 오래 걸리는 작업 진행 중…');
      // 오래 걸리는 작업…
    },
  },
};
```

위 예시에서는 제공된 정보를 포함하는 메시지를 기록합니다.

```shell
[long-process-plugin] 시간이 오래 걸리는 작업 진행 중…
```
