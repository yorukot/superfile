---
title: Référence des modules d'extension
description: Une vue d'ensemble de l'API des modules d'extension Starlight.
tableOfContents:
  maxHeadingLevel: 4
---

Les modules d'extension Starlight peuvent personnaliser la configuration, l'interface utilisateur et le comportement de Starlight, tout en étant faciles à partager et à réutiliser.
Cette page de référence documente l'API à laquelle ces modules d'extension ont accès.

Consultez la [référence de configuration](/fr/reference/configuration/#plugins) pour en savoir plus sur l'utilisation d'un module d'extension Starlight ou visitez la [vitrine des modules d'extension](/fr/resources/plugins/) pour voir une liste de modules d'extension disponibles.

## Référence rapide de l'API

Un module d'extension Starlight a la forme suivante.
Voir ci-dessous pour plus de détails sur les différentes propriétés et paramètres des hooks.

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

**Type :** `string`

Un module d'extension doit fournir un nom unique qui le décrit. Le nom est utilisé lors de [l'affichage des messages](#logger) liés à ce module d'extension et peut être utilisé par d'autres modules d'extension pour détecter la présence de ce dernier.

## `hooks`

Les hooks sont des fonctions que Starlight appelle pour exécuter le code du module d'extension à des moments spécifiques. Actuellement, Starlight prend en charge un seul hook `setup`.

### `hooks.setup`

La fonction de configuration du module d'extension appelée lorsque Starlight est initialisé (pendant le hook [`astro:config:setup`](https://docs.astro.build/fr/reference/integrations-reference/#astroconfigsetup) de l'intégration).
Le hook `setup` peut être utilisé pour mettre à jour la configuration de Starlight ou ajouter des intégrations Astro.

Ce hook est appelé avec les options suivantes :

#### `config`

**Type :** `StarlightUserConfig`

Une copie en lecture seule de la [configuration de Starlight](/fr/reference/configuration/) fournie par l'utilisateur.
Cette configuration peut avoir été mise à jour par d'autres modules d'extension configurés avant celui en cours.

#### `updateConfig`

**Type :** `(newConfig: StarlightUserConfig) => void`

Une fonction de rappel pour mettre à jour la [configuration de Starlight](/fr/reference/configuration/) fournie par l'utilisateur.
Spécifiez les clés de configuration de niveau racine que vous souhaitez remplacer.
Pour mettre à jour des valeurs de configuration imbriquées, vous devez fournir l'objet imbriqué entier.

Pour étendre une option de configuration existante sans la remplacer, étendez la valeur existante dans votre nouvelle valeur.
Dans l'exemple suivant, un nouveau compte de média [`social`](/fr/reference/configuration/#social) est ajouté à la configuration existante en étendant `config.social` dans le nouvel objet `social` :

```ts {6-11}
// module-extension.ts
export default {
  name: 'ajout-twitter-plugin',
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

**Type :** `(integration: AstroIntegration) => void`

Une fonction de rappel pour ajouter une [intégration Astro](https://docs.astro.build/fr/reference/integrations-reference/) requise par le module d'extension.

Dans l'exemple suivant, le module d'extension vérifie d'abord si [l'intégration React d'Astro](https://docs.astro.build/fr/guides/integrations-guide/react/) est configurée et, si ce n'est pas le cas, utilise `addIntegration()` pour l'ajouter :

```ts {14} "addIntegration,"
// module-extension.ts
import react from '@astrojs/react';

export default {
  name: 'plugin-utilisant-react',
  hooks: {
    setup({ addIntegration, astroConfig }) {
      const isReactLoaded = astroConfig.integrations.find(
        ({ name }) => name === '@astrojs/react'
      );

      // Ajoute seulement l'intégration React si elle n'est pas déjà chargée.
      if (!isReactLoaded) {
        addIntegration(react());
      }
    },
  },
};
```

#### `astroConfig`

**Type :** `AstroConfig`

Une copie en lecture seule de la [configuration d'Astro](https://docs.astro.build/fr/reference/configuration-reference/) fournie par l'utilisateur.

#### `command`

**Type :** `'dev' | 'build' | 'preview'`

La commande utilisée pour exécuter Starlight :

- `dev` - Le projet est exécuté avec `astro dev`
- `build` - Le projet est exécuté avec `astro build`
- `preview` - Le projet est exécuté avec `astro preview`

#### `isRestart`

**Type :** `boolean`

`false` lorsque le serveur de développement démarre, `true` lorsqu'un rechargement est déclenché.
Les raisons courantes d'un redémarrage incluent un utilisateur qui modifie son fichier `astro.config.mjs` pendant que le serveur de développement est en cours d'exécution.

#### `logger`

**Type :** `AstroIntegrationLogger`

Une instance du [journaliseur (logger) d'intégration Astro](https://docs.astro.build/fr/reference/integrations-reference/#astrointegrationlogger) que vous pouvez utiliser pour écrire des messages de journalisation.
Tous les messages seront préfixés par le nom du module d'extension.

```ts {6}
// module-extension.ts
export default {
  name: 'plugin-long-processus',
  hooks: {
    setup({ logger }) {
      logger.info("Démarrage d'un long processus…");
      // Un long processus…
    },
  },
};
```

L'exemple ci-dessus affichera un message qui inclut le message d'information fourni :

```plaintext frame="terminal"
[plugin-long-processus] Démarrage d'un long processus…
```
