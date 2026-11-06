// A check in the success role beside something finished: read, done or
// locked. The check itself is decorative; screen readers get its label.
export interface CheckMarkProps {
	/** What the mark means, for screen readers: read, done, locked. */
	label: string
}

export function CheckMark(props: CheckMarkProps) {
	const { label } = props

	return (
		<span className="inline-flex items-center font-mono text-[12px] text-status-success">
			<span aria-hidden="true">✓</span>
			<span className="sr-only">{label}</span>
		</span>
	)
}
