'use client'
import { useTranslations } from 'next-intl'

interface StubBannerProps {
	source: string
}
export function StubBanner(props: StubBannerProps) {
	const { source } = props

	const t = useTranslations('ui')
	if (source !== 'stub') return null

	return (
		<aside className="mb-4 text-[11px] text-muted-foreground">
			{t('stubBanner')}
		</aside>
	)
}
