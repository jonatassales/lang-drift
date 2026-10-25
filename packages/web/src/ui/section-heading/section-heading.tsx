import { cn } from '@langdrift/web/utilities'

export interface SectionHeadingProps {
	className?: string
	headingId: string
	kicker: string
	lead?: string
	title: string
}

// The head of a product-page section: a small mono kicker, the heading and an
// optional lead, in the website's section rhythm.
export function SectionHeading(props: SectionHeadingProps) {
	const { className, headingId, kicker, lead, title } = props

	return (
		<div className={cn('max-w-[760px]', className)}>
			<span className="font-mono text-[10px] tracking-[0.11em] text-muted-foreground uppercase">
				{kicker}
			</span>
			<h2
				className="my-[18px] text-[clamp(38px,5vw,62px)] leading-[1.02] font-normal tracking-[-0.035em] text-balance text-foreground"
				id={headingId}
			>
				{title}
			</h2>
			{lead ? (
				<p className="m-0 text-base leading-[1.55] text-pretty text-muted-foreground">
					{lead}
				</p>
			) : null}
		</div>
	)
}
