# 0084: Login and sign-up wear the website's background

**Spec:** .agents/specs/0013-console-areas.spec.md

**What to build:** The console's login and sign-up screens stand on the website's hero shader in its muted tone across the full screen, a faint purple and peach in the light scheme and grays in the dark one, with only the LangDrift brand above the form, the form inside the website's elevated glass frame, the language switcher at the top left and the theme toggle at the top right, as in the website's header. The title stays for screen readers only, and the intro line and the stub notice leave these screens. The shader's muted tone and the frame become shared components of the web package, so the website's hero and the console use the same ones.

**Blocked by:** None (can start immediately)

**Status:** wontfix. The owner closed it on 2026-11-05: superseded by ADR 0005 (Perfect Simulation Alignment).

- [ ] Login and sign-up show the muted shader and the elevated frame in both color schemes, on desktop and on a phone, and still sign in against the stub.
- [ ] The website's hero looks the same, now built from the shared pieces.
- [ ] Markup tests cover the shader's muted tone and the frame's merged classes.

### Context

- [The spec](/.agents/specs/0013-console-areas.spec.md)
