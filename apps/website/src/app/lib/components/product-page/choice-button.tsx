'use client'

import { cn } from '@langdrift/web/utilities'
import type { ReactNode } from 'react'

interface ChoiceButtonProps {
	children: ReactNode
	className?: string
	kind: 'toggle' | 'step'
	onSelect: () => void
	selected: boolean
}

// A choice in a product-page demo: a work item to read, a step to inspect.
// A toggle reports aria-pressed; a step reports aria-current="step".
export function ChoiceButton(props: ChoiceButtonProps) {
	const { children, className, kind, onSelect, selected } = props

	return (
		<button
			aria-current={kind === 'step' && selected ? 'step' : undefined}
			aria-pressed={kind === 'toggle' ? selected : undefined}
			className={cn(
				'w-full cursor-pointer rounded-xl border text-left text-foreground motion-safe:transition-colors',
				'hover:bg-foreground/5 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring',
				selected && 'border-primary/60 bg-primary/10 hover:bg-primary/10',
				className,
			)}
			onClick={onSelect}
			type="button"
		>
			{children}
		</button>
	)
}
