import { tools } from './integration-tools'

export function IntegrationBranches() {
	return (
		<div aria-hidden="true" className="integration-branches">
			<svg
				aria-hidden="true"
				className="integration-branch-spine"
				focusable="false"
				preserveAspectRatio="none"
				viewBox="0 0 2 100"
			>
				<path d="M1 0v100" />
			</svg>
			{tools.map((id, index) => (
				<svg
					aria-hidden="true"
					focusable="false"
					key={id}
					preserveAspectRatio="none"
					viewBox="0 0 100 12"
				>
					{/* the middle branch runs straight out of the spine */}
					<path d={index === 1 ? 'M0 6h100' : 'M50 6h50'} />
				</svg>
			))}
		</div>
	)
}
