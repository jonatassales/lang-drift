import type { TextMorphMode, TextMorphToken } from './text-morph'
import { tokenizeText } from './tokenize-text'

export function buildTokens(
	text: string,
	mode: TextMorphMode,
): TextMorphToken[] {
	const occurrences = new Map<string, number>()

	return tokenizeText(text, mode).map((value) => {
		const count = occurrences.get(value) ?? 0
		occurrences.set(value, count + 1)

		return {
			key: `${value}-${count}`,
			value,
		}
	})
}
