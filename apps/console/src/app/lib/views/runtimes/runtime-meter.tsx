'use client'
import { useFormatter, useTranslations } from 'next-intl'

interface RuntimeMeterProps {
	label: string
	value: number | null
	limit: number | null
}
export function RuntimeMeter(props: RuntimeMeterProps) {
	const { label, value, limit } = props

	const t = useTranslations('ui')
	const format = useFormatter()
	const known = value !== null && limit !== null

	return (
		<div>
			<div className="mb-1 flex flex-wrap justify-between gap-1 font-mono text-[10px]">
				<span>{label}</span>
				<span>
					{value === null
						? t('unknown')
						: format.number(value, { maximumFractionDigits: 3 })}{' '}
					/{' '}
					{limit === null
						? t('unknown')
						: format.number(limit, { maximumFractionDigits: 3 })}
				</span>
			</div>
			{known ? (
				<progress
					className="block h-1 w-full appearance-none accent-foreground [&::-webkit-progress-bar]:bg-muted [&::-webkit-progress-value]:bg-foreground [&::-moz-progress-bar]:bg-foreground"
					aria-label={label}
					value={value}
					max={limit}
				/>
			) : (
				<div className="h-1.5 rounded bg-muted" />
			)}
		</div>
	)
}
