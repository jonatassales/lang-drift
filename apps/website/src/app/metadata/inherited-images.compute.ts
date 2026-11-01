import type { ResolvingMetadata } from 'next'

export interface InheritedSocialImages {
	openGraph?: NonNullable<Awaited<ResolvingMetadata>['openGraph']>['images']
	twitter?: NonNullable<Awaited<ResolvingMetadata>['twitter']>['images']
}

export function inheritedImages(
	images: InheritedSocialImages['openGraph'] | InheritedSocialImages['twitter'],
) {
	if (!images) return
	if (Array.isArray(images) && images.length === 0) return

	return images
}
