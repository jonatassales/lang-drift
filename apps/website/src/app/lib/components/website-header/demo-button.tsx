'use client'

import { useTranslations } from 'next-intl'

import { Link } from '@i18n/navigation'

interface DemoButtonProps {
	/** The drawer closes itself when its button is tapped. */
	onClick?: () => void
}

export function DemoButton(props: DemoButtonProps) {
	const { onClick } = props

	const t = useTranslations('header')

	return (
		<Link
			className="inline-flex min-h-[38px] items-center justify-center rounded-[7px] bg-foreground px-[13px] py-2 text-xs font-semibold whitespace-nowrap text-background"
			href="/book-a-demo"
			onClick={onClick}
		>
			{t('bookDemo')}
		</Link>
	)
}
