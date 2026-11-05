import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'
import {
	type FinopsMessages,
	sourceIds,
} from '@routes/platform/finops/content/finops-content'
import { Eyebrow } from 'web/ui/eyebrow'

// What would feed the ledger, as a plain list under a label that says the
// sources are intended, closed by the line that says a missing one stays
// visible as a gap.
interface SourceListProps {
	copy: FinopsMessages['sources']
	locale: WebsiteLocale
}

export function SourceList(props: SourceListProps) {
	const { copy, locale } = props

	return (
		<div className="mt-11">
			<Eyebrow className="text-[10px]">{copy.intended}</Eyebrow>
			<ul className="m-0 mt-4 grid list-none gap-x-10 p-0 md:grid-cols-2">
				{sourceIds.map((id) => (
					<li
						className="grid content-start gap-1 border-t border-border py-4"
						key={id}
					>
						<strong className="text-[15px] leading-snug font-medium">
							{copy.items[id].title}
						</strong>
						<span className="text-[13px] leading-[1.55] text-muted-foreground">
							{copy.items[id].description}
						</span>
						{id === 'microruntimes' && (
							<Link
								className="smooth-text-link mt-1 justify-self-start focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
								href="/platform/microruntimes"
								locale={locale}
							>
								{copy.items.microruntimes.link}
							</Link>
						)}
					</li>
				))}
			</ul>
			<p className="mt-6 mb-0 max-w-[760px] text-[13px] leading-[1.6] text-muted-foreground">
				{copy.intendedModel}
			</p>
		</div>
	)
}
