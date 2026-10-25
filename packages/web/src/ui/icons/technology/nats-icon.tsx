import type { TechnologyIconProps } from './technology-icon-props'

// The NATS mark: an N inside its speech-bubble frame. The geometry is the official mark's, from
// simple-icons (CC0 1.0), drawn in currentColor.
export function NatsIcon(props: TechnologyIconProps) {
	const { className } = props

	return (
		<svg
			aria-hidden="true"
			className={className}
			fill="currentColor"
			viewBox="0 0 24 24"
		>
			<path d="M12.004 0H.404v18.807h9.938l1.714 1.602v-.026L15.966 24v-5.193h7.63V0H12.003zm7.578 14.45H15.38L6.898 6.519v7.93H4.116V4.376h4.349l8.344 7.784V4.375h2.773V14.45z" />
		</svg>
	)
}
