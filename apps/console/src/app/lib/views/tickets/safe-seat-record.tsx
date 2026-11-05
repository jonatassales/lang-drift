import { useTranslations } from 'next-intl'

import type { Ticket } from '@data/types'
import { Eyebrow } from 'web/ui/eyebrow'
import { cn } from 'web/utilities'

import { safeSeatView } from './safe-seat-view.map'
import type { SafeSeatView } from './safe-seat-view.type'

interface SafeSeatRecordProps {
	ticket: Ticket
}

const toneClasses: Record<SafeSeatView['tone'], string> = {
	success: 'text-status-success',
	warning: 'text-status-warning',
	danger: 'text-status-danger',
	muted: 'text-muted-foreground',
}

// The ticket's Safe Seat state, with the override reason when a lead
// recorded one, labeled when it is illustrative stub data.
export function SafeSeatRecord(props: SafeSeatRecordProps) {
	const { ticket } = props

	const translate = useTranslations('tickets')

	const view = safeSeatView(ticket)

	return (
		<section className="grid gap-1 border-t border-border pt-6">
			<p className="m-0 flex flex-wrap items-baseline gap-x-3 gap-y-1">
				<Eyebrow>{translate('safeSeat.title')}</Eyebrow>
				<span className={cn('text-xs', toneClasses[view.tone])}>
					{translate(view.stateKey)}
				</span>
				{view.illustrative && (
					<span className="font-mono text-xs text-muted-foreground">
						{translate('safeSeat.illustrative')}
					</span>
				)}
			</p>
			{view.overrideReason && (
				<p className="m-0 text-xs leading-snug text-muted-foreground">
					{translate('safeSeat.overrideReason', {
						reason: view.overrideReason,
					})}
				</p>
			)}
		</section>
	)
}
