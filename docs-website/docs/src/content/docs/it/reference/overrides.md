---
title: Riferimento Sostituzioni
description: Una panoramica dei componenti e delle proprietà dei componenti supportati dalle sostituzioni di Starlight.
tableOfContents:
  maxHeadingLevel: 4
---

Puoi sovrascrivere i componenti integrati di Starlight fornendo percorsi ai componenti sostitutivi nell'opzione di configurazione [`components`](/it/reference/configuration/#components) di Starlight.
Questa pagina elenca tutti i componenti disponibili per l'override e si collega alle loro implementazioni predefinite su GitHub.

Scopri di più nella [Guida alla sostituzione dei componenti](/it/guides/overriding-components/).

## Proprietà dei componenti

Tutti i componenti possono accedere a un oggetto standard `Astro.props` che contiene informazioni sulla pagina corrente.

Per aggiungere i tipi di dato ai tuoi componenti personalizzati, importa il tipo `Props` da Starlight:

```astro
---
// src/components/Custom.astro
import type { Props } from '@astrojs/starlight/props';

const { hasSidebar } = Astro.props;
//      ^ tipo: boolean
---
```

Questo ti darà il completamento automatico e i tipi di dato quando accedi a `Astro.props`.

### Proprietà

Starlight trasmetterà le seguenti proprietà ai tuoi componenti personalizzati.

#### `dir`

**tipo:** `'ltr' | 'rtl'`

Direzione di scrittura della pagina.

#### `lang`

**tipo:** `string`

Tag di lingua BCP-47 per le impostazioni internazionali di questa pagina, ad es. `en`, `zh-CN` o `pt-BR`.

#### `locale`

**tipo:** `string | undefined`

Il percorso di base in cui viene servita una lingua. `undefined` per gli slug della lingua di base.

#### `slug`

**tipo:** `string`

Lo slug per questa pagina generato dal nome del file di contenuto.

#### `id`

**tipo:** `string`

L'ID univoco per questa pagina in base al nome del file di contenuto.

#### `isFallback`

**tipo:** `true | undefined`

`true` se questa pagina non è tradotta nella lingua corrente e utilizza contenuti di riserva dalle impostazioni di lingua predefinite.
Utilizzato solo in siti multilingue.

#### `entryMeta`

**tipo:** `{ dir: 'ltr' | 'rtl'; lang: string }`

Metadati di lingua per il contenuto della pagina. Può essere diverso dai valori di lingua di livello superiore quando una pagina utilizza contenuti di fallback.

#### `entry`

La voce della raccolta dei contenuti Astro per la pagina corrente.
Include i valori frontmatter per la pagina corrente in `entry.data`.

```ts
entry: {
  data: {
    title: string;
    description: string | undefined;
    // ecc.
  }
}
```

