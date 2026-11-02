import { ScrollHint } from '@langdrift/web/ui/scroll-hint'
import type { ReactNode } from 'react'

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
