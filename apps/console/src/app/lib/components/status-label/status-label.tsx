'use client'

import { useTranslations } from 'next-intl'

interface StatusLabelProps {
	value: string
}

const statusTones: Record<string, string> = {
	connected: 'text-status-success',
	failed: 'text-status-danger',
	running: 'text-status-success',
	unavailable: 'text-status-danger',
	unknown: 'text-status-danger',
}

export function StatusLabel(props: StatusLabelProps) {
	const { value } = props

	const translate = useTranslations('status')
	const tone = statusTones[value] ?? 'text-muted-foreground'

	return (
		<span
			className={`inline-flex items-center gap-1.5 rounded-full border border-border px-2 py-1 text-[11px] ${tone}`}
		>
			<span className="h-1.5 w-1.5 rounded-full bg-current" />
			{translate.has(value) ? translate(value) : value}
		</span>
	)
}
