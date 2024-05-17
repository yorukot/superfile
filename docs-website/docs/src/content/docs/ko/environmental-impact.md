---
title: 친환경 문서
description: Starlight가 어떻게 친환경 문서 사이트를 구축하고 탄소 배출량을 줄이는 데 도움이 되는지 알아보세요.
---

웹 산업이 기후에 미치는 영향에 대한 추정치는 전 세계 탄소 배출량의 [2%][sf] ~ [4%][bbc] 범위이며, 이는 항공 산업의 배출량과 비슷합니다.
웹 사이트의 생태학적 영향을 계산하기 위한 여러 복잡한 요소가 존재하지만, 이 가이드는 문서 사이트의 환경적 영향을 줄이기 위한 몇 가지 팁을 포함합니다.

좋은 소식은 Starlight를 선택하는 것이 좋은 시작이라는 것입니다.
Website Carbon Calculator에 따르면 이 사이트는 [테스트된 웹페이지의 99%보다 깨끗하며][sl-carbon] 페이지 방문당 0.01g의 CO₂를 배출합니다.

## 페이지 크기

웹페이지가 전송하는 데이터가 많을수록 더 많은 에너지 자원이 필요합니다.
2023년 4월, [HTTP Archive의 데이터][http]에 따르면 웹페이지 중앙값은 2,000KB 이상을 다운로드해야 했습니다.

Starlight는 최대한 가벼운 페이지를 구축합니다.
예를 들어, 사용자는 첫 방문 시 50KB 미만의 압축 데이터를 다운로드하게 됩니다. 이는 HTTP Archive 중앙값의 2.5%에 불과합니다.
좋은 캐싱 전략을 사용하면 후속 탐색에서 10KB 정도만 다운로드할 수 있습니다.

### 이미지

Starlight는 좋은 기준을 제공하지만 문서 페이지에 이미지를 추가하면 페이지의 크기가 빠르게 증가할 수 있습니다.
Starlight는 Astro의 [최적화된 자산 지원][assets]을 사용하여 Markdown 및 MDX 파일에서 추가한 로컬 이미지를 최적화합니다.

### UI 컴포넌트

React 또는 Vue와 같은 UI 프레임워크로 구축된 컴포넌트는 페이지에 많은 양의 JavaScript를 추가합니다.
Starlight는 Astro 기반으로 구축되었고, Astro 자체 컴포넌트는 [Astro 아일랜드][islands] 덕분에 **기본적으로 클라이언트 측 JavaScript가 전혀 로드되지 않습니다**.

### 캐싱

캐싱은 브라우저가 이미 다운로드한 데이터를 저장하고 재사용하는 기간을 제어하는 ​​데 사용됩니다.
좋은 캐싱 전략은 콘텐츠가 변경될 때 새 콘텐츠를 최대한 빨리 얻을 수 있도록 하며, 변경되지 않은 동일한 콘텐츠를 무의미하게 반복해서 다운로드하는 것을 방지합니다.

캐싱을 구성하는 가장 일반적인 방법은 [`Cache-Control` HTTP 헤더][cache]를 사용하는 것입니다.
Starlight를 사용할 때 `/_astro/` 디렉터리의 모든 항목에 대해 긴 캐시 시간을 설정할 수 있습니다.
이 디렉터리에는 불필요한 다운로드를 줄여 영원히 안전하게 캐시할 수 있는 CSS, JavaScript 및 기타 번들 자산이 포함되어 있습니다.

```
Cache-Control: public, max-age=604800, immutable
```

캐싱을 구성하는 방법은 웹 호스트에 따라 다릅니다. 예를 들어 Vercel은 구성 없이 이 캐싱 전략을 적용하는 반면, 프로젝트에 `public/_headers` 파일을 추가하여 [Netlify용 사용자 정의 헤더][ntl-headers]를 설정할 수 있습니다.

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## 전력 소비

