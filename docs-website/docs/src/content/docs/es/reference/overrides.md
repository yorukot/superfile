---
title: Referencia de Personalización de Componentes
description: Una descripción general de los componentes y props de componentes compatibles con los reemplazos de Starlight.
tableOfContents:
  maxHeadingLevel: 4
---

Puedes reemplazar los componentes integrados de Starlight proporcionando rutas a los componentes en la opción de configuración [`components`](/es/reference/configuration/#components) de Starlight.

Esta página enumera todos los componentes disponibles para reemplazar y enlaces a sus implementaciones predeterminadas en GitHub.

Aprende más en la [Guía para Personalizar Componentes](/es/guides/overriding-components/).

## Props de Componentes

Todos los componentes pueden acceder a un objeto estándar `Astro.props` que contiene información sobre la página actual.

Para escribir los tipos de tus componentes personalizados, importa el tipo `Props` de Starlight:

```astro
---
// src/components/Custom.astro
import type { Props } from '@astrojs/starlight/props';

const { hasSidebar } = Astro.props;
//      ^ tipo: boolean
---
```

Esto te dará autocompletado y tipos al acceder a `Astro.props`.

### Props

Starlight pasará las siguientes props a tus componentes personalizados.

#### `dir`

**Tipo:** `'ltr' | 'rtl'`

La dirección de escritura de la página.

#### `lang`

**Tipo:** `string`

Etiqueta de idioma BCP-47 para la configuración regional de esta página, por ejemplo, `en`, `zh-CN` o `pt-BR`.

#### `locale`

**Tipo:** `string | undefined`

La ruta base en la que se sirve un idioma. `undefined` para los slugs de idioma raíz.

#### `siteTitle`

**Tipo:** `string`

El título del sitio para el idioma de esta página.

#### `slug`

**Tipo:** `string`

El slug se genera a partir del nombre de archivo de contenido.

#### `id`

**Tipo:** `string`

El ID único para esta página basado en el nombre de archivo de contenido.

#### `isFallback`

**Tipo:** `true | undefined`

`true` si esta página no está traducida en el idioma actual y está utilizando contenido de respaldo del idioma predeterminado.
Solo se usa en sitios multilingües.

#### `entryMeta`

**Tipo:** `{ dir: 'ltr' | 'rtl'; lang: string }`

Metadatos de configuración regional para el contenido de la página. Puede ser diferente de los valores de configuración regional de nivel superior cuando una página está utilizando contenido de respaldo.

#### `entry`

La entrada de la colección de contenido Astro para la página actual.
Incluye los valores de frontmatter para la página actual en `entry.data`.

```ts
entry: {
  data: {
    title: string;
    description: string | undefined;
    // etc.
  }
}
```

Aprende más sobre la forma de este objeto en la referencia de [Tipo de Entrada de la Colección de Astro](https://docs.astro.build/es/reference/api-reference/#tipo-de-entrada-de-la-colección).

#### `sidebar`

**Tipo:** `SidebarEntry[]`

Navegación de sitio entradas de barra lateral para esta página.

#### `hasSidebar`

**Tipo:** `boolean`

Si la barra lateral debe mostrarse o no en esta página.

#### `pagination`

**Tipo:** `{ prev?: Link; next?: Link }`

Los enlaces a la página anterior y siguiente en la barra lateral si están habilitados.

#### `toc`

**Tipo:** `{ minHeadingLevel: number; maxHeadingLevel: number; items: TocItem[] } | undefined`

Tabla de contenidos para esta página si está habilitada.

#### `headings`

**Tipo:** `{ depth: number; slug: string; text: string }[]`

Array de todos los encabezados Markdown extraídos de la página actual.
Usa [`toc`](#toc) en su lugar si deseas construir un componente de tabla de contenidos que respete las opciones de configuración de Starlight.

#### `lastUpdated`

**Tipo:** `Date | undefined`

Objeto `Date` de JavaScript que representa cuándo se actualizó por última vez esta página si está habilitado.

#### `editUrl`

**Tipo:** `URL | undefined`

Objeto `URL` para la dirección donde se puede editar esta página si está habilitado.

#### `labels`

**Tipo:** `Record<string, string>`

Un objecto que contiene cadenas de UI localizadas para la página actual. Consulta la guía ["Traducir la UI de Starlight"](/es/guides/i18n/#traduce-la-ui-de-starlight) para ver una lista de todas las claves disponibles.

---

## Componentes

### Head

Estos componentes son renderizados dentro del elemento `<head>` de cada página.
Solo deben incluir [elementos permitidos dentro de `<head>`](https://developer.mozilla.org/es/docs/Web/HTML/Element/head#see_also).

#### `Head`

**Componente por defecto:** [`Head.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Head.astro)

Componente renderizado dentro del elemento `<head>` de cada página.
Incluye etiquetas importantes como `<title>` y `<meta charset="utf-8">`.

Reemplaza este componente como último recurso.
Si es posible, prefiere la opción [`head`](/es/reference/configuration/#head) de la configuración de Starlight si es posible.

#### `ThemeProvider`

**Componente por defecto:** [`ThemeProvider.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeProvider.astro)

Componente renderizado dentro de `<head>` que configura el soporte de tema claro/oscuro.
La implementación predeterminada incluye un script en línea y una `<template>` utilizada por el script en [`<ThemeSelect />`](#themeselect).

---

### Accesibilidad

#### `SkipLink`

**Componente por defecto:** [`SkipLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SkipLink.astro)

El componente renderizado como el primer elemento dentro de `<body>` que enlaza al contenido principal de la página para la accesibilidad.
La implementación predeterminada está oculta hasta que un usuario la enfoca al presionar la tecla de tabulación con su teclado.

---

### Plantilla

Estos componentes son responsables de la ubicación de los componentes de Starlight y de la gestión de las vistas en diferentes tamaños de pantalla.
Reemplazar estos componentes viene con una complejidad significativa.
Cuando sea posible, es preferible reemplazar un componente de nivel inferior.

#### `PageFrame`

**Componente por defecto:** [`PageFrame.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageFrame.astro)

El componente plantilla que envuelve la mayor parte del contenido de la página.
La implementación predeterminada configura la plantilla header-sidebar-main e incluye slots nombrados `header` y `sidebar` con un slot predeterminado para el contenido principal.
También renderiza [`<MobileMenuToggle />`](#mobilemenutoggle) para alternar el renderizado de la navegación de la barra lateral en pantallas estrechas (móviles).

#### `MobileMenuToggle`

**Componente por defecto:** [`MobileMenuToggle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuToggle.astro)

Componente renderizado dentro de [`<PageFrame>`](#pageframe) que es responsable de alternar la navegación de la barra lateral en pantallas estrechas (móviles).

#### `TwoColumnContent`

**Componente por defecto:** [`TwoColumnContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TwoColumnContent.astro)

Componente plantilla que envuelve la columna de contenido principal y la barra lateral derecha (tabla de contenidos).
La implementación predeterminada maneja el cambio entre un diseño de una sola columna, en pantallas estrechas y un diseño de dos columnas en pantallas más grande.

---

### Header

Estos componentes renderizan la barra de navegación superior de Starlight.

#### `Header`

**Componente por defecto:** [`Header.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Header.astro)

Componente de encabezado que se muestra en la parte superior de cada página.
La implementación predeterminada muestra [`<SiteTitle />`](#sitetitle-1), [`<Search />`](#search), [`<SocialIcons />`](#socialicons), [`<ThemeSelect />`](#themeselect) y [`<LanguageSelect />`](#languageselect).

#### `SiteTitle`

**Componente por defecto:** [`SiteTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SiteTitle.astro)

Componente renderizado al comienzo del encabezado para renderizar el título de la web.
La implementación predeterminada incluye lógica para renderizar logotipos definidos en la configuración de Starlight.

#### `Search`

**Componente por defecto:** [`Search.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Search.astro)

Componente utilizado para renderizar la UI de búsqueda de Starlight.
La implementación predeterminada incluye el botón en el encabezado y el código para mostrar un modal de búsqueda cuando se hace clic y cargar la UI de [Pagefind](https://pagefind.app/).

Cuando [`pagefind`](/es/reference/configuration/#pagefind) está deshabilitado, el componente de búsqueda predeterminado no se renderizará. Sin embargo, si reemplazas `Search`, tu componente personalizado siempre se renderizará incluso si la opción de configuración `pagefind` es `false`. Esto te permite agregar una interfaz de usuario para proveedores de búsqueda alternativos cuando se deshabilita Pagefind.

#### `SocialIcons`

**Componente por defecto:** [`SocialIcons.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SocialIcons.astro)

Componente renderizado en el encabezado del sitio que incluye enlaces de iconos sociales.
La implementación predeterminada utiliza la opción [`social`](/es/reference/configuration/#social) en la configuración de Starlight para renderizar iconos y enlaces.

#### `ThemeSelect`

**Componente por defecto:** [`ThemeSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeSelect.astro)

Componente renderizado en el encabezado del sitio que permite a los usuarios seleccionar su esquema de color preferido.

#### `LanguageSelect`

**Componente por defecto:** [`LanguageSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LanguageSelect.astro)

Componente renderizado en el encabezado del sitio que permite a los usuarios cambiar a un idioma diferente.

---

### Barra lateral global

La barra lateral global de Starlight incluye la navegación principal del sitio.
En los pantallas estrechas esto está oculto detrás de un menú desplegable.

#### `Sidebar`

**Componente por defecto:** [`Sidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Sidebar.astro)

Componente renderizado antes del contenido de la página que contiene la navegación global.
La implementación predeterminada muestra una barra lateral lo suficientemente ancha y dentro de un menú desplegable en pantallas estrechas (móviles).
También renderiza [`<MobileMenuFooter />`](#mobilemenufooter) para mostrar elementos adicionales dentro del menú móvil.

#### `MobileMenuFooter`

**Componente por defecto:** [`MobileMenuFooter.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuFooter.astro)

Componente renderizado en la parte inferior del menú desplegable móvil.
La implementación por defecto renderiza [`<ThemeSelect />`](#themeselect) y [`<LanguageSelect />`](#languageselect).

---

### Barra lateral de la página

La barra lateral de la página de Starlight es responsable de mostrar una tabla de contenidos que describe los subtítulos de la página actual.
En pantallas estrechas esto se colapsa en un menú desplegable fijado.

#### `PageSidebar`

**Componente por defecto:** [`PageSidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageSidebar.astro)

El componente renderizado antes del contenido principal de la página para mostrar una tabla de contenidos.
La implementación renderiza [`<TableOfContents />`](#tableofcontents) y [`<MobileTableOfContents />`](#mobiletableofcontents)

#### `TableOfContents`

**Componente por defecto:** [`TableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TableOfContents.astro)

Componente que renderiza la tabla de contenidos de la página actual en pantallas más anchas.

#### `MobileTableOfContents`

**Componente por defecto:** [`MobileTableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileTableOfContents.astro)

Componente que renderiza la tabla de contenidos de la página actual en pantallas más estrechas (móviles).

---

### Contenido

Estos componentes se renderizan en la columna principal del contenido de la página.

#### `Banner`

**Componente por defecto:** [`Banner.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Banner.astro)

Componente Banner renderizado en la parte superior de cada página.
La implementación predeterminada usa el valor de frontmatter [`banner`](/es/reference/frontmatter/#banner) de la página para decidir si renderizar o no.

#### `ContentPanel`

**Componente por defecto:** [`ContentPanel.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ContentPanel.astro)

Componente plantilla utilizado para envolver secciones de la columna de contenido principal.

#### `PageTitle`

**Componente por defecto:** [`PageTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageTitle.astro)

Componente que contiene el elemento `<h1>` de la página actual.

Las implementaciones deben asegurarse de establecer `id="_top"` en el elemento `<h1>` como en la implementación predeterminada.

#### `DraftContentNotice`

**Componente por defecto:** [`DraftContentNotice.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/DraftContentNotice.astro)

Aviso mostrado a los usuarios durante el desarrollo cuando la página actual está marcada como borrador.

#### `FallbackContentNotice`

**Componente por defecto:** [`FallbackContentNotice.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/FallbackContentNotice.astro)

Aviso mostrado a los usuarios en páginas donde no está disponible una traducción para el idioma actual.

Solo se usa en sitios multilingües.

#### `Hero`

**Componente por defecto:** [`Hero.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Hero.astro)

Componente renderizado en la parte superior de la página cuando [`hero`](/es/reference/frontmatter/#hero) está establecido en frontmatter.
La implementación predeterminada muestra un título grande, un lema y enlaces de llamada a la acción junto con una imagen opcional.

#### `MarkdownContent`

**Componente por defecto:** [`MarkdownContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MarkdownContent.astro)

Componente renderizado alrededor del contenido principal de cada página.
La implementación predeterminada configura estilos básicos para aplicar al contenido de Markdown.

Los estilos de contenido Markdown también están expuestos en `@astrojs/starlight/style/markdown.css` y están limitados al ámbito de la clase CSS `.sl-markdown-content`.

---

### Pie de página

Estos componentes se renderizan en la parte inferior de la columna principal del contenido de la página.

#### `Footer`

**Componente por defecto:** [`Footer.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Footer.astro)

Componente de pie de página que se muestra en la parte inferior de cada página.
La implementación predeterminada muestra [`<LastUpdated />`](#lastupdated), [`<Pagination />`](#pagination) y [`<EditLink />`](#editlink).

#### `LastUpdated`

**Componente por defecto:** [`LastUpdated.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LastUpdated.astro)

Componente renderizado en el pie de página de la página para mostrar la fecha de la última actualización.

#### `EditLink`

**Componente por defecto:** [`EditLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/EditLink.astro)

Componente renderizado en el pie de página de la página para mostrar un enlace a donde se puede editar la página.

#### `Pagination`

**Componente por defecto:** [`Pagination.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Pagination.astro)

Componente renderizado en el pie de página de la página para mostrar flechas de navegación entre páginas anteriores/siguientes.
