---
title: Umweltfreundliche Dokumente
description: Erfahre, wie Starlight dir helfen kann, umweltfreundlichere Dokumentationen zu erstellen und deinen ökologischen Fußabdruck zu verringern.
---

Schätzungen über die Auswirkungen der Webindustrie auf das Klima liegen zwischen [2 %][sf] und [4 % der weltweiten Kohlenstoffemissionen][bbc], was in etwa den Emissionen der Luftfahrtindustrie entspricht.
Es gibt viele komplexe Faktoren bei der Berechnung der ökologischen Auswirkungen einer Website, aber dieser Leitfaden enthält einige Tipps, wie du den ökologischen Fußabdruck deiner Docs-Website verringern kannst.

Die gute Nachricht ist, dass die Wahl von Starlight ein guter Anfang ist.
Laut dem Website Carbon Calculator ist diese Website [sauberer als 99 % der getesteten Websiten][sl-carbon] und erzeugt 0,01g CO₂ pro Seitenbesuch.

## Seitengewicht

Je mehr Daten eine Website überträgt, desto mehr Energieressourcen benötigt sie.
Im April 2023 musste ein Nutzer laut [Daten aus dem HTTP-Archiv][http] für die durchschnittliche Website mehr als 2.000 KB herunterladen.

Starlight erstellt Seiten, die so leicht wie möglich sind.
So lädt ein Benutzer beim ersten Besuch weniger als 50 KB an komprimierten Daten herunter - nur 2,5 % des Medianwerts des HTTP-Archivs.
Mit einer guten Caching-Strategie können nachfolgende Besuche sogar nur 10 KB herunterladen.

### Bilder

Während Starlight eine gute Grundlage bietet, können Bilder, die du deinen Dokumentseiten hinzufügst, das Seitengewicht schnell erhöhen.
Starlight nutzt die [optimierte Asset-Unterstützung][Assets] von Astro, um lokale Bilder in deinen Markdown- und MDX-Dateien zu optimieren.

### UI-Komponenten

Komponenten, die mit UI-Frameworks wie React oder Vue erstellt wurden, können leicht große Mengen an JavaScript zu einer Seite hinzufügen.
Da Starlight auf Astro aufbaut, laden Komponenten wie diese dank [Astro Islands][islands] standardmäßig **kein clientseitiges JavaScript**.

### Caching

Caching wird verwendet, um zu kontrollieren, wie lange ein Browser Daten speichert und wiederverwendet, die er bereits heruntergeladen hat.
Eine gute Caching-Strategie stellt sicher, dass ein Benutzer neue Inhalte so schnell wie möglich erhält, wenn sich diese ändern, vermeidet aber auch, dass derselbe Inhalt unnötigerweise immer wieder heruntergeladen wird, wenn er sich nicht geändert hat.

Die gebräuchlichste Art, das Zwischenspeichern zu konfigurieren, ist der [`Cache-Control` HTTP-Header][cache].
Wenn du Starlight verwendest, kannst du eine lange Cache-Zeit für alles im Verzeichnis `/_astro/` einstellen.
Dieses Verzeichnis enthält CSS, JavaScript und andere gebündelte Inhalte, die sicher für immer zwischengespeichert werden können, wodurch unnötige Downloads vermieden werden:

```
Cache-Control: public, max-age=604800, immutable
```

Wie du das Caching konfigurierst, hängt von deinem Webhost ab. Zum Beispiel wendet Vercel diese Caching-Strategie für dich an, ohne dass eine Konfiguration erforderlich ist, während du [benutzerdefinierte Header für Netlify][ntl-headers] einstellen kannst, indem du eine `public/_headers`-Datei zu deinem Projekt hinzufügst:

```
/_astro/*
  Cache-Control: public
  Cache-Control: max-age=604800
  Cache-Control: immutable
```

[cache]: https://csswizardry.com/2019/03/cache-control-for-civilians/
[ntl-headers]: https://docs.netlify.com/routing/headers/

## Stromverbrauch

Die Art und Weise, wie eine Website aufgebaut ist, kann sich auf den Stromverbrauch auswirken, den sie auf dem Gerät des Benutzers benötigt.
Durch die Verwendung von minimalem JavaScript reduziert Starlight die Rechenleistung, die das Telefon, Tablet oder der Computer eines Nutzers zum Laden und Rendern von Seiten benötigt.

Sei jedoch vorsichtig, wenn du Funktionen wie analytische Tracking-Skripte oder JavaScript-lastige Inhalte wie Videoeinbettungen hinzufügst, da diese den Stromverbrauch der Seite erhöhen können.
Wenn du Analysen benötigst, solltest du eine schlanke Option wie [Cabin][cabin], [Fathom][fathom] oder [Plausible][plausible] wählen.
Einbettungen wie YouTube- und Vimeo-Videos können verbessert werden, indem man auf [Laden des Videos bei Benutzerinteraktion][lazy-video] wartet.
Pakete wie [`astro-embed`][embed] können bei gängigen Diensten helfen.

