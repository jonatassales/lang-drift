import type { PageId } from '@app/metadata/page-path.compute'
import { pageStructuredData } from '@app/metadata/page-structured-data.map'
import { serializeJsonLd } from '@app/metadata/serialize-json-ld.compute'
import type { WebsiteLocale } from '@i18n/routing'

interface WebsiteStructuredDataProps {
	page: PageId
	locale: WebsiteLocale
}

export function WebsiteStructuredData(props: WebsiteStructuredDataProps) {
	const { page, locale } = props

	const json = serializeJsonLd(pageStructuredData(page, locale))

	return <script type="application/ld+json">{json}</script>
}
