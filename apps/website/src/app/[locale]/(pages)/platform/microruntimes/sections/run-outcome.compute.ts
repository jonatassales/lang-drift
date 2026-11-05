import type { MicroruntimeSession } from '@app/atoms/microruntime-session.atom'
import type {
	RunLine,
	WorkflowChoices,
} from '@routes/platform/microruntimes/content/microruntimes-workflow'

export interface RunOutcome {
	choices: WorkflowChoices
	playedCount: number
	script: readonly RunLine[]
}

export function runOutcome(
	session: MicroruntimeSession,
	script: readonly RunLine[],
): RunOutcome | null {
	const { ran, stage } = session
	if (stage !== 'running' || !ran) return null

	return {
		choices: ran,
		playedCount: session.playedCount,
		script,
	}
}
