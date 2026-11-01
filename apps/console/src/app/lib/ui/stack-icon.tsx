import { Cube } from '@langdrift/web/ui/icons/client'

import icons from './stack-icons.json'

interface StackIconProps {
	stack: string
	size?: number
}

/** Locally bundled brand paths; unknown stacks retain the neutral component icon. */
export function StackIcon(props: StackIconProps) {
	const { stack, size = 16 } = props

	if (!Object.hasOwn(icons, stack)) return <Cube size={size} aria-hidden />
	const icon = icons[stack as keyof typeof icons]

	return (
		<svg
			viewBox="0 0 24 24"
			width={size}
			height={size}
			fill="currentColor"
			role="img"
			aria-label={icon.title}
			data-stack={stack}
			className="shrink-0"
		>
			<title>{icon.title}</title>
			<path d={icon.path} />
		</svg>
	)
}
