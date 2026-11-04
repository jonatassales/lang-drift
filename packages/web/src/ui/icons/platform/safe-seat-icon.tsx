import type { PlatformIconProps } from './platform-icon-props'

// A seat seen from the side: the backrest, the cushion and two legs.
export function SafeSeatIcon(props: PlatformIconProps) {
	const { className } = props

	return (
		<svg
			aria-hidden="true"
			className={className}
			fill="none"
			stroke="currentColor"
			strokeLinecap="round"
			strokeLinejoin="round"
			strokeWidth="1.6"
			viewBox="0 0 24 24"
		>
			<path d="M7 4a1 1 0 0 1 1-1h1a1 1 0 0 1 1 1v8" />
			<path d="M7 4v9a2 2 0 0 0 2 2h8a1 1 0 0 0 1-1v-1a1 1 0 0 0-1-1h-7" />
			<path d="m9 15-1 5" />
			<path d="m16 15 1 5" />
		</svg>
	)
}
