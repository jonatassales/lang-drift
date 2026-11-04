import { type MouseEvent, type ReactNode, useId } from 'react'

import { Eyebrow } from 'web/ui/eyebrow'
import { CaretDownIcon } from 'web/ui/icons/arrow'
import { cn } from 'web/utilities'

interface SessionCardProps {
	body: ReactNode
	expanded: boolean
	label: string
	onOpen: (event: MouseEvent<HTMLButtonElement>) => void
	phase?: string
	/** The card's one line while it is collapsed: its latest figure. */
	summary: ReactNode
}

// One of the cards the session fills: a header that names the act and
// carries its latest figure, over a body that slides open and shut. The
// body stays mounted either way, so its values keep their state while the
// card is collapsed; shut, it is inert and hidden from assistive
// technology.
export function SessionCard(props: SessionCardProps) {
	const { body, expanded, label, onOpen, phase, summary } = props

	const id = useId()
	const triggerId = `${id}-trigger`
	const regionId = `${id}-region`

	return (
		<article className="min-w-0 rounded-2xl border border-border bg-card text-left">
			<button
				aria-controls={regionId}
				aria-expanded={expanded}
				className="flex min-h-11 w-full items-center justify-between gap-4 px-5 py-3 text-left whitespace-nowrap focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary sm:px-6"
				data-phase={phase}
				id={triggerId}
				onClick={onOpen}
				type="button"
			>
				<Eyebrow className="shrink-0">{label}</Eyebrow>
				<span className="flex min-w-0 flex-1 items-center justify-end gap-3">
					<span className="min-w-0 truncate font-mono text-[12px] text-muted-foreground">
						{summary}
					</span>
					<CaretDownIcon
						className={cn(
							'shrink-0 text-muted-foreground motion-safe:transition-transform motion-safe:duration-300',
							expanded && 'rotate-180',
						)}
					/>
				</span>
			</button>
			<section
				aria-hidden={expanded ? undefined : 'true'}
				aria-labelledby={triggerId}
				className={cn(
					'grid motion-safe:transition-[grid-template-rows] motion-safe:duration-300 motion-safe:ease-out',
					expanded ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]',
				)}
				id={regionId}
				inert={!expanded}
			>
				<div className="min-h-0 overflow-hidden">
					<div className="px-5 pt-1 pb-5 sm:px-6 sm:pb-6">{body}</div>
				</div>
			</section>
		</article>
	)
}
