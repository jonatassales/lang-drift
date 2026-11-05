import {
	type RunLine,
	type WorkflowChoices,
	workStageIds,
} from './microruntimes-workflow'
import { stageLines } from './stage-lines.compute'

// The dev ticket: the terminal waits while the person writes the reserved
// part, then the ticket closes without an agent. No choice turns it off,
// because it is the person's own work, not a stage a person confirms.
const devTicket: readonly RunLine[] = [
	{
		delay: 700,
		id: 'devTicketWaiting',
		message: 'devTicketWaiting',
		stage: 'dev-ticket',
		type: 'output',
	},
	{
		delay: 2400,
		id: 'devTicketClosed',
		message: 'devTicketClosed',
		stage: 'dev-ticket',
		type: 'success',
	},
]

// The agent comes back after the person: its code review opens the review
// stage.
const codeReview: RunLine = {
	delay: 800,
	id: 'codeReview',
	message: 'codeReview',
	stage: 'review',
	type: 'output',
}

// The work between the gate and the microruntime: each stage through
// implement, waiting for a person where the choices ask for one, then the
// dev ticket the person writes, then the agent's code review.
export function workLines(choices: WorkflowChoices): RunLine[] {
	const throughImplement = workStageIds.flatMap((stage) =>
		stageLines(stage, choices.human[stage]),
	)

	return [...throughImplement, ...devTicket, codeReview]
}
