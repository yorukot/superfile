---
title: Frontmatter Referenz
description: Ein Überblick über die von Starlight unterstützten Standard-Frontmatter-Felder.
---

Du kannst einzelne Markdown- und MDX-Seiten in Starlight anpassen, indem du Werte in deren Frontmatter setzt. Zum Beispiel könnte eine normale Seite die Felder `title` und `description` setzen:

```md
---
title: Über dieses Projekt
description: Erfahre mehr über das Projekt, an dem ich gerade arbeite.
---

Willkommen auf der Info-Seite!
```

## Frontmatter-Felder

### `title` (erforderlich)

**type:** `string`

Du musst für jede Seite einen Titel angeben. Dieser wird oben auf der Seite, in Browser-Tabs und in den Seiten-Metadaten angezeigt.

### `description`

**type:** `string`

Die Seitenbeschreibung wird für die Metadaten der Seite verwendet und wird von Suchmaschinen und in der Vorschau von sozialen Medien angezeigt.

### `editUrl`

**type:** `string | boolean`

Überschreibt die [globale `editLink`-Konfiguration](/de/reference/configuration/#editlink). Setze die Konfiguration auf `false`, um den Link `Seite bearbeiten` für eine bestimmte Seite zu deaktivieren oder gibt eine alternative URL an, unter der der Inhalt dieser Seite bearbeitet werden kann.

### `head`

**type:** [`HeadConfig[]`](/de/reference/configuration/#headconfig)

Du kannst zusätzliche Tags zum `<head>` deiner Seite hinzufügen, indem du das Feld `head` Frontmatter verwendest. Dies bedeutet, dass du benutzerdefinierte Stile, Metadaten oder andere Tags zu einer einzelnen Seite hinzufügen kannst. Ähnlich wie bei der [globalen `head` Option](/de/reference/configuration/#head).

```md
---
title: Über uns
head:
  # Benutze einen eigenen <title> Tag
  - tag: title
    content: Benutzerdefinierter "Über uns"-Titel
---
```

### `tableOfContents`

**type:** `false | { minHeadingLevel?: number; maxHeadingLevel?: number; }`

Überschreibt die [globale `tableOfContents`-Konfiguration](/de/reference/configuration/#tableofcontents).
Passe die einzuschließenden Überschriftsebenen an oder setze sie auf `false`, um das Inhaltsverzeichnis auf dieser Seite auszublenden.

```md
---
title: Seite mit nur H2s im Inhaltsverzeichnis
tableOfContents:
  minHeadingLevel: 2
  maxHeadingLevel: 2
---
```

```md
---
title: Seite ohne Inhaltsverzeichnis
tableOfContents: false
---
```

### `template`

**type:** `'doc' | 'splash'`  
**default:** `'doc'`

Legt die Layoutvorlage für diese Seite fest.
Seiten verwenden standardmäßig das `'doc'`-Layout.
Setze den Typen auf `'splash'`, um ein breiteres Layout ohne Seitenleisten zu verwenden, welches spezifisch für Startseiten entwickelt wurde.

### `hero`

**type:** [`HeroConfig`](#heroconfig)

Fügt eine Hero-Komponente oben auf der Seite ein. Kann sehr gut mit `template: splash` kombiniert werden.

Zum Beispiel zeigt diese Konfiguration einige übliche Optionen, einschließlich des Ladens eines Bildes aus deinem Repository.

```md
---
title: Meine Website
template: splash
hero:
  title: 'Mein Projekt: Schnell ins All'
  tagline: Bringe deine Wertgegenstände im Handumdrehen auf den Mond und wieder zurück.
  image:
    alt: Ein glitzerndes, leuchtend farbiges Logo
    file: ../../assets/logo.png
  actions:
    - text: Erzähl mir mehr
      link: /getting-started/
      icon: right-arrow
      variant: primary
    - text: Schau mal auf GitHub vorbei
      link: https://github.com/astronaut/mein-projekt
      icon: external
---
```

Du kannst verschiedene Versionen der Hero-Komponente im hellen und dunklen Modus anzeigen.

```md
---
hero:
  image:
    alt: Ein glitzerndes, farbenfrohes Logo
    dark: ../../assets/logo-dark.png
    light: ../../assets/logo-light.png
---
```

#### `HeroConfig`

```ts
interface HeroConfig {
  title?: string;
  tagline?: string;
  image?:
    | {
        // Relativer Pfad zu einem Bild in deinem Repository.
        file: string;
        // Alt-Text, um das Bild für unterstützende Technologien zugänglich zu machen
        alt?: string;
      }
    | {
        // Relativer Pfad zu einem Bild in deinem Repository, das für den dunklen Modus verwendet werden soll.
        dark: string;
        // Relativer Pfad zu einem Bild in deinem Repository, das für den hellen Modus verwendet werden soll.
        light: string;
        // Alt-Text, um das Bild für unterstützende Technologien zugänglich zu machen
        alt?: string;
      }
    | {
        // HTML, welches im Bild-Slot verwendet werden soll.
        // Dies kann ein benutzerdefinierter `<img>`-Tag oder ein Inline-`<svg>` sein.
        html: string;
      };
  actions?: Array<{
    text: string;
    link: string;
    variant: 'primary' | 'secondary' | 'minimal';
    icon: string;
  }>;
}
```

### `banner`

**type:** `{ content: string }`

Zeigt ein Ankündigungsbanner oben auf dieser Seite an.

Der Wert `content` kann HTML für Links oder andere Inhalte enthalten.
Auf dieser Seite wird beispielsweise ein Banner mit einem Link zu `example.com` angezeigt.

```md
---
title: Seite mit Banner
banner:
  content: |
    Wir haben gerade etwas Cooles angefangen!
    <a href="https://example.com">Jetzt besuchen</a>
---
```

### `lastUpdated`

**type:** `Date | boolean`

Überschreibt die [globale Option `lastUpdated`](/de/reference/configuration/#lastupdated). Wenn ein Datum angegeben wird, muss es ein gültiger [YAML-Zeitstempel](https://yaml.org/type/timestamp.html) sein und überschreibt somit das im Git-Verlauf für diese Seite gespeicherte Datum.

```md
---
title: Seite mit einem benutzerdefinierten Datum der letzten Aktualisierung
lastUpdated: 2022-08-09
---
```

### `prev`

**type:** `boolean | string | { link?: string; label?: string }`

Überschreibt die [globale Option `pagination`](/de/reference/configuration/#pagination). Wenn eine Zeichenkette angegeben wird, wird der generierte Linktext ersetzt und wenn ein Objekt angegeben wird, werden sowohl der Link als auch der Text überschrieben.

```md
---
# Versteckt den Link zur vorherigen Seite
prev: false
---
```

```md
---
# Überschreibe den Linktext der vorherigen Seite
prev: Fortsetzung des Tutorials
---
```

```md
---
# Überschreibe sowohl den Link zur vorherigen Seite als auch den Text
prev:
  link: /unverwandte-seite/
  label: Schau dir diese andere Seite an
---
```

### `next`

**type:** `boolean | string | { link?: string; label?: string }`

Dasselbe wie [`prev`](#prev), aber für den Link zur nächsten Seite.

```md
---
# Versteckt den Link zur nächsten Seite
next: false
---
```

### `pagefind`

**type:** `boolean`  
**default:** `true`

Legt fest, ob diese Seite in den [Pagefind](https://pagefind.app/)-Suchindex aufgenommen werden soll. Setze das Feld auf `false`, um eine Seite von den Suchergebnissen auszuschließen:

```md
---
# Diese Seite aus dem Suchindex ausblenden
pagefind: false
---
```

### `sidebar`

**type:** [`SidebarConfig`](#sidebarconfig)

Steuert, wie diese Seite in der [Seitenleiste](/de/reference/configuration/#sidebar) angezeigt wird, wenn eine automatisch generierte Linkgruppe verwendet wird.

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

**type:** `string`  
**default:** the page [`title`](#title-erforderlich)

Legt die Bezeichnung für diese Seite in der Seitenleiste fest, wenn sie in einer automatisch erzeugten Linkgruppe angezeigt wird.

```md
---
title: Über dieses Projekt
sidebar:
  label: Infos
---
```

#### `order`

**type:** `number`

Steuere die Reihenfolge dieser Seite beim Sortieren einer automatisch erstellten Gruppe von Links.
Niedrigere Nummern werden in der Linkgruppe weiter oben angezeigt.

```md
---
title: Erste Seite
sidebar:
  order: 1
---
```

#### `hidden`

**type:** `boolean`  
**default:** `false`

Verhindert, dass diese Seite in eine automatisch generierte Seitenleistengruppe aufgenommen wird.

```md
---
title: Versteckte Seite
sidebar:
  hidden: true
---
```

#### `badge`

**type:** <code>string | <a href="/de/reference/configuration/#badgeconfig">BadgeConfig</a></code>

Füge der Seite in der Seitenleiste ein Abzeichen hinzu, wenn es in einer automatisch generierten Gruppe von Links angezeigt wird.
Bei Verwendung einer Zeichenkette wird das Abzeichen mit einer Standard-Akzentfarbe angezeigt.
Optional kann ein [`BadgeConfig` Objekt](/de/reference/configuration/#badgeconfig) mit den Feldern `text` und `variant` übergeben werden, um das Abzeichen anzupassen.

```md
---
title: Seite mit einem Badge
sidebar:
  # Verwendet die Standardvariante, die der Akzentfarbe deiner Website entspricht
  badge: Neu
---
```

```md
---
title: Seite mit einem Abzeichen
sidebar:
  badge:
    text: Experimentell
    variant: caution
---
```

#### `attrs`

**type:** `Record<string, string | number | boolean | undefined>`

HTML-Attribute, die dem Seitenlink in der Seitenleiste hinzugefügt werden, wenn er in einer automatisch generierten Gruppe von Links angezeigt wird.

```md
---
title: Seite im neuen Tab öffnen
sidebar:
  # Dies öffnet den Link in einem neuen Tab
  attrs:
    target: _blank
---
```
