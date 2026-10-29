# Search coverage and hardening for the website, the console, extub and ex

**Status:** ready-for-agent

## Problem Statement

LangDrift's public site is meant to be found, and its console, stub API and CLI are meant to be trusted. Neither holds today.

The website already carries most of what search engines read: localized titles and descriptions, self canonicals, hreflang, Open Graph and Twitter cards, a sitemap, and Organization and WebSite structured data, all ported from langdrift.md. Three of those signals contradict each other or depend on luck:

- Every page tells crawlers its x-default is the negotiating root, while the sitemap names the page's own English URL.
- Social images resolve against whatever origin the deployment environment reports, so the build cannot show they point at langdrift.md.
- Only the home page describes itself in structured data.

Nothing pins any of it. A change to the page list, the sitemap or the metadata helper can drift without a failing test. The console is rightly kept out of search, but its tab titles collide or mislead, its brand strings live in TypeScript, and a missing page renders Next's English fallback.

Neither app sets a Content-Security-Policy. The website sends no security headers at all and announces its framework. Both apps run Next.js releases with published critical advisories, including unauthenticated remote code execution, and the console pulls a zod release with a memory-exhaustion advisory through sinapsi. Nothing in continuous integration notices: the pinned govulncheck that spec 0002 promised never landed, and no job reads npm advisories.

The console has three request-handling flaws:

- Its login form sends credentials in the URL when it is submitted before hydration.
- Its server-rendered reads forward the browser's fetch metadata into a check meant for browser calls. A page reached from a link on another site, or rendered over HTTPS on the hosted preview, fails.
- A path with a dot skips its proxy, and with it the `no-store` the proxy sets.

extub falls short of a hardened API in four ways:

- It answers JSON without the headers an API should send.
- It leaks its token's length through timing.
- It answers an oversized body with 400.
- It accepts unbounded runtime fields.

Its WebRTC hub also lets pion open all-interface UDP sockets in a mode that promises loopback only. ex's judge lets the Anthropic SDK follow redirects with the API key attached, reads responses without a bound, and copies provider error bodies into warnings and evidence.

## Solution

Both apps ship patched Next.js, zod and vitest releases first, in their own pull request. Continuous integration gains a Security job that fails a pull request on any known high or critical advisory among the production npm packages, and on any reachable Go vulnerability.

The website answers every response with a static security header set:

- a Content-Security-Policy that keeps every page statically prerendered;
- HSTS ready for preloading;
- framing denied and MIME sniffing off;
- a strict referrer policy and powerful features disabled;
- no framework banner.

Its crawler-facing contract becomes one consistent story, pinned by tests:

- Every indexable page's x-default is its own English URL, as the sitemap says.
- Social images resolve on langdrift.md in the build output.
- Every indexable page carries a WebPage node in its structured data.
- robots.txt keeps search and AI answer crawlers in and AI training crawlers out.
- security.txt names the security contact.
- The 404 page carries a localized title.

Book a demo refuses cross-site and non-JSON requests, as early access does.

The console:

- sends a strict, per-request nonce Content-Security-Policy on every document and a full header set on everything else;
- stays out of search with a robots.txt as well as its noindex;
- gives every screen a unique localized title and renders its own localized 404 and error pages;
- posts its login form;
- reads extub on the server without trusting the browser's fetch metadata.

extub answers like a hardened API. ex's API clients refuse redirects, insist on HTTPS, bound what they read, and keep provider error bodies out of what they print and record.

## User Stories

