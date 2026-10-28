# The console wears the dashboard's shell, switches between Microruntimes and FinOps, and shows what the stub serves

**Status:** ready-for-agent

## Problem Statement

The owner removed the docs app, because the documentation is not due yet, and removed the console, believing the copy had drifted from Exphere's. The console is the platform's management surface, so it has to come back. It should carry no Exphere brand, and it should look like the old LangDrift dashboard at the point where its sidebar opened with a product switcher. The platform now has two areas, Microruntimes and FinOps, and the console should switch between them. It should show only what its stub API, extub, serves, and keep no code that nothing uses.

## Solution

The console returns from the copy last committed here, which matches Exphere's console as it still stands upstream plus this repository's adaptations. Every Exphere mark in it becomes LangDrift. Its sidebar opens with an area switcher in the dashboard's style, listing Microruntimes and FinOps; the URL decides the area. Microruntimes keeps every screen the console had: home, projects, runtimes, executions, history and secrets. FinOps opens a cost overview built only from the stub's runtime readings, where an unknown cost stays unknown. The header gains the dashboard's account menu. Routes that only redirected, the dashboard's notifications, voice and avatar, and the docs app with its token theme all go. The stub's seeded project is renamed so it no longer shares the Microruntimes area's name.

## User Stories

1. As a team member, I want the console to carry the LangDrift brand everywhere, so that nothing names a product we do not sell.
2. As a team member, I want the sidebar to open with an area switcher showing Microruntimes and FinOps, so that I always know which part of the platform I am in.
3. As a team member, I want the area to follow the URL, so that a link or a reload lands in the same area.
4. As a keyboard user, I want the area switcher to open, close with Escape and return focus, so that I can switch without a mouse.
5. As a team member in Microruntimes, I want every screen the console had, so that I lose nothing I used.
6. As a team member in FinOps, I want the known LLM cost of every runtime summed, so that I see what agent work has cost.
7. As a team member in FinOps, I want executions without a measured cost counted as unknown, so that missing evidence never reads as zero.
8. As a team member in FinOps, I want each runtime's CPU and memory against its envelope and its LLM budget, so that I see where resources and money go.
9. As a team member in FinOps, I want to ask for a runtime's recommendation from the overview, so that I act where I read.
10. As a team member, I want the account menu in the header to show my email and sign me out, so that the session is one click away.
11. As a team member, I want the theme and the language in the header, so that I set them without leaving the page.
12. As a team member, I want the stub's seeded project to have a name of its own, so that it does not read as the Microruntimes area.
13. As a developer, I want no route that only redirects and no component without data, so that the console holds only what it uses.
14. As a developer, I want the docs app and its token theme gone with their commands, so that the repository holds only what runs.

## Implementation Decisions

- **Starting point.** The console is restored from the last commit, which already maps Exphere's projects onto the stub's products and Exphere's design onto LangDrift tokens. The upstream console has no later change, so a fresh copy would only undo those adaptations.
- **Brand.** The Exphere logo file goes; the navigation and the sign-in form use the LangDrift brand from the shared web package. The app icon's title, the locale cookie and every catalog string say LangDrift, and the project folder the agent creates is `.langdrift`.
- **Area switcher.** A sidebar control styled after the dashboard's product switcher: a badge with the area's icon, a small "Area" label and the area name, opening a menu of the two areas. Routes under `/finops` are FinOps; every other route is Microruntimes, so today's URLs stay. The URL is the only state.
- **Navigation.** Microruntimes lists Home and Projects and, when a project is in context, its runtimes, executions, history and secrets. FinOps lists its cost overview. Settings sits at the foot. The project item that listed runtimes is labelled Runtimes.
- **Header.** The account control becomes a menu showing the session email, a Settings link and Sign out. Theme and language stay in the header. Notifications, voice and an avatar are not built, because the stub serves none of them.
- **FinOps overview.** A server-rendered page reads the runtime list and the products from the stub. It shows the summed known LLM cost, the count of executions with an unknown cost, the CPU and memory used against each envelope, each runtime's LLM budget when configured, each runtime's existing FinOps reading, and the existing recommendation action. Illustrative and stub notices stay.
- **Dead code.** Redirect-only routes go, and the sidebar links straight to the secrets page. Whatever the dead-code check reports afterwards goes or is justified in the allow list.
- **Docs.** The docs app, its launcher, its preview configuration, its token theme and export, and its mentions in AGENTS.md and the context file go. Two root launchers remain: the website and the console.
- **Stub.** Only the seeded project's display name changes, to LangDrift Studio; its id and every route stay.

## Testing Decisions

- A good test observes what a user sees or what a URL decides, not markup structure.
- **Seams.** The console's existing Node test runner for view logic, its type check and build, the stub's Go tests, the catalog parity of the console's four locales, and the repository's dead-code, localization and history checks.
- **New tests.** The area is derived from a pathname; the cost totals add only measured costs and count the unknown ones.
- **Prior art.** The console's runtime view tests.

## Out of Scope

- A sign-in service of its own, an auth provider, notifications and a voice assistant.
- New stub routes; the console shows only what the stub serves today.
- The documentation, which returns later in its own spec.
- Vision, which the dashboard's shell came from but which does not return.

## Further Notes

- The dashboard's shell was taken from langdrift commit 9f97e5e0, the last with a product switcher in its sidebar.
- Upstream Exphere's console was last changed on 2026-10-21; this repository imported it on 2026-10-28.
