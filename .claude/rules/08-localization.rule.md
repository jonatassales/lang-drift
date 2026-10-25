# Localization rules

Rules for copy and translations in every app. Loaded in every session.

1. **Copy is JSON.** Every translation catalog is a UTF-8 `.json` file under `src/app/template/i18n/messages/`, English included. TypeScript may import JSON, type messages, select a locale and format ICU; it holds no translated value, UI copy map, language-name map or navigation label.

2. **Parity.** Every locale carries the same keys, placeholders, rich-text tags and valid ICU as `en`, and the counts 0, 1, 2, 3 and 12 format in every locale. `drift check localization` and the app's catalog test enforce it; a change to a catalog or a loader runs both.

3. **Locales.** `en`, `pt-BR`, `zh-Hant`, `ja` and `vi`, served under lowercase prefixes (`/pt-br/`, `/zh-hant/`). Traditional Chinese never stands in for Simplified Chinese: a Simplified Chinese visitor gets English with an availability notice.

4. **"Product Vision" stays in English** as the metric name; the explanations around it are translated. Simulated previews, paused sign-up, illustrative data and product calculations are preserved in every locale.

5. **Locale behavior is code.** Locale ids, prefixes, the default locale and the cookie preference live in `template/i18n/routing.ts` and `src/proxy/`, not in the catalogs.
