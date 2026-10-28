export interface CaretDownIconProps {
	className?: string
	size?: number
}

// The caret a disclosure turns to show whether it is open. It draws in
// `currentColor` and hides from assistive technology: the control's
// `aria-expanded` carries the state.
export function CaretDownIcon(props: CaretDownIconProps) {
	const { className, size = 16 } = props

	return (
		<svg
			aria-hidden="true"
			className={className}
			fill="none"
			height={size}
			viewBox="0 0 24 24"
			width={size}
		>
			<path
				d="m6 9 6 6 6-6"
				stroke="currentColor"
				strokeLinecap="round"
				strokeLinejoin="round"
				strokeWidth="1.5"
			/>
		</svg>
	)
}
