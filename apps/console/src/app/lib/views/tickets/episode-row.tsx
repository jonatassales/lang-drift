import { ShareBar } from '@langdrift/web/ui/share-bar'
import { cn } from '@langdrift/web/utilities'

import type { TicketEpisode } from '@data/types'

interface EpisodeRowProps {
	episode: TicketEpisode
	/** The similarity in the viewer's locale. */
	similarity: string
	/** The verdict and its reason in the viewer's locale. */
	verdict: string
}

const verdictTones: Record<TicketEpisode['verdict'], string> = {
	kept: 'text-status-success',
	dropped: 'text-muted-foreground',
}

// One episode the retrieval ranked: its id and slug, a bar of its
// similarity with the figure, and the verdict, kept in the success role or
// dropped with its reason, faded back so the kept ones read as precedents.
export function EpisodeRow(props: EpisodeRowProps) {
	const { episode, similarity, verdict } = props

	return (
		<li
			className={cn(
				'grid grid-cols-[auto_minmax(0,1fr)_4rem_auto_minmax(6rem,auto)] items-center gap-x-3 py-1 max-sm:grid-cols-[minmax(0,1fr)_auto_auto]',
				episode.verdict === 'dropped' && 'text-muted-foreground',
			)}
		>
			<span className="font-mono text-[10px] text-muted-foreground max-sm:hidden">
				{episode.id}
			</span>
			<span className="truncate font-mono text-[12px]">{episode.slug}</span>
			<ShareBar className="max-sm:hidden" share={episode.similarity} />
			<span className="font-mono text-[12px] tabular-nums">{similarity}</span>
			<span
				className={cn(
					'text-right font-mono text-[10px]',
					verdictTones[episode.verdict],
				)}
			>
				{verdict}
			</span>
		</li>
	)
}
