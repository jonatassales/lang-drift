import assert from 'node:assert/strict'
import { describe, test } from 'node:test'

import type {
	SafeSeatState,
	Ticket,
	TicketPlan,
	TicketReservedPart,
	TicketSafeSeat,
} from '@data/types'

import { safeSeatView } from './safe-seat-view.map.ts'

const reservedPart: TicketReservedPart = {
	paths: [
		'services/signup-api/domain/admission.go',
		'services/signup-api/domain/admission_test.go',
	],
	goal: 'Write the admission rule, with its tests.',
	category: 'core-domain',
	chosenBy: 'judge',
}

const tones: Record<SafeSeatState, string> = {
	met: 'success',
	'needs-override': 'danger',
	overridden: 'warning',
	'not-applicable': 'muted',
	unknown: 'muted',
}

function ticket(plan: TicketPlan | null, safeSeat: TicketSafeSeat | null) {
	return {
		id: 'ticket-0031',
		plan,
		safeSeat,
	} as Ticket
}

describe('console/lib/views/tickets/safe-seat-view.map', () => {
	test('every state has its key and tone', () => {
		for (const [state, tone] of Object.entries(tones)) {
			const record: TicketSafeSeat = {
				state: state as SafeSeatState,
				overrideReason: null,
				illustrative: true,
			}

			const view = safeSeatView(ticket(null, record))

			assert.equal(view.stateKey, `safeSeat.states.${state}`)
			assert.equal(view.tone, tone)
			assert.equal(view.illustrative, true)
		}
	})

	test('a ticket without a record reads as unknown, never as met', () => {
		const view = safeSeatView(ticket(null, null))

		assert.equal(view.stateKey, 'safeSeat.states.unknown')
		assert.equal(view.tone, 'muted')
		assert.equal(view.overrideReason, null)
		assert.equal(view.illustrative, false)
	})

	test('an override shows its reason, or none when it has none', () => {
		const reason = 'The person was on call this week.'
		const cases: [string | null, string | null][] = [
			[reason, reason],
			[null, null],
			['  ', null],
		]

		for (const [overrideReason, shown] of cases) {
			const record: TicketSafeSeat = {
				state: 'overridden',
				overrideReason,
				illustrative: true,
			}

			assert.equal(safeSeatView(ticket(null, record)).overrideReason, shown)
		}
	})

	test('a reason appears only with an override', () => {
		const record: TicketSafeSeat = {
			state: 'met',
			overrideReason: 'A stray reason.',
			illustrative: true,
		}

		assert.equal(safeSeatView(ticket(null, record)).overrideReason, null)
	})

	test('a plan shows its reserved part, or none when it has none', () => {
		const plan = { reservedPart } as TicketPlan
		const withoutPart = { reservedPart: null } as TicketPlan

		assert.deepEqual(
			safeSeatView(ticket(plan, null)).reservedPart,
			reservedPart,
		)
		assert.equal(safeSeatView(ticket(withoutPart, null)).reservedPart, null)
		assert.equal(safeSeatView(ticket(null, null)).reservedPart, null)
	})
})
