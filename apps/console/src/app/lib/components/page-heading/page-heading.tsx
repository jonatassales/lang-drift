import type { ReactNode } from 'react'

interface PageHeadingProps {
	title: string
	description?: string
	children?: ReactNode
}
export function PageHeading(props: PageHeadingProps) {
	const { title, description, children } = props

	return (
		<div className="mb-6 flex flex-wrap items-start justify-between gap-4">
			<div className="grid gap-2">
				<h1 className="text-[24px] font-medium tracking-tight">{title}</h1>
				{description ? (
					<p className="max-w-2xl text-sm leading-relaxed text-muted-foreground">
						{description}
					</p>
				) : null}
			</div>
			{children}
		</div>
	)
}
