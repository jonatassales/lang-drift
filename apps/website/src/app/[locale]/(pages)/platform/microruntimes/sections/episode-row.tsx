import type {
	Episode,
	EpisodeVerdict,
} from '@routes/platform/microruntimes/content/microruntimes-content.data'
import { ShareBar } from 'web/ui/share-bar'
import { cn } from 'web/utilities'

import type { CardValueLabels, CardValueState } from './card-value'

const skeletonBar = 'h-2.5 rounded bg-foreground/10'
const rowGrid =
	'grid grid-cols-[auto_minmax(0,1fr)_4rem_auto_minmax(6rem,auto)] items-center gap-x-3'

const verdictTones: Record<EpisodeVerdict, string> = {
	kept: 'text-status-success',
	dropped: 'text-muted-foreground',
}

// One episode the retrieval ranked, on one line: a pulsing skeleton while
// the search runs, then its id and slug, a bar of its similarity with the
// figure, and once judged, kept in the success role or dropped with the
// reason. A dropped episode fades back so the kept ones read as the
// precedents.
interface EpisodeRowProps {
	episode: Episode
	labels: CardValueLabels
	/** The similarity in the visitor's locale. */
	similarity: string
	state: Exclude<CardValueState, 'unset'>
	/** The judgment of the episode, once the agent has made it. */
	verdict?: { kind: EpisodeVerdict; label: string; reason?: string }
}

export function EpisodeRow(props: EpisodeRowProps) {
	const { episode, labels, similarity, state, verdict } = props

	if (state === 'resolving') {
		return (
			<li className="py-1">
				<span
					aria-hidden="true"
					className={`${rowGrid} animate-pulse motion-reduce:animate-none`}
				>
					<span className={`${skeletonBar} w-12`} />
					<span className={`${skeletonBar} w-40`} />
					<span className={`${skeletonBar} w-full`} />
					<span className={`${skeletonBar} w-8`} />
					<span className={`${skeletonBar} w-16 justify-self-end`} />
				</span>
				<span className="sr-only">{labels.resolving}</span>
			</li>
		)
	}

	return (
		<li
			className={cn(
				rowGrid,
				'py-1 motion-safe:transition-opacity motion-safe:duration-300 starting:opacity-0',
				verdict?.kind === 'dropped' && 'text-muted-foreground',
			)}
		>
			<span className="font-mono text-[10px] text-muted-foreground">
				{episode.id}
			</span>
			<span className="truncate font-mono text-[12px]">{episode.slug}</span>
			<ShareBar share={episode.similarity} />
			<span className="font-mono text-[12px] tabular-nums">{similarity}</span>
			<span
				className={cn(
					'text-right font-mono text-[10px]',
					verdict && verdictTones[verdict.kind],
				)}
			>
				{verdict && (
					<span className="motion-safe:transition-opacity motion-safe:duration-300 starting:opacity-0">
						{verdict.reason
							? `${verdict.label} · ${verdict.reason}`
							: verdict.label}
					</span>
				)}
			</span>
		</li>
	)
}
