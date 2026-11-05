import { getTranslations } from 'next-intl/server'

import { RouteStateNotice } from '@components/route-state/route-state-notice'
import { Link } from '@i18n/navigation'

export { generateMetadata } from './not-found.metadata'

export default async function NotFound() {
	const [translateNotFound, translateUi] = await Promise.all([
		getTranslations('notFound'),
		getTranslations('ui'),
	])

	return (
		<RouteStateNotice
			title={translateNotFound('title')}
			description={translateNotFound('description')}
		>
			<Link
				className="w-fit text-primary underline underline-offset-4 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
				href="/home"
			>
				{translateUi('home')}
			</Link>
		</RouteStateNotice>
	)
}
