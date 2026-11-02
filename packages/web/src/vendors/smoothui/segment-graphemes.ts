export function segmentGraphemes(value: string): string[] {
	if (typeof Intl === 'undefined' || !('Segmenter' in Intl)) return [...value]

	const segmenter = new Intl.Segmenter(undefined, {
		granularity: 'grapheme',
	})

	return Array.from(segmenter.segment(value), (entry) => entry.segment)
}
