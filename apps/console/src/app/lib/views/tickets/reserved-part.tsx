import { useTranslations } from 'next-intl'

import type { TicketReservedPart } from '@data/types'
import { Eyebrow } from 'web/ui/eyebrow'

interface ReservedPartProps {
	reservedPart: TicketReservedPart | null
}

// The part of the ticket a person writes by hand under Safe Seat: its
// category, who chose it, its goal and the paths it covers. A plan that
// has not reserved one yet says so.
export function ReservedPart(props: ReservedPartProps) {
	const { reservedPart } = props

	const translate = useTranslations('tickets')

	if (!reservedPart) {
		return (
			<p className="m-0 text-xs text-muted-foreground">
				{translate('safeSeat.noReservedPart')}
			</p>
		)
	}

	return (
		<div className="grid gap-1.5">
			<p className="m-0 flex flex-wrap items-baseline gap-x-3 gap-y-1">
				<Eyebrow>{translate('safeSeat.reservedPart')}</Eyebrow>
				<span className="text-xs text-foreground">
					{translate(`safeSeat.categories.${reservedPart.category}`)}
				</span>
				<span className="text-xs text-muted-foreground">
					{translate(`safeSeat.chosenBy.${reservedPart.chosenBy}`)}
				</span>
			</p>
			<p className="m-0 text-xs leading-snug">{reservedPart.goal}</p>
			<ul className="m-0 grid list-none gap-0.5 p-0">
				{reservedPart.paths.map((path) => (
					<li key={path}>
						<code className="font-mono text-xs break-all text-muted-foreground">
							{path}
						</code>
					</li>
				))}
			</ul>
		</div>
	)
}
