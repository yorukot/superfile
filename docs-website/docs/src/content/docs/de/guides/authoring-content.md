---
title: Inhalte in Markdown verfassen
description: Ein Überblick über die von Starlight unterstützte Markdown-Syntax.
---

Starlight unterstützt die gesamte Bandbreite der [Markdown](https://daringfireball.net/projects/markdown/) Syntax in `.md` Dateien sowie Frontmatter [YAML](https://dev.to/paulasantamaria/introduction-to-yaml-125f) um Metadaten wie Titel und Beschreibung zu definieren.

Bitte prüfe die [MDX docs](https://mdxjs.com/docs/what-is-mdx/#markdown) oder [Markdoc docs](https://markdoc.dev/docs/syntax), wenn du diese Dateiformate verwendest, da die Unterstützung und Verwendung von Markdown unterschiedlich sein kann.

## Inline-Stile

Text kann **fett**, _italic_, oder ~~durchgestrichen~~ sein.

```md
Text kann **fett**, _italic_, oder ~~durchgestrichen~~ sein.
```

Du kannst [auf eine andere Seite](/de/getting-started/) verlinken.

```md
Du kannst [auf eine andere Seite](/de/getting-started/) verlinken.
```

Du kannst `inline code` mit Backticks hervorheben.

```md
Du kannst `inline code` mit Backticks hervorheben.
```

## Bilder

Bilder in Starlight verwenden [Astros eingebaute optimierte Asset-Unterstützung](https://docs.astro.build/de/guides/assets/).

Markdown und MDX unterstützen die Markdown-Syntax für die Anzeige von Bildern, einschließlich Alt-Text für Bildschirmleser und unterstützende Technologien.

![Eine Illustration von Planeten und Sternen mit dem Wort "Astro"](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-image.png)

```md
![Eine Illustration von Planeten und Sternen mit dem Wort "astro"](https://raw.githubusercontent.com/withastro/docs/main/public/default-og-imag)
```

Relative Bildpfade werden auch für lokal in Ihrem Projekt gespeicherte Bilder unterstützt.

```md
// src/content/docs/page-1.md

![Ein Raketenschiff im Weltraum](../../assets/images/rocket.svg)
```

## Überschriften

Mit einer Überschrift kannst du den Inhalt strukturieren. Überschriften in Markdown werden durch eine Reihe von `#` am Anfang der Zeile gekennzeichnet.

### Wie du Seiteninhalte in Starlight strukturierst

Starlight ist so konfiguriert, dass es automatisch den Seitentitel als Überschrift verwendet und eine "Übersicht"-Überschrift an den Anfang des Inhaltsverzeichnisses jeder Seite setzt. Wir empfehlen, jede Seite mit normalem Text zu beginnen und die Seitenüberschriften ab `<h2>` zu verwenden:

```md
---
title: Markdown Anleitung
description: Wie man Markdown in Starlight benutzt
---

Diese Seite beschreibt, wie man Markdown in Starlight benutzt.

## Inline-Stile

## Überschriften
```

### Automatische Überschriften-Ankerlinks

Wenn du Überschriften in Markdown verwendst, erhaltst du automatisch Ankerlinks, so dass du direkt auf bestimmte Abschnitte deiner Seite verlinken kannst:

```md
---
title: Meine Seite mit Inhalt
description: Wie man Starlight's eingebaute Ankerlinks benutzt
---

## Einleitung

Ich kann auf [meine Schlussfolgerung](#schlussfolgerung) weiter unten auf derselben Seite verlinken.

## Schlussfolgerung

`https://meine-site.com/seite1/#einleitung` navigiert direkt zu meiner Einleitung.
```

Überschriften der Ebene 2 (`<h2>`) und der Ebene 3 (`<h3>`) werden automatisch im Inhaltsverzeichnis der Seite angezeigt.

## Nebenbemerkungen

Nebenbemerkungen (auch bekannt als "Ermahnungen" oder "Callouts") sind nützlich, um sekundäre Informationen neben dem Hauptinhalt einer Seite anzuzeigen.

Starlight bietet eine eigene Markdown-Syntax für die Darstellung von Nebeninformationen. Seitenblöcke werden mit einem Paar dreifacher Doppelpunkte `:::` angezeigt, um den Inhalt zu umschließen, und können vom Typ `note`, `tip`, `caution` oder `danger` sein.

Sie können alle anderen Markdown-Inhaltstypen innerhalb einer Nebenbemerkung verschachteln, allerdings eignen sich diese am besten für kurze und prägnante Inhaltsstücke.

### Nebenbemerkung `note`

:::note
Starlight ist ein Toolkit für Dokumentations-Websites, das mit [Astro](https://astro.build/de) erstellt wurde. Du kannst mit diesem Befehl beginnen:

```sh
npm create astro@latest -- --template starlight
```

:::

````md
:::note
Starlight ist ein Toolkit für Dokumentations-Websites, das mit [Astro](https://astro.build/de) erstellt wurde. Du kannst mit diesem Befehl beginnen:

```sh
npm create astro@latest -- --template starlight
```

:::
````

### Benutzerdefinierte Nebenbemerkungstitel

Du kannst einen benutzerdefinierten Titel für die Nebenbemerkung in eckigen Klammern nach dem Typen angeben, z.B. `:::tip[Wusstest du schon?]`.

:::tip[Wusstest du schon?]
Astro hilft dir, schnellere Websites mit ["Islands Architecture"](https://docs.astro.build/de/concepts/islands/) zu erstellen.
:::

```md
:::tip[Wusstest du schon?]
Astro hilft dir, schnellere Websites mit ["Islands Architecture"](https://docs.astro.build/de/concepts/islands/) zu erstellen.
:::
```

### Weitere Typen

Vorsichts- und Gefahrenhinweise sind hilfreich, um die Aufmerksamkeit des Benutzers auf Details zu lenken, über die er stolpern könnte.
Wenn du diese häufig verwenden, kann das auch ein Zeichen dafür sein, dass die Sache, die Sie dokumentieren, von einem neuen Design profitieren könnte.

:::caution
Wenn du nicht sicher bist, ob du eine großartige Dokumentseite willst, überlege es dir zweimal, bevor du [Starlight](/de/) verwendest.
:::

:::danger
Deine Benutzer können dank hilfreicher Starlight-Funktionen produktiver sein und dein Produkt einfacher nutzen.

- Übersichtliche Navigation
- Benutzer-konfigurierbares Farbthema
- [i18n-Unterstützung](/de/guides/i18n/)

:::

```md
:::caution
Wenn du nicht sicher bist, ob du eine großartige Dokumentseite willst, überlege es dir zweimal, bevor du [Starlight](/de/) verwendest.
:::

:::danger
Deine Benutzer können dank hilfreicher Starlight-Funktionen produktiver sein und dein Produkt einfacher nutzen.

- Übersichtliche Navigation
- Benutzer-konfigurierbares Farbthema
- [i18n-Unterstützung](/de/guides/i18n/)

:::
```

## Blockzitate

> Dies ist ein Blockzitat, das üblicherweise verwendet wird, wenn eine andere Person oder ein Dokument zitiert wird.
>
> Blockzitate werden durch ein ">" am Anfang jeder Zeile gekennzeichnet.

```md
> Dies ist ein Blockzitat, das üblicherweise verwendet wird, wenn eine andere Person oder ein Dokument zitiert wird.
>
> Blockzitate werden durch ein ">" am Anfang jeder Zeile gekennzeichnet.
```

## Code blocks

Ein Codeblock wird durch einen Block mit drei Backticks <code>```</code> am Anfang und Ende gekennzeichnet. Du kannst die verwendete Programmiersprache nach den ersten drei Backticks angeben.

```js
// Javascript code with syntax highlighting.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```

````md
```js
// Javascript code with syntax highlighting.
var fun = function lang(l) {
  dateformat.i18n = require('./lang/' + l);
  return true;
};
```
````

```md
Lange, einzeilige Codeblöcke sollten nicht umgebrochen werden. Sie sollten horizontal scrollen, wenn sie zu lang sind. Diese Zeile sollte lang genug sein, um dies zu demonstrieren.
```

## Andere allgemeine Markdown-Funktionen

Starlight unterstützt alle anderen Markdown-Autorensyntaxen, wie Listen und Tabellen. Einen schnellen Überblick über alle Markdown-Syntaxelemente findest du im [Markdown Cheat Sheet von The Markdown Guide](https://www.markdownguide.org/cheat-sheet/).

## Erweiterte Markdown- und MDX-Konfiguration

Starlight verwendet Astros Markdown- und MDX-Renderer, der auf remark und rehype aufbaut. Du kannst eine Unterstützung für eigene Syntax und Verhalten hinzufügen, indem du `remarkPlugins` oder `rehypePlugins` in deiner Astro-Konfigurationsdatei hinzufügst. Weitere Informationen findest du unter ["Markdown konfigurieren"] (https://docs.astro.build/de/guides/markdown-content/#markdown-konfigurieren) in der Astro-Dokumentation.
