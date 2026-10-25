import organization from '@app/metadata/organization.json'
import site from '@app/metadata/site.json'

const organizationId = `${site.url}/#organization`

const structuredData = {
	'@context': 'https://schema.org',
	'@graph': [
		{
			'@id': organizationId,
			'@type': 'Organization',
			name: organization.name,
			url: organization.url,
			logo: organization.logo,
			sameAs: organization.sameAs,
		},
		{
			'@id': `${site.url}/#website`,
			'@type': 'WebSite',
			name: site.name,
			publisher: { '@id': organizationId },
			url: `${site.url}/`,
		},
	],
}

export function WebsiteStructuredData() {
	const json = JSON.stringify(structuredData).replaceAll('<', '\\u003c')

	return <script type="application/ld+json">{json}</script>
}
