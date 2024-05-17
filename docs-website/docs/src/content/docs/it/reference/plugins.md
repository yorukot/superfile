---
title: Riferimento ai plugin
description: Panoramica dell'API dei plugin di Starlight.
tableOfContents:
  maxHeadingLevel: 4
---

I plugin di Starlight possono personalizzare la configurazione, l'interfaccia utente e il comportamento di Starlight, rendendo allo stesso tempo facile la condivisione e il riutilizzo.
Questa pagina di riferimento documenta le API a cui i plugin hanno accesso.

Per saperne di più sull'utilizzo di un plugin Starlight, consulta il [riferimento alla configurazione](/it/reference/configuration/#plugins) o visita la [vetrina dei plugin](/it/resources/plugins/#plugins) per vedere un elenco dei plugin disponibili.

## Riferimento rapido per le API

Un plugin Starlight ha la seguente forma.
Vedi sotto i dettagli delle diverse proprietà e dei parametri degli hook.

```ts
interface StarlightPlugin {
  name: string;
  hooks: {
    setup: (options: {
      config: StarlightUserConfig;
      updateConfig: (newConfig: StarlightUserConfig) => void;
      addIntegration: (integration: AstroIntegration) => void;
      astroConfig: AstroConfig;
      command: 'dev' | 'build' | 'preview';
      isRestart: boolean;
      logger: AstroIntegrationLogger;
    }) => void | Promise<void>;
  };
}
```

## `name`

**name**: `string`

Un plugin deve fornire un nome univoco che lo descriva. Il nome viene utilizzato durante la [registrazione dei messaggi](#logger) relativi a questo plugin e può essere utilizzato da altri plugin per rilevare la presenza di questo plugin.

## `hooks`

Gli hook sono funzioni che Starlight richiama per eseguire il codice del plugin in momenti specifici. Al momento, Starlight supporta un singolo hook `setup`.

### `hooks.setup`

Funzione di configurazione del plugin chiamata quando Starlight viene inizializzato (durante l'hook di integrazione [`astro:config:setup`](https://docs.astro.build/it/reference/integrations-reference/#astroconfigsetup)).
L'hook `setup` può essere utilizzato per aggiornare la configurazione di Starlight o aggiungere integrazioni Astro.

Questo hook viene chiamato con le seguenti opzioni:

#### `config`

**tipo:** `StarlightUserConfig`

Una copia di sola lettura della [configurazione Starlight](/it/reference/configuration/) fornita all'utente.
Questa configurazione potrebbe essere stata aggiornata da altri plugin configurati prima di quello corrente.

#### `updateConfig`

**tipo:** `(newConfig: StarlightUserConfig) => void`

Una funzione di callback per aggiornare la [configurazione Starlight](/it/reference/configuration/) fornita all'utente.
Fornisci le chiavi di configurazione di livello root che vuoi sovrascrivere.
Per aggiornare valori di configurazione annidati, devi fornire l'intero oggetto annidato.

Per estendere un'opzione di configurazione esistente senza sovrascriverla, unisci il valore esistente nel tuo nuovo valore.
Nel seguente esempio, un nuovo account [`social`](/it/reference/configuration/#social) viene aggiunto alla configurazione esistente unendo `config.social` nel nuovo oggetto `social`:

```ts {6-11}
// plugin.ts
export default {
  name: 'add-twitter-plugin',
  hooks: {
    setup({ config, updateConfig }) {
      updateConfig({
        social: {
          ...config.social,
          twitter: 'https://twitter.com/astrodotbuild',
        },
      });
    },
  },
};
```

#### `addIntegration`

**tipo:** `(integration: AstroIntegration) => void`

Una funzione di callback per aggiungere un'[integrazione Astro](https://docs.astro.build/it/reference/integrations-reference/) richiesta dal plugin.

Nel seguente esempio, il plugin controlla prima se è configurata [l'integrazione React di Astro](https://docs.astro.build/it/guides/integrations-guide/react/) e, se non lo è, utilizza `addIntegration()` per aggiungerla:

```ts {14} "addIntegration,"
// plugin.ts
import react from '@astrojs/react';

export default {
  name: 'plugin-using-react',
  hooks: {
    setup({ addIntegration, astroConfig }) {
      const isReactLoaded = astroConfig.integrations.find(
        ({ name }) => name === '@astrojs/react'
      );

      // Aggiungi l'integrazione di React solo se non è stata già caricata.
      if (!isReactLoaded) {
        addIntegration(react());
      }
    },
  },
};
```

#### `astroConfig`

**tipo:** `AstroConfig`

Una copia di sola lettura della [configurazione Astro](https://docs.astro.build/it/reference/configuration-reference/) fornita all'utente.

#### `command`

**tipo:** `'dev' | 'build' | 'preview'`

Il comando utilizzato per eseguire Starlight:

- `dev` - Il progetto viene eseguito con `astro dev`
- `build` - Il progetto viene eseguito con `astro build`
- `preview` - Il progetto viene eseguito con `astro preview`

#### `isRestart`

**tipo:** `boolean`

`false` quando il server di sviluppo viene avviato, `true` quando viene attivato un riavvio.
Le ragioni comuni per un riavvio includono la modifica di `astro.config.mjs` da parte dell'utente mentre il server di sviluppo è in esecuzione.

#### `logger`

**tipo:** `AstroIntegrationLogger`

Un'istanza del [logger di integrazione Astro](https://docs.astro.build/it/reference/integrations-reference/#astrointegrationlogger) che puoi utilizzare per scrivere log.
Tutti i messaggi registrati saranno prefissati con il nome del plugin.

```ts {6}
// plugin.ts
export default {
  name: 'long-process-plugin',
  hooks: {
    setup({ logger }) {
      logger.info('Inizio di un processo lungo...');
      // Un processo lungo...
    },
  },
};
```

L'esempio sopra registrerà un messaggio che include il messaggio info fornito:

```shell
[long-process-plugin] Inizio di un processo lungo...
```
