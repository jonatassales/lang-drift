import { describe, expect, it } from 'vitest'

import {
	initialWorkflowChoices,
	type WorkflowChoices,
} from './microruntimes-workflow'
import { workLines } from './work-lines.compute'

const everyStageOff: WorkflowChoices = {
	...initialWorkflowChoices,
	human: {
		planning: false,
		spec: false,
		tickets: false,
		implement: false,
		review: false,
	},
}

const everyStageOn: WorkflowChoices = {
	...initialWorkflowChoices,
	human: {
		planning: true,
		spec: true,
		tickets: true,
		implement: true,
		review: true,
	},
}

const cases = [
	{
		name: 'the default choices',
		choices: initialWorkflowChoices,
	},
	{
		name: 'every switchable stage off',
		choices: everyStageOff,
	},
	{
		name: 'every switchable stage on',
		choices: everyStageOn,
	},
]

describe('website/platform/microruntimes/content/work-lines.compute', () => {
	it('waits for the planning owner, works through implement, then hands the dev ticket to the person before the code review', () => {
		const messages = workLines(initialWorkflowChoices).map(
			(line) => line.message,
		)

		expect(messages).toEqual([
			'waiting',
			'confirmed',
			'grilling',
			'toSpec',
			'toTickets',
			'hook',
			'implement',
			'devTicketWaiting',
			'devTicketClosed',
			'codeReview',
		])
	})

	it('runs straight through when no stage asks for a person, and still stops at the dev ticket', () => {
		const messages = workLines(everyStageOff).map((line) => line.message)

		expect(messages).toEqual([
			'grilling',
			'toSpec',
			'toTickets',
			'hook',
			'implement',
			'devTicketWaiting',
			'devTicketClosed',
			'codeReview',
		])
	})

	it('waits for a person before every stage that asks for one', () => {
		const ids = workLines(everyStageOn).map((line) => line.id)

		expect(ids).toEqual([
			'waiting-planning',
			'confirmed-planning',
			'grilling',
			'waiting-spec',
			'confirmed-spec',
			'toSpec',
			'waiting-tickets',
			'confirmed-tickets',
			'toTickets',
			'waiting-implement',
			'confirmed-implement',
			'hook',
			'implement',
			'devTicketWaiting',
			'devTicketClosed',
			'codeReview',
		])
	})

	it.each(cases)(
		'ends with the dev ticket right after implement and the code review after it, with $name',
		({ choices }) => {
			const messages = workLines(choices).map((line) => line.message)
			const implement = messages.indexOf('implement')

			expect(messages.slice(implement + 1)).toEqual([
				'devTicketWaiting',
				'devTicketClosed',
				'codeReview',
			])
		},
	)

	it.each(cases)(
		'tags every line with its stage, with $name',
		({ choices }) => {
			const lines = workLines(choices)
			const stageById = Object.fromEntries(
				lines.map((line) => [line.id, line.stage]),
			)

			expect(lines.every((line) => line.stage !== undefined)).toBe(true)
			expect(stageById).toMatchObject({
				implement: 'implement',
				devTicketWaiting: 'dev-ticket',
				devTicketClosed: 'dev-ticket',
				codeReview: 'review',
			})
		},
	)

	it.each(cases)('gives every line a unique id, with $name', ({ choices }) => {
		const ids = workLines(choices).map((line) => line.id)

		expect(new Set(ids).size).toBe(ids.length)
	})
})
