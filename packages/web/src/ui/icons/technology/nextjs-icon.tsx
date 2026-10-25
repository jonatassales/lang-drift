import type { TechnologyIconProps } from './technology-icon-props'

// The Next.js mark: an N inside a circle. The geometry is the official mark's, from
// simple-icons (CC0 1.0), drawn in currentColor.
export function NextjsIcon(props: TechnologyIconProps) {
	const { className } = props

	return (
		<svg
			aria-hidden="true"
			className={className}
			fill="currentColor"
			viewBox="0 0 24 24"
		>
			<path d="M18.665 21.978C16.758 23.255 14.465 24 12 24 5.377 24 0 18.623 0 12S5.377 0 12 0s12 5.377 12 12c0 3.583-1.574 6.801-4.067 9.001L9.219 7.2H7.2v9.596h1.615V9.251l9.85 12.727Zm-3.332-8.533 1.6 2.061V7.2h-1.6v6.245Z" />
		</svg>
	)
}
