import type { RunLine, RunMessageId } from './microruntimes-workflow.types'

// True once the run has completed a line with this message.
export function runHasPlayed(
	script: readonly RunLine[],
	playedCount: number,
	message: RunMessageId,
): boolean {
	const index = script.findIndex((line) => line.message === message)

	return index >= 0 && playedCount > index
}