Scopri di più sulla forma di questo oggetto nel riferimento [Tipo di voce della raccolta di Astro](https://docs.astro.build/it/reference/api-reference/#collection-entry-type).

#### `sidebar`

**tipo:** `SidebarEntry[]`

Voci della barra laterale di navigazione del sito per questa pagina.

#### `hasSidebar`

**tipo:** `boolean`

Se la barra laterale deve essere visualizzata o meno in questa pagina.

#### `pagination`

**tipo:** `{ prev?: Link; next?: Link }`

Collegamenti alla pagina precedente e successiva nella barra laterale, se abilitata.

#### `toc`

**tipo:** `{ minHeadingLevel: number; maxHeadingLevel: number; items: TocItem[] } | undefined`

Sommario per questa pagina se abilitato.

#### `headings`

**tipo:** `{ depth: number; slug: string; text: string }[]`

Matrice di tutte le intestazioni Markdown estratte dalla pagina corrente.
Utilizza invece [`toc`](#toc) se vuoi creare un sommario che rispetti le opzioni di configurazione di Starlight.

#### `lastUpdated`

**tipo:** `Date | undefined`

Oggetto JavaScript `Date` che rappresenta l'ultimo aggiornamento di questa pagina, se abilitato.

#### `editUrl`

**tipo:** `URL | undefined`

Oggetto `URL` per l'indirizzo in cui questa pagina può essere modificata se abilitata.

#### `labels`

**tipo:** `Record<string, string>`

Un oggetto contenente stringhe dell'interfaccia utente localizzate per la pagina corrente. Consulta la guida [“Tradurre l'interfaccia di Starlight”](/it/guides/i18n/#tradurre-linterfaccia-starlight) per un elenco di tutte le chiavi disponibili.

---

## Componenti

### Head

Questi componenti vengono renderizzati all'interno dell'elemento `<head>` di ciascuna pagina.
Dovrebbero includere solo [elementi consentiti all'interno di `<head>`](https://developer.mozilla.org/it/docs/Web/HTML/Element/head#see_also).

#### `Head`

**Componente standard:** [`Head.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Head.astro)

Componente renderizzato all'interno di `<head>` di ogni pagina.
Include tag importanti tra cui `<title>` e `<meta charset="utf-8">`.

Sostituisci questo componente come ultima risorsa.
Se possibile, preferisci l'opzione di configurazione [`head`](/it/reference/configuration/#head) di Starlight.

#### `ThemeProvider`

**Componente standard:** [`ThemeProvider.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeProvider.astro)

Componente renderizzato all'interno di `<head>` che imposta il supporto del tema scuro/chiaro.
L'implementazione predefinita include uno script in linea e un `<template>` utilizzato dallo script in [`<ThemeSelect />`](#themeselect).

---

### Accessibilità

#### `SkipLink`

**Componente standard:** [`SkipLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SkipLink.astro)

Componente renderizzato come primo elemento all'interno di `<body>` che si collega al contenuto della pagina principale per l'accessibilità.
L'implementazione predefinita è nascosta finché un utente non la focalizza premendo il tasto tab con la tastiera.

---

### Layout

Questi componenti sono responsabili del layout dei componenti di Starlight e della gestione delle visualizzazioni attraverso diversi punti di interruzione.
L'override di questi comporta una complessità significativa.
Quando possibile, prediligi sovrascrivere un componente di livello inferiore.

#### `PageFrame`

**Componente standard:** [`PageFrame.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageFrame.astro)

Componente di layout avvolto attorno alla maggior parte del contenuto della pagina.
L'implementazione predefinita imposta il layout header-sidebar-main e include slot denominati `header` e `sidebar` insieme a uno slot predefinito per il contenuto principale.
Renderizza inoltre [`<MobileMenuToggle />`](#mobilemenutoggle) per supportare l'attivazione/disattivazione della navigazione della barra laterale su piccole finestre (mobili).

#### `MobileMenuToggle`

**Componente standard:** [`MobileMenuToggle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuToggle.astro)

Componente reso all'interno di [`<PageFrame>`](#pageframe) responsabile dell'attivazione/disattivazione della navigazione della barra laterale su piccoli viewport (mobili).

#### `TwoColumnContent`

**Componente standard:** [`TwoColumnContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TwoColumnContent.astro)

Componente di layout avvolto attorno alla colonna del contenuto principale e alla barra laterale destra (sommario).
L'implementazione predefinita gestisce il passaggio da un layout a colonna singola con viewport piccolo a un layout a due colonne con viewport più grande.

---

### Intestazione

Questi componenti eseguono il rendering della barra di navigazione superiore di Starlight.

#### `Header`

**Componente standard:** [`Header.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Header.astro)

Componente dell'intestazione renderizzato nella parte superiore di ogni pagina.
L'implementazione predefinita renderizza [`<SiteTitle />`](#sitetitle), [`<Search />`](#search), [`<SocialIcons />`](#socialicons), [`<ThemeSelect />`](#themeselect) e [`<LanguageSelect />`](#languageselect).

#### `SiteTitle`

**Componente standard:** [`SiteTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SiteTitle.astro)

Componente renderizzato all'inizio dell'intestazione del sito per visualizzare il titolo del sito.
L'implementazione predefinita include la logica per il rendering dei loghi definita nella configurazione di Starlight.

#### `Search`

**Componente standard:** [`Search.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Search.astro)

Componente utilizzato per eseguire il rendering dell'interfaccia utente di ricerca di Starlight.
L'implementazione predefinita include il pulsante nell'intestazione e il codice per visualizzare una schermata di ricerca quando viene cliccata e caricare l'[interfaccia utente di Pagefind](https://pagefind.app/).

Quando [`pagefind`](/it/reference/configuration/#pagefind) è disabilitato, il componente di ricerca predefinito non verrà renderizzato.
Tuttavia, se si sovrascrive `Search`, il componente personalizzato verrà sempre renderizzato anche se l'opzione di configurazione `pagefind` è `false`.
Ciò consente di aggiungere un'interfaccia utente per i provider di ricerca alternativi quando si disabilita Pagefind.

#### `SocialIcons`

**Componente standard:** [`SocialIcons.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SocialIcons.astro)

Componente renderizzato nell'intestazione del sito, inclusi i collegamenti alle icone social.
L'implementazione predefinita utilizza l'opzione [`social`](/it/reference/configuration/#social) nella configurazione di Starlight per eseguire il rendering di icone e collegamenti.

#### `ThemeSelect`

**Componente standard:** [`ThemeSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeSelect.astro)

Componente renderizzato nell'intestazione del sito che consente agli utenti di selezionare la combinazione di colori preferita.

#### `LanguageSelect`

**Componente standard:** [`LanguageSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LanguageSelect.astro)

Componente renderizzato nell'intestazione del sito che consente agli utenti di passare a una lingua diversa.

---

### Barra Laterale Globale

La barra laterale globale di Starlight include la navigazione principale del sito.
Nelle finestre strette questa è nascosta dietro un menu a discesa.

#### `Sidebar`

**Componente standard:** [`Sidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Sidebar.astro)

Componente visualizzato prima del contenuto della pagina che contiene la navigazione globale.
L'implementazione predefinita viene visualizzata come barra laterale su viewport sufficientemente ampi e all'interno di un menu a discesa su viewport piccoli (mobili).
Visualizza inoltre [`<MobileMenuFooter />`](#mobilemenufooter) per mostrare elementi aggiuntivi all'interno del menu mobile.

#### `MobileMenuFooter`

**Componente standard:** [`MobileMenuFooter.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuFooter.astro)

Componente visualizzato nella parte inferiore del menu a discesa mobile.
L'implementazione predefinita visualizza [`<ThemeSelect />`](#themeselect) e [`<LanguageSelect />`](#languageselect).

---

### Barra Laterale della Pagina

La barra laterale della pagina di Starlight è responsabile della visualizzazione di un sommario che delinea i sottotitoli della pagina corrente.
Nelle finestre strette questo si comprime in un menu a discesa fisso.

#### `PageSidebar`

**Componente standard:** [`PageSidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageSidebar.astro)

Componente renderizzato prima del contenuto della pagina principale per visualizzare un sommario.
L'implementazione predefinita rende [`<TableOfContents />`](#tableofcontents) e [`<MobileTableOfContents />`](#mobiletableofcontents).

#### `TableOfContents`

**Componente standard:** [`TableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TableOfContents.astro)

Componente che renderizza il sommario della pagina corrente su finestre più ampie.

#### `MobileTableOfContents`

**Componente standard:** [`MobileTableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileTableOfContents.astro)

Componente che renderizza il sommario della pagina corrente su piccoli viewport (mobili).

---

### Contenuto

Questi componenti vengono visualizzati nella colonna principale del contenuto della pagina.

#### `Banner`

**Componente standard:** [`Banner.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Banner.astro)

Componente banner renderizzato nella parte superiore di ogni pagina.
L'implementazione predefinita utilizza il valore frontmatter [`banner`](/it/reference/frontmatter/#banner) della pagina per decidere se renderizzare o meno.

#### `ContentPanel`

**Componente standard:** [`ContentPanel.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ContentPanel.astro)

Componente di layout utilizzato per racchiudere le sezioni della colonna del contenuto principale.

#### `PageTitle`

**Componente standard:** [`PageTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageTitle.astro)

Componente contenente l'elemento `<h1>` per la pagina corrente.

Le implementazioni dovrebbero garantire di impostare `id="_top"` sull'elemento `<h1>` come nell'implementazione predefinita.

#### `FallbackContentNotice`

**Componente standard:** [`FallbackContentNotice.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/FallbackContentNotice.astro)

Avviso visualizzato agli utenti nelle pagine in cui non è disponibile una traduzione per la lingua corrente.
Utilizzato solo su siti multilingue.

#### `Hero`

**Componente standard:** [`Hero.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Hero.astro)

Componente renderizzato nella parte superiore della pagina quando [`hero`](/it/reference/frontmatter/#hero) è impostato in frontmatter.
L'implementazione predefinita mostra un titolo di grandi dimensioni, uno slogan e collegamenti di invito all'azione insieme a un'immagine facoltativa.

#### `MarkdownContent`

**Componente standard:** [`MarkdownContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MarkdownContent.astro)

Componente renderizzato attorno al contenuto principale di ogni pagina.
L'implementazione predefinita imposta gli stili di base da applicare al contenuto Markdown.

Anche gli stili del contenuto Markdown sono esposti in `@astrojs/starlight/style/markdown.css` e limitati alla classe CSS `.sl-markdown-content`.

---

### Piè di pagina

Questi componenti vengono visualizzati nella parte inferiore della colonna principale del contenuto della pagina.

#### `Footer`

**Componente standard:** [`Footer.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Footer.astro)

Componente piè di pagina renderizzato nella parte inferiore di ogni pagina.
L'implementazione predefinita visualizza [`<LastUpdated />`](#lastupdated), [`<Pagetion />`](#pagination) e [`<EditLink />`](#editlink).

#### `LastUpdated`

**Componente standard:** [`LastUpdated.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LastUpdated.astro)

Componente renderizzato nel piè di pagina per visualizzare la data dell'ultimo aggiornamento.

#### `EditLink`

**Componente standard:** [`EditLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/EditLink.astro)

Componente renderizzato nel piè di pagina per visualizzare un collegamento al punto in cui è possibile modificare la pagina.

#### `Pagination`

**Componente standard:** [`Pagination.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Pagination.astro)

Componente renderizzato nel piè di pagina per visualizzare le frecce di navigazione tra le pagine precedenti/successive.
