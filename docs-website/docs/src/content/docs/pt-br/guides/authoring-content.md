---
title: Escrevendo Conteúdo em Markdown
description: Uma visão geral da sintaxe Markdown suportada pelo Starlight.
---

Starlight suporta completamente a sintaxe [Markdown](https://daringfireball.net/projects/markdown/) em arquivos `.md` assim como frontmatter [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f) para definir metadados como o título e a descrição.

Por favor verifique a [documentação do MDX](https://mdxjs.com/docs/what-is-mdx/#markdown) ou a [documentação do Markdoc](https://markdoc.dev/docs/syntax) se estiver utilizando esses formatos de arquivo, já que o suporte e uso do Markdown podem variar.

## Frontmatter

Você pode customizar páginas individualmente no Starlight passando parâmetros no frontmatter.

Frontmatter é definido no topo de seus arquivos entre divisores `---`:

```md title="src/content/docs/exemplo.md"
---
title: Meu título
---

Conteúdo da página vem depois do segundo `---`.
```

Toda página deve incluir ao menos o `title` (título).
Veja a [referência do frontmatter](/pt-br/reference/frontmatter/) para todos os campos disponíveis e como adicionar campos customizados.

## Estilos Inline

Texto pode estar em **negrito**, _itálico_, ou ~~tachado~~.

```md
Texto pode estar em **negrito**, _itálico_, ou ~~tachado~~.
```

Você pode [fazer links para outras páginas](/pt-br/getting-started/).

```md
Você pode [fazer links para outras páginas](/pt-br/getting-started/).
```

Você pode destacar `código inline` com crases.

```md
Você pode destacar `código inline` com crases.
```

## Imagens

Imagens no Starlight usam o [suporte integrado a assets otimizados do Astro](https://docs.astro.build/pt-br/guides/assets/).

Markdown e MDX suportam a sintaxe do Markdown para mostrar imagens que incluem texto alternativo para leitores de tela e tecnologias assistivas.

![Uma ilustração de planetas e estrelas apresentando a palavra “astro”](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
![Uma ilustração de planetas e estrelas apresentando a palavra “astro”](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)
```

Caminhos de imagem relativos também são suportados para imagens armazenadas localmente no seu projeto.

```md
// src/content/docs/pagina-1.md

![Um foguete no espaço](../../assets/imagens/foguete.svg)
```

## Cabeçalhos

Você pode estruturar o conteúdo utilizando um cabeçalho. Cabeçalhos no Markdown são indicados pelo número de `#` no início de uma linha.

### Como estruturar o conteúdo da página no Starlight

Starlight é configurado para automaticamente utilizar o título da sua página como um cabeçalho superior e irá incluir um cabeçalho "Visão geral" no topo do índice de cada página. Nós recomendamos começar cada página com um parágrafo normal e utilizar cabeçalhos na página a partir de `<h2>` para baixo:

```md
---
title: Guia de Markdown
description: Como utilizar Markdown no Starlight
---

Esta página descreve como utilizar Markdown no Starlight.

## Estilos Inline

## Cabeçalhos
```

### Links de âncora automáticos de cabeçalho

Utilizar cabeçalhos no Markdown irá automaticamente dá-lo links de âncora para que você direcione diretamente a certas seções da sua página:

```md
---
title: Minha página de conteúdo
description: Como utilizar os links de âncora integrados do Starlight
---

## Introdução

Eu posso fazer um link para [minha conclusão](#conclusão) abaixo na mesma página.

## Conclusão

`https://meu-site.com/pagina1/#introdução` navega diretamente para minha Introdução.
```

Cabeçalhos de Nível 2 (`<h2>`) e Nível 3 (`<h3>`) vão aparecer automaticamente no índice da página.

Aprenda mais sobre como Astro processa `id`s de títulos na [documentação do Astro](https://docs.astro.build/pt-br/guides/markdown-content/#ids-de-t%C3%ADtulos)

## Asides

Asides (também conhecidos como “advertências” ou “frases de destaque”) são úteis para mostrar informações secundárias ao lado do conteúdo principal de uma página.

Starlight providencia uma sintaxe Markdown customizada para renderizar asides. Blocos Aside são indicados utilizando um par de dois pontos triplo `:::` para envolver o seu conteúdo, e podem ser do tipo `note`, `tip`, `caution` ou `danger`.

Você pode aninhar qualquer outras formas de conteúdo Markdown dentro de um aside, mas asides são mais adequados para blocos curtos e concisos de conteúdo.

### Aside de Nota

:::note
Starlight é um conjunto de ferramentas para websites de documentação feito com [Astro](https://astro.build/). Você pode começar com o comando:

```sh
npm create astro@latest -- --template starlight
```

:::

````md
:::note
Starlight é um conjunto de ferramentas para websites de documentação feito com [Astro](https://astro.build/). Você pode começar com o comando:

```sh
npm create astro@latest -- --template starlight
```

:::
````

### Títulos de aside customizados

Você pode especificar um título customizado para o aside em colchetes seguindo o tipo do aside, e.x. `:::tip[Você sabia?]`.

:::tip[Você sabia?]
Astro te ajuda a construir websites mais rápidos com a [“Arquitetura em Ilhas”](https://docs.astro.build/pt-br/concepts/islands/).
:::

```md
:::tip[Você sabia?]
Astro te ajuda a construir websites mais rápidos com a [“Arquitetura em Ilhas”](https://docs.astro.build/pt-br/concepts/islands/).
:::
```

### Mais tipos de aside

Asides de cuidado e perigo são úteis para chamar a atenção de um usuário a detalhes que podem o atrapalhar.
Se você anda os utilizando muito, pode ser um sinal de que o que você está documentando se beneficiaria com uma mudança.

:::caution
Se você não tem certeza de que você quer um site de documentação incrível, pense novamente antes de utilizar [Starlight](/pt-br/).
:::

:::danger
Seus usuários podem ser mais produtivos e considerar seu produto mais fácil de usar graças a funcionalidades úteis do Starlight.

- Navegação compreensível
- Tema de cores configurável pelo usuário
- [Suporte a internacionalização](/pt-br/guides/i18n/)

:::

```md
:::caution
Se você não tem certeza de que você quer um site de documentação incrível, pense novamente antes de utilizar [Starlight](/pt-br/).
:::

:::danger
Seus usuários podem ser mais produtivos e considerar seu produto mais fácil de usar graças a funcionalidades úteis do Starlight.

- Navegação compreensível
- Tema de cores configurável pelo usuário
- [Suporte a internacionalização](/pt-br/guides/i18n/)

:::
```

## Citações

> Esta é uma citação, que é comumente utilizada ao citar outra pessoa ou documento.
>
> Citações são indicadas com um `>` no começo de cada linha.

```md
> Esta é uma citação, que é comumente utilizada ao citar outra pessoa ou documento.
>
> Citações são indicadas com um `>` no começo de cada linha.
```

## Blocos de código

Um bloco de código é indicado por um bloco com três crases <code>```</code> no começo e fim. Você pode indicar a linguagem de programação sendo utilizada após as crases iniciais.

```js
// Código JavaScript com syntax highlighting.
var divertido = function lingua(l) {
  formatodata.i18n = require('./lingua/' + l);
  return true;
};
```

````md
```js
// Código JavaScript com syntax highlighting.
var divertido = function lingua(l) {
  formatodata.i18n = require('./lingua/' + l);
  return true;
};
```
````

### Funcionalidades do Expressive Code

Starlight usa [Expressive Code](https://github.com/expressive-code/expressive-code/tree/main/packages/astro-expressive-code) para aumentar as possibilidades de formatação em blocos de código.
Os plugins de marcadores de texto e moldura de janela do Expressive Code são habilitados por padrão.
A renderização de blocos de código pode ser configurada utilizando a [opção de configuração `expressiveCode`](/pt-br/reference/configuration/#expressivecode) do Starlight.

#### Marcadores de texto

Você pode destacar linhas ou partes específicas do seu bloco de código utilizando [marcadores de texto do Expressive Code](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#usage-in-markdown--mdx-documents) na linha de abertura do seu bloco de código.
Use chaves (`{ }`) para destacar linhas inteiras, e aspas para destacar segmentos do texto.

Existem três estilos de destaque: neutro para chamar a atenção para o código, verde para indicar código adicionado, e vermelho para indeicar código deletado.
Tanto texto quanto linhas inteiras podem ser marcados com o marcador padrão, ou combinados com `ins=` e `del=` para produzir o destaque desejado.

Expressive Code provê diversas opções para customizar a aparência visual dos seus exemplos de código.
Muitas dessas podem ser combinadas para exemplos de código altamente ilustrativos.
Explore a [documentação do Expressive Code](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md) para a lista extensiva de opções disponíveis.
Alguns dos exemplos mais comuns estão demonstrados abaixo:

- [Marque linhas inteiras e blocos de linhas usando o marcador `{ }`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-entire-lines--line-ranges):

  ```js {2-3}
  function demo() {
    // Esta linha (#2) e a próxima estão destacadas
    return 'Esta é a linha #3 do snippet';
  }
  ```

  ````md
  ```js {2-3}
  function demo() {
    // Esta linha (#2) e a próxima estão destacadas
    return 'Esta é a linha #3 do snippet';
  }
  ```
  ````

- [Marque partes do texto usando o marcador `" "` or expressões regulares](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#marking-individual-text-inside-lines):

  ```js "Termos individuais" /Até.*suportadas/
  // Termos individuais também podem ser destacados
  function demo() {
    return 'Até expressões regulares são suportadas';
  }
  ```

  ````md
  ```js "Termos individuais" /Até.*suportadas/
  // Termos individuais também podem ser destacados
  function demo() {
    return 'Até expressões regulares são suportadas';
  }
  ```
  ````

- [Marque texto ou linhas como inseridos ou deletados com `ins` ou `del`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#selecting-marker-types-mark-ins-del):

  ```js "return true;" ins="inserido" del="deletado"
  function demo() {
    console.log('Esses são os marcadores inserido e deletado');
    // A expressão de retorno usa o marcador padrão
    return true;
  }
  ```

  ````md
  ```js "return true;" ins="inserido" del="deletado"
  function demo() {
    console.log('Esses são os marcadores inserido e deletado');
    // A expressão de retorno usa o marcador padrão
    return true;
  }
  ```
  ````

- [Combine highlight de sintaxe com sintaxe de `diff`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-text-markers/README.md#combining-syntax-highlighting-with-diff-like-syntax):

  ```diff lang="js"
    function thisIsJavaScript() {
      // O bloco inteiro tem o highlight de JavaScript,
      // e ainda podemos colocar marcadores de diff nele!
  -   console.log('Código antigo a ser removido')
  +   console.log('Código novo e brilhante!')
    }
  ```

  ````md
  ```diff lang="js"
    function thisIsJavaScript() {
      // O bloco inteiro tem o highlight de JavaScript,
      // e ainda podemos colocar marcadores de diff nele!
  -   console.log('Código antigo a ser removido')
  +   console.log('Código novo e brilhante!')
    }
  ```
  ````

#### Molduras e títulos

Blocos de código podem ser renderizados dentro de molduras como se fossem janelas.
Uma moldura que parece com uma janela de terminal pode ser usada para linguagens de scripts em shell (e.g. `bash` e `sh`).
Outras linguages são exibidas dentro de uma moldura similar a um editor de código caso incluam um título.

O título opcional de um bloco de código pode ser definido tanto com um atributo `title="..."` seguindo as crases de abertura do bloco de código e o identificador da linguagem, ou com um comentário contendo o nome do arquivo na primeira linha do código.

- [Adicione uma aba com o nome do arquivo com um comentário](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```js
  // meu-arquivo-de-teste.js
  console.log('Olá Mundo!');
  ```

  ````md
  ```js
  // meu-arquivo-de-teste.js
  console.log('Olá Mundo!');
  ```
  ````

- [Adicione um título a uma janela de Terminal](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#adding-titles-open-file-tab-or-terminal-window-title)

  ```bash title="Instalando dependências…"
  npm install
  ```

  ````md
  ```bash title="Instalando dependências…"
  npm install
  ```
  ````

- [Desabilite molduras de janela com `frame="none"`](https://github.com/expressive-code/expressive-code/blob/main/packages/%40expressive-code/plugin-frames/README.md#overriding-frame-types)

  ```bash frame="none"
  echo "Isso não é exibido como um terminal apesar de usar a linguagem bash"
  ```

  ````md
  ```bash frame="none"
  echo "Isso não é exibido como um terminal apesar de usar a linguagem bash"
  ```
  ````

## Outras funcionalidades comuns do Markdown

Starlight suporta todo o resto da sintaxe de escrita do Markdown, como listas e tabelas. Veja a [Cheat Sheet de Markdown do The Markdown Guide](https://www.markdownguide.org/cheat-sheet/) para uma visão geral rápida de todos os elementos da sintaxe do Markdown.

## Markdown Avançado e Configurando MDX

O Starlight utiliza o mesmo rendizador Markdown e MDX do Astro, que suporta remark e rehype. Você pode adicionar sintaxe e comportamento personalizado adicionando `remarkPlugins` ou `rehypePlugins` no seu arquivo de configuração Astro. Visite [Configurando Markdown e MDX](https://docs.astro.build/pt-br/guides/markdown-content/#configurando-markdown-e-mdx) na documentação do Astro para saber mais.
