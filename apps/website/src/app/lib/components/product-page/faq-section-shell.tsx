import type { ReactNode } from 'react'

import { cn } from 'web/utilities'

interface FaqSectionShellProps {
	children: ReactNode
	/** Overrides the landmark's width where a page sets its own, as the home does. */
	className?: string
	description: string
	title: string
}

// The FAQ landmark the product pages and the home share: a heading, a lead
// and the answers that follow.
export function FaqSectionShell(props: FaqSectionShellProps) {
	const { children, className, description, title } = props

	return (
		<section
			aria-labelledby="faq-heading"
			className={cn(
				'mx-auto max-w-[1180px] scroll-mt-24 border-t border-border px-[clamp(18px,3vw,36px)] py-16 sm:py-20',
				className,
			)}
			id="faq"
		>
			<div className="max-w-[760px]">
				<h2
					className="m-0 text-[clamp(30px,4vw,46px)] leading-[1.08] font-normal tracking-[-0.03em] text-balance text-foreground"
					id="faq-heading"
				>
					{title}
				</h2>
				<p className="mt-4 mb-0 text-base leading-[1.55] text-pretty text-muted-foreground">
					{description}
				</p>
			</div>
			{children}
		</section>
	)
}
