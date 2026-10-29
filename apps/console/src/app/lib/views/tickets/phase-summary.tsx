import type { ReactNode } from 'react'

import type {
	Ticket,
	TicketAnalysis,
	TicketPlan,
	TicketWorkflow,
} from '@data/types'

import { type TicketFormat, useTicketFormat } from './ticket-format.hook'
import {
	analysisSummary,
	type SummaryPhase,
	workflowSummary,
} from './ticket-phases'
import { UnknownValue } from './unknown-value'

interface PhaseSummaryProps {
	phase: SummaryPhase
	ticket: Ticket
}

// The analysis's line: the latest thing it has settled.
function analysisLine(analysis: TicketAnalysis, format: TicketFormat) {
	const { t } = format
	const summary = analysisSummary(analysis)
	switch (summary.kind) {
		case 'suggestion':
			return t('analysis.suggestionSummary', {
				effort: t(`analysis.efforts.${summary.effort}`),
				model: summary.model,
			})
		case 'episodes':
			return t('analysis.episodesSummary', {
				kept: summary.kept,
				total: summary.total,
			})
		case 'documents':
			return t('analysis.documentsSummary', { count: summary.count })
		case 'unknown':
			return <UnknownValue />
	}
}

// The plan's line: a dash until it is drafted with totals, then the totals
// and whether the team has authorized them.
function planLine(plan: TicketPlan, format: TicketFormat): ReactNode {
	if (!plan.suggested) return <UnknownValue />

	return format.t('plan.summary', {
		figures: format.figures(plan.suggested),
		pending: format.authorization(plan),
	})
}

// The workflow's line: the model, effort and estimate against the budget
// while it runs, or how it ended.
function workflowLine(workflow: TicketWorkflow, format: TicketFormat) {
	const { t } = format
	const summary = workflowSummary(workflow)
	switch (summary.kind) {
		case 'paused':
			return t('workflow.summaryPaused')
		case 'done':
			return t('workflow.summaryDone', {
				number: summary.number,
				observed: format.money(summary.observedCents),
				state: t(`workflow.${summary.state}`),
			})
		case 'estimate':
			return t('workflow.summaryChoice', {
				budget: format.money(workflow.budgetCents),
				effort: t(`analysis.efforts.${workflow.effort}`),
				estimate: format.money(workflow.estimateCents),
				model: workflow.model.shortName,
			})
	}
}

// A phase's one-line summary, the header line of its panel, wherever the
// ticket appears: nothing while the ticket has not reached the phase.
export function PhaseSummary(props: PhaseSummaryProps) {
	const { phase, ticket } = props

	const format = useTicketFormat()

	switch (phase) {
		case 'analysis':
			return ticket.analysis ? analysisLine(ticket.analysis, format) : null
		case 'plan':
			return ticket.plan ? planLine(ticket.plan, format) : null
		case 'workflow':
			return ticket.workflow ? workflowLine(ticket.workflow, format) : null
	}
}
