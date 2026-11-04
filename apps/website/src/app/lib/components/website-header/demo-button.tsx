'use client'

import { useTranslations } from 'next-intl'

import { Link } from '@i18n/navigation'
import { cn } from 'web/utilities'

interface DemoButtonProps {
	/** The drawer gives its button the 44px touch height. */
	className?: string
	/** The drawer closes itself when its button is tapped. */
	onClick?: () => void
}

export function DemoButton(props: DemoButtonProps) {
	const { className, onClick } = props

	const translate = useTranslations('header')

	return (
		<Link
			className={cn(
				'inline-flex min-h-[38px] items-center justify-center rounded-[7px] bg-foreground px-[13px] py-2 text-xs font-semibold whitespace-nowrap text-background focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary',
				className,
			)}
			href="/book-a-demo"
			onClick={onClick}
		>
			{translate('bookDemo')}
		</Link>
	)
}
