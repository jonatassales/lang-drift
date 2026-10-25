export interface ArrowRightIconProps {
	className?: string
	size?: number
}

// The arrow the website's calls to action and flows point with. It draws in
// `currentColor` and hides from assistive technology: the text beside it
// carries the meaning.
export function ArrowRightIcon(props: ArrowRightIconProps) {
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
				d="M4 12h16m-6-6 6 6-6 6"
				stroke="currentColor"
				strokeLinecap="round"
				strokeLinejoin="round"
				strokeWidth="1.5"
			/>
		</svg>
	)
}
