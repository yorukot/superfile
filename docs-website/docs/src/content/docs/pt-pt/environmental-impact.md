---
title: Documentação ecológica
description: Aprenda como o Starlight te pode ajudar a construir sites de documentação mais ecológicos e reduzir a pegada de carbono.
---

As estimativas para o impacto climático da indústria web variam entre [2%][sf] e [4% das emissões globais de carbono][bbc], aproximadamente o equivalente às emissões da indústria aeronáutica.
Há vários fatores complexos no cálculo do impacto ecológico de um website, mas este guia inclui algumas dicas para reduzir a pegada ambiental do seu site de documentação.

A boa noticía é que, escolher o Starlight já é um bom começo!
De acordo com o "Website Carbon Calculator", este site é [mais limpo do que 99% das páginas web testadas][sl-carbon], produzindo 0.01g de CO₂ por cada visita à página.

## Peso da página

Quanto mais dados uma página web transfere mais recursos energéticos são necessários.
De acordo com [dados do HTTP Archive][http], em Abril de 2023, uma página web mediana necessitava que o utilizador baixasse mais de 2,000 KB.

O Starlight constrói páginas que são o mais leve possível.
Por exemplo, na primeira visita, um utilizador vai descarregar menos do que 50 KB de dados comprimidos, ou seja, apenas 2.5% da mediana indicada pelo HTTP Archive.
Com uma boa estratégia de _cache_, as navegações subsequentes podem descarregar tão pouco quanto 10 KB.

### Imagens

Enquanto o Starlight providencia uma boa base, as imagens que você adicionar à sua documentação podem aumentar o peso da sua página rapidamente.
O Starlight usa o [suporte a assets otimizados][assets] do Astro para otimizar imagens locais nos seus arquivos Markdown e MDX.

### Componentes de UI

Os componentes construídos com frameworks de UI como React ou Vue podem facilmente adicionar grandes quantidades de JavaScript a uma página.
Porque o Starlight é construído com o Astro, e graças às [Ilhas Astro][islands], esses componentes carregam, **por padrão, zero código JavaScript no lado do cliente**.

### _Cache_

A _Cache_ é usada para controlar por quanto tempo um navegador armazena e reutiliza os dados já descarregados.
Uma boa estratégia de _caching_ garante que um utilizador receba o conteúdo novo o mais cedo possível assim que ele muda, mas também evita descarregar inutil e repetidamente o mesmo conteúdo sem que ele mude.

A forma mais comum de configurar a _cache_ é com o [header HTTP `Cache-Control`][cache].
Ao utilizar o Starlight, você pode definir um grande tempo de cache para todo o conteúdo do diretório `/_astro/`.
Este diretório contém CSS, JavaScript e outros artefactos em _bundle_ que podem ser seguramente _cached_ para sempre, reduzindo assim downloads desnecessários:

```
Cache-Control: public, max-age=604800, immutable
```

A forma de configurar a _cache_ depende do seu alojamento web. Por exemplo, o Vercel aplica por você esta estratégia de _cache_ sem necessidade de configuração adicional, já a definição de [cabeçalhos customizados para Netlify][ntl-headers] necessita que adicione um arquivo `public/_headers` ao seu projeto:

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## Consumo de energia

A forma com que uma página web é construída pode ter impacto na energia necessária para executá-la no dispositivo de um utilizador.
Por utilizar JavaScript ao mínimo, o Starlight reduz a quantidade de energia de processamento que o celular, tablet ou computador de um utilizador precisa para carregar e renderizar páginas.

Tenha atenção ao adicionar funcionalidades como scripts de rastreamento de _Analytics_ ou conteúdo cheio de JavaScript como embeds de vídeo já que estes podem aumentar o consumo de energia da página.
Se você precisa de _Analytics_, considere escolher uma opção leve como [Cabin][cabin], [Fathom][fathom] ou [Plausible][plausible].
Embeds como vídeos do YouTube e Vimeo podem ser melhorados se [carregar o vídeo mediante a interação do usuário][lazy-video].
Pacotes como o [`astro-embed`][embed] podem ajudá-lo com alguns dos serviços comuns.

:::tip[Sabia que?]
Fazer parse e compilação de JavaScript é uma das tarefas mais caras que os navegadores tem que fazer.
Comparado com a renderização de uma imagem JPEG de mesmo tamanho, [o JavaScript pode levar mais do que 30 vezes o tempo para processar][cost-of-js].
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## Alojamento

O lugar onde uma página web é alojada pode ter um grande impacto no quão amigável ao ambiente o seu site de documentação é.
Os centros de dados e de servidores podem ter um grande impacto ecológico, incluindo alto consumo de eletricidade e uso intensivo de água.

Escolher um alojamento que utiliza energia renovável significará menos emissões de carbono para o seu site. A [Green Web Directory][gwb] é uma ferramenta que poderá ajudá-lo a encontrar empresas de alojamento.

[gwb]: https://www.thegreenwebfoundation.org/directory/

## Comparações

Está curioso para comparar com os outros frameworks de documentação?
Estes testes realizados com o [Website Carbon Calculator][wcc] comparam páginas semelhantes construídas com diferentes ferramentas.

| Framework                   | CO₂ por visita da página |
| --------------------------- | ------------------------ |
| [Starlight][sl-carbon]      | 0.01g                    |
| [VitePress][vp-carbon]      | 0.05g                    |
| [Docus][dc-carbon]          | 0.05g                    |
| [Sphinx][sx-carbon]         | 0.07g                    |
| [MkDocs][mk-carbon]         | 0.10g                    |
| [Nextra][nx-carbon]         | 0.11g                    |
| [docsify][dy-carbon]        | 0.11g                    |
| [Docusaurus][ds-carbon]     | 0.24g                    |
| [Read the Docs][rtd-carbon] | 0.24g                    |
| [GitBook][gb-carbon]        | 0.71g                    |

<small>Dados recolhidos a 14 de Maio de 2023. Clique num dos links para ver os valores atualizados.</small>

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

## Mais recursos

### Ferramentas

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Carbon Control](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### Artigos e palestras

- [“Building a greener web”](https://youtu.be/EfPoOt7T5lg), palestra por Michelle Barker
- [“Sustainable Web Development Strategies Within An Organization”](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/), artigo por Michelle Barker
- [“A sustainable web for everyone”](https://2021.stateofthebrowser.com/speakers/tom-greenwood/), palestra por Tom Greenwood
- [“How Web Content Can Affect Power Usage”](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/), artigo por Benjamin Poulain e Simon Fraser

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/pt-pt/guides/assets/
[islands]: https://docs.astro.build/pt-pt/concepts/islands/
[wcc]: https://www.websitecarbon.com/
