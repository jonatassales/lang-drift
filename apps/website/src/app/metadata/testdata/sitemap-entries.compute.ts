interface SitemapEntry {
	loc: string
	alternates: Record<string, string>
}

export function sitemapEntries(xml: string): SitemapEntry[] {
	return Array.from(xml.matchAll(/<url>([\s\S]*?)<\/url>/g), ([, body]) => ({
		loc: body.match(/<loc>(.*?)<\/loc>/)?.[1] ?? '',
		alternates: Object.fromEntries(
			Array.from(
				body.matchAll(/hreflang="([^"]+)" href="([^"]+)"/g),
				([, language, href]) => [language, href],
			),
		),
	}))
}
