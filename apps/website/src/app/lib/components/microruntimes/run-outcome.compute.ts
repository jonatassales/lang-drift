import type { MicroruntimeSession } from '@app/atoms/microruntime-session.atom'
import type { RunLine, WorkflowChoices } from '@data/microruntimes-workflow'

export interface RunOutcome {
	choices: WorkflowChoices
	playedCount: number
	script: readonly RunLine[]
}

export function runOutcome(
	session: MicroruntimeSession,
	script: readonly RunLine[],
): RunOutcome | null {
	if (session.stage !== 'running' || !session.ran) return null

	return {
		choices: session.ran,
		playedCount: session.playedCount,
		script,
	}
}
