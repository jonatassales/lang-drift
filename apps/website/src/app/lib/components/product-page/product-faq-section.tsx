import { BasicAccordion } from '@langdrift/web/vendors/smoothui/basic-accordion'

import { FaqSectionShell } from '@components/product-page/faq-section-shell'

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
		<FaqSectionShell description={description} title={title}>
			<BasicAccordion
				className="mt-10"
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
		</FaqSectionShell>
	)
}
