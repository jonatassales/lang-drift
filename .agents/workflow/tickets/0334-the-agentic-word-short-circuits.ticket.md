# 0334: The word "Agentic" short-circuits

**Spec:** none. The owner asked for the effect directly, as a visual experiment on the home hero, and approved the first version after trying two others.

**What to build:** On the home hero, the first word of the headline, "Agentic", reads as if it were in a short circuit. It trembles in place all the time, and at irregular moments (every 0.45 to 1.75 seconds) it jerks and flickers while one to three jagged bolts arc out of each side of it, redrawn three times in about 200 ms before they vanish. The bolts have a white core and a blue glow. The lines below the word keep their motion. The rotating word now cycles through DevEx, LLMOps, FinOps, Runtime and CI Gates: the owner added LLMOps and CI Gates to the same change.

**Blocked by:** none

**Status:** completed

**Alignment:** none. The effect changes the website's home hero, not the console's runtime view, so ADR 0005's simulation does not apply.

- [x] The first word of the home hero headline trembles and, at irregular moments, jerks while bolts arc out of both sides of it, in every locale.
- [x] The rotating word and "Platform" below it keep their motion, and no bolt reaches them.
- [x] The rotating word cycles DevEx, LLMOps, FinOps, Runtime and CI Gates in every locale, and "CI Gates Platform" fits a 375px screen without horizontal scroll.
- [x] The bolts use the tokens' roles (`rainbow-5`, `primary-foreground`, white in both themes) and no authored CSS.
- [x] Under reduced motion, or where the browser cannot animate an element, the word holds still and no bolt is drawn; a hidden tab skips the discharges.
- [x] Server and client render the word without bolts first, so hydration matches.

## Dev part

**Override:** No part qualifies: a presentational effect with no domain decision of its own. The owner asked for it and chose the version. 2026-11-05.
