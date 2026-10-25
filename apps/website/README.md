# Website

Public localized LangDrift product and demo presentation. Demo CTAs lead to the
return-soon section above the footer; requests are paused by default. The preserved
Early Access flow uses the private
`email-store` registry: admission → Website Siteverify → atomic contact/intent
persistence → NATS → `email-sender` → Resend. Website has no Resend credentials
or direct-send fallback. HTTP 202 confirms persistence, not provider acceptance
or inbox delivery. Preview signup remains disabled.

## Development and configuration

`pnpm --filter website dev` runs on port 3000. Tracked development defaults contain
nonsecret local configuration. Put private credentials in ignored
`.env.local`; never in client variables. Local integration requires the provisioned
runtime described in the [Early Access runbook](../../platform/runtimes/early-access/RUNTIME.md).

| Variable | Requirement |
| --- | --- |
| `WEBSITE_PREVIEW_SIGNUP_DISABLED` | Requests stay paused unless explicitly set to `false`; unset, `true`, and invalid values hide the form and block admission. Do not reactivate until the demo runtime is ready. |
| `DOCS_ORIGIN` | Server-only dedicated Docs deployment origin; HTTPS in hosted environments. Local development uses the Docs app on port 3004. Never use the redirecting legacy `docs.langdrift.md` host as this upstream. |
| `NEXT_PUBLIC_SSO_URL` | Explicit account-preview navigation origin; production HTTPS |
| `NEXT_PUBLIC_TURNSTILE_SITE_KEY` | Actual Invisible production widget; official Invisible test key locally |
| `TURNSTILE_SECRET_KEY` | Server-only widget secret; production rejects test keys |
| `TURNSTILE_EXPECTED_HOSTNAME` | Exact approved hostname; use localhost explicitly for local fixtures |
| `TURNSTILE_VERIFY_URL` | Production exactly `https://challenges.cloudflare.com/turnstile/v0/siteverify` |
| `EARLY_ACCESS_PUBLIC_ORIGIN` | Exact canonical browser origin, HTTPS in production |
| `EARLY_ACCESS_REGISTRY_URL` | Private authenticated registry origin, HTTPS in production |
| `EARLY_ACCESS_REGISTRY_KEY` | Server-only website role credential, at least 32 characters |
| `CRYPTO_SECRET` | Server-only 32-byte AES key in canonical base64. Replaces `COOKIE_ENCRYPTION_KEY`. Generate with `node -e "process.stdout.write(require('node:crypto').randomBytes(32).toString('base64'))"`. Put it in ignored `.env` or `.env.production.local`, never in `.env.development` or a `NEXT_PUBLIC_` variable. The same key must be present on every instance of a deploy. |

Resend sender fields and API key belong to Go services, not Website. Production
requires verified private connectivity from Website to registry. Source identity
uses Vercel-owned `x-vercel-forwarded-for` only on Vercel; the real proxy chain and
WAF still require isolated deployment proof. Self-hosted production is fail-closed
until its trusted ingress contract is implemented. Never expose registry admission
to browsers or accept client-supplied IP headers as proof.

The proxy stores the language preference in an encrypted `langdrift-website-locale`
cookie for 180 days. An older plaintext or `v1` cookie is ignored. A URL with a
locale prefix keeps that language. A URL without a prefix uses a valid cookie,
then `Accept-Language`, then English.

## Locales and Docs

Website supports English, Brazilian Portuguese, Traditional Chinese, Japanese and
Vietnamese at `/en/`, `/pt-br/`, `/zh-hant/`, `/ja/` and `/vi/`. Language changes
retain the current route, query and fragment. Simplified Chinese requests use
English with a visible availability notice. Localized missing paths, including
dotted paths, have recovery copy and return 404.

Translation values live in JSON under `src/app/template/i18n/messages/`; localized
metadata lives under `src/app/metadata/`. Some catalog families have direct
consumers outside the next-intl provider. See [LOCALIZATION.md](../../LOCALIZATION.md)
for ownership, terminology and review requirements.

The public Docs routes are `/{locale}/docs/`. Website forwards them and the
`/docs-static/` and `/_pagefind/` asset prefixes to the separate Docs app using
`DOCS_ORIGIN`. Deploy and verify that dedicated upstream, then verify forwarding
and search through Website, before enabling the Docs legacy-host redirect. Source
configuration alone does not establish a working hosted Docs deployment.

## User flow and boundaries

The voice interaction is a scripted preview with a prepared answer. Its visible
description and accessible state feedback disclose that it neither uses a
microphone nor generates a live answer. The following Early Access behavior is
preserved for a separately authorized activation; it remains paused here.

Focus lazily prepares Invisible Turnstile. Valid submission obtains a local lock,
executes one challenge and sends JSON to `/api/early-access/`. Spinner, live feedback,
`aria-busy`, disabled controls and reduced-motion behavior cover verification,
submission and navigation. Errors retain the address; retry executes a fresh
challenge. Unmount cancels client work. Cancellation after server persistence does
not undo the registration; repeated requests preserve the immutable intent.

The route checks exact origin/host, body≤4KiB, strict field types/properties and
address limits before registry admission. Rejected admission never calls Siteverify.
Siteverify has a 2.5s timeout and strict success/hostname/action validation. Only
verified tokens can reach registration. Tokens never enter the registry, broker,
logs, URL, browser persistence or Resend. Registry 429 includes Retry-After; 503 is
unavailable. No result reveals whether an address already exists.

## Verification

- `pnpm --filter website test`, `pnpm --filter website typecheck`
- `pnpm check:i18n`: repository JSON catalog guardrail
- `pnpm exec turbo run build --filter=website... --filter=docs...`: build both
  applications and workspace dependencies. Supply a valid
  `DOCS_ORIGIN` and use explicit local-build configuration for local origins.

End-to-end tests belong to `platform/runtimes/`; Website has no app-owned
end-to-end suite.

The historical `@langdrift/early-access` package commands are no longer present in
this workspace. Use the current Website tests and the Early Access runtime runbook
for the preserved service boundary. Local fixtures do not prove deployed WAF,
provider configuration or inbox delivery. Separately authorized real-widget smoke
on its approved hostname remains a release gate before signup activation.

## Privacy and release

Localized `/privacy/` and footer links disclose Cloudflare Invisible Turnstile and
reference its [Privacy Addendum](https://www.cloudflare.com/en-gb/turnstile-privacy-policy/).
The notice now describes private contact persistence and queued confirmation.
Retention/deletion and operator/legal details still require owner review before
release; no global compliance claim is made. Production needs approved widget
hostnames, pre-clearance off, TLS/private connectivity, WAF/proxy proof, Resend
permissions and controlled smoke. A future CSP must follow
[Cloudflare guidance](https://developers.cloudflare.com/turnstile/reference/content-security-policy/).

Rollback means pausing admission and deploying only a verified format2-compatible
image. Do not restore direct Resend, v1 writers or old data to resume sends. Follow
the runtime migration/restore procedure; this change has durable state.
