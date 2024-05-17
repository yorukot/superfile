---
title: Substituindo Componentes
description: Aprenda como substituir os componentes nativos do Starlight para adicionar elementos personalizados a UI do seu site de documentação.
---

<!---
TODO: Check all links and html anchors
-->

A UI e configuração padrão do Starlight foi projetada para ser flexível a adaptável a uma gama de conteúdos. Boa parte da customização da aparência padrão do Starlight pode ser feita via [CSS](/pt-br/guides/css-and-tailwind/) e [opções de configuração](/pt-br/guides/customization/).

Caso você precise de mais possibilidades, o Starlight suporta a criação dos seus próprios componentes para estender ou substituir completamente os componentes padrões.

## Em que casos substituir

Substituir os componentes padrões do Starlight pode ser útil nos seguintes casos:

- Você deseja mudar parte da UI do Starlight de forma que não é possível com [CSS personalizado](/pt-br/guides/css-and-tailwind/).
- Você deseja mudar o comportamento de parte da UI do Starlight.
- Você deseja adicionar elementos de UI junto da UI existente do Starlight.

## Como substituir

1. Escolha qual componente você deseja substituir.
   Você pode encontrar uma lista completa de componentes na [Referência de Substituições](/pt-br/reference/overrides/).

   Neste exemplo, substituiremos o componente do Starlight [`SocialIcons`](/pt-br/reference/overrides/#socialicons) que fica na barra de navegação.

2. Crie um componente Astro para substituir os componentes Starlight.
   O exemplo abaixo é de um link de contato.

   ```astro
   ---
   // src/components/LinkDeEmail.astro
   import type { Props } from '@astrojs/starlight/props';
   ---

   <a href="mailto:houston@exemplo.com.br">Nosso e-mail</a>
   ```

3. Diga ao Starlight para utilizar seu componente personalizado na opção [`components`](/pt-br/reference/configuration/#components) do arquivo `astro.config.mjs`:

   ```js {9-12}
   // astro.config.mjs
   import { defineConfig } from 'astro/config';
   import starlight from '@astrojs/starlight';

   export default defineConfig({
     integrations: [
       starlight({
         title: 'Minha Documentação com Substituições',
         components: {
           // Substitui o componente padrão `SocialIcons`.
           SocialIcons: './src/components/LinkDeEmail.astro',
         },
       }),
     ],
   });
   ```

## Reutilize um componente padrão

Você pode construir com os componentes de UI padrão do Starlight da mesma forma que faria ao criar seus próprios componentes: importando e renderizando-o dentro do seu componente personalizado. Isso permite que você mantenha toda a UI base do Starlight em seu design e, ao mesmo tempo, adicionar novos elementos a ela.

O exemplo a seguir mostra um componente personalizado que renderiza um link de e-mail junto do componente padrão `SocialIcons`:

```astro {4,8}
---
// src/components/LinkDeEmail.astro
import type { Props } from '@astrojs/starlight/props';
import Padrao from '@astrojs/starlight/components/SocialIcons.astro';
---

<a href="mailto:houston@exemplo.com.br">Nosso e-mail</a>
<Padrao {...Astro.props}><slot /></Padrao>
```

Quando estiver utilizando um componente padrão num componente personalizado:

- Utilize a sintaxe de espalhamentodo `Astro.props`: `{...Astro.props}`. Assim garante-se que o componente receberá todos os dados necessários para renderizar corretamente.
- Adicione um [`<slot />`](https://docs.astro.build/pt-br/core-concepts/astro-components/#slots) dentro do componente padrão. Isso é para garantir que o Astro saiba onde renderizar elementos-filho no componente, se algum for passado.

## Utilize dados da página

Quando estiver substituindo um componente Starlight, a sua implementação receberá um objeto padrão `Astro.props` contendo todas as informações da página atual.
Isso permite que você utilize esses valores para controlar como seu componente renderiza.

Por exemplo, você pode ler os valores do frontmatter a partir do `Astro.props.entry.data`. No exemplo a seguir, utilizamos [`PageTitle`](/pt-br/reference/overrides/#pagetitle) para exibir o título da página atual num componente substituto.

```astro {5} "{title}"
---
// src/components/Titulo.astro
import type { Props } from '@astrojs/starlight/props';

const { titulo } = Astro.props.entry.data;
---

<h1 id="_top">{titulo}</h1>

<style>
  h1 {
    font-family: 'Comic Sans';
  }
</style>
```

Aprenda mais sobre todos os props disponíveis na [Referência de Substituição](/pt-br/reference/overrides/#props-de-componentes).

### Substituindo apenas em páginas específicas

A substituição de componentes aplica-se a todas as páginas. Porém, você pode fazer o componente
renderizar condicionalmente utilizando `Astro.props` para determinar quando exibir a sua UI personaliza, ou a UI padrão do Starlight, ou até mesmo para exibir algo totalmente diferente.

No exemplo a seguir, um componente está substituindo o [`Footer`](/pt-br/reference/overrides/#footer) padrão do Starlight para exibir "Feito com Starlight 🌟" exclusivamente na página principal, e nas outras exibir o rodapé padrão.

```astro
---
// src/components/RodapeCondicional.astro
import type { Props } from '@astrojs/starlight/props';
import Padrao from '@astrojs/starlight/components/Footer.astro';

const isPaginaPrincial = Astro.props.slug === '';
---

{
  isPaginaPrincial ? (
    <footer>Feito com Starlight 🌟</footer>
  ) : (
    <Padrao {...Astro.props}>
      <slot />
    </Padrao>
  )
}
```

Aprenda mais sobre renderização condicional no [Guia de Sintaxe de Template Astro](https://docs.astro.build/pt-br/core-concepts/astro-syntax/#html-din%C3%A2mico).