:::tip[Wusstest du schon?]
Das Parsen und Kompilieren von JavaScript ist eine der aufwändigsten Aufgaben, die ein Browser zu erledigen hat.
Verglichen mit dem Rendern eines JPEG-Bildes derselben Größe kann die [Verarbeitung von JavaScript mehr als 30 Mal so lange dauern][cost-of-js].
:::

[cabin]: https://withcabin.com/
[fathom]: https://usefathom.com/
[plausible]: https://plausible.io/
[lazy-video]: https://web.dev/iframe-lazy-loading/
[embed]: https://www.npmjs.com/package/astro-embed
[cost-of-js]: https://medium.com/dev-channel/the-cost-of-javascript-84009f51e99e

## Hosting

Wo eine Website gehostet wird, kann einen großen Einfluss darauf haben, wie umweltfreundlich deine Dokumentationsseite ist.
Rechenzentren und Serveranlagen können große ökologische Auswirkungen haben, einschließlich eines hohen Stromverbrauchs und eines intensiven Wasserverbrauchs.

Wenn du dich für einen Hoster entscheidest, der erneuerbare Energien einsetzt, wird deine Website weniger Kohlenstoffemissionen verursachen. Das [Green Web Directory][gwb] ist ein Tool, das dir helfen kann, Hosting-Unternehmen zu finden.

[gwb]: https://www.thegreenwebfoundation.org/directory/

## Vergleiche

Bist du neugierig, wie andere Docs-Frameworks im Vergleich abschneiden?
Diese Tests mit dem [Website Carbon Calculator][wcc] vergleichen ähnliche Seiten, die mit verschiedenen Tools erstellt wurden.

| Framework                   | CO₂ pro Seitenaufruf |
| --------------------------- | -------------------- |
| [Starlight][sl-carbon]      | 0.01g                |
| [VitePress][vp-carbon]      | 0.05g                |
| [Docus][dc-carbon]          | 0.05g                |
| [Sphinx][sx-carbon]         | 0.07g                |
| [MkDocs][mk-carbon]         | 0.10g                |
| [Nextra][nx-carbon]         | 0.11g                |
| [docsify][dy-carbon]        | 0.11g                |
| [Docusaurus][ds-carbon]     | 0.24g                |
| [Read the Docs][rtd-carbon] | 0.24g                |
| [GitBook][gb-carbon]        | 0.71g                |

<small>Daten erhoben am 14. Mai 2023. Klicke auf einen Link, um aktuelle Zahlen zu sehen.</small>

[sl-carbon]: https://www.websitecarbon.com/website/starlight-astro-build-getting-started/
[vp-carbon]: https://www.websitecarbon.com/website/vitepress-dev-guide-what-is-vitepress/
[dc-carbon]: https://www.websitecarbon.com/website/docus-dev-introduction-getting-started/
[sx-carbon]: https://www.websitecarbon.com/website/sphinx-doc-org-en-master-usage-quickstart-html/
[mk-carbon]: https://www.websitecarbon.com/website/mkdocs-org-getting-started/
[nx-carbon]: https://www.websitecarbon.com/website/nextra-site-docs-docs-theme-start/
[dy-carbon]: https://www.websitecarbon.com/website/docsify-js-org/
[ds-carbon]: https://www.websitecarbon.com/website/docusaurus-io-docs/
[rtd-carbon]: https://www.websitecarbon.com/website/docs-readthedocs-io-en-stable-index-html/
[gb-carbon]: https://www.websitecarbon.com/website/docs-gitbook-com/

## Weitere Ressourcen

### Werkzeuge

- [Website Carbon Calculator][wcc]
- [GreenFrame](https://greenframe.io/)
- [Ecograder](https://ecograder.com/)
- [WebPageTest Kohlenstoffkontrolle](https://www.webpagetest.org/carbon-control/)
- [Ecoping](https://ecoping.earth/)

### Artikel und Vorträge

- ["Building a greener web"](https://youtu.be/EfPoOt7T5lg), Vortrag von Michelle Barker
- ["Sustainable Web Development Strategies Within An Organization"](https://www.smashingmagazine.com/2022/10/sustainable-web-development-strategies-organization/), Artikel von Michelle Barker
- ["A sustainable web for everyone"](https://2021.stateofthebrowser.com/speakers/tom-greenwood/), Vortrag von Tom Greenwood
- ["How Web Content Can Affect Power Usage"](https://webkit.org/blog/8970/how-web-content-can-affect-power-usage/), Artikel von Benjamin Poulain und Simon Fraser

[sf]: https://www.sciencefocus.com/science/what-is-the-carbon-footprint-of-the-internet/
[bbc]: https://www.bbc.com/future/article/20200305-why-your-internet-habits-are-not-as-clean-as-you-think
[http]: https://httparchive.org/reports/state-of-the-web
[assets]: https://docs.astro.build/en/guides/assets/
[islands]: https://docs.astro.build/en/concepts/islands/
[wcc]: https://www.websitecarbon.com/
