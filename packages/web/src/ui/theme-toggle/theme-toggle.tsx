'use client'

import { Moon, Sun } from '@phosphor-icons/react'

export interface ThemeToggleProps {
	label: string
	onToggle: () => void
	theme: 'dark' | 'light'
}

// The button that switches between the dark and the light scheme. It shows
// the scheme it switches to: the sun while dark, the moon while light.
export function ThemeToggle(props: ThemeToggleProps) {
	const { label, onToggle, theme } = props

	const Icon = theme === 'dark' ? Sun : Moon

	return (
		<button
			aria-label={label}
			className="inline-flex h-[34px] w-[34px] cursor-pointer items-center justify-center rounded-lg border-0 bg-transparent text-inherit hover:bg-muted focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
			onClick={onToggle}
			title={label}
			type="button"
		>
			<Icon aria-hidden size={17} />
		</button>
	)
}
