import {
	type RunLine,
	stageWork,
	type workStageIds,
} from './microruntimes-workflow'

// One stage: the wait for the owner when the stage asks for a person, then
// the stage's own work.
export function stageLines(
	stage: (typeof workStageIds)[number],
	human: boolean,
): RunLine[] {
	const confirmation: RunLine[] = human
		? [
				{
					delay: 600,
					id: `waiting-${stage}`,
					message: 'waiting',
					stage,
					type: 'output',
					values: { stage },
				},
				{
					delay: 1600,
					id: `confirmed-${stage}`,
					message: 'confirmed',
					stage,
					type: 'success',
					values: { stage },
				},
			]
		: []
	const work = stageWork[stage].map((line) => ({
		...line,
		id: line.message,
		stage,
	}))

	return [...confirmation, ...work]
}
