# 0146: The console's control flow and external input follow rule 00

**Spec:** .agents/workflow/specs/0026-console-rules-compliance.spec.md

**What to build:** Chained ternaries and the breadcrumb else-if chain become guards, a switch or a lookup, and compound predicates are named. Catches stop discarding errors and keep today's outcome. The try/catch around a pure function goes, and EXTUB_URL is read at its use without a fallback. No visual or behavioural change.

**Blocked by:** 0145.

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).
