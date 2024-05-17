---
title: Creazione di contenuti in Markdown
description: Una panoramica della sintassi Markdown supportata da Starlight.
---

Starlight supporta l'intera sintassi [Markdown](https://daringfireball.net/projects/markdown/) nei file `.md` insieme al frontmatter [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f) per definire metadati come il titolo e la descrizione.

Assicurarsi di guardare la [documentazione MDX](https://mdxjs.com/docs/what-is-mdx/#markdown) o la [documentazione Markdoc](https://markdoc.dev/docs/syntax) se si vogliono usare questi formati, dato che il supporto Markdown può variare.

## Frontmatter

È possibile personalizzare le singole pagine in Starlight impostando i valori nel loro frontmatter.
Il frontmatter è impostato nella parte superiore dei file tra i separatori `---`:

```md title="src/content/docs/example.md"
---
title: Titolo della mia pagina
---

Il contenuto della pagina segue il secondo `---`.
```

Ogni pagina deve includere almeno un `title`.
Consulta il [riferimento al frontmatter](/it/reference/frontmatter/) per tutti i campi disponibili e come aggiungere campi personalizzati.

## Stili in linea

Il testo può essere **grassetto**, _corsivo_, o ~~barrato~~.

```md
Il testo può essere **grassetto**, _corsivo_, o ~~barrato~~.
```

Puoi [aggiungere un link ad un'altra pagina](/it/getting-started/).

```md
Puoi [aggiungere un link ad un'altra pagina](/it/getting-started/).
```

Puoi evidenziare `codice in linea` con apici inversi.

```md
Puoi evidenziare `codice in linea` con apici inversi.
```

## Immagini

Le immagini in Starlight utilizzano [l'ottimizzazione degli asset di Astro](https://docs.astro.build/it/guides/assets/).

Markdown e MDX supportano la sintassi Markdown per rappresentare immagini che includono testo alternativo per le tecnologie assistive.

![Un'illustrazione di pianeti e stelle con la scritta "astro"](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
![Un'illustrazione di pianeti e stelle con la scritta "astro"](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

I percorsi relativi sono supportati per immagini salvate localmente nel tuo progetto.

```md
// src/content/docs/page-1.md

![Un'astronave nello spazio](../../assets/images/rocket.svg)
```

## Titoli

Puoi strutturare i contenuti utilizzando dei titoli. In Markdown sono indicati dal numero di `#` all'inizio della linea.

### Come strutturare i contenuti della pagina in Starlight

Starlight è configurato per utilizzare automaticamente il titolo della pagina come intestazione e includerà una "Panoramica" in alto per ogni tabella dei contenuti. Si raccomanda di iniziare ogni pagina con un paragrafo e di usare titoli a partire da `<h2>`:

```md
---
title: Guida Markdown
description: Come utilizzare Markdown in Starlight
---

Questa pagina descrive come utilizzare Markdown in Starlight.

## Stili in linea

## Titoli
```

### Link titoli automatici

Utilizzando titoli in Markdown verranno generati automaticamente i rispettivi link per navigare velocemente in certe sezioni della tua pagina:

```md
---
title: La mia pagina dei contenuti
description: Come utilizzare i link automatici di Starlight
---

## Introduzione

Posso collegarmi alla [mia conclusione](#conclusione) che si trova più in basso.

## Conclusione

`https://my-site.com/page1/#introduzione` porta direttamente all'introduzione.
```

Titoli di livello 2 (`<h2>`) e di livello 3 (`<h3>`) verranno inclusi automaticamente nella tabella dei contenuti.

Scopri come Astro processa gli `id` di heading nella [documentazione di Astro](https://docs.astro.build/it/guides/markdown-content/#heading-ids)

## Avvisi

Gli avvisi sono utili per indicare contenuti secondari insieme ai contenuti principali.

Starlight fornisce una sintassi Markdown personalizzata per indicarli. Gli avvisi sono indicati da `:::` per racchiudere i contenuti e possono essere di tipo `note`, `tip`, `caution` o `danger`.

Dentro un avviso puoi inserire qualsiasi altro contenuto Markdown anche se sono più indicati per contenere poche informazioni.

### Avviso come nota

:::note
Starlight è uno strumento per siti da documentazione con [Astro](https://astro.build/). Puoi iniziare con questo comando:

```sh
npm create astro@latest -- --template starlight
```

:::

````md
:::note
Starlight è uno strumento per siti da documentazione con [Astro](https://astro.build/). Puoi iniziare con questo comando:

```sh
npm create astro@latest -- --template starlight
```

:::
````

### Avvisi con titoli personalizzati

Si può specificare un titolo personalizzato per gli avvisi in parentesi quadre dopo aver specificato il tipo di avviso, per esempio `:::tip[Lo sapevi?]`.

:::tip[Lo sapevi?]
Astro ti aiuta a costruire siti più veloci con ["Islands Architecture"](https://docs.astro.build/it/concepts/islands/).
:::

```md
:::tip[Lo sapevi?]
Astro ti aiuta a costruire siti più veloci con ["Islands Architecture"](https://docs.astro.build/it/concepts/islands/).
:::
```

### Altri tipi di avvisi

Gli avvisi `caution` e `danger` sono d'aiuto per richiamare l'attenzione dell'utente a dettagli che potrebbero sorprenderli.
Se ti ritrovi ad usarli spesso, potrebbe essere segno che quelo che stai documentando potrebbe trarre beneficio da una riprogettazione.

:::caution
Se non sei sicuro di voler un sito per documentazione fantastico, pensaci due volte prima di usare [Starlight](/it/).
:::

:::danger
Gli utenti potrebbero essere più produttivi e trovare il tuo prodotto più facile da usare grazie alle utili funzioni di Starlight.

- Navigazione chiara
- Temi configurabili dall'utente
- [Supporto per i18n](/it/guides/i18n/)

:::

```md
:::caution
Se non sei sicuro di voler un sito per documentazione fantastico, pensaci due volte prima di usare [Starlight](/it/).
:::

:::danger
Gli utenti potrebbero essere più produttivi e trovare il tuo prodotto più facile da usare grazie alle utili funzioni di Starlight.

- Navigazione chiara
- Temi configurabili dall'utente
- [Supporto per i18n](/it/guides/i18n/)

:::
```

## Citazioni

> Questo è un blockquote, che di solito viene utilizzato per citazioni di persone o documenti.
>
> I blockquote sono indicati da `>` all'inizio di ogni riga.

```md
> Questo è un blockquote, che di solito viene utilizzato per citazioni di persone o documenti.
>
> I blockquote sono indicati da `>` all'inizio di ogni riga.
```

## Blocco di codice

Un blocco di codice è indicato da tre backtick <code>```</code> all'inizio e alla fine. Puoi indicare il linguaggio di programmazione dopo i primi backtick.

```js
// Codice JavaScript con sintassi evidenziata.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// Codice JavaScript con sintassi evidenziata.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

### Funzionalità di Expressive Code

Starlight utilizza [Expressive Code](https://github.com/expressive-code/expressive-code/tree/main/packages/astro-expressive-code) per estendere le possibilità di formattazione dei blocchi di codice.
I plugin di marcatori di testo e cornici per finestre di Expressive Code sono abilitati per impostazione predefinita.
La resa dei blocchi di codice può essere configurata utilizzando [l'opzione di configurazione `expressiveCode`](/it/reference/configuration/#expressivecode) di Starlight.

#### Marcatori di testo

È possibile evidenziare linee specifiche o parti dei blocchi di codice utilizzando i [marcatori di testo di Expressive Code](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#usage-in-markdown--mdx-documents) sulla riga di apertura del blocco di codice.
Utilizzare le parentesi graffe (`{ }`) per evidenziare intere linee e virgolette per evidenziare stringhe di testo.

Ci sono tre stili di evidenziazione: neutro per attirare l'attenzione sul codice, verde per indicare il codice inserito e rosso per indicare il codice eliminato.
Sia il testo che intere linee possono essere contrassegnati utilizzando il marcatore predefinito, o in combinazione con `ins=` e `del=` per produrre l'evidenziazione desiderata.

Expressive Code fornisce diverse opzioni per personalizzare l'aspetto visivo dei tuoi esempi di codice.
Molte di queste possono essere combinate, per esempi di codice altamente illustrativi.
Si prega di esplorare la [documentazione di Expressive Code](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md) per le numerose opzioni disponibili.
Di seguito sono mostrati alcuni degli esempi più comuni:

- [Contrassegnare linee intere e intervalli di linee utilizzando il marcatore `{ }`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-entire-lines--line-ranges):

  ```js {2-3}
  function demo() {
    // Questa linea (#2) e la successiva sono evidenziate
    return 'Questa è la linea #3 di questo snippet';
  }
  ```

  ````md
  ```js {2-3}
  function demo() {
    // Questa linea (#2) e la successiva sono evidenziate
    return 'Questa è la linea #3 di questo snippet';
  }
  ```
  ````

- [Contrassegnare selezioni di testo utilizzando il marcatore `" "` o espressioni regolari](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-individual-text-inside-lines):

  ```js "termini individuali" /Anche.*supportate/
  // Anche i termini individuali possono essere evidenziati
  function demo() {
    return 'Anche le espressioni regolari sono supportate';
  }
  ```

  ````md
  ```js "termini individuali" /Anche.*supportate/
  // Anche i termini individuali possono essere evidenziati
  function demo() {
    return 'Anche le espressioni regolari sono supportate';
  }
  ```
  ````

- [Contrassegnare testi o linee come inseriti o eliminati con `ins` o `del`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#selecting-marker-types-mark-ins-del):

  ```js "return true;" ins="inseriti" del="eliminati"
  function demo() {
    console.log('Questi sono tipi di marcatore inseriti ed eliminati');
    // La dichiarazione di ritorno utilizza il tipo di marcatore predefinito
    return true;
  }
  ```

  ````md
  ```js "return true;" ins="inseriti" del="eliminati"
  function demo() {
    console.log('Questi sono tipi di marcatore inseriti ed eliminati');
    // La dichiarazione di ritorno utilizza il tipo di marcatore predefinito
    return true;
  }
  ```
  ````

- [Combina l'evidenziazione della sintassi con una sintassi simile a `diff`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#combining-syntax-highlighting-with-diff-like-syntax):

  ```diff lang="js"
    function thisIsJavaScript() {
      // Questo intero blocco viene evidenziato come JavaScript,
      // e possiamo comunque aggiungere marcatori diff ad esso!
  -   console.log('Vecchio codice da rimuovere')
  +   console.log('Nuovo e splendido codice!')
    }
  ```

  ````md
  ```diff lang="js"
    function thisIsJavaScript() {
      // Questo intero blocco viene evidenziato come JavaScript,
      // e possiamo comunque aggiungere marcatori diff ad esso!
  -   console.log('Vecchio codice da rimuovere')
  +   console.log('Nuovo e splendido codice!')
    }
  ```
  ````

#### Frame e titoli

I blocchi di codice possono essere visualizzati all'interno di un frame simile a una finestra.
Un frame che assomiglia a una finestra del terminale verrà utilizzato per i linguaggi di scripting della shell (ad esempio `bash` o `sh`).
Altri linguaggi vengono visualizzati all'interno di un frame simile a un editor di codice se includono un titolo.

Il titolo opzionale di un blocco di codice può essere impostato sia con un attributo `title="..."` che segue l'identificatore del linguaggio e le backticks di apertura del blocco di codice, sia con un commento del nome del file nelle prime righe del codice.

- [Aggiungi una scheda con il nome del file con un commento](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```js
  // my-test-file.js
  console.log('Ciao mondo!');
  ```

  ````md
  ```js
  // my-test-file.js
  console.log('Ciao mondo!');
  ```
  ````

- [Aggiungi un titolo a una finestra del terminale](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```bash title="Installando dipendenze..."
  npm install
  ```

  ````md
  ```bash title="Installando dipendenze..."
  npm install
  ```
  ````

- [Disabilita i frame delle finestre con `frame="none"`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#overriding-frame-types)

  ```bash frame="none"
  echo "Questo non verrà mostrato come un terminale nonostante stia usando bash come linguaggio"
  ```

  ````md
  ```bash frame="none"
  echo "Questo non verrà mostrato come un terminale nonostante stia usando bash come linguaggio"
  ```
  ````

## Altre funzionalità Markdown utili

Starlight supporta tutte le altre funzionalità Markdown, come liste e tabelle. Guarda la [Markdown Cheat Sheet da The Markdown Guide](https://www.markdownguide.org/cheat-sheet/) per una panoramica veloce su tutte le funzionalità Markdown.

## Configurazione avanzata di Markdown e MDX

Starlight utilizza il renderer Markdown e MDX di Astro costruito su remark e rehype. Puoi aggiungere supporto per la sintassi e il comportamento personalizzati aggiungendo `remarkPlugins` o `rehypePlugins` nel file di configurazione di Astro. Vedi [“Configuring Markdown and MDX”](https://docs.astro.build/it/guides/markdown-content/#configuring-markdown-and-mdx) nella documentazione di Astro per saperne di più.
