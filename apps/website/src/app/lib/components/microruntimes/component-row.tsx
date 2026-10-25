import {
	GoIcon,
	NatsIcon,
	NextjsIcon,
	NodejsIcon,
	RedisIcon,
	type TechnologyIconProps,
} from '@langdrift/web/ui/icons/technology'
import type { ComponentType } from 'react'

import type {
	MicroruntimesMessages,
	PlanComponent,
	Technology,
} from '@data/microruntimes-content'

interface ComponentRowProps {
	component: PlanComponent
	cpu: string
	labels: Pick<MicroruntimesMessages['session']['card'], 'resolving'>
	memory: string
	/** The component's share of the CPU envelope, from 0 to 1. */
	share: number
	state: 'resolving' | 'set'
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

// One row of the ticket card's components manifest: a pulsing skeleton while
// the console resolves it, then the technology's mark and name, the path,
// the kind, the figures, and a thin bar with its share of the CPU envelope.
export function ComponentRow(props: ComponentRowProps) {
	const { component, cpu, labels, memory, share, state } = props

	if (state === 'resolving') {
		return (
			<li className="grid gap-2 py-2">
				<span
					aria-hidden="true"
					className="grid animate-pulse grid-cols-[auto_minmax(0,1fr)_auto_auto] items-center gap-x-3 motion-reduce:animate-none"
				>
					<span className="size-4 rounded bg-foreground/10" />
					<span className={`${skeletonBar} w-36`} />
					<span className={`${skeletonBar} w-10`} />
					<span className={`${skeletonBar} w-16`} />
				</span>
				<span
					aria-hidden="true"
					className="block h-0.5 rounded bg-foreground/5"
				/>
				<span className="sr-only">{labels.resolving}</span>
			</li>
		)
	}

	const Icon = technologyIcons[component.technology]

	return (
		<li className="grid gap-1.5 py-2">
			<div className="grid grid-cols-[auto_minmax(0,1fr)_auto_auto] items-center gap-x-3">
				<Icon className="size-4 shrink-0" />
				<span className="grid min-w-0 gap-0.5">
					<span className="truncate font-mono text-[12px]">
						{component.path}
					</span>
					<span className="text-[11px] text-muted-foreground">
						{component.technology} · {component.kind}
					</span>
				</span>
				<span className="font-mono text-[12px] tabular-nums">{cpu}</span>
				<span className="font-mono text-[12px] tabular-nums">{memory}</span>
			</div>
			<span
				aria-hidden="true"
				className="block h-0.5 overflow-hidden rounded bg-foreground/5"
			>
				<span
					className="block h-full rounded bg-primary motion-safe:transition-[width] motion-safe:duration-700 motion-safe:ease-out"
					style={{ width: `${share * 100}%` }}
				/>
			</span>
		</li>
	)
}
