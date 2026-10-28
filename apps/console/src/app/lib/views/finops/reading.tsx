import type { ReactNode } from 'react'

interface ReadingProps {
	children: ReactNode
	label: string
}

// One labelled value in the cost overview's runtime readings.
export function Reading(props: ReadingProps) {
	const { children, label } = props

	return (
		<div className="grid gap-1">
			<dt className="text-muted-foreground">{label}</dt>
			<dd>{children}</dd>
		</div>
	)
}
