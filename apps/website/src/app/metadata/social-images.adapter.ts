import type { ResolvingMetadata } from 'next'

export async function socialImages(parent: ResolvingMetadata) {
	const resolved = await parent

	return {
		openGraph: resolved.openGraph?.images,
		twitter: resolved.twitter?.images,
	}
}
