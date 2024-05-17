---
title: 마크다운으로 콘텐츠 작성
description: Starlight가 지원하는 Markdown 구문의 개요입니다.
---

Starlight는 `.md` 파일에서 제목 및 설명과 같은 메타데이터를 정의하기 위해 [Markdown](https://daringfireball.net/projects/markdown/)의 모든 구문과 프론트매터 [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f)을 지원합니다.

해당 파일 형식을 사용하는 경우 Markdown 지원 및 사용법이 다를 수 있으므로 [MDX 문서](https://mdxjs.com/docs/what-is-mdx/#markdown) 또는 [Markdoc 문서](https://markdoc.dev/docs/syntax)를 확인하세요.

## 프런트매터

프런트매터의 값을 설정하여 Starlight의 개별 페이지를 사용자 정의할 수 있습니다.
프런트매터는 다음과 같이 파일 상단에 `---` 구분 기호를 사용하여 설정합니다.

```md title="src/content/docs/example.md"
---
title: 페이지 제목
---

페이지의 콘텐츠는 두 번째 `---` 뒤에 옵니다.
```

모든 페이지에는 최소 하나의 `title`이 포함되어야 합니다.
사용 가능한 모든 필드를 확인하고 사용자 정의 필드를 추가하는 방법을 알아보기 위해 [프런트매터 참조](/ko/reference/frontmatter/)를 확인하세요.

## 인라인 스타일

텍스트는 **굵게**, _기울임꼴_ 또는 ~~취소선~~으로 표시할 수 있습니다.

```md
텍스트는 **굵게**, _기울임꼴_ 또는 ~~취소선~~으로 표시할 수 있습니다.
```

[다른 페이지로 링크](/ko/getting-started/)할 수 있습니다.

```md
[다른 페이지로 링크](/ko/getting-started/)할 수 있습니다.
```

백틱을 사용하여 `인라인 코드`를 강조 표시할 수 있습니다.

```md
백틱을 사용하여 `인라인 코드`를 강조 표시할 수 있습니다.
```

## 이미지

Starlight의 이미지는 [Astro에 내장된 최적화된 자산 지원](https://docs.astro.build/ko/guides/assets/)을 사용합니다.

Markdown 및 MDX는 스크린 리더 및 보조 기술에서 사용되는 대체 텍스트가 포함된 이미지를 표시하기 위한 Markdown 구문을 지원합니다.

!["astro"라는 단어가 포함된 행성과 별 그림](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
!["astro"라는 단어가 포함된 행성과 별 그림](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

프로젝트 내 로컬 이미지 파일에 대한 상대 경로도 지원합니다.

```md
// src/content/docs/page-1.md

![우주에 있는 로켓](../../assets/images/rocket.svg)
```

## 제목

제목을 사용하여 콘텐츠를 구조화할 수 있습니다. Markdown의 제목은 줄 시작 부분에 `#` 개수로 나타냅니다.

### Starlight에서 페이지 콘텐츠를 구성하는 방법

Starlight는 페이지 제목을 최상위 제목으로 사용하도록 구성되어 있으며 각 페이지 목차 상단에 "개요" 제목을 포함합니다. 각 페이지를 일반 단락 텍스트 콘텐츠로 시작하고 `<h2>`부터 아래로 페이지 제목을 사용하는 것이 좋습니다.

```md
---
title: Markdown 가이드
description: Starlight에서 Markdown을 사용하는 방법
---

이 페이지는 Starlight에서 Markdown을 사용하는 방법을 설명합니다.

## 인라인 스타일

## 제목
```

### 제목 링크

Markdown에서 제목을 사용하면 자동으로 링크가 제공되므로 페이지의 특정 섹션에 직접 연결할 수 있습니다.

```md
---
title: 내 콘텐츠 페이지
description: Starlight에 내장된 링크를 사용하는 방법
---

## 서론

[나의 결론](#결론)은 같은 페이지 하단에 링크될 수 있습니다.

## 결론

`https://my-site.com/page1/#서론` 서론으로 바로 이동합니다.
```

레벨 2 (`<h2>`) 및 레벨 3 (`<h3>`) 제목이 페이지 목차에 자동으로 나타납니다.

[Astro 공식 문서](https://docs.astro.build/ko/guides/markdown-content/#heading-ids)에서 Astro가 제목의 `id`를 처리하는 방법에 대해 자세히 알아보세요.

## 주석

주석은 "admonitions" 또는 "callouts" 라고도 하며, 페이지의 기본 콘텐츠 주변에 보조 정보를 표시하는 데 유용합니다.

Starlight는 주석 렌더링을 위한 사용자 정의 Markdown 구문을 제공합니다. 주석 블록은 내용을 감싸기 위해 세 개의 콜론 `:::`을 사용하며 `note`, `tip`, `caution` 또는 `danger` 타입일 수 있습니다.

다른 Markdown 콘텐츠를 주석 안에 중첩시킬 수도 있지만 짧고 간결한 콘텐츠 덩어리에 가장 적합합니다.

### Note 주석

:::note

Starlight는 [Astro](https://astro.build/)로 구축된 문서 웹사이트 툴킷입니다. 다음 명령으로 시작할 수 있습니다.

```sh
npm create astro@latest -- --template starlight
```

:::

````md
:::note
Starlight는 [Astro](https://astro.build/)로 구축된 문서 웹사이트 툴킷입니다. 다음 명령으로 시작할 수 있습니다.

```sh
npm create astro@latest -- --template starlight
```

:::
````

### 사용자 정의 주석 제목

주석 타입 다음에 대괄호를 사용해 주석의 제목을 지정할 수 있습니다. `:::tip[알고 계셨나요?]`

:::tip[알고 계셨나요?]

Astro는 ["Islands Architecture"](https://docs.astro.build/ko/concepts/islands/)를 사용하여 더 빠른 웹사이트를 구축할 수 있도록 도와줍니다.
:::

```md
:::tip[알고 계셨나요?]
Astro는 ["Islands Architecture"](https://docs.astro.build/ko/concepts/islands/)를 사용하여 더 빠른 웹사이트를 구축할 수 있도록 도와줍니다.
:::
```

### 더 많은 주석 타입

Caution과 Danger 주석은 실수하기 쉬운 세부 사항에 대해 사용자를 집중시키는 데 도움이 됩니다. 이러한 기능을 많이 사용하고 있다면, 문서화중인 내용을 다시 디자인하는 것이 좋습니다.

:::caution
당신이 멋진 문서 사이트를 원하지 않는다면 [Starlight](/ko/)는 필요하지 않을 수도 있습니다.
:::

:::danger
Starlight의 유용한 기능 덕분에 사용자의 생산성이 향상되고 제품을 더 쉽게 사용할 수 있습니다.

- 쉬운 탐색
- 사용자 구성 가능한 색상 테마
- [i18n 지원](/ko/guides/i18n/)

:::

```md
:::caution
당신이 멋진 문서 사이트를 원하지 않는다면 [Starlight](/ko/)는 필요하지 않을 수도 있습니다.
:::

:::danger
Starlight의 유용한 기능 덕분에 사용자의 생산성이 향상되고 제품을 더 쉽게 사용할 수 있습니다.

- 쉬운 탐색
- 사용자 구성 가능한 색상 테마
- [i18n 지원](/ko/guides/i18n/)

:::
```

## 인용

> 이것은 인용 구문입니다. 다른 사람의 말이나 문서를 인용할 때 자주 사용됩니다.
>
> 인용은 각 줄의 시작 부분에 `>`를 사용하여 나타낼 수 있습니다.

```md
> 이것은 인용 구문입니다. 다른 사람의 말이나 문서를 인용할 때 자주 사용됩니다.
>
> 인용은 각 줄의 시작 부분에 `>`를 사용하여 나타낼 수 있습니다.
```

## 코드 블록

코드 블록은 시작과 끝 부분에 세 개의 백틱 <code>```</code>이 있는 블록으로 나타냅니다. 시작하는 백틱 뒤에 프로그래밍 언어를 명시할 수 있습니다.

```js
// 구문 강조 기능이 있는 Javascript 코드입니다.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// 구문 강조 기능이 있는 Javascript 코드입니다.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

### Expressive Code 기능

Starlight는 [Expressive Code](https://github.com/expressive-code/expressive-code/tree/main/packages/astro-expressive-code)를 사용하여 코드 블록의 형식 지정 가능성을 확장합니다.
기본적으로 Expressive Code의 텍스트 마커와 창 프레임 플러그인은 활성화되어 있습니다.
코드 블록 렌더링은 Starlight의 [`expressiveCode` 구성 옵션](/ko/reference/configuration/#expressivecode)을 사용하여 구성할 수 있습니다.

#### 텍스트 마커

코드 블록의 시작 줄에 [Expressive Code 텍스트 마커](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#usage-in-markdown--mdx-documents)를 사용하여 코드 블록의 특정 줄이나 부분을 강조 표시할 수 있습니다. 전체 줄을 강조 표시하려면 중괄호(`{ }`)를 사용하고, 텍스트 문자열을 강조 표시하려면 따옴표를 사용하세요.

세 가지 강조 스타일이 있습니다. 코드에 주의를 환기시키는 중립, 삽입된 코드를 나타내는 녹색, 삭제된 코드를 나타내는 빨간색입니다.
텍스트와 전체 줄 모두 기본 마커를 사용하거나 `ins=` 및 `del=`과 함께 표시하여 원하는 강조 표시를 생성할 수 있습니다.

Expressive Code는 코드 샘플의 시각적 모습을 사용자 정의하기 위한 여러 옵션을 제공합니다. 이들 중 다수는 예시적인 코드 샘플을 위해 결합될 수 있습니다. 사용 가능한 광범위한 옵션을 보려면 [Expressive Code 문서](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md)를 살펴보세요. 가장 일반적인 예시 중 일부는 다음과 같습니다.

- [`{ }` 마커를 사용하여 전체 줄과 줄 범위를 표시합니다.](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-entire-lines--line-ranges)

  ```js {2-3}
  function demo() {
    // 이 줄(#2)과 다음 줄이 강조 표시됩니다.
    return '이 줄은 이 스니펫의 라인 #3입니다.';
  }
  ```

  ````md
  ```js {2-3}
  function demo() {
    // 이 줄(#2)과 다음 줄이 강조 표시됩니다.
    return '이 줄은 이 스니펫의 라인 #3입니다.';
  }
  ```
  ````

- [`" "` 마커 또는 정규 표현식을 사용하여 텍스트 선택 항목을 표시합니다.](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-individual-text-inside-lines)

  ```js "Individual terms" /정규.*지원됩니다./
  // 개별 용어도 강조 표시할 수 있습니다.
  function demo() {
    return '정규 표현식도 지원됩니다.';
  }
  ```

  ````md
  ```js "Individual terms" /정규.*지원됩니다./
  // 개별 용어도 강조 표시할 수 있습니다.
  function demo() {
    return '정규 표현식도 지원됩니다.';
  }
  ```
  ````

- [`ins` 또는 `del`을 사용하여 텍스트나 줄을 삽입 또는 삭제된 것으로 표시합니다.](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#selecting-marker-types-mark-ins-del)

  ```js "return true;" ins="삽입" del="삭제된"
  function demo() {
    console.log('삽입 및 삭제된 마커 타입입니다.');
    // return 문은 기본 마커 타입을 사용합니다.
    return true;
  }
  ```

  ````md
  ```js "return true;" ins="삽입" del="삭제된"
  function demo() {
    console.log('삽입 및 삭제된 마커 타입입니다.');
    // return 문은 기본 마커 타입을 사용합니다.
    return true;
  }
  ```
  ````

- [구문 강조와 `diff` 유사 구문을 결합합니다.](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#combining-syntax-highlighting-with-diff-like-syntax)

  ```diff lang="js"
    function thisIsJavaScript() {
      // 이 전체 블록은 JavaScript로 강조표시됩니다.
      // 그리고 여전히 diff 마커를 추가할 수 있습니다!
  -   console.log('제거할 이전 코드')
  +   console.log('새롭고 빛나는 코드!')
    }
  ```

  ````md
  ```diff lang="js"
    function thisIsJavaScript() {
      // 이 전체 블록은 JavaScript로 강조표시됩니다.
      // 그리고 여전히 diff 마커를 추가할 수 있습니다!
  -   console.log('제거할 이전 코드')
  +   console.log('새롭고 빛나는 코드!')
    }
  ```
  ````

#### Frames 및 titles

코드 블록은 창과 같은 프레임 내부에서 렌더링될 수 있습니다.
터미널 창처럼 보이는 프레임은 쉘 스크립팅 언어(예: `bash` 또는 `sh`)에 사용됩니다.
title이 포함된 다른 언어는 코드 편집기 스타일의 프레임에 표시됩니다.

코드 블록의 선택적 제목은 코드 블록을 여는 백틱 및 언어 식별자 뒤에 `title="..."` 속성을 추가하거나 코드 첫 번째 줄에 파일 이름 주석을 추가하여 설정할 수 있습니다.

- [설명과 함께 파일 이름 탭 추가](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```js
  // my-test-file.js
  console.log('안녕하세요!');
  ```

  ````md
  ```js
  // my-test-file.js
  console.log('안녕하세요!');
  ```
  ````

- [터미널 창에 제목 추가](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```bash title="종속성 설치 중…"
  npm install
  ```

  ````md
  ```bash title="종속성 설치 중…"
  npm install
  ```
  ````

- [`frame="none"`으로 창 프레임 비활성화](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#overriding-frame-types)

  ```bash frame="none"
  echo "bash 언어를 사용해도 터미널로 렌더링되지 않습니다."
  ```

  ````md
  ```bash frame="none"
  echo "bash 언어를 사용해도 터미널로 렌더링되지 않습니다."
  ```
  ````

## 기타 일반적인 Markdown 기능

Starlight는 목록 및 테이블과 같은 다른 모든 Markdown 작성 구문을 지원합니다. 모든 Markdown 구문 요소에 대한 간략한 개요는 [Markdown Guide의 Markdown 치트 시트](https://www.markdownguide.org/cheat-sheet/)를 참조하세요.

## 고급 Markdown 및 MDX 구성

Starlight는 remark 및 rehype를 기반으로 구축된 Astro의 Markdown 및 MDX 렌더러를 사용합니다. Astro 구성 파일에 `remarkPlugins` 또는 `rehypePlugins`를 추가하여 사용자 정의 구문 및 동작에 대한 지원을 추가할 수 있습니다. 자세한 내용은 Astro 문서의 [Markdown 및 MDX 구성](https://docs.astro.build/ko/guides/markdown-content/#configuring-markdown-and-mdx)을 참조하세요.
