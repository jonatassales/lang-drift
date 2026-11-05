'use client'

import { useSearchParams } from 'next/navigation'
import { useTranslations } from 'next-intl'

export function LocaleAvailabilityNotice() {
	const params = useSearchParams()
	const translate = useTranslations('routeState')
	if (params.get('availability') !== 'zh-CN') return null

	return (
		<output className="relative z-50 mx-auto block max-w-5xl bg-background px-6 py-3 text-center text-sm text-foreground">
			{translate('unavailable')}
		</output>
	)
}
