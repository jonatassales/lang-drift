import assert from 'node:assert/strict'
import { describe, test } from 'node:test'

import type { TicketStage } from '@data/types'

import { stageMarker } from './stage-marker.compute.ts'

function stage(values: Partial<TicketStage>): TicketStage {
	return {
		id: 'implement',
		skills: null,
		human: false,
		drivenBy: 'agent',
		progress: 'pending',
		...values,
	}
}

describe('console/lib/views/tickets/stage-marker.compute', () => {
	test('a person-driven stage is marked apart from an approved one', () => {
		const written = stage({
			id: 'dev-ticket',
			drivenBy: 'person',
		})

		assert.equal(stageMarker(written), 'person')
	})

	test('a human-approved stage carries the approval mark', () => {
		const review = stage({
			id: 'review',
			human: true,
		})

		assert.equal(stageMarker(review), 'approval')
	})

	test('an agent stage carries no mark, and merge keeps its own lock', () => {
		const merge = stage({
			id: 'merge',
			human: true,
		})

		assert.equal(stageMarker(stage({})), null)
		assert.equal(stageMarker(merge), null)
	})
})
