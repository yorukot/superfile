---
title: Creación de contenido en Markdown
description: Una descripción general de la sintaxis Markdown que soporta Starlight.
---

Starlight admite la gama completa de la sintaxis [Markdown](https://daringfireball.net/projects/markdown/) en archivos `.md`, así como el frontmatter en [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f) para definir metadatos como el título y la descripción.

Por favor, asegúrate de consultar la documentación de [MDX](https://mdxjs.com/docs/what-is-mdx/#markdown) o [Markdoc](https://markdoc.dev/docs/syntax) si estás utilizando esos formatos de archivo, ya que el soporte y el uso de Markdown pueden variar.

## Frontmatter

Puedes personalizar individualmente las páginas en Starlight estableciendo valores en el frontmatter.
El frontmatter se establece en la parte superior de tus archivos entre separadores `---`:

```md title="src/content/docs/example.md"
---
title: Mi título de página
---

El contenido de la página sigue luego de los `---`.
```

Cada página debe incluir al menos un `title`.
Consulta la [referencia de frontmatter](/es/reference/frontmatter/) para ver todos los campos disponibles y cómo añadir campos personalizados.

## Estilos en línea

El texto puede estar **en negrita**, _en cursiva_, o ~~tachado~~.

```md
El texto puede estar **en negrita**, _en cursiva_, o ~~tachado~~.
```

Puedes [enlazar a otra página](/es/getting-started/).

```md
Puedes [enlazar a otra página](/es/getting-started/).
```

Puedes resaltar `código en línea` con comillas invertidas.

```md
Puedes resaltar `código en línea` con comillas invertidas.
```

## Imágenes

Las imágenes en Starlight utilizan el [soporte de assets optimizados incorporado en Astro](https://docs.astro.build/en/guides/assets/).

Markdown y MDX admiten la sintaxis Markdown para mostrar imágenes, que incluye alt-text para lectores de pantalla y tecnología de asistencia.

![Una ilustración de planetas y estrellas con la palabra “astro”](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
![Una ilustración de planetas y estrellas con la palabra “astro”](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

También se admiten rutas de imágenes relativas para imágenes almacenadas localmente en tu proyecto.

```md
// src/content/docs/page-1.md

![Una nave espacial en el espacio](../../assets/images/rocket.svg)
```

## Encabezados

Puedes estructurar el contenido utilizando encabezados. Los encabezados en Markdown se indican con uno o más `#` al comienzo de la línea.

### Cómo estructurar el contenido de la página en Starlight

Starlight está configurado para utilizar automáticamente el título de tu página como un encabezado de nivel superior y se incluirá un encabezado "Visión general" en la parte superior de la tabla de contenido de cada página. Recomendamos comenzar cada página con contenido de texto de párrafo regular y utilizar encabezados dentro de la página a partir de `<h2>` en adelante:

```md
---
title: Guía de Markdown
description: Cómo utilizar Markdown en Starlight
---

Esta página describe cómo utilizar Markdown en Starlight.

## Estilos en línea

## Encabezados
```

### Enlaces de anclaje automáticos para encabezados.

Al utilizar encabezados en Markdown, se generan automáticamente enlaces de anclaje para que puedas vincular directamente a ciertas secciones de tu página:

```md
---
title: Mi página de contenido
description: Cómo utilizar los enlaces de anclaje integrados de Starlight.
---

## Introducción

Puedo enlazar a [mi conclusión](#conclusión) más abajo en la misma página.

## Conclusión

`https://mi-sitio.com/page1/#introduction` navega directamente a mi Introducción.
```

Los encabezados de nivel 2 (`<h2>`) y nivel 3 (`<h3>`) aparecerán automáticamente en la tabla de contenido de la página.

Aprende más sobre cómo Astro procesa los `id`s de los encabezados en [la documentación de Astro](https://docs.astro.build/es/guides/markdown-content/#ids-de-encabezado)

## Apartados

Los apartados (también conocidos como “apartados” o ”contenido destacado”) son útiles para mostrar información secundaria junto al contenido principal de una página.

Starlight proporciona una sintaxis personalizada de Markdown para renderizar apartados. Los bloques de apartados se indican utilizando un par de triples dos puntos `:::` para envolver tu contenido, y pueden ser de tipo `note`, `tip`, `caution` o `danger`.

Puedes anidar cualquier otro tipo de contenido Markdown dentro de un apartado, pero los apartados son más adecuados para fragmentos de contenido cortos y concisos.

### Nota de apartados

:::note
Starlight es un conjunto de herramientas para crear sitios web de documentación construido con [Astro](https://astro.build/). Puedes comenzar con este comando:

```sh
npm run create astro@latest --template starlight
```

:::

````md
:::note
Starlight es un conjunto de herramientas para sitios de documentación construido con [Astro](https://astro.build/). Puedes comenzar con este comando:

```sh
npm run create astro@latest --template starlight
```

:::
````

### Títulos personalizados para los apartados

Puedes especificar un título personalizado para el apartado utilizando corchetes cuadrados después del tipo del apartado, por ejemplo, `:::tip[¿Sabías esto?]`.

:::tip[¿Sabías esto?]
Astro te ayuda a construir sitios web más rápidos con la[“Arquitectura de Islas”](https://docs.astro.build/es/concepts/islands/).
:::

```md
:::tip[¿Sabías esto?]
Astro te ayuda a construir sitios web más rápidos con la[“Arquitectura de Islas”](https://docs.astro.build/es/concepts/islands/).
:::
```

### Más tipos de apartados

Los apartados de caution y danger son útiles para llamar la atención del usuario sobre detalles que podrían generar problemas. Si te encuentras utilizando estos tipos de apartados con frecuencia, también puede ser una señal de que lo que estás documentando podría beneficiarse de una reestructuración o rediseño.

:::caution
Si no estás seguro de si deseas un sitio de documentación increíble, piénsalo dos veces antes de usar [Starlight](/es/).
:::

:::danger
Tus usuarios pueden ser más productivos y encontrar más fácil de usar tu producto gracias a las útiles características de Starlight.

- Navegación clara
- Tema de color configurable por el usuario
- [Soporte de i18n](/es/guides/i18n/)

:::

```md
:::caution
Si no estás seguro de si deseas un sitio de documentación increíble, piénsalo dos veces antes de usar [Starlight](/es/).
:::

:::danger
Tus usuarios pueden ser más productivos y encontrar más fácil de usar tu producto gracias a las útiles características de Starlight.

- Navegación clara
- Tema de color configurable por el usuario
- [Soporte de i18n](/es/guides/i18n/)

:::
```

## Citas en bloque

> Esto es una cita en bloque, que se utiliza comúnmente para citar a otra persona o documento.
>
> Las citas en bloque se indican con un `>` al inicio de cada línea.

```md
> Esto es una cita en bloque, que se utiliza comúnmente para citar a otra persona o documento.
>
> Las citas en bloque se indican con un `>` al inicio de cada línea.
```

## Bloques de código

Un bloque de código se indica con un bloque de tres comillas invertidas <code>```</code> al inicio y al final. Puedes indicar el lenguaje de programación que se está utilizando después de las comillas invertidas de apertura.

```js
// Código JavaScript con resaltado de sintaxis.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// Código JavaScript con resaltado de sintaxis.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

### Características de Expressive Code

Starlight usa [Expressive Code](https://github.com/expressive-code/expressive-code/tree/main/packages/astro-expressive-code) para ampliar las posibilidades de formato de los bloques de código.
Los marcadores de texto de Expressive Code y los plugins de marcos de ventana están habilitados de forma predeterminada.
El renderizado de los bloques de código se puede configurar utilizando la opción de configuración [`expressiveCode`](/es/reference/configuration/#expressivecode) de Starlight.

#### Marcadores de texto

Puedes resaltar líneas específicas o partes de tus bloques de código utilizando [los marcadores de texto de Expressive Code](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#usage-in-markdown--mdx-documents) en la línea de apertura de tu bloque de código.
Usa llaves (`{ }`) para resaltar líneas enteras, y comillas para resaltar cadenas de texto.

Hay tres estilos de resaltado: neutral para llamar la atención sobre el código, verde para indicar código insertado y rojo para indicar código eliminado.
Tanto el texto como las líneas enteras pueden marcarse con el marcador predeterminado, o en combinación con `ins=` y `del=` para producir el resaltado deseado.

Expressive Code proporciona varias opciones para personalizar la apariencia visual de tus ejemplos de código.
Muchas de estas opciones se pueden combinar, para obtener ejemplos de código altamente ilustrativos.
Por favor, explora la [documentación de Expressive Code](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md) para ver las extensas opciones disponibles.
Algunos de los ejemplos más comunes se muestran a continuación:

- [Marca líneas enteras y rangos de líneas usando el marcador `{ }`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-entire-lines--line-ranges):

  ```js {2-3}
  function demo() {
    // Esta línea (#2) y la siguiente están resaltadas
    retrun 'Esta es la línea #3 de este fragmento'
  }
  ```

  ````md
  ```js {2-3}
  function demo() {
    // Esta línea (#2) y la siguiente están resaltadas
    return 'Esta es la línea #3 de este fragmento';
  }
  ```
  ````

- [Marca selecciones de texto usando el marcador `" "` o expresiones regulares](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-individual-text-inside-lines):

  ```js "Términos individuales" /También.*compatibles/
  // Términos individuales también pueden ser resaltados
  function demo() {
    return 'También las expresiones regulares son compatibles';
  }
  ```

  ````md
  ```js "Términos individuales" /También.*compatibles/
  // Términos individuales también pueden ser resaltados
  function demo() {
    return 'También las expresiones regulares son compatibles';
  }
  ```
  ````

- [Marca texto o líneas como insertadas o eliminadas con `ins` o `del`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#selecting-marker-types-mark-ins-del):

  ```js "return true;" ins="insertados" del="eliminados"
  function demo() {
    console.log('Estos son tipos de marcadores insertados y eliminados');
    // La declaración de retorno utiliza el tipo de marcador predeterminado
    return true;
  }
  ```

  ````md
  ```js "return true;" ins="insertados" del="eliminados"
  function demo() {
    console.log('Estos son tipos de marcadores insertados y eliminados');
    // La declaración de retorno utiliza el tipo de marcador predeterminado
    return true;
  }
  ```
  ````

- [Combina el resaltado de sintaxis con la sintaxis similar a `diff`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#combining-syntax-highlighting-with-diff-like-syntax):

  ```diff lang="js"
    function thisIsJavaScript() {
      // ¡El bloque completo se resalta como JavaScript,
      // y aún podemos añadir marcadores de diferencias a él!
  -   console.log('Código antiguo a eliminar')
  +   console.log('¡Nuevo y brillante código!')
    }
  ```

  ````md
  ```diff lang="js"
    function thisIsJavaScript() {
      // ¡El bloque completo se resalta como JavaScript,
      // y aún podemos añadir marcadores de diferencias a él!
  -   console.log('Código antiguo a eliminar')
  +   console.log('¡Nuevo y brillante código!')
    }
  ```
  ````

#### Marcos y títulos

Los bloques de código se pueden representar dentro de un marco similar a una ventana.
Un marco que se parece a una ventana de código se utilizará para todos los demás lenguajes de programación (por ejemplo, `bash`o `sh`).
Otros lenguajes se muestran dentro de un marco de estilo de editor de código si incluyen un título.

Un título opcional del bloque de código se puede establecer con un atributo `title="..."` después de las comillas invertidas de apertura del bloque de código y el identificador del lenguaje, o con un comentario del nombre del archivo en las primeras líneas del código.

- [Añade una pestaña con el nombre del archivo con un comentario](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```js
  // mi-archivo-de-prueba.js
  console.log('¡Hola mundo!');
  ```

  ````md
  ```js
  // mi-archivo-de-prueba.js
  console.log('¡Hola mundo!');
  ```
  ````

- [Agrega un título a la ventana Terminal](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```bash title="Instalando dependencias…"
  npm install
  ```

  ````md
  ```bash title="Instalando dependencias…"
  npm install
  ```
  ````

- [Desactiva los marcos de ventana con `frame="none"`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#overriding-frame-types)

  ```bash frame="none"
  echo "Esto no se renderiza como una terminal a pesar de usar el lenguaje bash"
  ```

  ````md
  ```bash frame="none"
  echo "Esto no se renderiza como una terminal a pesar de usar el lenguaje bash"
  ```
  ````

## Otras características comunes de Markdown

Starlight admite todas las demás sintaxis de autoría de Markdown, como listas y tablas. Puedes consultar la [Guía de referencia de Markdown](https://www.markdownguide.org/cheat-sheet/) para obtener una descripción general rápida de todos los elementos de sintaxis de Markdown.

## Configuración avanzada de Markdown y MDX

Starlight utiliza el motor de renderizado de Markdown y MDX de Astro, construido sobre remark y rehype. Puedes añadir soporte para sintaxis y comportamientos personalizados añadiendo `remarkPlugins` o `rehypePlugins` en tu archivo de configuración de Astro. Consulta la sección ["Configuración de Markdown y MDX"](https://docs.astro.build/es/guides/markdown-content/#configuraci%C3%B3n-de-markdown-y-mdx) en la documentación de Astro para obtener más información.
