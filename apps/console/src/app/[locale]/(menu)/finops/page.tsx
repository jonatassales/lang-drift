import { getTranslations } from 'next-intl/server'

import { CostOverview } from '@views/finops/cost-overview'

export async function generateMetadata() {
	const t = await getTranslations('ui')

	return { title: t('costOverview') }
}

export default function Page() {
	return <CostOverview />
}
