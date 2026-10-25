'use client'

interface PillarIconProps {
	type: 'vision' | 'loop' | 'evidence'
}

export function PillarIcon(props: PillarIconProps) {
	const { type } = props

	return (
		<svg
			aria-hidden="true"
			className="vle-pillar-icon"
			width="24"
			height="24"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			strokeWidth="1.4"
			strokeLinecap="round"
			strokeLinejoin="round"
		>
			{type === 'vision' ? (
				<>
					<circle cx="12" cy="12" r="8" />
					<circle cx="12" cy="12" r="3" />
					<path d="M12 2v3m0 14v3M2 12h3m14 0h3" />
				</>
			) : type === 'loop' ? (
				<path d="M19.4 8A8 8 0 0 0 5 5L3 8m0-5v5h5M4.6 16A8 8 0 0 0 19 19l2-3m0 5v-5h-5" />
			) : (
				<path d="M7 3h7l4 4v14H6V3h1m7 0v5h4M9 12h6m-6 4h6" />
			)}
		</svg>
	)
}
