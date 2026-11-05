import { createHash, timingSafeEqual } from 'node:crypto'

// Compares the presented secret with the deployment's in constant time: the
// digests have one length whatever the secrets' lengths.
export function secretMatches(presented: string, expected: string): boolean {
	const given = createHash('sha256').update(presented).digest()
	const wanted = createHash('sha256').update(expected).digest()

	return timingSafeEqual(given, wanted)
}
