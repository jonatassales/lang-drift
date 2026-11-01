'use client'
import { useTranslations } from 'next-intl'

interface StubBannerProps {
	source: string
}
export function StubBanner(props: StubBannerProps) {
	const { source } = props

	const translate = useTranslations('ui')
	if (source !== 'stub') return null

	return (
		<aside className="mb-4 text-[11px] text-muted-foreground">
			{translate('stubBanner')}
		</aside>
	)
}
