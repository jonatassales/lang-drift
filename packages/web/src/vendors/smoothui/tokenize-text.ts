import { segmentGraphemes } from './segment-graphemes'
import type { TextMorphMode } from './text-morph'

export function tokenizeText(text: string, mode: TextMorphMode): string[] {
	if (mode === 'characters') return segmentGraphemes(text)

	return text.split(/(\s+)/).filter((token) => token.length > 0)
}
