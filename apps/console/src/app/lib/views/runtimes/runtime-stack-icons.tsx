import { StackIcon } from '@ui/stack-icon'
import { Cube } from 'web/ui/icons/client'

interface RuntimeStackIconsProps {
	stacks?: string[]
	size?: number
}

export function RuntimeStackIcons(props: RuntimeStackIconsProps) {
	const { stacks, size = 16 } = props

	if (!stacks?.length)
		return <Cube size={size} aria-hidden className="shrink-0" />

	return (
		<span className="inline-flex shrink-0 items-center gap-1.5">
			{Array.from(new Set(stacks))
				.slice(0, 4)
				.map((stack) => (
					<StackIcon key={stack} stack={stack} size={size} />
				))}
		</span>
	)
}
