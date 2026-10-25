import type { PlatformIconProps } from './platform-icon-props'

// A voice, as five bars of sound.
export function AiVoiceIcon(props: PlatformIconProps) {
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
			<path d="M5 10v4" />
			<path d="M8.5 7v10" />
			<path d="M12 4v16" />
			<path d="M15.5 8v8" />
			<path d="M19 10v4" />
		</svg>
	)
}
