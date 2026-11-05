import type { ReactNode } from 'react'

// The first focusable element of a page: hidden until focused, it jumps past
// the header to the page's first section.
interface SkipLinkProps {
	children: ReactNode
	href: string
}

export function SkipLink(props: SkipLinkProps) {
	const { children, href } = props

	return (
		<a
			className="fixed top-0 left-3 z-200 max-w-[calc(100vw-1.5rem)] -translate-y-full bg-foreground px-2.5 py-2 whitespace-normal text-card focus-visible:top-3 focus-visible:translate-y-0 focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
			href={href}
		>
			{children}
		</a>
	)
}