1. As the owner, I want both apps on Next.js releases without published critical advisories, so that no known remote code execution reaches production.
2. As the owner, I want sinapsi to resolve the patched zod release until sinapsi republishes, so that the console carries no known high advisory.
3. As a reviewer, I want a pull request to fail when a production npm dependency has a known high or critical advisory, so that the next one is caught before merge.
4. As a reviewer, I want a pull request to fail when a Go module calls a known vulnerable function, so that govulncheck finally guards the Go code as spec 0002 promised.
5. As a website visitor, I want every response to deny framing, so that no other site can clickjack the book-a-demo form.
6. As a website visitor, I want scripts, styles, images, fonts, workers and connections limited to langdrift.md, so that an injected tag cannot load code from or send data to another origin.
7. As the owner, I want the website's pages to stay statically prerendered and CDN-cached under that policy, so that security does not cost Core Web Vitals.
8. As the owner, I want HSTS with includeSubDomains and preload on both apps, so that langdrift.md can be submitted to the browsers' preload lists.
9. As a website visitor, I want the site to stop announcing its framework, so that one response tells an attacker less.
10. As a website visitor, I want the camera, microphone, geolocation, payment, USB and topics features disabled, so that no script can ask for them.
11. As a website visitor, I want cross-origin windows kept apart from the site, so that an opener cannot reach into it.
12. As a reviewer of a preview deployment, I want the Vercel Toolbar to keep working, while production never allows its origin.
13. As a crawler, I want a page's x-default to be the same URL in its HTML and in the sitemap, so that each language cluster is unambiguous.
14. As a crawler, I want og:image and twitter:image to be absolute URLs on langdrift.md, so that link previews never point at a preview host or at localhost.
15. As a crawler, I want every indexable page to describe itself as a WebPage, in its language, part of the LangDrift website, so that page and site form one entity graph.
16. As the owner, I want the home page's Organization and WebSite nodes, sameAs included, kept as they are.
17. As the owner, I want search and AI answer crawlers allowed and AI training crawlers disallowed, so that LangDrift is found and cited without feeding training sets.
18. As a security researcher, I want /.well-known/security.txt to tell me where to report a vulnerability, in which languages and until when.
19. As a visitor who reaches a missing page, I want a localized 404 title and a noindex, so that the error reads clearly and is never indexed.
20. As a frontend developer, I want a test to fail when a page is in the metadata JSON but not in the sitemap, or the other way round, so that the sitemap never drifts.
21. As a frontend developer, I want a test to fail when a page loses its unique title, description, self canonical, hreflang, Open Graph or Twitter card or robots directive, so that coverage stays complete.
22. As a frontend developer, I want a test to pin the header set on a page, an API route and a static file, so that a config change cannot silently drop a header.
23. As the owner, I want book a demo to refuse cross-site requests, foreign origins, encoded or non-JSON bodies and oversized declared lengths before it reads anything, so that it already holds the early-access contract when it reaches a provider.
24. As a console user, I want every document to carry a policy whose scripts need this response's nonce, so that an injected script never runs.
25. As a console user, I want that policy on every document, a path with a dot and the 404 included, so that no route escapes it.
26. As a console user, I want every console response to deny framing, sniffing, cross-origin embedding and cross-origin openers, so that the console cannot be framed or embedded.
27. As a console user, I want the console's API responses to forbid any rendering context, so that JSON is never treated as a document.
28. As a console user, I want my login form to post, so that my password never lands in a URL, a history entry or a log.
29. As a console user arriving from a link on another site, I want the page to load, so that the cross-site check guards browser calls, not my navigation.
30. As a console user on the hosted preview, I want server-rendered reads to work over HTTPS.
31. As a platform developer, I want every id that reaches an extub path checked at the boundary, so that paths are built from known-safe segments and nothing is encoded twice.
32. As a platform developer, I want the build to fail if a client component imports the extub bridge, so that the token can never reach a browser bundle.
33. As a console user, I want each screen's tab to say what it shows, with the project, runtime or ticket it is about, so that tabs, history and bookmarks can be told apart.
34. As a console user in any of its locales, I want the titles, the site name and the 404 and error pages in my language.
35. As the owner, I want the console disallowed in robots.txt as well as noindex, so that no compliant crawler fetches it.
36. As a console user, I want a favicon, an Apple touch icon, the application name and the theme color, so that the console reads as LangDrift in the browser.
37. As an extub client, I want every response to forbid rendering, framing, referrers and cross-origin reads, so that its JSON is only ever data.
38. As an extub client, I want a missing or wrong token answered with WWW-Authenticate, and the Bearer scheme matched case-insensitively, as RFC 6750 says.
39. As the owner, I want the token comparison to reveal neither the token nor its length through timing.
40. As an extub client, I want an oversized body answered 413, so that the error says what went wrong.
41. As the owner, I want component ports, images and readiness paths and folder paths bounded, so that the stub's contract matches the API it stands in for.
42. As the owner, I want extub's WebRTC peers to open loopback sockets only, so that the local contract is really loopback-only.
43. As a platform developer, I want extub's tests to run with the race detector and to pin the body limit, the fetch-metadata check, the headers and the caps.
44. As the owner, I want ex's judge to refuse redirects, so that the API key is never sent to another host.
45. As the owner, I want ex to refuse a model base URL that is not HTTPS, loopback aside, so that the key never travels in clear text.
46. As the owner, I want ex's judge to bound the response it reads, so that a misbehaving endpoint cannot exhaust memory.
47. As a team reading evidence, I want a provider failure recorded as its status and request id, never the provider's body, so that evidence stays token-free and small.
48. As a CLI user, I want a malformed console URL reported without the URL itself, so that credentials pasted into it never reach the terminal or a log.
49. As the owner, I want every contract above in the READMEs, together with a post-deploy checklist, so that what code cannot verify (domains, preload, Deployment Protection) is not forgotten.

