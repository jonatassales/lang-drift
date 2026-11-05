import type { ReactNode } from 'react'

import { ScrollHint } from 'web/ui/scroll-hint'

// Shared by every `platform/*` product page, so the bottom scroll hint is not
// wired up page by page.
interface PlatformLayoutProps {
	children: ReactNode
}

export default function PlatformLayout(props: PlatformLayoutProps) {
	const { children } = props

	return (
		<>
			{children}
			<ScrollHint />
		</>
	)
}
