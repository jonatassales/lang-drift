import { TextLink } from 'web/ui/text-link'

// The closing line of a product-page section, with a link to the section it
// leads to.
interface SectionOutroProps {
	fragment: string
	link: string
	text: string
}

export function SectionOutro(props: SectionOutroProps) {
	const { fragment, link, text } = props

	return (
		<div className="mt-8 flex flex-wrap items-center justify-between gap-x-8 gap-y-3">
			<p className="m-0 max-w-[640px] text-[13px] leading-[1.6] text-muted-foreground">
				{text}
			</p>
			<TextLink href={`#${fragment}`}>
				{link}
				<span aria-hidden="true"> ↗</span>
			</TextLink>
		</div>
	)
}
