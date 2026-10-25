import type { PlatformIconProps } from './platform-icon-props'

// Three bars under one budget line.
export function FinopsIcon(props: PlatformIconProps) {
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
			<path d="M4 6h16" strokeDasharray="2.6 2.6" />
			<path d="M7 19v-6" />
			<path d="M12 19V9.5" />
			<path d="M17 19v-4" />
		</svg>
	)
}
