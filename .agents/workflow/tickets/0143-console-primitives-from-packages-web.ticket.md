# 0143: The console uses shadcn primitives and icons from packages/web

**Spec:** .agents/workflow/specs/0026-console-rules-compliance.spec.md

**What to build:** packages/web owns shadcn Button, Card, Dialog and Skeleton, styled with token role utilities, plus ThemeToggle, ListEnd and the Phosphor icons the console uses. The console imports them in place of lib/ui, drops @phosphor-icons/react, and its focus ring moves from the body to each interactive element. This is the one approved visual change.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent
