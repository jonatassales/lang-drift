# 0050: The components resolve one by one on the ticket card

**Spec:** .agents/workflow/specs/0007-microruntimes-hero.spec.md

**What to build:** while the terminal resolves the spec, the card's components section goes from a dash to five pulsing skeleton rows, then fills one row per terminal line: the technology mark, the path, the kind, CPU and memory, and a thin bar under the row showing its share of the CPU envelope. A total line counts resolved rows and sums their figures as they arrive, and reads five components with the suggested totals at the end.

**Blocked by:** 0047, 0049.

**Status:** ready-for-agent

### Content

- [ ] The ticket plan in the content module carries the five components (path, technology, kind, CPU millicores, memory MiB): Go api 250m and 256 MiB, Node.js api 300m and 512 MiB, NATS broker 200m and 256 MiB, Redis cache 200m and 256 MiB, Next.js web 500m and 1,024 MiB; the bounded envelope of 2,000m and 3,072 MiB; and the suggested totals of 1,450m and 2,304 MiB. The reveal table maps each row to its line.

### Card

- [ ] `ComponentRow`: a skeleton row (a square, a name bar, a kind bar, two figure bars and an empty track) or a filled row with the technology icon at 16px in `currentColor`, the path, a kind chip (`api`, `broker`, `cache`, `web`), the CPU and memory figures formatted for the locale, and a two-pixel share bar whose fill grows to the row's share of the CPU envelope with a motion-safe width transition; skeletons and bars are plain Tailwind spans local to the card.
- [ ] The components section: a header with the components, CPU and memory labels; a dash row before the resolve line; five skeleton rows after it; each row filled after its line; a total that reads `n of 5 resolved` with the running sums while filling and `5 components` with the suggested totals when done, as ICU messages in five locales.
- [ ] The technology name sits in HTML beside the mark, so the icon is decorative.

### Tests

- [ ] The content module test: the component figures sum to the suggested totals; the suggested totals stay within the envelope; every component has a session line.
- [ ] The `TicketCard` test grows: after the resolve line five skeleton rows show and the total reads none resolved; after the first component line the row shows its path and figures and the total reads one of five; at the last component line the total reads five components with the suggested totals.

### Verification

- `pnpm turbo run check`, `pnpm turbo run typecheck test --filter=website --filter=@langdrift/web`, `drift check localization`, `drift check dead-code`.
- `/en/platform/microruntimes/` in light and dark: the marks follow the text color; rows fill one by one and the bars grow; 375px keeps the figures readable under the names.
