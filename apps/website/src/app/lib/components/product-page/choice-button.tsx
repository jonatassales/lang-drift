'use client'

import { cn } from '@langdrift/web/utilities'
import type { MouseEvent, ReactNode } from 'react'

interface ChoiceButtonProps {
	children: ReactNode
	className?: string
	'data-id'?: string
	kind: 'toggle' | 'step'
	onSelect: (event: MouseEvent<HTMLButtonElement>) => void
	selected: boolean
}

// A choice in a product-page demo: a work item to read, a step to inspect.
// A toggle reports aria-pressed; a step reports aria-current="step".
export function ChoiceButton(props: ChoiceButtonProps) {
	const {
		children,
		className,
		'data-id': dataId,
		kind,
		onSelect,
		selected,
	} = props

	return (
		<button
			aria-current={kind === 'step' && selected ? 'step' : undefined}
			aria-pressed={kind === 'toggle' ? selected : undefined}
			className={cn(
				'w-full cursor-pointer rounded-xl border text-left whitespace-nowrap text-foreground motion-safe:transition-colors',
				'hover:bg-foreground/5 focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary',
				selected && 'border-primary/60 bg-primary/10 hover:bg-primary/10',
				className,
			)}
			data-id={dataId}
			onClick={onSelect}
			type="button"
		>
			{children}
		</button>
	)
}
