# 0167: Adopt Apache 2.0 license, add NOTICE and SECURITY.md, and update website and console copyright

**What to build:**
1. Transition repository root and packages/ex licenses to Apache License Version 2.0.
2. Add root NOTICE file per Apache 2.0 section 4(d) with third-party attributions and component license mappings.
3. Add root SECURITY.md policy tailored for LangDrift FinOps governance and microruntime sandboxing.
4. Update website footer with Security and License links, retaining Brazilian attribution.
5. Update console navigation sidebar with subtle copyright notice and Apache 2.0 declaration.
6. Update platform quality detectors and filesystem rules to recognize NOTICE and SECURITY.md.

**Blocked by:** none

**Status:** completed

- [x] LICENSE and packages/ex/LICENSE carry Apache License Version 2.0.
- [x] packages/ex/package.json declares Apache-2.0.
- [x] NOTICE file added at repository root.
- [x] SECURITY.md added at repository root with security@langdrift.md and FinOps/Microruntimes scope.
- [x] Website footer renders Security and License links across all locales.
- [x] Console navigation renders copyright and Apache 2.0 notice across all locales.
- [x] Quality detectors and 01-filesystem rule recognize NOTICE and SECURITY.md.
- [x] All checks, typecheck, and test suites pass.
