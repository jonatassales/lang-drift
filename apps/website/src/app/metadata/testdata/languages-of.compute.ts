import type { Metadata } from 'next'

export function languagesOf(metadata: Metadata) {
	return Object.fromEntries(
		Object.entries(metadata.alternates?.languages ?? {}).map(
			([language, href]) => [language, String(href)],
		),
	)
}
