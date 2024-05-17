---
title: Referência de plugins
description: Uma visão geral da API de plugins do Starlight
tableOfContents:
  maxHeadingLevel: 4
---

Os plugins do Starlight podem personalizar a configuração, UI, e comportamento, além de serem fáceis de compartilhar e reutilizar.
Essa página de referência documenta a API que os plugins tem acesso.

Aprenda mais sobre como usar um plugin do Starlight na [Referência da Configuração](/pt-br/reference/configuration/#plugins) ou visite o [mostruário de plugins](/pt-br/resources/plugins/#plugins) para ver uma lista de plugins disponíveis.

## Referência rápida da API

Um plugin do Starlight segue o seguinte formato.
Veja abaixo os detalhes das diferentes propriedades e parâmetros de hooks.

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

**tipo:** `string`

Um plugin deve fornecer um nome único que o descreve. O nome é usado nas [mensagens de log](#logger) relacionadas ao plugin e pode ser utilizado por outros plugins para detectar a sua presença.

## `hooks`

Hooks são funções que o Starlight chama para executar um código do plugin em momentos específicos. Atualmente, o Starlight suporta apenas o hook `setup`.

### `hooks.setup`

Função de setup chamada quando o Starlight é iniciado (durante a execução do hook de integração [`astro:config:setup`](https://docs.astro.build/pt-br/reference/integrations-reference/#astroconfigsetup)).
O hook `setup` pode ser utilizado para atualizar a configuração do Starlight ou adicionar integrações ao Astro.

Este hook é chamado com as seguintes opções:

#### `config`

**tipo:** `StarlightUserConfig`

Uma cópia somente leitura da [configuração do Starlight](/pt-br/reference/configuration/) fornecida pelo usuário.
Essa configuração pode ter sido atualizada por outros plugins configurados antes desse.

#### `updateConfig`

**tipo:** `(newConfig: StarlightUserConfig) => void`

Uma função de callback para atualizar a [configuração do Starlight](/pt-br/reference/configuration/) fornecida pelo usuário.
Forneça as chaves de configuração no nível-raiz que você quer sobrescrever.
Para atualizar valores de configuração aninhados, você precisa fornecer todo o objeto aninhado.

Para estender uma configuração existente sem sobrescrevê-la, espalhe os valores existentes no novo valor.
No seguinte exemplo, uma nova conta de mídia [`social`](/pt-br/reference/configuration/#social) é adicionada à configuração existente espalhando o objeto `config.social` no novo objeto `social`:

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

Uma função de callback para adicionar uma [Integração Astro](https://docs.astro.build/pt-br/reference/integrations-reference/) requerida pelo plugin.

No seguinte exemplo, o plugin primeiro verifica se a [Integração de React do Astro](https://docs.astro.build/pt-br/guides/integrations-guide/react/) está configurada, caso não esteja, usa `addIntegration()` para adicioná-la:

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

      // Só adiciona a integração com o React se ela não estiver carregada.
      if (!isReactLoaded) {
        addIntegration(react());
      }
    },
  },
};
```

#### `astroConfig`

**tipo:** `AstroConfig`

Uma cópia somente leitura da [configuração do Astro](https://docs.astro.build/pt-br/reference/configuration-reference/) fornecida pelo usuário.

#### `command`

**tipo:** `'dev' | 'build' | 'preview'`

O comando usado para executar o Starlight:

- `dev` - O projeto é executado com `astro dev`
- `build` - O projeto é executado com `astro build`
- `preview` - O projeto é executado com `astro preview`

#### `isRestart`

**tipo:** `boolean`

`false` quando o servidor de desenvolvimento é iniciado, `true` quando uma atualização é disparada.
Razões comuns para um reinicio incluem um usuário editando o arquivo `astro.config.mjs` enquanto o servidor de desenvolvimento está sendo executado.

#### `logger`

**tipo:** `AstroIntegrationLogger`

Uma instância do [Astro logger](https://docs.astro.build/pt-br/reference/integrations-reference/#astrointegrationlogger) que você pode usar para escrever logs.
Todas as mensagens de log serão prefixadas com o nome do plugin.

```ts {6}
// plugin.ts
export default {
  name: 'long-process-plugin',
  hooks: {
    setup({ logger }) {
      logger.info('Iniciando um processo longo…');
      // Algum processo longo…
    },
  },
};
```

O exemplo acima registrará um log que inclui a mensagem `info` fornecida:

```shell
[long-process-plugin] Iniciando um processo longo…
```
