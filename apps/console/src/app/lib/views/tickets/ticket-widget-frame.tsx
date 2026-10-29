import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { useTranslations } from 'next-intl'
import type { ReactNode } from 'react'

import type { TicketPhase } from '@data/types'

interface TicketWidgetFrameProps {
	children: ReactNode
	phase: TicketPhase
	/** The panel's one line: its latest figure. */
	summary: ReactNode
}

// A ticket widget as the website draws its panel: the phase's eyebrow and
// one-line summary over the full body, never collapsed.
export function TicketWidgetFrame(props: TicketWidgetFrameProps) {
	const { children, phase, summary } = props

	const t = useTranslations('tickets')

	return (
		<article className="grid min-w-0 gap-3 text-left">
			<header className="flex min-w-0 items-center justify-between gap-4">
				<Eyebrow className="shrink-0">{t(`labels.${phase}`)}</Eyebrow>
				<span className="min-w-0 truncate font-mono text-[12px] text-muted-foreground">
					{summary}
				</span>
			</header>
			<div>{children}</div>
		</article>
	)
}
