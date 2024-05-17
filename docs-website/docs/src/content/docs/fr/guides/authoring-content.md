---
title: Création de contenu en Markdown
description: Un aperçu de la syntaxe Markdown prise en charge par Starlight.
---

Starlight prend en charge l'ensemble de la syntaxe [Markdown](https://daringfireball.net/projects/markdown/) dans les fichiers `.md` ainsi que la syntaxe frontale [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f) pour définir des métadonnées telles qu'un titre et une description.

Veillez à consulter les [MDX docs](https://mdxjs.com/docs/what-is-mdx/#markdown) ou les [Markdoc docs](https://markdoc.dev/docs/syntax) si vous utilisez ces formats de fichiers, car la prise en charge et l'utilisation de Markdown peuvent varier.

## Frontmatter

Vous pouvez personnaliser chaque page individuellement en définissant des valeurs dans leur frontmatter.
Le frontmatter se situe en haut de vos fichiers entre les séparateurs `---` :

```md title="src/content/docs/exemple.md"
---
title: Le titre de ma page
---

Le contenu de la page suit le second `---`.
```

Chaque page doit inclure au moins un titre (`title`).
Consultez la [référence du frontmatter](/fr/reference/frontmatter/) pour connaître tous les champs disponibles et comment ajouter des champs personnalisés.

## Styles en ligne

Le texte peut être **gras**, _italique_, ou ~~barré~~.

```md
Le texte peut être **gras**, _italique_, ou ~~barré~~.
```

Vous pouvez [faire un lien vers une autre page](/fr/getting-started/).

```md
Vous pouvez [faire un lien vers une autre page](/fr/getting-started/).
```

Vous pouvez mettre en évidence le `code en ligne` à l'aide d'un astérisque.

```md
Vous pouvez mettre en évidence le `code en ligne` à l'aide de barres de défilement.
```

## Images

Les images dans Starlight utilisent [la prise en charge intégrée des ressources optimisées d'Astro](https://docs.astro.build/en/guides/assets/).

Markdown et MDX supportent la syntaxe Markdown pour l'affichage des images qui inclut le texte alt pour les lecteurs d'écran et les technologies d'assistance.

![Une illustration de planètes et d'étoiles avec le mot "astro"](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
![Une illustration de planètes et d'étoiles avec le mot "astro"](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

Les chemins d'accès relatifs aux images sont également supportés pour les images stockées localement dans votre projet.

```md
// src/content/docs/page-1.md

![Une fusée dans l'espace](../../assets/images/rocket.svg)
```

## En-têtes

Vous pouvez structurer le contenu à l'aide d'un titre. En Markdown, les titres sont indiqués par un nombre de `#` en début de ligne.

### Comment structurer le contenu d'une page dans Starlight

Starlight est configuré pour utiliser automatiquement le titre de votre page comme titre de premier niveau et inclura un titre "Aperçu" en haut de la table des matières de chaque page. Nous vous recommandons de commencer chaque page par un paragraphe de texte normal et d'utiliser des titres de page à partir de `<h2>` :

```md
---
title: Guide Markdown
description: Comment utiliser Markdown dans Starlight
---

Cette page décrit comment utiliser Markdown dans Starlight.

## Styles en ligne

## Titres
```

### Liens d'ancrage automatiques pour les titres

L'utilisation de titres en Markdown vous donnera automatiquement des liens d'ancrage afin que vous puissiez accéder directement à certaines sections de votre page :

```md
---
title: Ma page de contenu
description: Comment utiliser les liens d'ancrage intégrés de Starlight
---

## Introduction

Je peux faire un lien vers [ma conclusion](#conclusion) plus bas sur la même page.

## Conclusion

`https://my-site.com/page1/#introduction` renvoie directement à mon Introduction.
```

Les titres de niveau 2 (`<h2>`) et de niveau 3 (`<h3>`) apparaissent automatiquement dans la table des matières de la page.

Pour en apprendre davantage sur la façon dont Astro traite les attributs `id` des titres de section, consultez la [documentation d'Astro](https://docs.astro.build/fr/guides/markdown-content/#identifiants-den-t%C3%AAte).

## Encarts

Les encarts (également connus sous le nom de « admonitions » ou « asides » en anglais) sont utiles pour afficher des informations secondaires à côté du contenu principal d'une page.

Starlight fournit une syntaxe Markdown personnalisée pour le rendu des encarts. Les blocs d'encarts sont indiqués en utilisant une paire de triples points `:::` pour envelopper votre contenu, et peuvent être de type `note`, `tip`, `caution` ou `danger`.

Vous pouvez imbriquer n'importe quel autre type de contenu Markdown à l'intérieur d'un aparté, mais les aparté sont mieux adaptés à des morceaux de contenu courts et concis.

### Encart de type note

:::note
Starlight est une boîte à outils pour sites web de documentation construite avec [Astro](https://astro.build/). Vous pouvez démarrer avec cette commande :

```sh
npm run create astro@latest --template starlight
```

:::

````md
:::note
Starlight est une boîte à outils pour sites web de documentation construite avec [Astro](https://astro.build/). Vous pouvez démarrer avec cette commande :

```sh
npm run create astro@latest --template starlight
```

:::
````

### Titres personnalisés dans les encarts

Vous pouvez spécifier un titre personnalisé pour l'encart entre crochets après le type d'encarts, par exemple `:::tip[Le saviez-vous ?]`.

:::tip[Le saviez-vous ?]
Astro vous aide à construire des sites Web plus rapides grâce à ["Islands Architecture"](https://docs.astro.build/fr/concepts/islands/).
:::

```md
:::tip[Le saviez-vous ?]
Astro vous aide à construire des sites Web plus rapides grâce à ["Islands Architecture"](https://docs.astro.build/fr/concepts/islands/).
:::
```

### Plus de types d'encarts

Les encarts de type Attention et Danger sont utiles pour attirer l'attention de l'utilisateur sur des détails qui pourraient le perturber. Si vous vous retrouvez à utiliser ces derniers fréquemment, cela pourrait aussi être un signe que ce que vous documentez pourrait bénéficier d'une refonte.

:::caution
Si vous n'êtes pas sûr de vouloir un site de documentation génial, réfléchissez à deux fois avant d'utiliser [Starlight](/fr/).
:::

:::danger
Vos utilisateurs peuvent être plus productifs et trouver votre produit plus facile à utiliser grâce aux fonctionnalités utiles de Starlight.

- Navigation claire
- Thème de couleurs configurable par l'utilisateur
- [Support i18n](/fr/guides/i18n/)

:::

```md
:::caution
Si vous n'êtes pas sûr de vouloir un site de documentation génial, réfléchissez à deux fois avant d'utiliser [Starlight](/fr/).
:::

:::danger
Vos utilisateurs peuvent être plus productifs et trouver votre produit plus facile à utiliser grâce aux fonctionnalités utiles de Starlight.

- Navigation claire
- Thème de couleurs configurable par l'utilisateur
- [Support i18n](/fr/guides/i18n/)

:::
```

## Blockquotes

> Il s'agit d'une citation en bloc, couramment utilisée pour citer une autre personne ou un document.
>
> Les guillemets sont indiqués par un `>` au début de chaque ligne.

```md
> Il s'agit d'une citation en bloc, couramment utilisée pour citer une autre personne ou un document.
>
> Les guillemets sont indiqués par un `>` au début de chaque ligne.
```

## Blocs de code

Un bloc de code est indiqué par un bloc avec trois accents graves <code>```</code> au début et à la fin. Vous pouvez indiquer le langage de programmation utilisé après les premiers accents graves.

```js
// Code Javascript avec coloration syntaxique.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// Code Javascript avec coloration syntaxique.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

### Fonctionnalités d'Expressive Code

Starlight utilise [Expressive Code](https://github.com/expressive-code/expressive-code/tree/main/packages/astro-expressive-code) pour étendre les possibilités de formatage des blocs de code.
Les plugins Expressive Code de marqueurs de texte et de cadres de fenêtre sont activés par défaut.
L'affichage des blocs de code peut être configuré à l'aide de [l'option de configuration `expressiveCode`](/fr/reference/configuration/#expressivecode) de Starlight.

#### Marqueurs de texte

Vous pouvez mettre en évidence des lignes ou des portions spécifiques de vos blocs de code à l'aide des [marqueurs de texte d'Expressive Code](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#usage-in-markdown--mdx-documents) sur la première ligne de votre bloc de code.
Utilisez des accolades (`{ }`) pour mettre en évidence des lignes entières, et des guillemets pour mettre en évidence des chaînes de texte.

Il existe trois styles de mise en évidence : neutre pour attirer l'attention sur le code, vert pour indiquer du code inséré, et rouge pour indiquer du code supprimé.
Du texte et des lignes entières peuvent être marqués à l'aide du marqueur par défaut, ou en combinaison avec `ins=` et `del=` pour produire la mise en évidence souhaitée.

Expressive Code fournit plusieurs options pour personnaliser l'apparence visuelle de vos exemples de code.
Beaucoup d'entre elles peuvent être combinées pour obtenir des exemples de code très illustratifs.
Merci d'explorer la [documentation d'Expressive Code](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md) pour obtenir une liste complète des options disponibles.
Certaines des options les plus courantes sont présentées ci-dessous :

- [Marquer des lignes entières et des plages de lignes à l'aide du marqueur `{ }`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-entire-lines--line-ranges) :

  ```js {2-3}
  function demo() {
    // Cette ligne (#2) et la suivante sont mises en évidence
    return 'Ceci est la ligne #3 de cet exemple';
  }
  ```

  ````md
  ```js {2-3}
  function demo() {
    // Cette ligne (#2) et la suivante sont mises en évidence
    return 'Ceci est la ligne #3 de cet exemple';
  }
  ```
  ````

- [Marquer des sélections de texte à l'aide du marqueur `" "` ou d'expressions régulières](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-individual-text-inside-lines) :

  ```js "termes individuels" /Même.*charge/
  // Des termes individuels peuvent également être mis en évidence
  function demo() {
    return 'Même les expressions régulières sont prises en charge';
  }
  ```

  ````md
  ```js "termes individuels" /Même.*charge/
  // Des termes individuels peuvent également être mis en évidence
  function demo() {
    return 'Même les expressions régulières sont prises en charge';
  }
  ```
  ````

- [Marquer du texte ou des lignes comme insérés ou supprimés avec `ins` ou `del`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#selecting-marker-types-mark-ins-del) :

  ```js "return true;" ins="insertion" del="suppression"
  function demo() {
    console.log("Voici des marqueurs d'insertion et de suppression");
    // La déclaration return utilise le type de marqueur par défaut
    return true;
  }
  ```

  ````md
  ```js "return true;" ins="insertion" del="suppression"
  function demo() {
    console.log("Voici des marqueurs d'insertion et de suppression");
    // La déclaration return utilise le type de marqueur par défaut
    return true;
  }
  ```
  ````

- [Combiner coloration syntaxique et syntaxe de type `diff`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#combining-syntax-highlighting-with-diff-like-syntax) :

  ```diff lang="js"
    function ceciEstDuJavaScript() {
      // Ce bloc entier utilise la coloration syntaxique JavaScript,
      // et nous pouvons toujours y ajouter des marqueurs de différence !
  -   console.log('Ancien code à supprimer')
  +   console.log('Nouveau code brillant !')
    }
  ```

  ````md
  ```diff lang="js"
    function ceciEstDuJavaScript() {
      // Ce bloc entier utilise la coloration syntaxique JavaScript,
      // et nous pouvons toujours y ajouter des marqueurs de différence !
  -   console.log('Ancien code à supprimer')
  +   console.log('Nouveau code brillant !')
    }
  ```
  ````

#### Cadres et titres

Les blocs de code peuvent être affichés dans un cadre ressemblant à une fenêtre.
Un cadre ressemblant à une fenêtre de terminal sera utilisé pour les langages de script shell (par exemple `bash` ou `sh`).
Les autres langages s'affichent dans un cadre de style éditeur de code s'ils incluent un titre.

Le titre optionnel d'un bloc de code peut être défini soit avec un attribut `title="..."` après les accents graves d'ouverture et l'identifiant de langage, ou avec un nom de fichier en commentaire sur la première ligne du bloc de code.

- [Ajouter un nom de fichier avec un commentaire](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title) :

  ```js
  // mon-fichier-de-test.js
  console.log('Hello World!');
  ```

  ````md
  ```js
  // mon-fichier-de-test.js
  console.log('Hello World!');
  ```
  ````

- [Ajouer un title à une fenêtre de terminal](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title) :

  ```bash title="Installation des dépendances…"
  npm install
  ```

  ````md
  ```bash title="Installation des dépendances…"
  npm install
  ```
  ````

- [Désactiver les cadres de fenêtre avec `frame="none"`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#overriding-frame-types) :

  ```bash frame="none"
  echo "Ceci n'est pas affiché comme un terminal malgré l'utilisation du langage bash"
  ```

  ````md
  ```bash frame="none"
  echo "Ceci n'est pas affiché comme un terminal malgré l'utilisation du langage bash"
  ```
  ````

## Autres fonctionnalités courantes de Markdown

Starlight prend en charge toutes les autres syntaxes de rédaction Markdown, telles que les listes et les tableaux. Voir [Markdown Cheat Sheet from The Markdown Guide](https://www.markdownguide.org/cheat-sheet/) pour un aperçu rapide de tous les éléments de la syntaxe Markdown.

## Configuration avancée de Markdown et MDX

Starlight utilise le moteur de rendu Markdown et MDX d'Astro basé sur remark et rehype. Vous pouvez ajouter la prise en charge de syntaxe et comportement personnalisés en ajoutant `remarkPlugins` ou `rehypePlugins` dans votre fichier de configuration Astro. Pour en savoir plus, consultez [« Configuration de Markdown et MDX »](https://docs.astro.build/fr/guides/markdown-content/#configuration-de-markdown-et-mdx) dans la documentation d'Astro.