## Implementation Decisions

### Delivery

- Two pull requests. The first carries only the dependency fix:
  - `next` 16.4.3 in both apps and `vitest` 5.0.5;
  - a workspace override that resolves sinapsi's zod to 4.6.7;
  - removal of the orbo-voice release-age exclusion, whose stated expiry has passed.

  Every version clears the repository's 24-hour minimum release age.
- The second pull request is branched from the first and carries everything else.

### Website headers

- One source: the Next config's headers, for every path. The proxy sets no security header.
- Policy:
  - `default-src 'self'`
  - `script-src 'self' 'unsafe-inline'`, `style-src 'self' 'unsafe-inline'`
  - `img-src 'self' data: blob:`, `font-src 'self' data:`
  - `connect-src 'self'`, `worker-src 'self' blob:`, `manifest-src 'self'`
  - `frame-src 'none'`, `object-src 'none'`, `base-uri 'self'`, `form-action 'self'`, `frame-ancestors 'none'`
- Allowing `'unsafe-inline'` for scripts is the owner's decision. Next inlines its page payload as scripts, nonces would force every page to render on demand, and the site renders no user-supplied HTML.
- Variants by environment:
  - Development adds `'unsafe-eval'` and the development server's websocket.
  - Preview builds add `https://vercel.live`, for the Vercel Toolbar.
  - Production adds nothing.

  The environment is read where the policy is built.
- No `upgrade-insecure-requests`: HSTS already forces HTTPS on the real origin, and the directive breaks plain-HTTP loopback runs.
- The other headers:
  - `Strict-Transport-Security: max-age=63072000; includeSubDomains; preload`
  - `X-Content-Type-Options: nosniff`
  - `Referrer-Policy: strict-origin-when-cross-origin`
  - `X-Frame-Options: DENY`
  - `Permissions-Policy: camera=(), microphone=(), geolocation=(), payment=(), usb=(), browsing-topics=()`
  - `Cross-Origin-Opener-Policy: same-origin`
  - The powered-by header is off.
- Turnstile's origins enter the policy only when sign-up reopens. The README's reopen steps say so.

### Book a demo

- Before reading the body, the route applies the early-access request contract:
  - POST only.
  - `Sec-Fetch-Site: cross-site` and an `Origin` other than the request's own origin answer 403.
  - Any `Content-Encoding`, or a content type other than JSON, answers 415.
  - A declared `Content-Length` over the 4 KiB limit answers 413.
- A same-origin JSON request keeps today's behavior: field errors answer 400, and a valid request answers 503 while the demo flow is paused.

### Website crawler contract

- The metadata helper gives each page an x-default equal to its own English URL, the URL the sitemap already lists.
- The social image files move under the localized segment whose layout sets `metadataBase`, so their absolute URLs resolve on langdrift.md at build time. If the build shows that Next still resolves them elsewhere, they stay where they are and the post-deploy checklist keeps the origin check. Each image keeps its alt text file.
- The localized 404 carries its title from the route-state catalog.
- Structured data:
  - Every indexable page renders one JSON-LD graph holding a WebPage node:
    - `@id` is the canonical with `#webpage` appended;
    - `url` is the canonical;
    - `name` and `description` are the page's meta title and description;
    - `inLanguage` is the locale's BCP 47 tag;
    - `isPartOf` references the WebSite node by its `@id`.
  - The home page's graph keeps its Organization and WebSite nodes, sameAs included, unchanged.
  - The graph is built from the same JSON the meta tags read, and serialized with `<` escaped.
  - Noindex pages carry no WebPage node.
- robots.txt keeps its group for every agent (allow all, disallow `/api/`) and its sitemap line. It adds one group for GPTBot, ClaudeBot, Google-Extended, CCBot, Applebot-Extended and Bytespider that disallows everything.
- security.txt (RFC 9116) is a static public file served at `/.well-known/security.txt` with these fields:
  - `Contact: mailto:security@langdrift.md`
  - `Expires` one year after the commit
  - `Preferred-Languages: en, pt`
  - `Canonical` on langdrift.md

  The proxy never matches it, because its path has a dot.

### Console headers

