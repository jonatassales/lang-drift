import { BasicAccordion } from '@langdrift/web/vendors/smoothui/basic-accordion'

export interface ProductFaqItem {
	answer: string
	id: string
	question: string
}

interface ProductFaqSectionProps {
	description: string
	items: readonly ProductFaqItem[]
	title: string
}

// The FAQ of a product page: a heading, a lead and an accordion of questions.
export function ProductFaqSection(props: ProductFaqSectionProps) {
	const { description, items, title } = props

	return (
		<section
			aria-labelledby="faq-heading"
			className="mx-auto max-w-[1180px] scroll-mt-24 border-t border-border px-[clamp(18px,3vw,36px)] py-16 sm:py-20"
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
			<BasicAccordion
				className="mt-10 [&_button]:whitespace-normal!"
				headingLevel={3}
				items={items.map((item) => ({
					content: (
						<p className="m-0 text-[13px] leading-[1.6] text-muted-foreground">
							{item.answer}
						</p>
					),
					id: item.id,
					title: item.question,
				}))}
				layoutId="product-faq"
			/>
		</section>
	)
}
