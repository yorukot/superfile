---
title: Экологичная документация
description: Узнайте, как Starlight может помочь создавать экологически сайты с документацией и снижать ваш углеродный след.
---

Оценки влияния интернет-индустрии на климат варьируются от [2%][sf] до [4% глобальных выбросов углерода][bbc],
что близко к выбросам авиационной индустрии.
Расчет экологического воздействия веб-сайта включает в себя множество сложных факторов,
но в данном руководстве перечислены несколько советов по снижению экологического следа вашего сайта с документацией.

Хорошая новость в том, что выбор Starlight — отличное начало.
Согласно Website Carbon Calculator, этот сайт [чище, чем 99% протестированных веб-страниц][sl-carbon],
производя 0,01 г CO₂ за каждое посещение страницы.

## Вес страницы

Чем больше данных передает веб-страница, тем больше энергетических ресурсов она требует.
В апреле 2023 года, медианная веб-страница требовала от пользователя скачать более 2 000 КБ данных, согласно [данным из HTTP Archive][http].

Starlight создает страницы лёгкими, настолько, насколько это возможно
Например, при первом посещении пользователь загрузит менее 50 КБ сжатых данных, что составляет всего 2,5% от медианного значения HTTP архива.
При хорошей стратегии кэширования последующие навигации могут загружать всего 10 КБ.

### Изображения

Хоть Starlight и предлагает лёгкие страницы по умолчанию, изображения, которые вы добавляете на страницы документации, могут быстро увеличивать вес вашей страницы.
Starlight использует [оптимизировацию ресурсов][assets] Astro для оптимизации локальных изображений в ваших файлах Markdown и MDX.

### UI-компоненты

Компоненты, на UI-фреймворках, как React или Vue, могут легко добавлять большие объемы JavaScript на страницу.
Поскольку Starlight основан на Astro, эти компоненты по умолчанию **не загружают клиентский JavaScript** благодаря [Островам Astro][islands].

### Кэширование

Кэширование управляет тем, как долго браузер хранит и повторно использует данные, которые он уже загрузил.
Хорошая стратегия кэширования гарантирует, что пользователь получает новое содержание как можно быстрее,
когда оно меняется, но также избегает бесполезной повторной загрузки одного и того же содержания снова и снова, когда оно не изменилось.

Самым распространённым способом настройки кэширования является использование [HTTP-заголовка `Cache-Control`][cache].
При использовании Starlight вы можете установить длительное время кэширования для всего, что находится в каталоге `/_astro/`.
Этот каталог содержит CSS, JavaScript и другие ресурсы, которые можно безопасно кэшировать навсегда, что позволяет снизить избыточные загрузки:

```
Cache-Control: public, max-age=604800, immutable
```

Как настроить кэширование зависит от вашего веб-хоста. Например, Vercel автоматически применяет эту стратегию кэширования без необходимости настройки,
в то же время вы можете установить [заголовки для Netlify][ntl-headers], добавив файл `public/_headers` в ваш проект:

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## Потребление энергии

То, как реализована веб-страница может влиять на потребление энергии при её запуске на устройстве пользователя.
За счет минимального использования JavaScript, Starlight снижает объем вычислительных ресурсов, необходимых телефону,
планшету или компьютеру пользователя для загрузки и отображения страниц.

Будьте внимательны при добавлении функций, таких как скрипты отслеживания аналитики или контент, зависящий от JavaScript,
например, встроенные видео, так как они могут увеличить энергопотребление страницы.
Если вам необходима аналитика, рассмотрите выбор легковесного варианта,
такого как [Cabin][cabin], [Fathom][fathom] или [Plausible][plausible].
Встроенные видео, такие как YouTube и Vimeo, можно улучшить, ожидая [взаимодействие пользователя для загрузки видео][lazy-video].
Пакеты, такие как [`astro-embed`][embed], могут помочь с часто используемыми сервисами.

:::tip
Разбор и компиляция JavaScript являются одной из самых ресурсоемких задач, которые браузерам приходится выполнять.
По сравнению с отображением изображения JPEG того же размера, [обработка JavaScript может занять более чем в 30 раз больше времени][cost-of-js].
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## Хостинг

Место, где размещена веб-страница, может иметь большое влияние на то, насколько экологичен ваш сайт с документацией.
Дата-центры и серверные фермы могут оказывать значительное экологическое воздействие, включая высокий энергопотребление и интенсивное использование воды.

Выбор хостинга, использующего возобновляемую энергию, снизит выбросы углерода для вашего сайта.
[Справочник Green Web][gwb] - один из инструментов, который может помочь вам найти хостинговые компании, работающие с экологически чистой энергией.

[gwb]: https://www.thegreenwebfoundation.org/directory/

## Сравнения

Хотите сравнивить другие фреймворки для документации?
Эти тесты с использованием [Website Carbon Calculator][wcc] сравнивают аналогичные страницы, созданные с помощью разных инструментов.

| Фреймворк                   | CO₂ на каждое посещение стр. |
| --------------------------- | ---------------------------- |
| [Starlight][sl-carbon]      | 0.01g                        |
| [VitePress][vp-carbon]      | 0.05g                        |
| [Docus][dc-carbon]          | 0.05g                        |
| [Sphinx][sx-carbon]         | 0.07g                        |
| [MkDocs][mk-carbon]         | 0.10g                        |
| [Nextra][nx-carbon]         | 0.11g                        |
| [docsify][dy-carbon]        | 0.11g                        |
| [Docusaurus][ds-carbon]     | 0.24g                        |
| [Read the Docs][rtd-carbon] | 0.24g                        |
| [GitBook][gb-carbon]        | 0.71g                        |

<small>Данные собраны 14 мая 2023 года. Чтобы увидеть актуальные цифры, перейдите по ссылке.</small>

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

## Дополнительные ресурсы

### Инструменты

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Carbon Control](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### Статьи и выступления

- [Построение более экологичного веба](https://youtu.be/EfPoOt7T5lg), выступление Мишель Баркер
- [Стратегии устойчивого веб-развития в организации](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/), статья Мишель Баркер
- [Экологически устойчивый веб для каждого](https://2021.stateofthebrowser.com/speakers/tom-greenwood/), выступление Тома Гринвуда
- [Как веб-контент может влиять на энергопотребление](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/), статья Бенджамина Пулена и Саймона Фрейзера.

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/ru/guides/assets/
[islands]: https://docs.astro.build/ru/concepts/islands/
[wcc]: https://www.websitecarbon.com/
