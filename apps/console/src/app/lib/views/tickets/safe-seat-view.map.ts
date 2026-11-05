import type { SafeSeatState, Ticket } from '@data/types'

import type { SafeSeatView } from './safe-seat-view.type.ts'

const tones: Record<SafeSeatState, SafeSeatView['tone']> = {
	met: 'success',
	'needs-override': 'danger',
	overridden: 'warning',
	'not-applicable': 'muted',
	unknown: 'muted',
}

// How a ticket's Safe Seat record reads on its cards. A ticket without a
// record reads as unknown, never as met; only an override shows its
// reason, and only when it has one; a plan without a reserved part shows
// none, not an empty one.
export function safeSeatView(ticket: Ticket): SafeSeatView {
	const record = ticket.safeSeat
	const state = record?.state ?? 'unknown'
	const reason = record?.overrideReason?.trim() || null
	const overrideReason = state === 'overridden' ? reason : null

	return {
		stateKey: `safeSeat.states.${state}`,
		tone: tones[state],
		overrideReason,
		reservedPart: ticket.plan?.reservedPart ?? null,
		illustrative: record?.illustrative ?? false,
	}
}
