---
title: Référence des redéfinitions
description: Une vue d'ensemble de tous les composants et les props des composants supportés par les redéfinitions de Starlight.
tableOfContents:
  maxHeadingLevel: 4
---

Vous pouvez redéfinir les composants intégrés à Starlight en spécifiant des chemins vers des composants de remplacement avec l'option de configuration [`components`](/fr/reference/configuration/#components) de Starlight.
Cette page répertorie tous les composants disponibles qui peuvent être redéfinis et fournit des liens vers leurs implémentations par défaut sur GitHub.

Pour en savoir plus, consultez le [guide des redéfinitions de composants](/fr/guides/overriding-components/).

## Props des composants

Tous les composants peuvent accéder à un objet `Astro.props` standard qui contient des informations concernant la page courante.

Pour typer vos composants personnalisés, importez le type `Props` depuis Starlight :

```astro
---
// src/components/Custom.astro
import type { Props } from '@astrojs/starlight/props';

const { hasSidebar } = Astro.props;
//      ^ type: boolean
---
```

Cela vous permettra d'obtenir de l'autocomplétion et un typage lors de l'utilisation de `Astro.props`.

### Props

Starlight passera les props suivantes à vos composants personnalisés.

#### `dir`

**Type :** `'ltr' | 'rtl'`

Le sens d'écriture de la page.

#### `lang`

**Type :** `string`

L’étiquette d’identification BCP-47 pour la langue de la page, par exemple `en`, `zh-CN` ou `pt-BR`.

#### `locale`

**Type :** `string | undefined`

Le chemin de base utilisé pour servir une langue. `undefined` pour les slugs de la locale racine.

#### `siteTitle`

**Type :** `string`

Le titre du site pour la langue de cette page.

#### `siteTitleHref`

**Type :** `string`

La valeur de l’attribut `href` du titre du site, renvoyant à la page d'accueil, par exemple `/`.
Pour les sites multilingues, cette valeur inclura la locale actuelle, par exemple `/fr/` ou `/zh-cn/`.

#### `slug`

**Type :** `string`

Le slug de la page généré à partir du nom du fichier du contenu.

#### `id`

**Type :** `string`

L'identifiant unique de cette page basé sur le nom du fichier du contenu.

#### `isFallback`

**Type :** `true | undefined`

`true` si cette page n'est pas traduite dans la langue actuelle et utilise le contenu de la langue par défaut en tant que repli.
Utilisé uniquement dans les sites multilingues.

#### `entryMeta`

**Type :** `{ dir: 'ltr' | 'rtl'; lang: string }`

Métadonnées de la locale pour le contenu de la page. Peut être différent des valeurs de locale de premier niveau lorsque la page utilise un contenu de repli.

#### `entry`

L'entrée de la collection de contenu Astro pour la page courante.
Inclut les valeurs du frontmatter pour la page courante dans `entry.data`.

```ts
entry: {
  data: {
    title: string;
    description: string | undefined;
    // etc.
  }
}
```

