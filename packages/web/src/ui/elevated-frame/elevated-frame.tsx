import type { ReactNode } from 'react'

import { cn } from 'web/utilities'

export interface ElevatedFrameProps {
	children: ReactNode
	className?: string
}

// Glass that sits above the page: a translucent, blurred surface with a lit
// top edge and a deep, soft shadow under it, in the foreground's gray over a
// light field and the background's black over a dark one, where the lit edge
// does the lifting.
const elevated =
	'rounded-2xl border border-border bg-background/70 backdrop-blur-xl shadow-[inset_0_1px_0_color-mix(in_srgb,var(--foreground)_14%,transparent),0_32px_80px_-20px_color-mix(in_srgb,var(--foreground)_38%,transparent)] dark:shadow-[inset_0_1px_0_color-mix(in_srgb,var(--foreground)_16%,transparent),0_32px_80px_-16px_color-mix(in_srgb,var(--background)_95%,transparent)]'

export function ElevatedFrame(props: ElevatedFrameProps) {
	const { children, className } = props

	return <div className={cn(elevated, className)}>{children}</div>
}
