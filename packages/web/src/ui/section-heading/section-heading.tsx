import { Eyebrow } from 'web/ui/eyebrow'

export interface SectionHeadingProps {
	headingId: string
	kicker?: string
	lead?: string
	title: string
}

// The head of a page section: an optional mono kicker, the heading and an
// optional lead.
export function SectionHeading(props: SectionHeadingProps) {
	const { headingId, kicker, lead, title } = props

	return (
		<div className="max-w-[760px]">
			{kicker ? (
				<Eyebrow className="text-[10px] tracking-[0.11em]">{kicker}</Eyebrow>
			) : null}
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
