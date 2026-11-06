import { dollars } from '@app/template/formatters/dollars.fmt'
import type { WebsiteLocale } from '@i18n/routing'
import { cn } from 'web/utilities'

export interface EnvelopeSegment {
	amount: number
	color: string
	id: string
	label: string
}

// An authorized envelope as a stacked bar, then a legend that gives each
// segment its amount and its share of the envelope.
interface EnvelopeBarProps {
	authorized: number
	locale: WebsiteLocale
	segments: readonly EnvelopeSegment[]
}

export function EnvelopeBar(props: EnvelopeBarProps) {
	const { authorized, locale, segments } = props

	const money = dollars(locale)
	const percent = new Intl.NumberFormat(locale, {
		maximumFractionDigits: 0,
		style: 'percent',
	})
	const barLabel = segments
		.map((segment) => `${segment.label} ${money.format(segment.amount)}`)
		.join(', ')

	return (
		<>
			<div
				aria-label={barLabel}
				className="mt-6 flex h-3 overflow-hidden rounded-full bg-foreground/5"
				role="img"
			>
				{segments.map((segment) => (
					<span
						className={cn('h-full', segment.color)}
						key={segment.id}
						style={{
							flexBasis: `${(segment.amount / authorized) * 100}%`,
						}}
					/>
				))}
			</div>
			<ul className="m-0 mt-4 grid list-none grid-cols-2 gap-x-6 gap-y-3 p-0 sm:grid-cols-4">
				{segments.map((segment) => (
					<li className="grid gap-1" key={segment.id}>
						<span className="flex items-center gap-2 text-[12px] text-muted-foreground">
							<span
								aria-hidden="true"
								className={cn('size-2 rounded-full', segment.color)}
							/>
							{segment.label}
						</span>
						<span className="font-mono text-[14px] tabular-nums">
							{money.format(segment.amount)}
							<span className="ml-2 text-[11px] text-muted-foreground">
								{percent.format(segment.amount / authorized)}
							</span>
						</span>
					</li>
				))}
			</ul>
		</>
	)
}
