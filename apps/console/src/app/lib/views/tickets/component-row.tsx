import {
	GoIcon,
	NatsIcon,
	NextjsIcon,
	NodejsIcon,
	RedisIcon,
	type TechnologyIconProps,
} from '@langdrift/web/ui/icons/technology'
import { ShareBar } from '@langdrift/web/ui/share-bar'
import type { ComponentType } from 'react'

import type { Technology, TicketComponent } from '@data/types'

interface ComponentRowProps {
	component: TicketComponent
	cpu: string
	/** The kind's label in the viewer's locale. */
	kind: string
	memory: string
	/** The component's share of the CPU envelope, from 0 to 1. */
	share: number
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

// One row of the plan's components manifest: the technology's mark, the
// path, the technology and kind, the figures, and under it a bar with its
// share of the CPU envelope.
export function ComponentRow(props: ComponentRowProps) {
	const { component, cpu, kind, memory, share } = props

	const Icon = technologyIcons[component.technology]

	return (
		<li className="grid gap-0.5 py-1">
			<div className="grid grid-cols-[auto_minmax(0,1fr)_auto_auto_auto] items-center gap-x-3">
				<Icon className="size-4 shrink-0" />
				<span className="truncate font-mono text-[12px]">{component.path}</span>
				<span className="flex items-center gap-1.5 text-[11px] text-muted-foreground max-sm:hidden">
					{component.technology}
					<span className="rounded-full border border-border px-1.5 text-[10px]">
						{kind}
					</span>
				</span>
				<span className="font-mono text-[12px] tabular-nums">{cpu}</span>
				<span className="font-mono text-[12px] tabular-nums">{memory}</span>
			</div>
			<ShareBar share={share} />
		</li>
	)
}
