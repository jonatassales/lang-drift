import { ScrollHint } from '@langdrift/web/ui/scroll-hint'
import type { ReactNode } from 'react'

// Shared by every `platform/*` product page, so the bottom scroll hint is not
// wired up page by page.
export default function PlatformLayout({ children }: { children: ReactNode }) {
	return (
		<>
			{children}
			<ScrollHint />
		</>
	)
}
