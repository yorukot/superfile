---
title: Riferimenti frontmatter
description: Una panoramica sui campi predefiniti del frontmatter Starlight.
---

Puoi personalizzare pagine Markdown e MDX in Starlight definendo i valori nel frontmatter. Per esempio, una pagina potrebbe definire `title` e `description` :

```md {3-4}
---
# src/content/docs/example.md
title: A proposito del progetto
description: Scopri di più sul progetto a cui sto lavorando.
---

Benvenuto alla pagina "a proposito del progetto"!
```

## Campi del frontmatter

### `title` (obbligatorio)

**tipo:** `string`

Devi fornire un titolo ad ogni pagina. Questo sarà usato in testa alla pagina, nelle finestre del browser e nei metadati della pagina.

### `description`

**tipo:** `string`

La descrizione è utilizzata nei metadati e sarà utilizzata dai motori di ricerca e nelle anteprime nei social.

### `slug`

**tipo:**: `string`

Sovrascrivi lo slug della pagina. Vedi [“Definizione degli slug personalizzati”](https://docs.astro.build/it/guides/content-collections/#defining-custom-slugs) nella documentazione di Astro per ulteriori dettagli.

### `editUrl`

**tipo:** `string | boolean`

Sovrascrive la [configurazione globale `editLink`](/it/reference/configuration/#editlink). Metti a `false` per disabilitare "Modifica la pagina" per quella pagina specifica oppure fornisci un link alternativo.

### `head`

**tipo:** [`HeadConfig[]`](/it/reference/configuration/#headconfig)

Puoi aggiungere tag aggiuntivi nell'`<head>` della pagina utilizzando la chiave `head` nel frontmatter. Questo significa che puoi aggiungere stili personalizzati, metadati o altri tag in una pagina. Il funzionamento è simile [all'opzione globale `head`](/it/reference/configuration/#head).

```md
---
# src/content/docs/example.md
title: Chi siamo
head:
  # Utilizza un <title> personalizzato
  - tag: title
    content: Titolo personalizzato
---
```

### `tableOfContents`

**tipo:** `false | { minHeadingLevel?: number; maxHeadingLevel?: number; }`

Sovrascrive la [configurazione globale `tableOfContents`](/it/reference/configuration/#tableofcontents).
Cambia i livelli di titoli inclusi o, se messo a `false`, nasconde la tabella dei contenuti della pagina.

```md
---
# src/content/docs/example.md
title: Pagina con solo H2 nella tabella dei contenuti della pagina
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 2
---
```

```md
---
# src/content/docs/example.md
title: Pagina senza tabella dei contenuti della pagina
tableOfContents: false
---
```

### `template`

**tipo:** `'doc' | 'splash'`  
**predefinito:** `'doc'`

Definisce il layout per la pagina.
Le pagine utilizzano `'doc'` come predefinita.
Se valorizzato a `'splash'` viene utilizzato un layout senza barre laterali ottimale per la pagina iniziale.

### `hero`

**tipo:** [`HeroConfig`](#heroconfig)

Aggiunge un componente hero all'inizio della pagina. Funziona bene con `template: splash`.

Per esempio, questa configurazione illustra comuni opzioni, incluso il caricamento di un'immagine.

```md
---
# src/content/docs/example.md
title: La mia pagina principale
template: splash
hero:
  title: 'Il mio progetto: Stellar Stuff Sooner'
  tagline: Porta le tue cose sulla Luna e torna indietro in un battito d'occhio.
  image:
    alt: Un logo brillante e luminoso
    file: ../../assets/logo.png
  actions:
    - text: Dimmi di più
      link: /getting-started/
      icon: right-arrow
      variant: primary
    - text: Vedi su GitHub
      link: https://github.com/astronaut/my-project
      icon: external
      attrs:
        rel: me
---
```

Puoi mostrare diverse versioni dell'immagine in base alla modalità chiara o scura.

```md
---
# src/content/docs/example.md
hero:
  image:
    alt: Un logo brillante e luminoso
    dark: ~/assets/logo-dark.png
    light: ~/assets/logo-light.png
---
```

#### `HeroConfig`

```ts
interface HeroConfig {
  title?: string;
  tagline?: string;
  image?:
    | {
        // Percorso relativo a un'immagine nella tua repository.
        file: string;
        // Testo alternativo per rendere l'immagine accessibile alla tecnologia assistiva
        alt?: string;
      }
    | {
        // Percorso relativo a un'immagine nella tua repository da utilizzare per la modalità scura.
        dark: string;
        // Percorso relativo a un'immagine nella tua repository da utilizzare per la modalità chiara.
        light: string;
        // Testo alternativo per rendere l'immagine accessibile alla tecnologia assistiva
        alt?: string;
      }
    | {
        // HTML grezzo da utilizzare nello slot dell'immagine.
        // Potrebbe essere un tag `<img>` personalizzato o `<svg>` inline.
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

Visualizza un banner di annuncio nella parte superiore di questa pagina.

Il valore `content` può includere HTML per collegamenti o altri contenuti.
Ad esempio, questa pagina visualizza un banner che include un collegamento a `example.com`.

```md
---
# src/content/docs/example.md
title: Pagina con un banner
banner:
  content: |
    Abbiamo appena lanciato qualcosa di interessante!
    <a href="https://example.com">Dai un'occhiata</a>
---
```

### `lastUpdated`

**tipo:** `Date | boolean`

Sostituisce l'[opzione globale `lastUpdated`](/it/reference/configuration/#lastupdated). Se viene specificata una data, deve essere un [timestamp YAML](https://yaml.org/type/timestamp.html) valido e sovrascriverà la data archiviata nella cronologia Git per questa pagina.

```md
---
# src/content/docs/example.md
title: Pagina con una data di ultimo aggiornamento personalizzata
lastUpdated: 2022-08-09
---
```

### `prev`

**tipo:** `boolean | string | { link?: string; label?: string }`

Sostituisce l'[opzione globale `paginazione`](/it/reference/configuration/#pagination). Se viene specificata una stringa, il testo del collegamento generato verrà sostituito e se viene specificato un oggetto, sia il collegamento che il testo verranno sovrascritti.

```md
---
# src/content/docs/example.md
# Nascondi il collegamento alla pagina precedente
prev: false
---
```

```md
---
# src/content/docs/example.md
# Sostituisci il testo del collegamento della pagina precedente
prev: Continua il tutorial
---
```

```md
---
# src/content/docs/example.md
# Sostituisci sia il collegamento che il testo della pagina precedente
prev:
  link: /pagina-non-correlata/
  label: Dai un'occhiata a quest'altra pagina
---
```

### `next`

**tipo:** `boolean | string | { link?: string; label?: string }`

Uguale a [`prev`](#prev) ma per il collegamento alla pagina successiva.

```md
---
# src/content/docs/example.md
# Nascondi il collegamento alla pagina successiva
next: false
---
```

### `pagefind`

**tipo:** `boolean`  
**predefinito:** `true`

Imposta se questa pagina deve essere inclusa nell'indice di ricerca [Pagefind](https://pagefind.app/). Imposta su `false` per escludere una pagina dai risultati di ricerca:

```md
---
# src/content/docs/example.md
# Nascondi questa pagina dai risultati di ricerca
pagefind: false
---
```

### `sidebar`

**tipo:** [`SidebarConfig`](#sidebarconfig)

Controlla il modo in cui questa pagina viene visualizzata nella [barra laterale](/it/reference/configuration/#sidebar), quando si utilizza un gruppo di collegamenti generato automaticamente.

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
**predefinito:** la pagina [`title`](#title-obbligatorio)

Imposta l'etichetta per questa pagina nella barra laterale quando viene visualizzata in un gruppo di collegamenti generato automaticamente.

```md
---
# src/content/docs/example.md
title: Informazioni su questo progetto
sidebar:
  label: Informazioni
---
```

#### `order`

**tipo:** `number`

Controlla l'ordine di questa pagina quando ordini un gruppo di collegamenti generato automaticamente.
I numeri più bassi vengono visualizzati più in alto nel gruppo di collegamenti.

```md
---
# src/content/docs/example.md
title: Pagina da visualizzare per prima
sidebar:
  order: 1
---
```

#### `hidden`

**tipo:** `boolean`  
**predefinito:** `false`

Impedisce che questa pagina venga inclusa in un gruppo della barra laterale generato automaticamente.

```md
---
# src/content/docs/example.md
title: Pagina da nascondere dalla barra laterale generata automaticamente
sidebar:
  hidden: true
---
```

#### `badge`

**tipo:** <code>string | <a href="/it/reference/configuration/#badgeconfig">BadgeConfig</a></code>

Aggiungi un badge alla pagina nella barra laterale quando viene visualizzata in un gruppo di collegamenti generato automaticamente.
Quando si utilizza una stringa, il badge verrà visualizzato con un colore in risalto predefinito.
Facoltativamente, passa un [oggetto `BadgeConfig`](/it/reference/configuration/#badgeconfig) con i campi `text` e `variant` per personalizzare il badge.

```md
---
# src/content/docs/example.md
title: Pagina con un badge
sidebar:
  # Utilizza la variante predefinita corrispondente al colore principale del tuo sito
  badge: nuovo
---
```

```md
---
# src/content/docs/example.md
title: Pagina con un badge
sidebar:
  badge:
    text: Sperimentale
    variant: caution
---
```

#### `attrs`

**tipo:** `Record<string, string | number | boolean | undefined>`

Attributi HTML da aggiungere al collegamento della pagina nella barra laterale quando viene visualizzato in un gruppo di collegamenti generato automaticamente.

```md
---
# src/content/docs/example.md
title: Pagina che si aprirà in una nuova scheda
sidebar:
  # Apre la pagina in una nuova scheda
  attrs:
    target: _blank
---
```

## Personalizza lo schema del frontmatter

Lo schema del frontmatter per la raccolta di contenuti `docs` di Starlight è configurato in `src/content/config.ts` utilizzando l'helper `docsSchema()`:

```ts {3,6}
// src/content/config.ts
import { defineCollection } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({ schema: docsSchema() }),
};
```

Per saperne di più sugli schemi di raccolta dei contenuti, consulta [“Definizione di uno schema di raccolta”](https://docs.astro.build/it/guides/content-collections/#defining-a-collection-schema) nella documentazione di Astro.

`docsSchema()` accetta le seguenti opzioni:

### `extend`

**tipo:** Schema Zod o funzione che restituisce uno schema Zod  
**predefinito:** `z.object({})`

Estendi lo schema di Starlight con campi aggiuntivi impostando `extend` nelle opzioni di `docsSchema()`.
Il valore dovrebbe essere uno [schema Zod](https://docs.astro.build/it/guides/content-collections/#defining-datatypes-with-zod).

Nell'esempio seguente, forniamo un tipo più restrittivo per `description` per renderlo obbligatorio e aggiungiamo un nuovo campo opzionale `category`:

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: z.object({
        // Rendi un campo built-in obbligatorio anziché opzionale.
        description: z.string(),
        // Aggiungi un nuovo campo allo schema.
        category: z.enum(['tutorial', 'guide', 'reference']).optional(),
      }),
    }),
  }),
};
```

Per sfruttare l'[helper `image()` di Astro,](https://docs.astro.build/it/guides/images/#images-in-content-collections), utilizza una funzione che restituisce l'estensione dello schema:

```ts {8-13}
// src/content/config.ts
import { defineCollection, z } from 'astro:content';
import { docsSchema } from '@astrojs/starlight/schema';

export const collections = {
  docs: defineCollection({
    schema: docsSchema({
      extend: ({ image }) => {
        return z.object({
          // Aggiungi un campo che deve risolvere in un'immagine locale.
          cover: image(),
        });
      },
    }),
  }),
};
```
