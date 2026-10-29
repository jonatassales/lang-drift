import assert from 'node:assert/strict'
import test from 'node:test'

import type { Ticket, TicketPlan, TicketWorkflow } from '@data/types'

import { compactTokens, dollars } from './ticket-format.ts'
import {
	analysisSummary,
	componentShare,
	keptEpisodes,
	listedPhase,
	phaseReached,
	phaseStep,
	planFigures,
	summaryPhases,
	ticketPath,
	ticketPhases,
	ticketSummary,
	workflowSummary,
} from './ticket-phases.ts'

function ticket(values: Partial<Ticket>): Ticket {
	return {
		id: 'ticket-0031',
		number: '0031',
		slug: 'early-access-runtime',
		phase: 'ticket',
		analysis: null,
		plan: null,
		workflow: null,
		...values,
	} as Ticket
}

const plan: TicketPlan = {
	bounded: {
		cpuMillis: 2000,
		memoryMiB: 3072,
	},
	components: [],
	suggested: {
		cpuMillis: 1450,
		memoryMiB: 2304,
	},
	validationErrors: 0,
	revision: 1,
	authorization: 'pending',
}

test('a ticket has reached its own phase and every phase before it', () => {
	const planned = ticket({ phase: 'plan' })
	assert.equal(phaseReached(planned, 'ticket'), true)
	assert.equal(phaseReached(planned, 'analysis'), true)
	assert.equal(phaseReached(planned, 'plan'), true)
	assert.equal(phaseReached(planned, 'workflow'), false)
})

test('an entry summarizes its listing phase, or the current one in Tickets', () => {
	const triaged = ticket({ phase: 'ticket' })
	const analyzed = ticket({ phase: 'analysis' })
	const running = ticket({ phase: 'workflow' })
	assert.equal(listedPhase(triaged, 'ticket'), null)
	assert.equal(listedPhase(analyzed, 'ticket'), 'analysis')
	assert.equal(listedPhase(running, 'ticket'), 'workflow')
	assert.equal(listedPhase(running, 'analysis'), 'analysis')
	assert.equal(listedPhase(running, 'plan'), 'plan')
	assert.equal(listedPhase(running, 'workflow'), 'workflow')
})

test('only the analysis, plan and workflow phases have a summary of their own', () => {
	assert.deepEqual(summaryPhases, ['analysis', 'plan', 'workflow'])
	assert.deepEqual(ticketPhases, ['ticket', ...summaryPhases])
})

test('a ticket counts the phases it has reached out of four', () => {
	assert.deepEqual(phaseStep(ticket({ phase: 'ticket' })), {
		reached: 1,
		total: 4,
	})
	assert.deepEqual(phaseStep(ticket({ phase: 'workflow' })), {
		reached: 4,
		total: 4,
	})
})

test('an entry opens the ticket page, at its phase in a phase listing', () => {
	const planned = ticket({
		id: 'ticket-0102',
		phase: 'plan',
		productId: 'seed-microruntimes',
	})
	const page = '/products/seed-microruntimes/tickets/ticket-0102'
	assert.equal(ticketPath(planned, 'ticket'), page)
	assert.equal(ticketPath(planned, 'analysis'), `${page}#analysis`)
	assert.equal(ticketPath(planned, 'plan'), `${page}#plan`)
	assert.equal(ticketPath(planned, 'workflow'), `${page}#workflow`)
})

test('the analysis summary names its latest settled step, or nothing yet', () => {
	assert.deepEqual(
		analysisSummary({
			documentsRead: ['prd', 'designDoc'],
			episodes: [],
			model: null,
			effort: null,
		} as unknown as Ticket['analysis'] & object),
		{
			kind: 'documents',
			count: 2,
		},
	)
	assert.deepEqual(
		analysisSummary({
			documentsRead: ['prd'],
			episodes: [],
			model: {
				id: 'claude-sonnet-5-5',
				name: 'Claude Sonnet 5.5',
			},
			effort: 'medium',
		} as unknown as Ticket['analysis'] & object),
		{
			kind: 'suggestion',
			model: 'Claude Sonnet 5.5',
			effort: 'medium',
		},
	)
	assert.deepEqual(
		analysisSummary({
			documentsRead: [],
			episodes: null,
			model: null,
			effort: null,
		} as unknown as Ticket['analysis'] & object),
		{ kind: 'unknown' },
	)
})

test('a plan without suggested totals has no headroom, never a zero one', () => {
	assert.deepEqual(planFigures(plan), {
		suggested: {
			cpuMillis: 1450,
			memoryMiB: 2304,
		},
		headroom: {
			cpuMillis: 550,
			memoryMiB: 768,
		},
		cpuShare: 0.725,
		memoryShare: 0.75,
	})
	assert.deepEqual(
		planFigures({
			...plan,
			suggested: null,
		}),
		{
			suggested: null,
			headroom: null,
			cpuShare: 0,
			memoryShare: 0,
		},
	)
})

test('the workflow summary follows the outcome', () => {
	const workflow = {
		outcome: 'readyForReview',
		observedCents: 980,
		pullRequest: {
			number: 19,
		},
	} as TicketWorkflow
	assert.deepEqual(workflowSummary(workflow), {
		kind: 'done',
		number: 19,
		observedCents: 980,
		state: 'readyForReview',
	})
	assert.deepEqual(
		workflowSummary({
			...workflow,
			outcome: 'paused',
		}),
		{ kind: 'paused' },
	)
	assert.deepEqual(
		workflowSummary({
			...workflow,
			outcome: 'running',
			pullRequest: null,
			observedCents: null,
		}),
		{ kind: 'estimate' },
	)
})

test('the ticket summary is its number and slug', () => {
	assert.equal(
		ticketSummary({
			number: '0031',
			slug: 'early-access-runtime',
		}),
		'0031 · early-access-runtime',
	)
})

test('a component takes its share of the CPU envelope', () => {
	assert.equal(componentShare({ cpuMillis: 500 }, plan.bounded), 0.25)
	assert.equal(componentShare({ cpuMillis: 250 }, plan.bounded), 0.125)
})

test('kept episodes are counted against all of them', () => {
	assert.deepEqual(
		keptEpisodes([
			{ verdict: 'kept' },
			{ verdict: 'dropped' },
			{ verdict: 'kept' },
		] as never),
		{
			kept: 2,
			total: 3,
		},
	)
})

test('money comes from cents and tokens from thousands', () => {
	assert.equal(dollars('en', 950), '$9.50')
	assert.equal(dollars('pt-br', 2500).replace(/\s/g, ' '), 'US$ 25,00')
	assert.equal(compactTokens('en', 3500), '3.5M')
	assert.equal(compactTokens('en', 250), '250K')
})
