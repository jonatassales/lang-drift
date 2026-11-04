# Resend Book a Demo Integration

**Status:** ready-for-review

**Superseded:** spec 0031 moves email dispatch out of the website into the Nest microservices; the sent page's auto-redirect and the toast introduced here stay.

## Tickets

- [x] 0172: Resend email adapters for confirmation and notification
- [x] 0173: Integrate email sending in book-a-demo API route
- [x] 0174: Five-second auto-redirect on sent confirmation page
- [x] 0175: Toast error feedback and outcome wiring in BookDemoForm

## Problem Statement

Visitors who submit the "Book a demo" form on `apps/website` receive a paused notice (HTTP 503 `UNAVAILABLE`) because demo request forwarding was not implemented. Prospective clients who wish to schedule a walkthrough cannot have their contact information recorded or delivered, and internal team members are not alerted of incoming requests.

## Solution

Directly within the Next.js website application (App Router API route handler without spinning up an extra microservice), integrate with Resend to dispatch two emails upon valid form submission:
1. An applicant confirmation email sent from `no-reply@langdrift.md` notifying the requester that their submission was received and the team will be in touch shortly.
2. An internal notification email sent from `no-reply@langdrift.md` to `jonatas.langdrift@gmail.com` with the lead's submission details (company email, company name, position, locale, and submission timestamp).

On successful submission, the user is navigated to the sent confirmation page (`/early-access-sent`, the cat/bird sending confirmation screen), which displays the confirmation, waits for 5 seconds, and automatically redirects the visitor to the site's home page (`/`).

If the submission fails (network failure or server error), the form remains on screen and displays an error toast informing the user that something went wrong for now.

## User Stories

1. As a prospective customer, I want to submit the Book a Demo form with my company email, company name, and position so that I can schedule a platform walkthrough.
2. As a prospective customer, I want to receive an immediate email confirmation from `no-reply@langdrift.md` letting me know my request was received and that the LangDrift team will reach out soon.
3. As a prospective customer, upon successful demo request submission, I want to see the sent confirmation screen for 5 seconds before being smoothly redirected back to the home page.
4. As a prospective customer, if something goes wrong during submission, I want to see an error toast informing me that something went wrong for now so that I know to try again.
5. As a LangDrift team member, I want to receive an email notification at `jonatas.langdrift@gmail.com` containing all details of the demo lead whenever a form is submitted.
6. As a platform operator, I want the email dispatch to run inside the website's App Router API route using Resend and authenticated by `RESEND_API_KEY`, without maintaining an external microservice.

## Implementation Decisions

- **Resend SDK and Client:** Add `resend` to `apps/website`. Read `process.env.RESEND_API_KEY` at the boundary in `apps/website/src/infra/resend/`.
- **Single Responsibility Adapters:**
  - `send-demo-confirmation.adapter.ts`: Formats and delivers the confirmation email to the applicant from `no-reply@langdrift.md`.
  - `send-demo-lead-notification.adapter.ts`: Formats and delivers the notification email with the lead's details to `jonatas.langdrift@gmail.com` from `no-reply@langdrift.md`.
- **API Route Handling:**
  - In `apps/website/src/app/api/book-a-demo/route.ts`, upon successful request validation (`bookDemoSubmissionSchema`), call both email sending adapters.
  - Return HTTP 200 with `{ ok: true }` when emails are dispatched.
  - Return HTTP 500 with `{ ok: false, code: 'DELIVERY_FAILED' }` if email dispatch fails.
- **Sent Page Auto-Redirect:**
  - In `apps/website/src/app/[locale]/(pages)/early-access-sent/`, add a client component (`auto-redirect.tsx`) that schedules a 5-second timer (`setTimeout` with cleanup) to navigate the visitor back to `/`.
- **Toast Feedback on Error:**
  - Create a reusable, accessible Toast component in `@langdrift/web/ui/toast` adhering to design tokens.
  - In `BookDemoForm` (`apps/website/src/app/lib/components/book-a-demo/form.tsx`), on submission error, trigger the toast displaying the localized error message (`copy.form.failed`).
  - On submission success, push navigation to `/early-access-sent`.

## Testing Decisions

- Unit tests for both email adapters under `apps/website/src/infra/resend/` with mocked Resend SDK client calls.
- Unit tests for `apps/website/src/app/api/book-a-demo/route.ts` in `route.unit.ts`, testing successful 200 response upon email dispatch and 500 on delivery failure, alongside existing validation cases.
- Unit tests for `AutoRedirect` in `auto-redirect.unit.tsx` checking that `router.push('/')` is called after 5000ms.
- Unit tests for `BookDemoForm` in `form.unit.tsx` testing that success redirects to `/early-access-sent` and failure renders the error toast.
- All checks pass: `pnpm turbo run check typecheck test build`, `platform/cli/drift quality check`, `drift check localization`, `drift check biome`.

## Out of Scope

- CRM integrations or databases beyond Resend email delivery.
- Creating an external microservice.
- Modifying other forms or pages.
