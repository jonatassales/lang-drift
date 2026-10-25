import { ChromaText } from '@langdrift/web/vendors/smoothui/chroma-text'
import type { ReactNode } from 'react'

export interface ProductHeroCta {
	href: string
	label: string
}

interface ProductHeroProps {
	caption?: string
	children?: ReactNode
	headingId: string
	lead: string
	primaryCta: ProductHeroCta
	secondaryCta: ProductHeroCta
	summary?: string
	titleEnd: string
	titleHighlight: string
	titleStart: string
}

// The hero of a product page, in the Vision hero's layout: a centered
// question headline with one highlighted phrase, a lead, two calls to action,
// then the illustrative demo described by its caption and a summary for
// screen readers.
export function ProductHero(props: ProductHeroProps) {
	const {
		caption,
		children,
		headingId,
		lead,
		primaryCta,
		secondaryCta,
		summary,
		titleEnd,
		titleHighlight,
		titleStart,
	} = props

	const captionId = `${headingId}-demo-caption`
	const summaryId = `${headingId}-demo-summary`

	return (
		<section
			aria-labelledby={headingId}
			className="mx-auto max-w-[1280px] px-[clamp(18px,3vw,36px)] pt-[58px] pb-[70px] max-md:pt-10 max-md:pb-12"
		>
			<div className="mx-auto text-center">
				<h1
					className="m-0 text-[clamp(28px,3.2vw,43px)] leading-[1.16] font-medium tracking-[-0.045em] text-balance"
					id={headingId}
				>
					{titleStart}
					<ChromaText>{titleHighlight}</ChromaText>
					{titleEnd}
				</h1>
				<p className="mx-auto mt-[23px] mb-0 max-w-[750px] text-base leading-[1.65] text-pretty text-muted-foreground">
					{lead}
				</p>
				<div className="mt-[25px] flex flex-wrap items-center justify-center gap-x-[26px] gap-y-4">
					<a className="smooth-primary-button min-h-11" href={primaryCta.href}>
						{primaryCta.label}
						<svg
							aria-hidden="true"
							fill="none"
							height="16"
							viewBox="0 0 24 24"
							width="16"
						>
							<path
								d="M4 12h16m-6-6 6 6-6 6"
								stroke="currentColor"
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth="1.5"
							/>
						</svg>
					</a>
					<a
						className="smooth-text-link inline-flex min-h-11 items-center"
						href={secondaryCta.href}
					>
						{secondaryCta.label}
					</a>
				</div>
			</div>
			{children ? (
				<div aria-describedby={`${captionId} ${summaryId}`} className="mt-12">
					{children}
					<p
						className="mx-auto mt-5 mb-0 max-w-[880px] text-center text-xs leading-[1.65] text-muted-foreground"
						id={captionId}
					>
						{caption}
					</p>
					<p className="sr-only" id={summaryId}>
						{summary}
					</p>
				</div>
			) : null}
		</section>
	)
}
