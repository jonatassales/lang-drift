'use client'

import { useTranslations } from 'next-intl'

import { Link } from '@i18n/navigation'

interface AccessLinksProps {
	onNavigate?: () => void
}

export function AccessLinks(props: AccessLinksProps) {
	const { onNavigate } = props

	const t = useTranslations('header')

	return (
		<Link
			className="website-get-started"
			href="/#early-access"
			onClick={onNavigate}
		>
			{t('getStarted')}
		</Link>
	)
}
