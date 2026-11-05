import type { QuestionKind } from './judge.types'

export interface JudgmentLog {
	cents: number
	inputTokens: number
	kind: QuestionKind
	model: string
	outputTokens: number
	promptVersion: string
	refused: boolean
}

// One line per judgment on the console's log, which Vercel keeps: what it
// consumed and cost, never the text of the sources or of the answer.
export function logJudgment(entry: JudgmentLog): void {
	console.info(
		JSON.stringify({
			event: 'judgment',
			...entry,
		}),
	)
}
