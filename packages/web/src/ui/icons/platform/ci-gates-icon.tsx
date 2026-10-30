import type { PlatformIconProps } from './platform-icon-props'

// Two gate posts with horizontal rails and cross braces.
export function CiGatesIcon(props: PlatformIconProps) {
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
			<path d="M5 4v16" />
			<path d="M19 4v16" />
			<path d="M5 8h14" />
			<path d="M5 16h14" />
			<path d="m5 8 14 8" />
			<path d="m5 16 14-8" />
		</svg>
	)
}
