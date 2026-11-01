import type { MicroruntimeSession } from '@app/atoms/microruntime-session.atom'

import type { TicketCardCopy } from './ticket-card-format'

export function playLabel(
	workflow: TicketCardCopy['workflow'],
	session: MicroruntimeSession,
	running: boolean,
) {
	if (session.stage === 'launching') return workflow.starting
	if (running) return workflow.playing
	if (session.stage === 'running') return workflow.playAgain

	return workflow.play
}
