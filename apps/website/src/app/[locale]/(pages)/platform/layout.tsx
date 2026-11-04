import type { ReactNode } from 'react'

import { ScrollHint } from 'web/ui/scroll-hint'

interface PlatformLayoutProps {
	children: ReactNode
}

// Shared by every `platform/*` product page, so the bottom scroll hint is not
// wired up page by page.
export default function PlatformLayout(props: PlatformLayoutProps) {
	const { children } = props

	return (
		<>
			{children}
			<ScrollHint />
		</>
	)
}
