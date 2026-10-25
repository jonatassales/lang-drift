import type { PlatformIconProps } from './platform-icon-props'

// A target and the point that moved away from its center.
export function VisionIcon(props: PlatformIconProps) {
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
			<circle cx="12" cy="12" r="8" />
			<path d="m12 12 3.5-3.5" />
			<circle cx="15.5" cy="8.5" fill="currentColor" r="1.4" stroke="none" />
		</svg>
	)
}
