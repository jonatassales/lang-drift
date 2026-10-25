import {
	GoIcon,
	NatsIcon,
	NextjsIcon,
	NodejsIcon,
	RedisIcon,
	type TechnologyIconProps,
} from '@langdrift/web/ui/icons/technology'
import type { ComponentType } from 'react'

import type { PlanComponent, Technology } from '@data/microruntimes-content'

import type { CardValueLabels, CardValueState } from './card-value'
import { ShareBar } from './share-bar'

interface ComponentRowProps {
	component: PlanComponent
	cpu: string
	/** The kind's label from the catalog. */
	kind: string
	labels: CardValueLabels
	memory: string
	/** The component's share of the CPU envelope, from 0 to 1. */
	share: number
	state: Exclude<CardValueState, 'unset'>
}

const technologyIcons: Record<
	Technology,
	ComponentType<TechnologyIconProps>
> = {
	Go: GoIcon,
	'Node.js': NodejsIcon,
	NATS: NatsIcon,
	Redis: RedisIcon,
	'Next.js': NextjsIcon,
}

const skeletonBar = 'h-2.5 rounded bg-foreground/10'
const rowGrid =
	'grid grid-cols-[auto_minmax(0,1fr)_auto_auto] items-center gap-x-3'

// One row of the ticket card's components manifest: a pulsing skeleton while
// the console resolves it, then the technology's mark and name, the path,
// the kind, the figures, and a bar with its share of the CPU envelope. The
// bar stays mounted through both states so its growth is a transition.
export function ComponentRow(props: ComponentRowProps) {
	const { component, cpu, kind, labels, memory, share, state } = props

	const Icon = technologyIcons[component.technology]

	return (
		<li className="grid gap-1.5 py-2">
			{state === 'resolving' ? (
				<span
					aria-hidden="true"
					className={`${rowGrid} animate-pulse motion-reduce:animate-none`}
				>
					<span className="size-4 rounded bg-foreground/10" />
					<span className="grid gap-1">
						<span className={`${skeletonBar} w-36`} />
						<span className={`${skeletonBar} w-20`} />
					</span>
					<span className={`${skeletonBar} w-10`} />
					<span className={`${skeletonBar} w-16`} />
				</span>
			) : (
				<div
					className={`${rowGrid} motion-safe:transition-opacity motion-safe:duration-300 starting:opacity-0`}
				>
					<Icon className="size-4 shrink-0" />
					<span className="grid min-w-0 gap-0.5">
						<span className="truncate font-mono text-[12px]">
							{component.path}
						</span>
						<span className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
							{component.technology}
							<span className="rounded-full border border-border px-1.5 text-[10px]">
								{kind}
							</span>
						</span>
					</span>
					<span className="font-mono text-[12px] tabular-nums">{cpu}</span>
					<span className="font-mono text-[12px] tabular-nums">{memory}</span>
				</div>
			)}
			{state === 'resolving' && (
				<span className="sr-only">{labels.resolving}</span>
			)}
			<ShareBar share={state === 'set' ? share : 0} />
		</li>
	)
}
