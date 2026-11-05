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
	'rounded-2xl border border-border bg-background/70 backdrop-blur-xl inset-shadow-[0_1px_0] inset-shadow-foreground/14 shadow-[0_32px_80px_-20px] shadow-foreground/38 dark:inset-shadow-foreground/16 dark:shadow-[0_32px_80px_-16px] dark:shadow-background/95'

export function ElevatedFrame(props: ElevatedFrameProps) {
	const { children, className } = props

	return <div className={cn(elevated, className)}>{children}</div>
}
