import { describe, expect, it } from 'vitest'

import { microruntimesMessages } from '@i18n/messages/index'
import { choices } from '@routes/platform/microruntimes/content/testdata/choices.mock'

import {
	budgetCents,
	initialWorkflowChoices,
	type RunMessageId,
	runSession,
	stageProgress,
	workflowEfforts,
	workflowFigures,
	workflowModels,
} from './microruntimes-workflow'

describe('website/platform/microruntimes/content/microruntimes-workflow', () => {
	it('prices each model and effort at list prices, in cents', () => {
		const estimates = workflowModels.map((model) =>
			workflowEfforts.map(
				(effort) =>
					workflowFigures(
						choices({
							effort,
							model: model.id,
						}),
					).estimate,
			),
		)
		expect(estimates).toEqual([
			[1040, 1900, 3400],
			[520, 950, 1700],
			[26, 48, 85],
		])
	})

	it('reserves a fifth above the estimate and releases what the usage record did not use', () => {
		expect(workflowFigures(initialWorkflowChoices)).toEqual({
			estimate: 950,
			reserved: 1140,
			observed: 980,
			released: 160,
			overBudget: false,
		})
	})

	it('breaks the budget only with Opus 5.5 at high effort', () => {
		const over = workflowModels.flatMap((model) =>
			workflowEfforts
				.filter(
					(effort) =>
						workflowFigures(
							choices({
								effort,
								model: model.id,
							}),
						).overBudget,
				)
				.map((effort) => `${model.id} ${effort}`),
		)
		expect(over).toEqual(['claude-opus-5-5 high'])
		expect(budgetCents).toBe(2500)
	})

	it('waits for the owner only at the stages chosen for a person', () => {
		const ids = runSession(initialWorkflowChoices).map((line) => line.id)
		expect(ids).toContain('waiting-planning')
		expect(ids).toContain('confirmed-planning')
		expect(ids).not.toContain('waiting-implement')
		expect(ids.at(-2)).toBe('ready')

		const everyone = runSession(
			choices({
				human: {
					implement: true,
					planning: true,
					review: true,
					spec: true,
					tickets: true,
				},
			}),
		).map((line) => line.id)
		expect(everyone).toContain('waiting-implement')
	})

	it('stops at the gate with nothing pushed when the reservation breaks the budget', () => {
		const script = runSession(
			choices({
				effort: 'high',
				model: 'claude-opus-5-5',
			}),
		)
		const messages = script.map((line) => line.message)
		expect(messages).toContain('gatePaused')
		expect(messages).not.toContain('pushed')
		expect(messages).not.toContain('pullRequest')
		expect(script.at(-1)?.message).toBe('paused')
	})

	it('ends to merge when the review is not the owner’s', () => {
		const script = runSession(
			choices({
				human: {
					...initialWorkflowChoices.human,
					review: false,
				},
			}),
		)
		expect(script.map((line) => line.message)).toContain('readyToMerge')
	})

	it('types one command and gives every line a unique id', () => {
		const script = runSession(initialWorkflowChoices)
		expect(script.filter((line) => line.type === 'command')).toHaveLength(1)
		expect(new Set(script.map((line) => line.id)).size).toBe(script.length)
	})

	it('never names a FIN criterion as the gate', () => {
		for (const messages of Object.values(microruntimesMessages)) {
			expect(messages.session.terminal.run.gateOpened).not.toMatch(/FIN-/)
			expect(messages.session.terminal.run.gatePaused).not.toMatch(/FIN-/)
		}
	})

	it('stops at the dev ticket after implement and reviews after the person', () => {
		const messages = runSession(initialWorkflowChoices).map(
			(line) => line.message,
		)
		const at = (message: RunMessageId) => messages.indexOf(message)

		expect(at('devTicketWaiting')).toBeGreaterThan(at('implement'))
		expect(at('devTicketClosed')).toBe(at('devTicketWaiting') + 1)
		expect(at('codeReview')).toBeGreaterThan(at('devTicketClosed'))
		expect(at('safeSeatMet')).toBeGreaterThan(at('checksPassed'))
		expect(at('safeSeatMet')).toBeLessThan(at('readyForReview'))
	})

	it('keeps the dev ticket when no stage asks for a person', () => {
		const messages = runSession(
			choices({
				human: {
					implement: false,
					planning: true,
					review: false,
					spec: false,
					tickets: false,
				},
			}),
		).map((line) => line.message)

		expect(messages).toContain('devTicketWaiting')
		expect(messages).toContain('devTicketClosed')
		expect(messages).toContain('safeSeatMet')
	})

	it('follows the dev ticket from pending to done', () => {
		const script = runSession(initialWorkflowChoices)
		const first = script.findIndex((line) => line.stage === 'dev-ticket')

		expect(stageProgress(script, first, 'dev-ticket')).toBe('pending')
		expect(stageProgress(script, first + 1, 'dev-ticket')).toBe('active')
		expect(stageProgress(script, script.length, 'dev-ticket')).toBe('done')
	})

	it('never prints a share in the Safe Seat lines', () => {
		for (const messages of Object.values(microruntimesMessages)) {
			const { run } = messages.session.terminal

			for (const line of [
				run.devTicketWaiting,
				run.devTicketClosed,
				run.safeSeatMet,
			]) {
				expect(line).not.toMatch(/%/)
			}
		}
	})

	it('follows each stage from pending to done', () => {
		const script = runSession(initialWorkflowChoices)
		const first = script.findIndex((line) => line.stage === 'planning')
		expect(stageProgress(script, first, 'planning')).toBe('pending')
		expect(stageProgress(script, first + 1, 'planning')).toBe('active')
		expect(stageProgress(script, script.length, 'planning')).toBe('done')
	})
})