- A pure policy builder takes a fresh nonce (16 random bytes, base64) and returns the document policy:
  - `default-src 'self'`
  - `script-src 'self' 'nonce-…' 'strict-dynamic'`, `style-src 'self' 'unsafe-inline'`
  - `img-src 'self' data: blob:`, `font-src 'self'`
  - `connect-src 'self'`, `worker-src 'self' blob:`, `manifest-src 'self'`
  - `frame-src 'none'`, `object-src 'none'`, `base-uri 'self'`, `form-action 'self'`, `frame-ancestors 'none'`

  It takes the same development and preview variants as the website.
- Styles keep `'unsafe-inline'`: React style props in shared primitives, and sinapsi's shadow-root style, cannot carry a nonce.
- The proxy builds the policy and does two things with it:
  - It hands the request to next-intl's middleware as a new request carrying the policy. next-intl copies the request headers into its rewrite and next responses, so Next stamps the nonce on its scripts.
  - It sets the same policy on the response.
- The proxy matcher covers every path except `/api`, `/_next`, and the root files (`favicon.ico`, `icon.svg`, `apple-icon.png`, `robots.txt`). A path with a dot under a locale, and the 404, therefore get the policy and `Cache-Control: private, no-store`.
- Config headers, on every path:
  - the same HSTS, nosniff and Permissions-Policy as the website;
  - `Referrer-Policy: same-origin`, `X-Frame-Options: DENY`;
  - `Cross-Origin-Opener-Policy: same-origin`, `Cross-Origin-Resource-Policy: same-origin`;
  - `X-Robots-Tag: noindex, nofollow`.
- API paths add `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`. The proxy's document policy and the config's API policy never apply to the same path, so every response carries at most one policy.

### Console data path

- Server-rendered reads build their in-process request from the incoming host alone. They carry none of the browser's `Origin`, `Sec-Fetch-*` or cookie headers. The cross-site rule keeps guarding browser calls to the route handlers.
- One predicate defines the entity id alphabet: letters, digits, hyphen and underscore, 1 to 128 characters.
  - A page whose route id falls outside it answers 404 without reading extub.
  - Every client call that interpolates an id checks it the same way.
  - No id is encoded twice. Next already encodes route params.
- The extub bridge imports `server-only`.
- The login and sign-up form declares `method="post"`.

### Console metadata

- The layout's metadata comes from the catalogs: the site name, the default title, the title template and the application name. The viewport declares both color schemes and the theme color.
- Entity pages title themselves after their entity:
  - a project page after the project's name and section;
  - a runtime page after the runtime's name and section;
  - a ticket after its number and title.

  The metadata reads go through the request-cached server read, so metadata and page share one extub call. A failed read falls back to the section label.
- The new-runtime page and the project's runtime list take the keys their headings already use. The zh-Hant label for configuration stops repeating the one for settings.
- New convention files:
  - a static robots.txt that disallows everything;
  - a favicon and an Apple touch icon from the LangDrift mark the website ships;
  - a localized not-found;
  - a localized global error.

  The not-found and global error use Tailwind classes only, so no inline style escapes the policy.

### extub

- Every response adds four headers:
  - `Content-Security-Policy: default-src 'none'; frame-ancestors 'none'`
  - `X-Frame-Options: DENY`
  - `Referrer-Policy: no-referrer`
  - `Cross-Origin-Resource-Policy: same-origin`
- Authentication:
  - 401 carries `WWW-Authenticate: Bearer realm="extub"`.
  - The Bearer scheme is matched case-insensitively.
  - The token is compared as SHA-256 digests in constant time, so neither the value nor its length leaks.
- A body over the limit answers 413 `request_too_large`.
- Bounds:
  - the `local/folders` path must be absolute, at most 1024 bytes, without NUL, like product paths;
  - a component port, when set, is 1 to 65535;
  - an image is at most 256 bytes;
  - a readiness path starts with `/` and is at most 256 bytes.
- The WebRTC setting engine disables multicast DNS, so peers open loopback sockets only.
- The test script runs the race detector.

### ex

- The judge's SDK client gets four guards:
  - an HTTP client that refuses every redirect;
  - a transport that bounds response bodies at 4 MiB, as the console client does;
  - a check that the model base URL is HTTPS, or HTTP on loopback, before any request;
  - provider failures reported as status and request id only, in warnings and in the judge's stop reason.
- The console client reports a malformed console URL without quoting it.
- The README states the SDK's credential sources, the redirect refusal and the order in which a credential's private file and record are written.

