---
title: 프론트매터 참조
description: Starlight가 지원하는 기본 프론트매터 필드에 대한 개요입니다.
---

프론트매터의 값을 설정하여 Starlight에서 개별 Markdown 및 MDX 페이지를 변경할 수 있습니다. 예를 들어 일반 페이지에서는 `title` 및 `description` 필드를 설정할 수 있습니다.

```md {3-4}
---
# src/content/docs/example.md
title: 이 프로젝트에 대하여
description: 내가 진행 중인 프로젝트에 대해 자세히 알아보세요.
---

나를 소개하는 페이지에 오신 것을 환영합니다!
```

## 프론트매터 필드

### `title` (필수)

**타입:** `string`

모든 페이지에 제목을 제공해야 합니다. 이는 페이지 상단, 브라우저 탭 및 페이지 메타데이터에 표시됩니다.

### `description`

**타입:** `string`

페이지 설명은 페이지 메타데이터에 사용되며 검색 엔진과 소셜 미디어 미리 보기에서 선택됩니다.

### `slug`

**타입**: `string`

페이지의 슬러그를 재정의합니다. 자세한 내용은 Astro 공식문서의 [“사용자 정의 슬러그 정의”](https://docs.astro.build/ko/guides/content-collections/#defining-custom-slugs)를 참조하세요.

### `editUrl`

**타입:** `string | boolean`

[전역 editLink 구성](/ko/reference/configuration/#editlink)을 변경합니다. 특정 페이지에 대한 "페이지 편집" 링크를 비활성화하거나 이 페이지의 콘텐츠를 편집할 수 있는 대체 URL을 제공하려면 `false`로 설정합니다.

### `head`

**타입:** [`HeadConfig[]`](/ko/reference/configuration/#headconfig)

`head` 프론트매터 필드를 사용하여 페이지의 `<head>`에 태그를 추가할 수 있습니다. 이는 사용자 정의 스타일, 메타데이터 또는 기타 태그를 단일 페이지에 추가할 수 있음을 의미합니다. [전역 `head` 옵션](/ko/reference/configuration/#head)과 유사합니다.

```md
---
# src/content/docs/example.md
title: 회사 소개
head:
  # 사용자 정의 <title> 태그 사용
  - tag: title
    content: title에 대한 사용자 정의
---
```

### `tableOfContents`

**타입:** `false | { minHeadingLevel?: number; maxHeadingLevel?: number; }`

[전역 tableOfContents 구성](/ko/reference/configuration/#tableofcontents)을 변경합니다.
포함된 제목의 레벨을 변경하거나 값을 `false`로 설정하여 페이지에서 목차를 숨길 수 있습니다.

```md
---
# src/content/docs/example.md
title: 목차에 H2만 있는 페이지
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 2
---
```

```md
---
# src/content/docs/example.md
title: 목차가 없는 페이지
tableOfContents: false
---
```

### `template`

**타입:** `'doc' | 'splash'`  
**기본값:** `'doc'`

이 페이지의 레이아웃 템플릿을 설정합니다. 페이지는 기본적으로 `'doc'` 레이아웃을 사용합니다. 랜딩 페이지용으로 설계된 사이드바 없이 더 넓은 레이아웃을 사용하려면 `'splash'`로 설정하세요.

### `hero`

**타입:** [`HeroConfig`](#heroconfig)

페이지 상단에 hero 컴포넌트를 추가합니다. `template: splash`와 잘 작동합니다.

예를 들어 이 구성은 저장소에서 이미지를 로드하는 것을 포함하여 몇 가지 일반적인 옵션을 보여줍니다.

```md
---
# src/content/docs/example.md
title: 나의 홈페이지
template: splash
hero:
  title: '내 프로젝트: Stellar Stuff Sooner'
  tagline: 물건을 달에 가져갔다가 눈 깜짝할 사이에 다시 가져올 수 있습니다.
  image:
    alt: 반짝이는 밝은 색상의 로고
    file: ../../assets/logo.png
  actions:
    - text: 더보기
      link: /getting-started/
      icon: right-arrow
      variant: primary
    - text: Github에서 보기
      link: https://github.com/astronaut/my-project
      icon: external
      attrs:
        rel: me
---
```

밝은 모드와 어두운 모드에서 다양한 버전의 hero 이미지를 표시할 수 있습니다.

```md
---
# src/content/docs/example.md
hero:
  image:
    alt: 반짝이는 밝은 색상의 로고
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
        // 저장소에 있는 이미지의 상대 경로입니다.
        file: string;
        // 보조 기술이 이미지에 접근할 수 있도록 하는 대체 텍스트입니다.
        alt?: string;
      }
    | {
        // 어두운 모드에 사용할 저장소의 이미지에 대한 상대 경로입니다.
        dark: string;
        // 밝은 모드에 사용할 저장소의 이미지에 대한 상대 경로입니다.
        light: string;
        // 보조 기술이 이미지에 접근할 수 있도록 하는 대체 텍스트입니다.
        alt?: string;
      }
    | {
        // 이미지 슬롯에 사용할 원시 HTML입니다.
        // 사용자 정의 `<img>` 태그 또는 인라인 `<svg>` 태그일 수 있습니다.
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

**타입:** `{ content: string }`

이 페이지 상단에 공지 배너를 표시합니다.

`content` 값에는 링크나 다른 콘텐츠에 대한 HTML이 포함될 수 있습니다.
예를 들어, 이 페이지에서는 `example.com`으로 이동하는 링크가 포함된 배너가 표시됩니다.

```md
---
# src/content/docs/example.md
title: 배너가 포함된 페이지
banner:
  content: |
    방금 멋진 것을 출시했습니다!
    <a href="https://example.com">확인하러 가기</a>
---
```

### `lastUpdated`

**타입:** `Date | boolean`

[전역 `lastUpdated` 옵션](/ko/reference/configuration/#lastupdated)을 변경합니다. 날짜가 지정된 경우 유효한 [YAML 타임스탬프](https://yaml.org/type/timestamp.html)여야 하며 이 페이지의 Git 기록에 저장된 날짜를 변경합니다.

```md
---
# src/content/docs/example.md
title: 수정된 최종 업데이트 날짜가 포함된 페이지
lastUpdated: 2022-08-09
---
```

### `prev`

**타입:** `boolean | string | { link?: string; label?: string }`

[전역 `pagination` 옵션](/ko/reference/configuration/#pagination)을 변경합니다. 문자열로 설정하면 생성된 링크 텍스트가 대체되고, 객체로 설정하면 링크와 텍스트가 모두 변경됩니다.

```md
---
# src/content/docs/example.md
# 이전 페이지 링크 숨기기
prev: false
---
```

```md
---
# src/content/docs/example.md
# 이전 페이지 링크의 텍스트 변경
prev: 튜토리얼 계속하기
---
```

```md
---
# src/content/docs/example.md
# 이전 페이지 링크와 텍스트 모두 변경
prev:
  link: /unrelated-page/
  label: 다른 페이지를 확인하세요.
---
```

### `next`

**타입:** `boolean | string | { link?: string; label?: string }`

[`prev`](#prev)와 동일하지만 다음 페이지 링크용입니다.

```md
---
# src/content/docs/example.md
# 다음 페이지 링크 숨기기
next: false
---
```

### `pagefind`

**타입:** `boolean`  
**기본값:** `true`

이 페이지를 [Pagefind](https://pagefind.app/) 검색 색인에 포함할지 여부를 설정합니다. 검색 결과에서 페이지를 제외하려면 값을 `false`로 설정하세요.

```md
---
# src/content/docs/example.md
# 검색 색인에서 이 페이지 숨기기
pagefind: false
---
```

### `draft`

**타입:** `boolean`  
**기본값:** `false`

이 페이지를 초안으로 간주하여 [프로덕션 빌드](https://docs.astro.build/ko/reference/cli-reference/#astro-build) 및 [자동 생성된 링크 그룹](/ko/guides/sidebar/#자동-생성-그룹)에 포함하지 않을지 여부를 설정합니다. 페이지를 초안으로 표시하고 개발 중에만 표시하려면 `true`로 설정하세요.

```md
---
# src/content/docs/example.md
# 프로덕션 빌드에서 이 페이지 제외
draft: true
---
```

### `sidebar`

**타입:** [`SidebarConfig`](#sidebarconfig)

자동 생성된 링크 그룹을 사용할 때 이 페이지가 [사이드바](/ko/reference/configuration/#sidebar)에 표시되는 방식을 제어합니다.

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

**타입:** `string`  
**기본값:** 페이지의 [`title`](#title-필수)

자동 생성된 링크 그룹에 표시될 때 사이드바에서 이 페이지에 대한 라벨을 설정합니다.

```md
---
# src/content/docs/example.md
title: 이 프로젝트에 대하여
sidebar:
  label: 소개
---
```

#### `order`

**타입:** `number`

자동 생성된 링크 그룹을 정렬할 때 이 페이지의 순서를 제어합니다. 링크 그룹에서는 낮은 숫자가 위쪽에 표시됩니다.

```md
---
# src/content/docs/example.md
title: 첫 번째로 표시될 페이지
sidebar:
  order: 1
---
```

#### `hidden`

**타입:** `boolean`
**기본값:** `false`

이 페이지가 자동 생성된 사이드바 그룹에 포함되지 않도록 합니다.

```md
---
# src/content/docs/example.md
title: 자동 생성된 사이드바에서 숨길 페이지
sidebar:
  hidden: true
---
```

#### `badge`

**타입:** <code>string | <a href="/ko/reference/configuration/#badgeconfig">BadgeConfig</a></code>

자동 생성된 링크 그룹에 표시될 때 사이드바의 페이지에 배지를 추가합니다. 문자열을 사용하면 배지가 기본 강조 색상으로 표시됩니다. 선택적으로, `text` 및 `variant`필드가 포함된 [BadgeConfig 객체](/ko/reference/configuration/#badgeconfig)를 전달하여 배지를 사용자가 원하는대로 변경할 수 있습니다.

```md
---
# src/content/docs/example.md
title: 배지를 사용하는 페이지
sidebar:
  # 사이트의 강조 색상과 일치하는 기본 변형을 사용합니다.
  badge: New
---
```

```md
---
# src/content/docs/example.md
title: 배지를 사용하는 페이지
sidebar:
  badge:
    text: 실험적 기능
    variant: caution
---
```

#### `attrs`

**타입:** `Record<string, string | number | boolean | undefined>`

사이드바에서 자동 생성된 링크 그룹을 사용할 때, 이 페이지의 링크에 추가할 HTML 속성을 설정합니다.

```md
---
# src/content/docs/example.md
title: 새 탭에서 열리는 페이지
sidebar:
  # 새 탭에서 페이지를 엽니다.
  attrs:
    target: _blank
---
```

## 프런트매터 스키마 맞춤설정

Starlight의 `docs` 콘텐츠 컬렉션에 대한 프런트매터 스키마는 `docsSchema()` 도우미를 사용하여 `src/content/config.ts`에 구성됩니다.

```ts {3,6}
// src/content/config.ts
import { defineCollection } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({ schema: docsSchema() }),
};
```

Astro 공식문서의 ["컬렉션 스키마 정의"](https://docs.astro.build/ko/guides/content-collections/#defining-a-collection-schema)에서 콘텐츠 컬렉션 스키마에 대해 자세히 알아보세요.

`docsSchema()`는 다음 옵션을 사용합니다:

### `extend`

**타입:** Zod 스키마 또는 Zod 스키마를 반환하는 함수  
**기본값:** `z.object({})`

`docsSchema()` 옵션에서 `extend`를 설정하여 추가 필드로 Starlight의 스키마를 확장하세요.
값은 [Zod 스키마](https://docs.astro.build/ko/guides/content-collections/#defining-datatypes-with-zod)여야 합니다.

다음 예시에서는 `description` 필드에 더 엄격한 타입을 제공하여 필수 항목으로 만들고, 새로운 선택적 필드인 `category`를 추가합니다.

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: z.object({
        // 기본 제공 필드를 선택 사항이 아닌 필수 항목으로 변경합니다.
        description: z.string(),
        // 스키마에 새 필드를 추가합니다.
        category: z.enum(['tutorial', 'guide', 'reference']).optional(),
      }),
    }),
  }),
};
```

[Astro `image()` 도우미](https://docs.astro.build/ko/guides/images/#images-in-content-collections)를 활용하려면 스키마 확장을 반환하는 함수를 사용하세요.

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: ({ image }) => {
        return z.object({
          // 로컬 이미지로 확인되어야 하는 필드를 추가합니다.
          cover: image(),
        });
      },
    }),
  }),
};
```
