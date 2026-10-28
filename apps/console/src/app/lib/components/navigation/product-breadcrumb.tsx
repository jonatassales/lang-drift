import { getTranslations } from 'next-intl/server'

import { readConsole } from '@data/server'
import type {
	ProductNavigation,
	RuntimeDetailSnapshot,
	TicketDetailSnapshot,
} from '@data/types'
import { Link } from '@i18n/navigation'

const ticketSections = ['tickets', 'analysis', 'plans', 'workflows']

interface ProductBreadcrumbProps {
	locale: string
	segments: string[]
}

export async function ProductBreadcrumb(props: ProductBreadcrumbProps) {
	const { locale, segments } = props

	const [root, id, section, ticketId] = segments
	if (!id || !['products', 'runtimes'].includes(root) || id === 'new')
		return null
	const [t, navigation, detail, ticket] = await Promise.all([
		getTranslations({
			locale,
			namespace: 'ui',
		}),
		readConsole<ProductNavigation>('products'),
		root === 'runtimes'
			? readConsole<RuntimeDetailSnapshot>(`runtimes/${encodeURIComponent(id)}`)
			: null,
		section === 'tickets' && ticketId
			? readConsole<TicketDetailSnapshot>(
					`tickets/${encodeURIComponent(ticketId)}`,
				)
			: null,
	])

	const productId = detail?.runtime.productId ?? id
	const product = navigation.products.find((item) => item.id === productId)
	if (!product) return null
	const base = `/products/${product.id}`
	const items = [
		{
			key: 'project',
			label: product.name,
			href: base,
		},
	]
	if (detail) {
		items.push(
			{
				key: 'runtimes',
				label: t('runtimes'),
				href: base,
			},
			{
				key: 'runtime',
				label: detail.runtime.name,
				href: `/runtimes/${id}`,
			},
		)
		if (section === 'configuration' || section === 'history')
			items.push({
				key: section,
				label: t(section),
				href: `/runtimes/${id}/${section}`,
			})
	} else if (section === 'config') {
		items.push({
			key: 'secrets',
			label: t('secrets'),
			href: `${base}/config/secrets`,
		})
	} else if (section === 'runtimes' && segments[3] === 'new') {
		items.push(
			{
				key: 'runtimes',
				label: t('runtimes'),
				href: base,
			},
			{
				key: 'new',
				label: t('createRuntime'),
				href: `${base}/runtimes/new`,
			},
		)
	} else if (section && ticketSections.includes(section)) {
		items.push({
			key: section,
			label: t(section),
			href: `${base}/${section}`,
		})
		if (ticket)
			items.push({
				key: 'ticket',
				label: `${ticket.ticket.number} · ${ticket.ticket.slug}`,
				href: `${base}/tickets/${ticket.ticket.id}`,
			})
	} else {
		items.push({
			key: section || 'runtimes',
			label: t(
				section === 'executions' || section === 'history'
					? section
					: 'runtimes',
			),
			href: section ? `${base}/${section}` : base,
		})
	}

	return (
		<nav
			aria-label={t('breadcrumb')}
			className="min-w-0 overflow-x-auto text-[11px] text-muted-foreground"
		>
			<ol className="flex min-w-max items-center gap-2 py-2 pr-3">
				{items.map((item, index) => (
					<li key={item.key} className="flex items-center gap-2">
						{index > 0 ? <span aria-hidden>/</span> : null}
						{index === items.length - 1 ? (
							<span aria-current="page" className="font-medium text-foreground">
								{item.label}
							</span>
						) : (
							<Link
								href={item.href}
								locale={locale}
								className="rounded hover:text-foreground hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
							>
								{item.label}
							</Link>
						)}
					</li>
				))}
			</ol>
		</nav>
	)
}
