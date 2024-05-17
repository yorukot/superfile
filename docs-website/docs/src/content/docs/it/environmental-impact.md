---
title: Documentazione ecologica
description: Impara come Starlight può aiutarti a costruire siti per documentazione più verdi e ridurre la tua impronta ecologica.
---

Le stime dell'impatto ambientale del web si aggirano tra il [2%][sf] e il [4% delle emissioni globali di gas serra][bbc], equivalente circa alle emissioni dell'industria aerea.
Ci sono molti complessi fattori per calcolare l'impatto ambientale di un sito web, ma questa guida include dei consigli per ridurre l'impronta ecologica del tuo sito.

La buona notizia è che scegliere Starlight è un ottimo punto di partenza.
Secondo il Website Carbon Calculator, questo sito è [più pulito del 99% delle pagine web testate][sl-carbon], producendo 0,01 g di CO₂ per pagina visitata.

## Peso per pagina

Più dati una pagina web trasferisce, più risorse energetiche sono necessarie.
Nell'aprile 2023, la pagina web media richiede all'utente di scaricare più di 2.000 KB secondo i [dati dell'HTTP Archive][http].

Starlight costruisce le pagine nel modo più leggero possibile.
Per esempio, alla prima visita, l'utente scaricherà meno di 50 KB di dati compressi — soltanto il 2,5 % della media dell'HTTP Archive.
Inoltre, con una strategia appropriata di cache, successive visite potranno richiedere solamente 10 KB.

### Immagini

Anche se Starlight fornisce un buon punto di partenza, le immagini possono velocemente aumentare il peso della pagina.
Starlight usa il [supporto ottimizzato degli asset][assets] di Astro per ottimizzare le immagini nei file Markdown e MDX.

### Componenti UI

I componenti costruiti con framework UI come React o Vue possono facilmente aggiungere grandi quantità di JavaScript nella pagina.
Dato che Starlight è costruito su Astro, i componenti come questo caricano **zero JavaScript all'utente per default** grazie alle [Isole Astro][islands].

### Caching

La cache è utilizzata per salvare per un periodo di tempo dati già scaricati in modo che il browser possa riutilizzarli.
Una buona strategia di cache permette all'utente di scaricare nuovi contenuti il prima possibile quando cambiano, ma anche evitare di scaricare nuovamente gli stessi più volte quando non sono cambiati.

Il modo più comune per configurarla è grazie al [`Cache-Control` header HTTP][cache].
Quando si utilizza Starlight, puoi controllare per quanto tempo salvare in cache nella cartella `/_astro/`.
Questa cartella contiene CSS, JavaScript e altri asset che possono essere salvati in cache per sempre, riducendo download non necessari:

```
Cache-Control: public, max-age=604800, immutable
```

La configurazione della cache dipende dal tuo host. Per esempio, Vercel applica questa strategia di cache per te senza configurazione richiesta, mentre puoi specificare degli [headers personalizzati per Netlify][ntl-headers] aggiungendo il file `public/_headers` nel progetto:

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## Consumo energetico

Come una pagina web sia costruita può impattare l'energia richiesta per visualizzarla nel dispositivo dell'utente.
Starlight riduce il consumo energetico del cellulare, tablet o computer che l'utente utilizza sfruttando pochissimo JavaScript.

Bisogna fare attenzione quando si vuole aggiungere funzionalità come script analitici o contenuti che utilizzano grandi quantità di JavaScript come video incorporati, dato che questi possono aumentare il consumo energetico.
Se necessario, si considerino [Cabin][cabin], [Fathom][fathom], o [Plausible][plausible] per funzionalità analitiche siccome non sono pesanti.
Video integrati come YouTube e Vimeo possono essere migliorati aspettando di [caricare il video dopo interazione][lazy-video].
Pacchetti come [`astro-embed`][embed] possono aiutare per servizi comuni.

:::tip[Lo sapevi ?]
L'analisi e la compilazione di JavaScript è una delle operazioni più costose che il browser deve fare.
Se si confronta con la visualizzazione di un'immagine JPEG dello stesso peso, [JavaScript può richiedere anche più di 30 volte il tempo necessario per essere processato][cost-of-js].
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## Hosting

La piattaforma utilizzata per hostare un sito può avere un impatto significativo per l'impronta ecologica.
I data center e server farm possono impattare di molto l'ambiente, usando grandi quantità di energia elettrica e d'acqua.

Scegliere un host che usi energia da fonti rinnovabili significa ridurre le emissioni di gas serra per il tuo sito. Il [Green Web Directory][gwb] è uno degli strumenti che si possono utilizzare per trovare host di questo tipo.

[gwb]: https://www.thegreenwebfoundation.org/directory/

## Comparazioni

Curioso di come altri framework per documentazioni si comparano?
Questi test eseguiti con [Website Carbon Calculator][wcc] confrontano pagine simili costruite con diversi framework.

| Framework                   | CO₂ per visita |
| --------------------------- | -------------- |
| [Starlight][sl-carbon]      | 0,01 g         |
| [VitePress][vp-carbon]      | 0,05 g         |
| [Docus][dc-carbon]          | 0,05 g         |
| [Sphinx][sx-carbon]         | 0,07 g         |
| [MkDocs][mk-carbon]         | 0,10 g         |
| [Nextra][nx-carbon]         | 0,11 g         |
| [docsify][dy-carbon]        | 0,11 g         |
| [Docusaurus][ds-carbon]     | 0,24 g         |
| [Read the Docs][rtd-carbon] | 0,24 g         |
| [GitBook][gb-carbon]        | 0,71 g         |

<small>Dati collezionati il 14 Maggio 2023. Clicca i link per vedere i dati aggiornati.</small>

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

## Altre risorse

### Strumenti

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Carbon Control](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### Articoli e discussioni

- [“Building a greener web”](https://youtu.be/EfPoOt7T5lg), conferenza di Michelle Barker
- [“Sustainable Web Development Strategies Within An Organization”](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/), articolo di Michelle Barker
- [“A sustainable web for everyone”](https://2021.stateofthebrowser.com/speakers/tom-greenwood/), conferenza di Tom Greenwood
- [“How Web Content Can Affect Power Usage”](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/), articolo di Benjamin Poulain e Simon Fraser

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/en/guides/assets/
[islands]: https://docs.astro.build/en/concepts/islands/
[wcc]: https://www.websitecarbon.com/
