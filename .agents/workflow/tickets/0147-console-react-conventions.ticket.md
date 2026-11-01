# 0147: The console's components follow rule 05

**Spec:** .agents/workflow/specs/0026-console-rules-compliance.spec.md

**What to build:** Props interfaces sit directly above their components, scope follows the 05.5 order, and handlers are named. Loading flags become transitions, and use client sits at the narrowest boundary. Derivable effects move into render or handlers, and the remaining effects say why. Login and signup compose their own views, and template.tsx files return their children directly. No visual or behavioural change.

**Blocked by:** 0146.

**Status:** ready-for-agent
