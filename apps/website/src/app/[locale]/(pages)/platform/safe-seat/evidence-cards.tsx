import { InfoCards } from 'web/ui/info-cards'

import {
	evidenceFollowUps,
	evidenceIds,
	evidenceSources,
	sourceLinks,
} from './safe-seat-content.data'
import type { SafeSeatEvidenceCopy } from './safe-seat-messages.types'

const sourceLinkClass =
	'justify-self-start text-[13px] underline decoration-current/25 underline-offset-5 focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary'

// One card per finding: its source, what it found, where it stops and a
// link to the primary source, and to a later source the card also reports,
// so a reader can check every claim.
interface EvidenceCardsProps {
	copy: SafeSeatEvidenceCopy
}

export function EvidenceCards(props: EvidenceCardsProps) {
	const { copy } = props

	return (
		<InfoCards
			className="mt-11"
			items={evidenceIds.map((id) => {
				const followUp = evidenceFollowUps[id]

				return {
					description: copy.items[id].description,
					eyebrow: copy.items[id].source,
					footer: (
						<div className="grid gap-3">
							<p className="m-0 text-xs leading-[1.55] text-muted-foreground">
								<strong className="font-medium text-foreground">
									{copy.limitLabel}
								</strong>{' '}
								{copy.items[id].limit}
							</p>
							<a
								className={sourceLinkClass}
								href={sourceLinks[evidenceSources[id]]}
								rel="noopener noreferrer"
								target="_blank"
							>
								{copy.sourceLink}
							</a>
							{followUp && (
								<a
									className={sourceLinkClass}
									href={sourceLinks[followUp]}
									rel="noopener noreferrer"
									target="_blank"
								>
									{copy.followUpLink}
								</a>
							)}
						</div>
					),
					id,
					title: copy.items[id].title,
				}
			})}
		/>
	)
}
