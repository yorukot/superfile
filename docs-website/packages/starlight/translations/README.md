# Starlight UI translation files

This directory contains translation data for Starlight’s UI.
Each language has its own JSON file and follows the [translation structure described in Starlight’s docs](https://starlight.astro.build/guides/i18n/#translate-starlights-ui).

## Add a new language

1. Create a JSON file named using the BCP-47 tag for the language, e.g. `en.json` or `ja.json`.

2. Fill the file with translations for each UI string. You can base your translations on [`en.json`](./en.json). Translate only the values, leaving the keys in English (e.g. `"search.label": "Buscar"`).

3. Import your file in [`index.ts`](./index.ts) and add your language to the `Object.entries`.

4. Open a pull request on GitHub to add your file to Starlight!
