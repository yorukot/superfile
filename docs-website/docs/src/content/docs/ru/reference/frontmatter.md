---
title: Метаданные
description: Обзор полей метаданных, поддерживаемых Starlight по умолчанию.
---

Вы можете настроить отдельные страницы Markdown и MDX в Starlight, задав значения в их метаданных. Например, на обычной странице можно задать поля `title` и `description`:

```md {3-4}
---
# src/content/docs/example.md
title: Об этом проекте
description: Узнайте больше о проекте, над которым я работаю.
---

Добро пожаловать на страницу «О сайте»!
```

## Поля метаданных

### `title` (обязательно)

**тип:** `string`

Вы должны указать заголовок для каждой страницы. Это будет отображаться в верхней части страницы, на вкладках браузера и в метаданных страницы.

### `description`

**тип:** `string`

Описание страницы используется в качестве метаданных страницы и будет воспринято поисковыми системами и в превью социальных сетей.

### `slug`

**type**: `string`

Переопределите slug страницы. Более подробную информацию вы найдете в разделе [Определение пользовательских слагов](https://docs.astro.build/ru/guides/content-collections/#defining-custom-slugs) в документации Astro.

### `editUrl`

**тип:** `string | boolean`

Переопределяет [глобальную конфигурацию `editLink`](/ru/reference/configuration/#editlink). Установите значение `false`, чтобы отключить ссылку «Редактировать страницу» для конкретной страницы или предоставить альтернативный URL, по которому можно редактировать содержимое этой страницы.

### `head`

**тип:** [`HeadConfig[]`](/ru/reference/configuration/#headconfig)

Вы можете добавить дополнительные теги в `<head>` вашей страницы, используя поле `head` метаданных. Это означает, что вы можете добавлять пользовательские стили, метаданные и другие теги на одну страницу. Аналогично [глобальной опции `head`](/ru/reference/configuration/#head).

```md
---
# src/content/docs/example.md
title: О нас
head:
  # Используем свой тег <title>
  - tag: title
    content: Пользовательский заголовок
---
```

### `tableOfContents`

**тип:** `false | { minHeadingLevel?: number; maxHeadingLevel?: number; }`

Переопределяет [глобальную конфигурацию `tableOfContents`](/ru/reference/configuration/#tableofcontents).
Настройте уровни заголовков, которые будут включены, или установите значение `false`, чтобы скрыть оглавление на этой странице.

```md
---
# src/content/docs/example.md
title: Страница, содержащая только H2 в оглавлении
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 2
---
```

```md
---
# src/content/docs/example.md
title: Страница без оглавления
tableOfContents: false
---
```

### `template`

**тип:** `'doc' | 'splash'`  
**по умолчанию:** `'doc'`

Установите шаблон макета для этой страницы.
Страницы используют макет `'doc'` по умолчанию.
Установите значение `'splash'`, чтобы использовать более широкий макет без боковых панелей, предназначенный для целевых страниц.

### `hero`

**тип:** [`HeroConfig`](#heroconfig)

Добавьте компонент hero в верхнюю часть этой страницы. Хорошо сочетается с `template: splash`.

Например, в этом конфиге показаны некоторые общие опции, включая загрузку изображения из вашего репозитория.

```md
---
# src/content/docs/example.md
title: Моя домашняя страница
template: splash
hero:
  title: 'Мой проект: Быстрая доставка в космосе'
  tagline: Доставьте свои вещи на Луну и обратно в мгновение ока.
  image:
    alt: Сверкающий, ярко раскрашенный логотип
    file: ~/assets/logo.png
  actions:
    - text: Расскажите мне больше
      link: /getting-started/
      icon: right-arrow
      variant: primary
    - text: Просмотр на GitHub
      link: https://github.com/astronaut/my-project
      icon: external
      attrs:
        rel: me
---
```

Вы можете отображать разные версии главного изображения в светлом и тёмном режимах.

```md
---
# src/content/docs/example.md
hero:
  image:
    alt: Сверкающий, ярко раскрашенный логотип
    dark: ~/assets/logo-dark.png
    light: ~/assets/logo-light.png
---
```

#### `HeroConfig`

```ts
interface HeroConfig {
  title?: string;
  tagline?: string;
  image?:
    | {
        // Относительный путь к изображению в вашем репозитории.
        file: string;
        // Alt-текст, чтобы сделать изображение доступным для вспомогательных технологий
        alt?: string;
      }
    | {
        // Относительный путь к изображению в вашем репозитории, которое будет использоваться для тёмного режима.
        dark: string;
        // Относительный путь к изображению в вашем репозитории, которое будет использоваться для светлого режима.
        light: string;
        // Alt-текст, чтобы сделать изображение доступным для вспомогательных технологий
        alt?: string;
      }
    | {
        // Необработанный HTML для использования в слоте изображения.
        // Это может быть пользовательский тег `<img>` или встроенный `<svg>`.
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

**тип:** `{ content: string }`

Отображает баннер объявления в верхней части этой страницы.

Значение `content` может включать HTML для ссылок или другого содержимого.
Например, на этой странице отображается баннер со ссылкой на `example.com`.

```md
---
# src/content/docs/example.md
title: Страница с баннером
banner:
  content: |
    Мы только что запустили нечто крутое!
    <a href="https://example.com">Проверьте</a>
---
```

### `lastUpdated`

**тип:** `Date | boolean`

Переопределяет [глобальную опцию `lastUpdated`](/ru/reference/configuration/#lastupdated). Если указана дата, она должна быть действительной [временной меткой YAML](https://yaml.org/type/timestamp.html) и будет переопределять дату, хранящуюся в истории Git для этой страницы.

```md
---
# src/content/docs/example.md
title: Страница с пользовательской датой последнего обновления
lastUpdated: 2022-08-09
---
```

### `prev`

**тип:** `boolean | string | { link?: string; label?: string }`

Переопределяет [глобальную опцию `pagination`](/ru/reference/configuration/#pagination). Если указана строка, будет заменен сгенерированный текст ссылки, а если указан объект, будут переопределены и ссылка, и текст.

```md
---
# src/content/docs/example.md
# Скрываем ссылку на предыдущую страницу
prev: false
---
```

```md
---
# src/content/docs/example.md
# Переопределяем текст ссылки на предыдущую страницу
prev: Продолжить обучение
---
```

```md
---
# src/content/docs/example.md
# Переопределяем ссылку и текст предыдущей страницы
prev:
  link: /unrelated-page/
  label: Загляните на другую страницу
---
```

### `next`

**тип:** `boolean | string | { link?: string; label?: string }`

То же самое, что и [`prev`](#prev), но для ссылки на следующую страницу.

```md
---
# src/content/docs/example.md
# Скрываем ссылку на следующую страницу
next: false
---
```

### `pagefind`

**тип:** `boolean`  
**по умолчанию:** `true`

Установите, должна ли эта страница быть включена в поисковый индекс [Pagefind](https://pagefind.app/). Установите значение `false`, чтобы исключить страницу из результатов поиска:

```md
---
# src/content/docs/example.md
# Скрываем эту страницу из поискового индекса
pagefind: false
---
```

### `draft`

**тип:** `boolean`  
**по умолчанию:** `false`

Установите, следует ли считать эту страницу черновиком и не включать её в [производственные сборки](https://docs.astro.build/ru/reference/cli-reference/#astro-build) и [группы автогенерируемых ссылок](/ru/guides/sidebar/#автогенерируемые-группы). Установите значение `true`, чтобы пометить страницу как черновик и сделать её видимой только во время разработки.

```md
---
# src/content/docs/example.md
# Исключить эту страницу из производственных сборок
draft: true
---
```

### `sidebar`

**тип:** [`SidebarConfig`](#sidebarconfig)

Управление отображением этой страницы в [боковой панели](/ru/reference/configuration/#sidebar) при использовании автогенерируемой группы ссылок.

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

**тип:** `string`  
**по умолчанию:** [`title`](#title-обязательно) страницы

Устанавливает метку для этой страницы в боковой панели при отображении в автогенерируемой группе ссылок.

```md
---
# src/content/docs/example.md
title: Об этом проекте
sidebar:
  label: О сайте
---
```

#### `order`

**тип:** `number`

Управляйте порядком этой страницы при сортировке автоматически созданной группы ссылок.
Страницы с меньшим значением параметра `order` отображаются выше в группе ссылок.

```md
---
# src/content/docs/example.md
title: Страница, которая будет отображаться первой
sidebar:
  order: 1
---
```

#### `hidden`

**тип:** `boolean`  
**по умолчанию:** `false`

Запрещает включать эту страницу в автоматически создаваемую группу боковой панели.

```md
---
# src/content/docs/example.md
title: Страница, которую нужно скрыть из автоматически созданной боковой панели
sidebar:
  hidden: true
---
```

#### `badge`

**тип:** <code>string | <a href="/ru/reference/configuration/#badgeconfig">BadgeConfig</a></code>

Добавьте значок на страницу в боковой панели, если она отображается в автогенерируемой группе ссылок.
При использовании строки значок будет отображаться с акцентным цветом по умолчанию.
В качестве опции передайте объект [`BadgeConfig`](/ru/reference/configuration/#badgeconfig) с полями `text` и `variant` для настройки значка.

```md
---
# src/content/docs/example.md
title: Страница со значком
sidebar:
  # Используется вариант по умолчанию, соответствующий акцентному цвету вашего сайта
  badge: Новое
---
```

```md
---
# src/content/docs/example.md
title: Страница со значком
sidebar:
  badge:
    text: Экспериментально
    variant: caution
---
```

#### `attrs`

**тип:** `Record<string, string | number | boolean | undefined>`

Атрибуты HTML для добавления к ссылке на страницу в боковой панели при отображении в автогенерируемой группе ссылок.

```md
---
# src/content/docs/example.md
title: Открытие страницы в новой вкладке
sidebar:
  # Открывает страницу в новой вкладке
  attrs:
    target: _blank
---
```

## Настройка схемы метаданных

Схема метаданных для коллекции контента Starlight `docs` настраивается в файле `src/content/config.ts` с помощью помощника `docsSchema()`:

```ts {3,6}
// src/content/config.ts
import { defineCollection } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({ schema: docsSchema() }),
};
```

Подробнее о схемах коллекций содержимого читайте в разделе [Определение схемы коллекции](https://docs.astro.build/ru/guides/content-collections/#defining-a-collection-schema) в документации Astro.

`docsSchema()` принимает следующие параметры:

### `extend`

**тип:** Схема Zod или функция, возвращающая схему Zod  
**по умолчанию:** `z.object({})`

Расширьте схему Starlight дополнительными полями, задав `extend` в опциях `docsSchema()`.
Значение должно быть [схемой Zod](https://docs.astro.build/ru/guides/content-collections/#defining-datatypes-with-zod).

В следующем примере мы задаем более строгий тип для `description`, чтобы сделать его обязательным, и добавляем новое необязательное поле `category`:

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: z.object({
        // Делаем встроенное поле обязательным
        description: z.string(),
        // Добавляем новое поле в схему
        category: z.enum(['tutorial', 'guide', 'reference']).optional(),
      }),
    }),
  }),
};
```

Чтобы воспользоваться преимуществами [хелпера `image()`](https://docs.astro.build/ru/guides/images/#images-in-content-collections), используйте функцию, которая возвращает расширение вашей схемы:

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: ({ image }) => {
        return z.object({
          // Добавляем поле, которое должно разрешаться в локальное изображение
          cover: image(),
        });
      },
    }),
  }),
};
```