Pour en savoir plus sur le format de cet objet, consultez la [référence du type d'entrée de collection](https://docs.astro.build/fr/reference/api-reference/#collection-entry-type).

#### `sidebar`

**Type :** `SidebarEntry[]`

Les entrées de la barre latérale de navigation du site pour cette page.

#### `hasSidebar`

**Type :** `boolean`

Indique si la barre latérale est affichée sur cette page.

#### `pagination`

**Type :** `{ prev?: Link; next?: Link }`

Liens vers la page précédente et suivante dans la barre latérale si celle-ci est activée.

#### `toc`

**Type :** `{ minHeadingLevel: number; maxHeadingLevel: number; items: TocItem[] } | undefined`

Table des matières de la page courante si celle-ci est activée.

#### `headings`

**Type :** `{ depth: number; slug: string; text: string }[]`

Un tableau de toutes les en-têtes Markdown extraites de la page courante.
Utilisez [`toc`](#toc) à la place si vous souhaitez construire un composant de table des matières qui respecte les options de configuration de Starlight.

#### `lastUpdated`

**Type :** `Date | undefined`

Un objet JavaScript de type `Date` représentant la date de dernière mise à jour de cette page si cette fonctionnalité est activée.

#### `editUrl`

**Type :** `URL | undefined`

Un objet `URL` de l'adresse où cette page peut être modifiée si cette fonctionnalité est activée.

#### `labels`

**Type :** `Record<string, string>`

Un objet contenant les chaînes de l’interface utilisateur localisées pour la page courante. Consultez le guide [« Traduire l’interface utilisateur de Starlight »](/fr/guides/i18n/#traduire-linterface-utilisateur-de-starlight) pour obtenir une liste de toutes les clés disponibles.

---

## Composants

### Métadonnées

Ces composants sont utilisés dans l'élément `<head>` de chaque page.
Ils ne doivent inclure que des [éléments autorisés à l'intérieur de `<head>`](https://developer.mozilla.org/fr/docs/Web/HTML/Element/head).

#### `Head`

**Composant par défaut :** [`Head.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Head.astro)

Composant utilisé à l'intérieur de l'élément `<head>` de chaque page.
Inclut des balises importantes comme `<title>` et `<meta charset="utf-8">`.

Redéfinissez ce composant en dernier recours.
Préférez l'option [`head`](/fr/reference/configuration/#head) de la configuration de Starlight si possible.

#### `ThemeProvider`

**Composant par défaut :** [`ThemeProvider.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeProvider.astro)

Composant utilisé à l'intérieur de l'élément `<head>` qui configure la prise en charge du thème sombre/clair.
L'implémentation par défaut inclut un script en ligne et un élément `<template>` utilisé par le script situé dans [`<ThemeSelect />`](#themeselect).

---

### Accessibilité

#### `SkipLink`

**Composant par défaut :** [`SkipLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SkipLink.astro)

Composant utilisé comme premier élément à l'intérieur du `<body>` qui relie au contenu principal de la page pour des raisons d'accessibilité.
L'implémentation par défaut est masquée jusqu'à ce qu'il reçoive le focus d'un utilisateur utilisant la navigation au clavier.

---

### Mise en page

Ces composants sont responsables de la mise en page des composants de Starlight et de la gestion des vues pour différents points d'arrêt.
Redéfinir ceux-ci implique une complexité significative.
Lorsque cela est possible, préférez redéfinir un composant de plus bas niveau.

#### `PageFrame`

**Composant par défaut :** [`PageFrame.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageFrame.astro)

Composant de mise en page contenant la plupart du contenu de la page.
L'implémentation par défaut configure la mise en page de l'en-tête, de la barre latérale et du contenu principal et inclut des emplacements (slots) nommés `header` et `sidebar` en plus de l'emplacement par défaut pour le contenu principal.
Il affiche également [`<MobileMenuToggle />`](#mobilemenutoggle) qui prend en charge l'affichage de la barre latérale de navigation sur petits écrans (mobiles).

#### `MobileMenuToggle`

**Composant par défaut :** [`MobileMenuToggle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuToggle.astro)

Composant utilisé à l'intérieur de [`<PageFrame>`](#pageframe) qui est responsable de l'affichage de la barre latérale de navigation sur petits écrans (mobiles).

#### `TwoColumnContent`

**Composant par défaut :** [`TwoColumnContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TwoColumnContent.astro)

Composant de mise en page enveloppant le contenu principal de la page et la barre latérale de droite (table des matières).
L'implémentation par défaut prend en charge le changement entre une mise en page à une seule colonne pour petits écrans et une mise en page à deux colonnes pour écrans plus larges.

---

### En-tête

Ces composants affichent la barre de navigation supérieure de Starlight.

#### `Header`

**Composant par défaut :** [`Header.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Header.astro)

Composant d'en-tête affiché en haut de chaque page.
L'implémentation par défaut affiche [`<SiteTitle />`](#sitetitle-1), [`<Search />`](#search), [`<SocialIcons />`](#socialicons), [`<ThemeSelect />`](#themeselect) et [`<LanguageSelect />`](#languageselect).

#### `SiteTitle`

**Composant par défaut :** [`SiteTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SiteTitle.astro)

Composant utilisé au début de l'en-tête du site pour afficher le titre du site.
L'implémentation par défaut inclut la logique pour afficher les logos définis dans la configuration de Starlight.

#### `Search`

**Composant par défaut :** [`Search.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Search.astro)

Composant utilisé pour afficher l'interface de recherche de Starlight.
L'implémentation par défaut inclut le bouton dans l'en-tête et le code pour afficher une fenêtre modale de recherche lorsqu'il est cliqué et charger [l'interface utilisateur de Pagefind](https://pagefind.app/).

Lorsque [`pagefind`](/fr/reference/configuration/#pagefind) est désactivé, le composant de recherche par défaut ne sera pas affiché.
Cependant, si vous redéfinissez `Search`, votre composant personnalisé sera toujours affiché même si l'option de configuration `pagefind` est `false`.
Cela vous permet d'ajouter une interface de recherche alternative lorsque vous désactivez Pagefind.

#### `SocialIcons`

**Composant par défaut :** [`SocialIcons.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/SocialIcons.astro)

Composant utilisé dans l'en-tête du site qui inclut des liens avec des icônes vers différents médias sociaux.
L'implémentation par défaut utilise l'option [`social`](/fr/reference/configuration/#social) de la configuration de Starlight pour afficher les icônes et les liens.

#### `ThemeSelect`

**Composant par défaut :** [`ThemeSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ThemeSelect.astro)

Composant utilisé dans l'en-tête du site qui permet aux utilisateurs de sélectionner leur thème de couleur préféré.

#### `LanguageSelect`

**Composant par défaut :** [`LanguageSelect.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LanguageSelect.astro)

Component utilisé dans l'en-tête du site qui permet aux utilisateurs de changer de langue.

---

### Barre latérale globale

La barre latérale globale de Starlight contient la navigation principale du site.
Sur des écrans peu larges, elle est masquée derrière un menu déroulant.

#### `Sidebar`

**Composant par défaut :** [`Sidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Sidebar.astro)

Composant utilisé avant le contenu de la page qui contient la navigation globale.
L'implémentation par défaut est affichée comme une barre latérale sur des écrans suffisamment larges et à l'intérieur d'un menu déroulant sur des écrans plus petits (mobiles).
Il utilise aussi [`<MobileMenuFooter />`](#mobilemenufooter) pour afficher des éléments supplémentaires dans le menu mobile.

#### `MobileMenuFooter`

**Composant par défaut :** [`MobileMenuFooter.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileMenuFooter.astro)

Composant utilisé à la fin du menu déroulant mobile.
L'implémentation par défaut affiche [`<ThemeSelect />`](#themeselect) et [`<LanguageSelect />`](#languageselect).

---

### Barre latérale de page

La barre latérale de page de Starlight est responsable d'afficher une table des matières mettant en avant les titres de section de la page courante.
Sur des écrans peu larges, elle est remplacée par un menu déroulant adhérant.

#### `PageSidebar`

**Composant par défaut :** [`PageSidebar.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageSidebar.astro)

Composant affiché avant le contenu de la page et contenant la table des matières.
L'implémentation par défaut affiche [`<TableOfContents />`](#tableofcontents) et [`<MobileTableOfContents />`](#mobiletableofcontents).

#### `TableOfContents`

**Composant par défaut :** [`TableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/TableOfContents.astro)

Composant qui affiche la table des matières de la page courante sur des écrans suffisamment larges.

#### `MobileTableOfContents`

**Composant par défaut :** [`MobileTableOfContents.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MobileTableOfContents.astro)

Composant qui affiche la table des matières de la page courante sur des petits écrans (mobiles).

---

### Contenu

Ces composants sont utilisés dans la colonne principale de contenu de la page.

#### `Banner`

**Composant par défaut :** [`Banner.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Banner.astro)

Composant représentant une bannière affichée en haut de chaque page.
L'implémentation par défaut utilise la valeur du champ [`banner`](/fr/reference/frontmatter/#banner) du frontmatter de la page pour décider de l'affichage ou non.

#### `ContentPanel`

**Composant par défaut :** [`ContentPanel.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/ContentPanel.astro)

Composant de mise en page utilisé pour envelopper les section de la colonne principale de contenu.

#### `PageTitle`

**Composant par défaut :** [`PageTitle.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/PageTitle.astro)

Composant contenant l'élement `<h1>` de la page courante.

Les implémentations personnalisées doivent s'assurer qu'elles définissent `id="_top"` sur l'élément `<h1>` comme dans l'implémentation par défaut.

#### `DraftContentNotice`

**Composant par défaut :** [`DraftContentNotice.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/DraftContentNotice.astro)

Note affichée aux utilisateurs durant le développement lorsque la page actuelle est marquée comme une ébauche.

#### `FallbackContentNotice`

**Composant par défaut :** [`FallbackContentNotice.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/FallbackContentNotice.astro)

Note affichée aux utilisateurs sur les pages où une traduction pour la langue courante n'est pas disponible.
Utilisé uniquement sur les sites multilingues.

#### `Hero`

**Composant par défaut :** [`Hero.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Hero.astro)

Composant affiché en haut de la page lorsque le champ [`hero`](/fr/reference/frontmatter/#hero) est défini dans le frontmatter.
L'implémentation par défaut affiche un large titre, une accroche et des liens d'appel à l'action à côté d'une image facultative.

#### `MarkdownContent`

**Composant par défaut :** [`MarkdownContent.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/MarkdownContent.astro)

Composant affiché autour du contenu principal de chaque page.
L'implémentation par défaut définit les styles de base à appliquer au contenu Markdown.

Les styles de contenu Markdown sont également exposés dans `@astrojs/starlight/style/markdown.css` avec une portée limitée à la classe CSS `.sl-markdown-content`.

---

### Pied de page

Ces composants sont affichés en bas de la colonne principale de contenu de la page.

#### `Footer`

**Composant par défaut :** [`Footer.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Footer.astro)

Composant pied de page affiché en bas de chaque page.
L'implémentation par défaut affiche [`<LastUpdated />`](#lastupdated), [`<Pagination />`](#pagination) et [`<EditLink />`](#editlink).

#### `LastUpdated`

**Composant par défaut :** [`LastUpdated.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/LastUpdated.astro)

Composant utilisé dans le pied de page pour afficher la date de dernière mise à jour de la page.

#### `EditLink`

**Composant par défaut :** [`EditLink.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/EditLink.astro)

Composant utilisé dans le pied de page pour afficher un lien vers l'emplacement où la page peut être modifiée.

#### `Pagination`

**Composant par défaut :** [`Pagination.astro`](https://github.com/withastro/starlight/blob/main/packages/starlight/components/Pagination.astro)

Composant utilisé dans le pied de page pour afficher des flèches de navigation entre les pages précédentes/suivantes.
