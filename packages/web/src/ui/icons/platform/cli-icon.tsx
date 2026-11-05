import type { PlatformIconProps } from './platform-icon-props'

// A terminal window with a prompt and a cursor.
export function CliIcon(props: PlatformIconProps) {
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
			<rect height="16" rx="2.5" width="18" x="3" y="4" />
			<path d="m7.5 9.5 3 2.5-3 2.5" />
			<path d="M12.5 15h4" />
		</svg>
	)
}
