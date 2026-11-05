import type {
	RunLine,
	StageProgress,
	WorkflowStageId,
} from './microruntimes-workflow.data'

// Where a stage stands in a playing run: pending until its first line
// prints, active while its lines print, done after its last.
export function stageProgress(
	script: readonly RunLine[],
	playedCount: number,
	stage: WorkflowStageId,
): StageProgress {
	const indexes = script
		.map((line, index) => (line.stage === stage ? index : -1))
		.filter((index) => index >= 0)
	const first = indexes.at(0)
	const last = indexes.at(-1)

	const missingBounds = first === undefined || last === undefined

	if (missingBounds) return 'pending'
	if (playedCount > last) return 'done'
	if (playedCount > first) return 'active'

	return 'pending'
}
