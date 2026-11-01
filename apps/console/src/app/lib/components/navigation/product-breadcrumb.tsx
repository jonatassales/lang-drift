import { getTranslations } from 'next-intl/server'

import { isEntityId } from '@data/entity-id'
import { readConsole } from '@data/server'
import type {
	ProductNavigation,
	RuntimeDetailSnapshot,
	TicketDetailSnapshot,
} from '@data/types'
import { Link } from '@i18n/navigation'

import { breadcrumbItems } from './breadcrumb-items.compute'

interface ProductBreadcrumbProps {
	locale: string
	segments: string[]
}

export async function ProductBreadcrumb(props: ProductBreadcrumbProps) {
	const { locale, segments } = props

	const [root, id, section, ticketId] = segments
	const validRoot = ['products', 'runtimes'].includes(root)
	const validId = Boolean(id) && id !== 'new'
	if (!validRoot || !validId) return null

	// An id outside the alphabet names nothing, and its page answers 404.
	const ids = ticketId ? [id, ticketId] : [id]
	if (!ids.every(isEntityId)) return null

	const [translate, navigation, detail, ticket] = await Promise.all([
		getTranslations({
			locale,
			namespace: 'ui',
		}),
		readConsole<ProductNavigation>('products'),
		root === 'runtimes'
			? readConsole<RuntimeDetailSnapshot>(`runtimes/${id}`)
			: null,
		section === 'tickets' && ticketId
			? readConsole<TicketDetailSnapshot>(`tickets/${ticketId}`)
			: null,
	])

	const productId = detail?.runtime.productId ?? id
	const product = navigation.products.find((item) => item.id === productId)
	if (!product) return null

	const base = `/products/${product.id}`
	const items = [
		{
			href: base,
			key: 'project',
			label: product.name,
		},
		...breadcrumbItems({
			base,
			detail,
			id,
			section,
			segments,
			ticket,
			translate,
		}),
	]

	return (
		<nav
			aria-label={translate('breadcrumb')}
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
