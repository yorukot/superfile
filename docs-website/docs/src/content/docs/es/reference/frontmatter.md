---
title: Referencia de Frontmatter
description: Una visión general de los campos de frontmatter predeterminados que admite Starlight.
---

Puedes personalizar individualmente las páginas Markdown y MDX en Starlight estableciendo valores en su frontmatter. Por ejemplo, una página regular podría establecer los campos `title` y `description`:

```md {3-4}
---
# src/content/docs/example.md
title: Acerca de este proyecto
description: Aprende más sobre el proyecto en el que estoy trabajando.
---

¡Bienvenido a la página Acerca de!
```

## Campos de frontmatter

### `title` (requerido)

**tipo:** `string`

Debes proporcionar un título para cada página. Este se mostrará en la parte superior de la página, en las pestañas del navegador y en los metadatos de la página.

### `description`

**tipo:** `string`

La descripción de la página es usada para los metadatos de la página y será recogida por los motores de búsqueda y en las vistas previas de las redes sociales.

### `slug`

**tipo**: `string`

Sobreescribe el slug de la página. Consulta [“Definiendo slugs personalizados”](https://docs.astro.build/es/guides/content-collections/#definiendo-slugs-personalizados) en la documentación de Astro para más detalles.

### `editUrl`

**tipo:** `string | boolean`

Reemplaza la [configuración global `editLink`](/es/reference/configuration/#editlink). Establece a `false` para deshabilitar el enlace "Editar página" para una página específica o proporciona una URL alternativa donde el contenido de esta página es editable.

### `head`

**tipo:** [`HeadConfig[]`](/es/reference/configuration/#headconfig)

Puedes agregar etiquetas adicionales a la etiqueta `<head>` de tu página usando el campo `head` del frontmatter. Esto significa que puedes agregar estilos personalizados, metadatos u otras etiquetas a una sola página. Similar a la [opción global `head`](/es/reference/configuration/#head).

```md
---
# src/content/docs/example.md
title: Acerca de nosotros
head:
  # Usa una etiqueta <title> personalizada
  - tag: title
    content: Título personalizado sobre nosotros
---
```

### `tableOfContents`

**tipo:** `false | { minHeadingLevel?: number; maxHeadingLevel?: number; }`

Reemplaza la [configuración global `tableOfContents`](/es/reference/configuration/#tableofcontents).
Personaliza los niveles de encabezado que se incluirán o establece en `false` para ocultar la tabla de contenidos en esta página.

```md
---
# src/content/docs/example.md
title: Página con solo encabezados H2 en la tabla de contenidos
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 2
---
```

```md
---
# src/content/docs/example.md
title: Página sin tabla de contenidos
tableOfContents: false
---
```

### `template`

**tipo:** `'doc' | 'splash'`  
**por defecto:** `'doc'`

Establece la plantilla de diseño para esta página.
Las páginas usan el diseño `'doc'` por defecto.
Establece `'splash'` para usar un diseño más amplio sin barras laterales diseñado para las landing pages.

### `hero`

**tipo:** [`HeroConfig`](#heroconfig)

Agrega un componente hero en la parte superior de esta página. Funciona bien con `template: splash`.

Por ejemplo, esta configuración muestra algunas opciones comunes, incluyendo la carga de una imagen desde tu repositorio.

```md
---
# src/content/docs/example.md
title: Mi página de inicio
template: splash
hero:
  title: 'Mi proyecto: Cosas estelares más pronto'
  tagline: Lleva tus cosas a la luna y de vuelta en un abrir y cerrar de ojos.
  image:
    alt: Un logotipo brillante, de colores brillantes
    file: ../../assets/logo.png
  actions:
    - text: Cuéntame más
      link: /getting-started/
      icon: right-arrow
      variant: primary
    - text: View on GitHub
      link: https://github.com/astronaut/my-project
      icon: external
      attrs:
        rel: me
---
```

Puedes mostrar diferentes versiones de la imagen hero en los modos claro y oscuro.

```md
---
# src/content/docs/example.md
hero:
  image:
    alt: Un logotipo brillante, de colores brillantes
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
        // Ruta relativa a una imagen en tu repositorio.
        file: string;
        // Texto alternativo para hacer que la imagen sea accesible a la tecnología de asistencia
        alt?: string;
      }
    | {
        // Ruta relativa a una imagen en tu repositorio para usar en el modo oscuro.
        dark: string;
        // Ruta relativa a una imagen en tu repositorio para usar en el modo claro.
        light: string;
        // Texto alternativo para hacer que la imagen sea accesible a la tecnología de asistencia
        alt?: string;
      }
    | {
        // HTML crudo para usar en el espacio de la imagen.
        // Podría ser una etiqueta `<img>` personalizada o un `<svg>` en línea.
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

**tipo:** `{ content: string }`

Muestra un banner de anuncio en la parte superior de esta página.

El valor `content` puede incluir HTML para enlaces u otro contenido.
Por ejemplo, esta página muestra un banner que incluye un enlace a `example.com`.

```md
---
# src/content/docs/example.md
title: Página con un banner
banner:
  content: |
    ¡Acabamos de lanzar algo genial!
    <a href="https://example.com">Checalo</a>
---
```

### `lastUpdated`

**type:** `Date | boolean`

Sobrescribe la [opción global `lastUpdated`](/es/reference/configuration/#lastupdated). Si se especifica una fecha, debe ser una [marca de tiempo YAML](https://yaml.org/type/timestamp.html) válida y sobrescribirá la fecha almacenada en el historial de Git para esta página.

```md
---
# src/content/docs/example.md
title: Página con una fecha de última actualización personalizada
lastUpdated: 2022-08-09
---
```

### `prev`

**tipo:** `boolean | string | { link?: string; label?: string }`

Anula la [opción global de `pagination`](/es/reference/configuration/#pagination). Si se especifica un string, el texto del enlace generado se reemplazará, y si se especifica un objeto, tanto el enlace como el texto serán anulados.

```md
---
# src/content/docs/example.md
# Ocultar el enlace de la página anterior
prev: false
---
```

```md
---
# src/content/docs/example.md
# Sobrescribir el texto del enlace de la página anterior
prev: Continuar con el tutorial
---
```

```md
---
# src/content/docs/example.md
# Sobrescribir tanto el enlace de la página anterior como el texto
prev:
  link: /página-no-relacionada/
  label: Echa un vistazo a esta otra página
---
```

### `next`

**tipo:** `boolean | string | { link?: string; label?: string }`

Lo mismo que [`prev`](#prev), pero para el enlace de la página siguiente.

```md
---

# src/content/docs/example.md

# Ocultar el enlace de la página siguiente

next: false
```

### `pagefind`

**tipo:** `boolean`  
**por defecto:** `true`

Establece si esta página debe incluirse en el índice de búsqueda de [Pagefind](https://pagefind.app/). Establece en `false` para excluir una página de los resultados de búsqueda:

```md
---
# src/content/docs/example.md
# Ocultar esta página del índice de búsqueda
pagefind: false
---
```

### `draft`

**tipo:** `boolean`  
**por defecto:** `false`

Establece si esta página debe considerarse como un borrador y no incluirse en las [compilaciones de producción](https://docs.astro.build/es/reference/cli-reference/#astro-build) y [grupos de enlaces autogenerados](/es/guides/sidebar/#grupos-autogenerados). Establece en `true` para marcar una página como borrador y hacerla visible solo durante el desarrollo.

```md
---
# src/content/docs/example.md
# Excluye esta página de las compilaciones de producción
draft: true
---
```

### `sidebar`

**tipo:** [`SidebarConfig`](#sidebarconfig)

Controla cómo se muestra esta página en el [sidebar](/es/reference/configuration/#sidebar) al utilizar un grupo de enlaces generado automáticamente.

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

**tipo:** `string`  
**por defecto:** El [`title`](#title-requerido) de la página

Establece la etiqueta para esta página en la barra lateral cuando se muestra en un grupo de enlaces generado automáticamente.

```md
---
# src/content/docs/example.md
title: Acerca de este proyecto
sidebar:
  label: Acerca de
---
```

#### `order`

**tipo:** `number`

Controla el orden de esta página al ordenar un grupo de enlaces generado automáticamente.
Los números más bajos se muestran más arriba en el grupo de enlaces.

```md
---
# src/content/docs/example.md
title: Página para mostrar primero
sidebar:
  order: 1
---
```

#### `hidden`

**tipo:** `boolean`
**por defecto:** `false`

Previene que esta página se incluya en un grupo de enlaces generado automáticamente en la barra lateral.

```md
---
# src/content/docs/example.md
title: Página para ocultar de la barra lateral autogenerada
sidebar:
  hidden: true
---
```

#### `badge`

**tipo:** <code>string | <a href="/es/reference/configuration/#badgeconfig">BadgeConfig</a></code>

Agrega una insignia a la página en la barra lateral cuando se muestra en un grupo de enlaces generado automáticamente.
Cuando se usa un string, la insignia se mostrará con el color de acento predeterminado.
Opcionalmente, pasa un objeto [`BadgeConfig`](/es/reference/configuration/#badgeconfig) con los campos `text` y `variant` para personalizar la insignia.

```md
---
# src/content/docs/example.md
title: Página con una insignia
sidebar:
  # Usa la variante predeterminada que coincide con el color de acento de tu sitio
  badge: Nuevo
---
```

```md
---
# src/content/docs/example.md
title: Página con una insignia
sidebar:
  badge:
    text: Experimental
    variant: caution
---
```

#### `attrs`

**type:** `Record<string, string | number | boolean | undefined>`

Atributos HTML para agregar al enlace de la página en la barra lateral cuando se muestra en un grupo de enlaces generado automáticamente.

```md
---
# src/content/docs/example.md
title: Página que se abre en una nueva pestaña
sidebar:
  # Abre la página en una nueva pestaña
  attrs:
    target: _blank
---
```

## Personaliza el esquema del frontmatter

El esquema del frontmatter para la colección de contenido `docs` de Starlight se configura en `src/content/config.ts` usando el auxiliar `docsSchema()`:

```ts {3,6}
// src/content/config.ts
import { defineCollection } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({ schema: docsSchema() }),
};
```

Aprende más sobre los esquemas de colección de contenido en [“Definir un esquema de colección”](https://docs.astro.build/es/guides/content-collections/#definiendo-un-esquema-de-colección) en la documentación de Astro.

`docsSchema()` toma las siguientes opciones:

### `extend`

**tipo:** esquema Zod o función que devuelve un esquema Zod
**por defecto:** `z.object({})`

Extiende el esquema de Starlight con campos adicionales estableciendo `extend` en las opciones de `docsSchema()`.
El valor debe ser un [esquema Zod](https://docs.astro.build/es/guides/content-collections/#definiendo-tipos-de-datos-con-zod).

En el siguiente ejemplo, proporcionamos un tipo más estricto para `description` para hacerlo requerido y agregamos un nuevo campo opcional `category`:

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: z.object({
        // Hacer un campo integrado requerido en lugar de opcional.
        description: z.string(),
        // Agrega un nuevo campo al esquema.
        category: z.enum(['tutorial', 'guide', 'reference']).optional(),
      }),
    }),
  }),
};
```

Para tomar ventaja del [auxiliar `image()` de Astro](https://docs.astro.build/es/guides/images/#imágenes-en-colecciones-de-contenido), usa una función que devuelva tu extensión de esquema:

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: ({ image }) => {
        return z.object({
          // Agrega un campo que debe resolverse a una imagen local.
          cover: image(),
        });
      },
    }),
  }),
};
```
