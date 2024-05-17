---
title: Referencia de Plugins
description: Una descripción general de la API de plugins de Starlight.
tableOfContents:
  maxHeadingLevel: 4
---

Los plugins de Starlight pueden personalizar la configuración, la UI y el comportamiento, mientras que son fáciles de compartir y reutilizar.
Esta página de referencia documenta la API a la que tienen acceso los plugins.

Aprende más sobre el uso de un plugin de Starlight en la [Referencia de Configuración](/es/reference/configuration/#plugins) o visita la [exhibición de plugins](/es/resources/plugins/#plugins) para ver una lista de los plugins disponibles.

## Referencia rápida de la API

Un plugin de Starlight tiene la siguiente forma.
Consulta a continuación los detalles de las diferentes propiedades y parámetros del hook.

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

Un plugin debe proporcionar un nombre único que lo describa. El nombre se utiliza cuando se [registran mensajes](#logger) relacionados con este plugin y puede ser utilizado por otros plugins para detectar la presencia de este plugin.

## `hooks`

Los hooks son funciones que Starlight llama para ejecutar código de plugin en momentos específicos. Actualmente, Starlight admite un único hook `setup`.

### `hooks.setup`

La función de configuración es llamada cuando se inicializa Starlight (durante el hook de integración [`astro:config:setup`](https://docs.astro.build/es/reference/integrations-reference/#astroconfigsetup)).

El hook `setup` se puede utilizar para actualizar la configuración de Starlight o añadir integraciones de Astro.

Este hook es llamado con las siguientes opciones:

#### `config`

**tipo:** `StarlightUserConfig`

Una copia de lectura de la [configuración de Starlight](/es/reference/configuration/) proporcionada por el usuario.
Esta configuración puede haber sido actualizada por otros plugins configurados antes del actual.

#### `updateConfig`

**tipo:** `(newConfig: StarlightUserConfig) => void`

Una función callback para actualizar la [configuración de Starlight](/es/reference/configuration/).
Proporciona las claves de configuración de nivel raíz que deseas sobreescribir.
Para actualizar los valores de configuración anidados, debes proporcionar el objeto anidado completo.

Para extender una opción de configuración existente sin sobreescribirla, extiende el valor existente en tu nuevo valor.
En el siguiente ejemplo, se agrega una nueva cuenta en [`social`](/es/reference/configuration/#social) a la configuración existente extendiendo 'config.social' en el nuevo objeto social:

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

Una función callback para añadir una [integración de Astro](https://docs.astro.build/es/reference/integrations-reference/) requerida por el plugin.

En el siguiente ejemplo, el plugin primero comprueba si la [integración de React de Astro](https://docs.astro.build/es/guides/integrations-guide/react/) está configurada y, si no lo está, utiliza `addIntegration()` para añadirla:

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

      // Solo agrega la integración de React si aún no está cargada.
      if (!isReactLoaded) {
        addIntegration(react());
      }
    },
  },
};
```

#### `astroConfig`

**tipo:** `AstroConfig`

Una copia de lectura de la [configuración de Astro](https://docs.astro.build/es/reference/configuration-reference/) proporcionada por el usuario.

#### `command`

**tipo:** `'dev' | 'build' | 'preview'`

El comando usado para ejecutar Starlight:

- `dev` - El proyecto se ejecuta con `astro dev`
- `build` - El proyecto se ejecuta con `astro build`
- `preview` - El proyecto se ejecuta con `astro preview`

#### `isRestart`

**tipo:** `boolean`

`false` cuando el servidor de desarrollo se inicia, `true` cuando se activa una recarga.
Common reasons for a restart include a user editing their `astro.config.mjs` while the dev server is running.

#### `logger`

**tipo:** `AstroIntegrationLogger`

Una instancia del [logger de integración de Astro](https://docs.astro.build/es/reference/integrations-reference/#astrointegrationlogger) que puedes utilizar para escribir logs.
Todos los mensajes de registro se prefijarán con el nombre del plugin.

```ts {6}
// plugin.ts
export default {
  name: 'long-process-plugin',
  hooks: {
    setup({ logger }) {
      logger.info('Empezando un proceso largo…');
      // Algun proceso largo…
    },
  },
};
```

El ejemplo anterior registrará un mensaje que incluye el mensaje de información proporcionado:

```shell
[long-process-plugin] Empezando un proceso largo…
```
