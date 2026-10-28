'use client'
import { useTranslations } from 'next-intl'

interface StatusLabelProps {
	value: string
}
export function StatusLabel(props: StatusLabelProps) {
	const { value } = props

	const t = useTranslations('status')

	return (
		<span
			className={`inline-flex items-center gap-1.5 rounded-full border border-border px-2 py-1 text-[11px] ${value === 'running' || value === 'connected' ? 'text-status-success' : value === 'failed' || value === 'unknown' || value === 'unavailable' ? 'text-status-danger' : 'text-muted-foreground'}`}
		>
			<span className="h-1.5 w-1.5 rounded-full bg-current" />
			{t.has(value) ? t(value) : value}
		</span>
	)
}
