import type { PageId } from '@app/metadata/page-metadata'
import {
	pageStructuredData,
	serializeJsonLd,
} from '@app/metadata/structured-data'
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
