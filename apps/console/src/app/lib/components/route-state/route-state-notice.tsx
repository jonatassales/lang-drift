import type { ReactNode } from 'react'

interface RouteStateNoticeProps {
	title: string
	description: string
	children: ReactNode
}

// The full-screen notice the 404 and the global error share: what happened,
// and the ways out the caller gives.
export function RouteStateNotice(props: RouteStateNoticeProps) {
	const { title, description, children } = props

	return (
		<main className="mx-auto flex min-h-screen max-w-xl flex-col justify-center gap-5 px-6">
			<h1 className="text-3xl font-semibold">{title}</h1>
			<p className="text-muted-foreground">{description}</p>
			{children}
		</main>
	)
}
