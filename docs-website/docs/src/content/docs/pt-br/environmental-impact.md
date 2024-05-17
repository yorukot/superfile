---
title: Documentação amigável ao meio ambiente
description: Aprenda como Starlight pode te ajudar a construir sites de documentação mais verdes e reduzir sua pegada de carbono.
---

Estimativas para o impacto climático da indústria web variam entre [2%][sf] e [4% da emissão global de carbono][bbc], aproximadamente equivalente as emissões da indústria aeronáutica.
Há vários fatores complexos no cálculo do impacto ecológico de um website, mas este guia inclui algumas dicas para reduzir a pegada ambiental do seu site de documentação.

A boa noticía é, escolher Starlight é um bom começo.
De acordo com o "Website Carbon Calculator", este site é [mais limpo que 99% das páginas web testadas][sl-carbon], produzindo 0.01g de CO₂ por visita da página.

## Peso da página

Quanto mais dados uma página web transfere, mais recursos energéticos são necessários.
Em Abril de 2023, uma página web mediana necessitava que o usuário baixasse mais que 2,000 KB de acordo com [dados do HTTP Archive][http].

Starlight constrói páginas que são o mais leve possível.
Por exemplo, em uma primeira visita, um usuário vai baixar menos do que 50 KB de dados comprimidos — apenas 2.5% da mediana do HTTP Archive.
Com uma boa estratégia de cacheamento, navegações subsequentes podem baixar tão pouco quanto 10 KB.

### Imagens

Enquanto Starlight providencia uma boa base, imagens que você adiciona a sua documentação podem rapidamente aumentar o peso da sua página.
Starlight usa o [suporte a assets otimizados][assets] do Astro para otimizar imagens locais em seus arquivos Markdown e MDX.

### Componentes de UI

Componentes construídos com frameworks de UI como React ou Vue podem facilmente adicionar grandes quantidades de JavaScript a uma página.
Pelo Starlight ser construído com Astro, componentes assim carregam **zero JavaScript no lado do cliente por padrão** graças a [Ilhas Astro][islands].

### Cacheamento

Cacheamento é usado para controlar por quanto tempo um navegador armazena e reutiliza dados já baixados.
Uma boa estratégia de cacheamento garante que um usuário receba conteúdo novo o mais cedo possível quando ele muda, mas também evita baixar inutilmente o mesmo conteúdo de novo e de novo enquanto ele não mudou.

A forma mais comum de configurar cacheamento é com o [header HTTP `Cache-Control`][cache].
Enquanto utiliza Starlight, você pode definir um grande tempo de cache para tudo no diretório `/_astro/`.
Esse diretório contém CSS, JavaScript e outros assets em bundle que podem ser seguramente cacheados para sempre, reduzindo downloads desnecessários:

```
Cache-Control: public, max-age=604800, immutable
```

Como configurar cacheamento depende da sua hospedagem web. Por exemplo, a Vercel aplica essa estratégia de cacheamento para você sem nenhuma configuração necessária, enquanto você pode definir [cabeçalhos customizados para Netlify][ntl-headers] ao adicionar um arquivo `public/_headers` ao seu projeto:

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## Consumo de energia

A forma com que uma página web é construída pode impactar a energia necessária para executá-la no dispositivo de um usuário.
Por utilizar JavaScript ao mínimo, Starlight reduz a quantidade de poder de processamento que o celular, tablet ou computador de um usuário precisa para carregar e renderizar páginas.

Seja cuidadoso ao adicionar funcionalidades como scripts de rastreamento de analytics ou conteúdo cheio de JavaScript como embeds de vídeo já que estes podem aumentar o consumo de energia da página.
Se você precisa de analytics, considere escolher uma opção leve como [Cabin][cabin], [Fathom][fathom] ou [Plausible][plausible].
Embeds como vídeos do YouTube e Vimeo podem ser melhorados por esperar para [carregar o vídeo conforme interação do usuário][lazy-video].
Pacotes como [`astro-embed`][embed] podem ajudar com serviços comuns.

:::tip[Você sabia?]
Fazer parse e compilação de JavaScript é uma das tarefas mais caras que navegadores tem que fazer.
Comparado a renderizar uma imagem JPEG de mesmo tamanho, [JavaScript pode levar mais do que 30 vezes mais tempo para processar][cost-of-js].
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## Hospedagem

O lugar onde uma página web é hospedada podem ter um grande impacto no quão amigável ao ambiente seu site de documentação é.
Centro de dados e fazendas de servidores podem ter um grande impacto ecológico, incluindo alto consumo de eletricidade e uso intensivo de água.

Escolher uma hospedagem que utiliza energia renovável significará menos emissões de carbono para o seu site. A [Green Web Directory][gwb] é uma ferramenta que pode ajudá-lo a encontrar empresas de hospedagem.

[gwb]: https://www.thegreenwebfoundation.org/directory/

## Comparações

Curioso para saber como outros frameworks de documentação se comparam?
Esses testes com o [Website Carbon Calculator][wcc] comparam páginas similares construídas com diferentes ferramentas.

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

<small>Dados coletados em 14 de Maio de 2023. Clique num dos links para ver valores atualizados.</small>

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
[assets]: https://docs.astro.build/pt-br/guides/assets/
[islands]: https://docs.astro.build/pt-br/concepts/islands/
[wcc]: https://www.websitecarbon.com/
