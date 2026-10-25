import { ArrowRightIcon } from '@langdrift/web/ui/icons/arrow'
import { ChromaText } from '@langdrift/web/vendors/smoothui/chroma-text'
import type { ReactNode } from 'react'

export interface ProductHeroCopy {
	titleStart: string
	titleHighlight: string
	titleEnd: string
	lead: string
	primaryCta: string
	secondaryCta: string
	caption?: string
	summary: string
}

interface ProductHeroProps {
	children: ReactNode
	copy: ProductHeroCopy
	headingId: string
	primaryFragment: string
	secondaryFragment: string
}

// The hero of a product page, in the Vision hero's layout: a centered
// question headline with one highlighted phrase, a lead, two calls to action
// to sections of the same page, then the demo with an optional visible caption
// and a summary for screen readers.
export function ProductHero(props: ProductHeroProps) {
	const { children, copy, headingId, primaryFragment, secondaryFragment } =
		props

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
					{copy.titleStart}
					<ChromaText>{copy.titleHighlight}</ChromaText>
					{copy.titleEnd}
				</h1>
				<p className="mx-auto mt-[23px] mb-0 max-w-[750px] text-base leading-[1.65] text-pretty text-muted-foreground">
					{copy.lead}
				</p>
				<div className="mt-[25px] flex flex-wrap items-center justify-center gap-x-[26px] gap-y-4">
					<a
						className="smooth-primary-button min-h-11"
						href={`#${primaryFragment}`}
					>
						{copy.primaryCta}
						<ArrowRightIcon />
					</a>
					<a
						className="smooth-text-link inline-flex min-h-11 items-center"
						href={`#${secondaryFragment}`}
					>
						{copy.secondaryCta}
					</a>
				</div>
			</div>
			<div
				aria-describedby={
					copy.caption ? `${captionId} ${summaryId}` : summaryId
				}
				className="mt-12"
			>
				{children}
				{copy.caption && (
					<p
						className="mx-auto mt-5 mb-0 max-w-[880px] text-center text-xs leading-[1.65] text-muted-foreground"
						id={captionId}
					>
						{copy.caption}
					</p>
				)}
				<p className="sr-only" id={summaryId}>
					{copy.summary}
				</p>
			</div>
		</section>
	)
}