### Continuous integration

- A Security job in the pull-request workflow has two steps:
  - `pnpm audit --prod --audit-level=high`;
  - govulncheck at a pinned version, run over each Go module.

  It follows the workflow's conventions: actions pinned to full SHAs, a read-only token, no persisted credentials. The context file's CI table lists the job.

## Testing Decisions

- A good test asserts what the outside sees: the headers a response carries, the metadata a page exports, the files crawlers read, the status a request gets. It never asserts the shape of helper code.
- Website, in the vitest unit suite:
  - **Config headers.** Through Next's config testing utility, a page, an API route and a static file each carry every header with its value, and the powered-by header stays off. Prior art: langdrift's Next config unit test.
  - **SEO contract.** For every page in the metadata JSON and every locale:
    - title and description unique within the locale;
    - canonical absolute on langdrift.md with a trailing slash;
    - hreflang for every locale, with x-default equal to the English URL;
    - Open Graph url equal to the canonical, with locale and alternates;
    - the Twitter card and the index directive.

    The sitemap lists exactly those canonicals with the same alternates. robots.txt keeps `/api/` out, names the sitemap and disallows the six training agents. security.txt carries its four fields. The early-access confirmation page stays noindex and out of the sitemap.
  - **Structured data.** The builder yields a WebPage per page with the expected `@id`, url, name, description, `inLanguage` and `isPartOf`, and escapes `<`.
  - **Book a demo.** Each refused request answers its status without reading the body, and a same-origin JSON request keeps today's statuses. Prior art: the book-a-demo route test.
- Console, under Node's test runner, with relative imports only (prior art: the extub-link test):
  - The policy builder embeds a fresh nonce per call, holds `'strict-dynamic'` and no `'unsafe-inline'` for scripts, and covers its three environments.
  - The entity id predicate accepts the seed ids and rejects `.`, `..`, `/`, `%`, empty and over-long values.
  - The cross-site rule refuses foreign browser calls and admits an in-process read.
- extub, with `httptest`, in the existing handler and hosted tests:
  - every response carries the headers;
  - 401 carries WWW-Authenticate, and a lowercase scheme is admitted;
  - an oversized body answers 413;
  - `Sec-Fetch-Site: cross-site` and a foreign `Origin` answer 403;
  - the caps hold over HTTP;
  - the folder and runtime bounds answer 400.

  The suite runs with the race detector.
- ex, with `httptest` servers (prior art: the console client's tests):
  - a redirecting endpoint never sees the key on the second host;
  - an oversized body fails bounded;
  - a provider error body never reaches the error text;
  - a non-HTTPS base URL is refused before any request;
  - a malformed console URL's error does not contain the URL.
- Verification beyond tests:
  - build-output inspection of every prerendered page;
  - curl against both apps' production servers;
  - a browser pass over every website page and console route that shows no CSP violation, with one documented exception (zod's caught probe inside sinapsi);
  - an `lsof` check during a WebRTC session.

## Out of Scope

- Real authentication and per-user sessions. The auth provider is an open owner decision; the hosted console relies on Vercel Deployment Protection.
- The console's missing `vi` locale and its catalog parity test.
- Republishing sinapsi. It is the owner's package.
- Vercel project settings, DNS, Search Console, HSTS preload submission and the redirect of the production `vercel.app` alias. These are the owner's steps, listed in the checklist.
- FAQPage, SoftwareApplication and BreadcrumbList structured data, and llms.txt.
- Edge rate limiting (the Vercel Firewall is the owner's setting) and a nonce-based policy on the website (the owner chose a static one).
- The console's light-theme flash, and error codes lost across the server-client boundary. Both are functional issues, not security or search.

## Further Notes

- The findings come from these sources, all on 2026-10-29:
  - static reading of `c2906ff`;
  - production builds of both apps;
  - `pnpm audit`.

  The live Vercel settings were not verified: the connector had no access to the team.
- zod's `new Function` probe inside sinapsi reports a CSP violation where the runtime graph renders. zod catches it and falls back to its interpreter. The report disappears once sinapsi configures zod as jitless.
- The owner's post-deploy checklist:
  - confirm that langdrift.md is the website's production domain and that `og:image` resolves there;
  - redirect the production `vercel.app` alias to langdrift.md;
  - keep Deployment Protection on the console and the stubs;
  - submit langdrift.md to hstspreload.org once every subdomain serves HTTPS;
  - renew security.txt before its Expires date;
  - keep `security@langdrift.md` monitored.
