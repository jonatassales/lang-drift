import type { SafeSeatState, TicketReservedPart } from '@data/types'

// What the ticket cards show of a ticket's Safe Seat record: the state's
// message key under tickets and its status tone, the override reason, the
// reserved part and whether the record is illustrative stub data.
export interface SafeSeatView {
	stateKey: `safeSeat.states.${SafeSeatState}`
	tone: 'success' | 'warning' | 'danger' | 'muted'
	overrideReason: string | null
	reservedPart: TicketReservedPart | null
	illustrative: boolean
}
