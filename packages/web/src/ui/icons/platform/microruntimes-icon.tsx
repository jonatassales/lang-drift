import type { PlatformIconProps } from './platform-icon-props'

// A bounded frame with the workload it runs inside it.
export function MicroruntimesIcon(props: PlatformIconProps) {
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
			<rect height="16" rx="3.5" width="16" x="4" y="4" />
			<rect height="6" rx="1.5" width="6" x="9" y="9" />
		</svg>
	)
}