웹 페이지가 구축되는 방식은 사용자 장치를 실행하는 데 필요한 전력에 영향을 미칠 수 있습니다. Starlight는 최소한의 JavaScript를 사용하여 사용자의 휴대폰, 태블릿 또는 컴퓨터가 페이지를 로드하고 렌더링하는 데 필요한 전력을 감소시킵니다.

분석 추적 스크립트나 동영상 삽입과 같은 JavaScript 중심 콘텐츠를 추가할 때 페이지의 전력 사용량이 증가할 수 있으므로 주의하세요.
분석이 필요한 경우 [Cabin][cabin], [Fathom][fathom] 또는 [Plausible][plausible]과 같은 가벼운 옵션을 선택하는 것이 좋습니다.
YouTube 및 Vimeo와 같은 동영상 삽입은 [상호 작용 시 동영상 로드][lazy-video]를 통해 개선될 수 있습니다.
[`astro-embed`][embed]와 같은 패키지는 일반적인 서비스에 도움이 될 수 있습니다.

:::tip[알고 계셨나요?]

JavaScript를 분석하고 컴파일하는 것은 브라우저가 수행해야 하는 비용이 가장 많이 드는 작업 중 하나입니다.
동일한 크기의 JPEG 이미지를 렌더링하는 것을 다른 언어들과 비교해보면 [JavaScript는 처리하는 데 30배 이상 더 오래 걸릴 수 있습니다][cost-of-js].

:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## 호스팅

웹 페이지가 호스팅되는 위치는 문서 사이트가 얼마나 친환경적인지에 큰 영향을 미칠 수 있습니다.
데이터 센터와 서버 팜은 높은 전력 소비와 물의 집중적 사용 등으로 인해 생태학적으로 큰 영향을 미칠 수 있습니다.

재생 가능 에너지를 사용하는 호스트를 선택하면 사이트의 탄소 배출량이 줄어듭니다. [Green Web Directory][gwb]는 호스팅 회사를 찾는 데 도움이 되는 도구 중 하나입니다.

[gwb]: https://www.thegreenwebfoundation.org/directory/

## 비교

다른 문서 프레임워크와 어떻게 비교되는지 궁금하십니까?
[Website Carbon Calculator][wcc]를 사용한 이러한 테스트는 서로 다른 도구로 작성된 유사한 페이지를 비교합니다.

| 프레임워크                  | 페이지 방문당 CO₂ |
| --------------------------- | ----------------- |
| [Starlight][sl-carbon]      | 0.01g             |
| [VitePress][vp-carbon]      | 0.05g             |
| [Docus][dc-carbon]          | 0.05g             |
| [Sphinx][sx-carbon]         | 0.07g             |
| [MkDocs][mk-carbon]         | 0.10g             |
| [Nextra][nx-carbon]         | 0.11g             |
| [docsify][dy-carbon]        | 0.11g             |
| [Docusaurus][ds-carbon]     | 0.24g             |
| [Read the Docs][rtd-carbon] | 0.24g             |
| [GitBook][gb-carbon]        | 0.71g             |

<small>2023년 5월 14일에 수집된 데이터. 최신 수치를 보려면 링크를 클릭하세요.</small>

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

## 더 많은 자료

### 도구

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Carbon Control](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### 기사 및 강연

- [“더 친환경적인 웹 구축”](https://youtu.be/EfPoOt7T5lg), Michelle Barker의 강연
- [“조직 내의 지속 가능한 웹 개발 전략”](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/), Michelle Barker의 기사
- [“모두를 위한 지속 가능한 웹”](https://2021.stateofthebrowser.com/speakers/tom-greenwood/), Tom Greenwood의 강연
- [“웹 콘텐츠가 전력 사용량에 미치는 영향”](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/), Benjamin Poulain 및 Simon Fraser의 기사

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/ko/guides/assets/
[islands]: https://docs.astro.build/ko/concepts/islands/
[wcc]: https://www.websitecarbon.com/
