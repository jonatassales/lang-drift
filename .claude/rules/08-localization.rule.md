# Localization rules

Copy and translations in every app. Loaded in every session.

1. **Copy is JSON.** Every translation catalog is a UTF-8 `.json` file under `src/app/template/i18n/messages/`, English included. Translated message values, UI copy maps, language-name maps and navigation labels never live in TypeScript or JavaScript. This applies to every locale and to dependencies maintained as part of this repository's work.
2. **TypeScript loads, it does not hold copy.** TypeScript may import JSON, type messages, select a locale and format ICU; those loaders hold no translated value.
3. **Parity.** Every locale carries the same keys, placeholders, rich-text tags, valid ICU and numeric formatter inputs as `en`, and the counts 0, 1, 2, 3 and 12 format in every locale with the production catalogs.
4. **Locales.** Every app serves `en`, `pt-BR`, `zh-Hant`, `ja` and `vi`, under lowercase prefixes (`/en/`, `/pt-br/`, `/zh-hant/`, `/ja/`, `/vi/`). Traditional Chinese never stands in for Simplified Chinese: a Simplified Chinese visitor gets English with an availability notice.
5. **Meaning survives translation.** Simulated previews, paused sign-up, illustrative data and product calculations are preserved in every locale.
6. **Locale behavior is code.** Locale ids, public prefixes, the default locale, language negotiation, fallback resolution and the cookie preference live in `template/i18n/routing.ts` and `src/proxy/`, never in the catalogs.
7. **Verification.** `drift check localization` keeps every catalog valid JSON and every loader free of copy; the app's catalog test enforces parity. A change to a catalog or a loader runs both, and review catches visible copy outside those locations.
