import type { Ticket } from '@data/types'

import { analysisLine } from './analysis-line.compute'
import { planLine } from './plan-line.compute'
import { useTicketFormat } from './ticket-format.hook'
import type { SummaryPhase } from './ticket-phases'
import { workflowLine } from './workflow-line.compute'

interface PhaseSummaryProps {
	phase: SummaryPhase
	ticket: Ticket
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
