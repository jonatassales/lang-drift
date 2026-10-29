import type { ReactNode } from 'react'

interface TicketSectionProps {
	children: ReactNode
	/** The section's fragment, which a listing entry links to. */
	id: string
	/** The section's one line, such as a phase's latest figure. */
	summary?: ReactNode
	title: string
}

// One section of the ticket page under a full-width hairline: its heading
// with its one-line summary beside it, then its body. It lands below the
// console's sticky header when the URL names it.
export function TicketSection(props: TicketSectionProps) {
	const { children, id, summary, title } = props

	const headingId = `${id}-heading`

	return (
		<section
			aria-labelledby={headingId}
			className="grid min-w-0 scroll-mt-20 gap-4 border-t border-border pt-6"
			id={id}
		>
			<header className="flex min-w-0 flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
				<h2 className="m-0 text-[15px] font-medium" id={headingId}>
					{title}
				</h2>
				{summary ? (
					<span className="min-w-0 font-mono text-[12px] text-muted-foreground">
						{summary}
					</span>
				) : null}
			</header>
			<div className="min-w-0">{children}</div>
		</section>
	)
}
