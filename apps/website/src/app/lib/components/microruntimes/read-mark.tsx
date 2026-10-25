interface ReadMarkProps {
	label: string
}

// The mark the ticket card puts beside an artifact once the agent has read
// it: read, not approved.
export function ReadMark(props: ReadMarkProps) {
	const { label } = props

	return (
		<span className="inline-flex items-center font-mono text-[12px] text-status-success">
			<span aria-hidden="true">✓</span>
			<span className="sr-only">{label}</span>
		</span>
	)
}
