import type { ComponentType } from 'react'

import type {
	PlanComponent,
	Technology,
} from '@routes/platform/microruntimes/content/microruntimes-content.types'
import {
	GoIcon,
	NatsIcon,
	NextjsIcon,
	NodejsIcon,
	RedisIcon,
	type TechnologyIconProps,
} from 'web/ui/icons/technology'
import type { RevealValueLabels, RevealValueState } from 'web/ui/reveal-value'
import { ShareBar } from 'web/ui/share-bar'

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
	'grid grid-cols-[auto_minmax(0,1fr)_auto_auto_auto] items-center gap-x-3'

// One row of the plan's components manifest: a pulsing skeleton while the
// console resolves it, then on one line the technology's mark, the path,
// the technology and kind, the figures, and under it a bar with its share
// of the CPU envelope. The bar stays mounted through both states so its
// growth is a transition.
interface ComponentRowProps {
	component: PlanComponent
	cpu: string
	/** The kind's label from the catalog. */
	kind: string
	labels: RevealValueLabels
	memory: string
	/** The component's share of the CPU envelope, from 0 to 1. */
	share: number
	state: Exclude<RevealValueState, 'unset'>
}

export function ComponentRow(props: ComponentRowProps) {
	const { component, cpu, kind, labels, memory, share, state } = props

	const Icon = technologyIcons[component.technology]

	return (
		<li className="grid gap-0.5 py-1">
			{state === 'resolving' ? (
				<span
					aria-hidden="true"
					className={`${rowGrid} animate-pulse motion-reduce:animate-none`}
				>
					<span className="size-4 rounded bg-foreground/10" />
					<span className={`${skeletonBar} w-36`} />
					<span className={`${skeletonBar} w-16`} />
					<span className={`${skeletonBar} w-10`} />
					<span className={`${skeletonBar} w-14`} />
				</span>
			) : (
				<div
					className={`${rowGrid} motion-safe:transition-opacity motion-safe:duration-300 starting:opacity-0`}
				>
					<Icon className="size-4 shrink-0" />
					<span className="truncate font-mono text-[12px]">
						{component.path}
					</span>
					<span className="flex items-center gap-1.5 text-[11px] text-muted-foreground max-sm:hidden">
						{component.technology}
						<span className="rounded-full border border-border px-1.5 text-[10px]">
							{kind}
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
