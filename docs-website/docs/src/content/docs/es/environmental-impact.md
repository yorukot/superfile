---
title: Documentación ecológica
description: Aprende cómo Starlight puede ayudarte a construir sitios de documentación más ecológicos y reducir tu huella de carbono.
---

Las estimaciones del impacto climático de la industria web oscilan entre el [2%][sf] y el [4% de las emisiones globales de carbono][bbc], equivalente aproximadamente a las emisiones de la industria de la aviación. Hay muchos factores complejos en el cálculo del impacto ecológico de un sitio web, pero esta guía incluye algunos consejos para reducir la huella ambiental de tu sitio de documentación.

La buena noticia es que elegir Starlight es un excelente comienzo. Según el Website Carbon Calculator, este sitio es [más limpio que el 99% de las páginas web analizadas][sl-carbon], produciendo 0.01g de CO₂ por visita a la página.

## Peso de la página

Cuanto más datos transfiera una página web, más recursos energéticos requerirá. En abril de 2023, la mediana de una página web requería que un usuario descargara más de 2.000 KB, según [los datos del HTTP Archive][http].

Starlight construye páginas que son lo más livianas posible. Por ejemplo, en una primera visita, un usuario descargará menos de 50 KB de datos comprimidos, lo que representa solo el 2.5% de la mediana del archivo HTTP. Con una buena estrategia de almacenamiento en caché, las navegaciones posteriores pueden descargar tan solo 10 KB.

### Imágenes

Si bien Starlight proporciona una buena base, las imágenes que agregas a tus páginas de documentación pueden aumentar rápidamente el peso de la página. Starlight utiliza el [soporte de assets optimizados][assets] de Astro para optimizar las imágenes locales en tus archivos Markdown y MDX.

### Componentes UI

Los componentes construidos con frameworks UI como React o Vue pueden añadir fácilmente grandes cantidades de JavaScript a una página. Sin embargo, debido a que Starlight está construido sobre Astro, los componentes como estos no cargan **ningún JavaScript del lado del cliente de forma predeterminada**, gracias a las [islas de Astro][islands].

### Caché

La caché se utiliza para controlar cuánto tiempo un navegador almacena y reutiliza los datos que ha estado descargando. Una buena estrategia de caché asegura que un usuario obtenga nuevo contenido tan pronto como sea posible cuando está cambiando, pero también evita descargar innecesariamente el mismo contenido una y otra vez cuando no ha estado cambiando.

La forma más común de configurar la caché es mediante la [cabecera HTTP `Cache-Control`][cache]. Al utilizar Starlight, puedes establecer un tiempo de caché prolongado para todo lo que se encuentra en el directorio `/_astro/`. Este directorio contiene CSS, JavaScript y otros activos empaquetados que se pueden almacenar en caché de forma segura para siempre, reduciendo las descargas innecesarias.

```
Cache-Control: public, max-age=604800, immutable
```

Cómo configurar la caché depende de tu proveedor de alojamiento web. Por ejemplo, Vercel aplica automáticamente esta estrategia de caché sin necesidad de configuración, mientras que en Netlify puedes establecer [cabeceras personalizadas][ntl-headers] añadiendo un archivo `public/_headers` a tu proyecto:

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## Consumo de energía

La forma en que se construye una página web puede afectar la cantidad de energía necesaria para que funcione en el dispositivo de un usuario. Al utilizar un JavaScript mínimo, Starlight reduce la cantidad de potencia de procesamiento que necesita el teléfono, la tableta o la computadora de un usuario para cargar y renderizar las páginas.

Ten en cuenta que al agregar funciones como scripts de seguimiento de análisis o contenido pesado en JavaScript como incrustaciones de video, esto puede aumentar el consumo de energía de la página. Si necesitas analíticas, considera elegir una opción liviana como [Cabin][cabin], [Fathom][fathom] o [Plausible][plausible]. Las incrustaciones de videos de servicios como YouTube y Vimeo se pueden mejorar al [cargar el video cuando haya interacción del usuario][lazy-video]. Paquetes como [astro-embed][embed] pueden ser útiles para servicios comunes.

:::tip[¿Sabías esto?]
El análisis y compilación de JavaScript es una de las tareas más costosas que los navegadores deben realizar. En comparación con el renderizado de una imagen JPEG del mismo tamaño, [el procesamiento de JavaScript puede llevar más de 30 veces más tiempo][cost-of-js].
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## Hospedaje

El lugar donde se hospeda una página web puede tener un gran impacto en la ecología de tu sitio de documentación. Los centros de datos y las granjas de servidores pueden tener un alto consumo de electricidad y un uso intensivo del agua.

Elegir un proveedor de alojamiento que utilice energía renovable significa tener emisiones de carbono más bajas para tu sitio. El [Directorio de la Web Ecológica][gwb] es una herramienta que puede ayudarte a encontrar empresas de alojamiento.

[gwb]: https://www.thegreenwebfoundation.org/directory/

## Comparaciones

¿Curioso por saber cómo se comparan otros frameworks de documentación?
Estas pruebas con el [Calculadora de Carbono de Sitios Web][wcc] comparan páginas similares construidas con diferentes herramientas.

| Framework                   | CO₂ por visita a la página |
| --------------------------- | -------------------------- |
| [Starlight][sl-carbon]      | 0.01g                      |
| [VitePress][vp-carbon]      | 0.05g                      |
| [Docus][dc-carbon]          | 0.05g                      |
| [Sphinx][sx-carbon]         | 0.07g                      |
| [MkDocs][mk-carbon]         | 0.10g                      |
| [Nextra][nx-carbon]         | 0.11g                      |
| [docsify][dy-carbon]        | 0.11g                      |
| [Docusaurus][ds-carbon]     | 0.24g                      |
| [Read the Docs][rtd-carbon] | 0.24g                      |
| [GitBook][gb-carbon]        | 0.71g                      |

<small>Datos recopilados el 14 de mayo de 2023. Haz clic en un enlace para ver cifras actualizadas.</small>

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

## Más recursos

### Herramientas

- [Calculadora de Carbono de Sitios Web][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [Control de Carbono WebPageTest](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### Artículos y charlas

- [“Construyendo una web más ecológica”](https://youtu.be/EfPoOt7T5lg), charla de Michelle Barker
- [“Estrategias de desarrollo web sostenible dentro de una organización”](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/), artículo de Michelle Barker
- [“Una web sostenible para todos”](https://2021.stateofthebrowser.com/speakers/tom-greenwood/), charla de Tom Greenwood
- [“Cómo el contenido web puede afectar el consumo de energía”](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/), artículo de Benjamin Poulain y Simon Fraser

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/en/guides/assets/
[islands]: https://docs.astro.build/en/concepts/islands/
[wcc]: https://www.websitecarbon.com/
