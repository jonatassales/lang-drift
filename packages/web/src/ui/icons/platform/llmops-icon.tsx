import type { PlatformIconProps } from './platform-icon-props'

// A spark of reasoning, the work a model call is paid for.
export function LlmopsIcon(props: PlatformIconProps) {
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
			<path d="M12 4c.6 4.2 3.8 7.4 8 8-4.2.6-7.4 3.8-8 8-.6-4.2-3.8-7.4-8-8 4.2-.6 7.4-3.8 8-8Z" />
		</svg>
	)
}
