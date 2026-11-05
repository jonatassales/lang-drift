import type { Metadata } from 'next'

export function canonicalOf(metadata: Metadata) {
	return String(metadata.alternates?.canonical)
}
