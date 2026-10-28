'use client'

import { Moon, Sun } from './icons'

interface ThemeToggleProps {
	label: string
	onToggle: () => void
	theme: 'dark' | 'light'
}

export function ThemeToggle(props: ThemeToggleProps) {
	const { label, onToggle, theme } = props

	return (
		<button
			aria-label={label}
			title={label}
			className="inline-flex h-[34px] w-[34px] cursor-pointer items-center justify-center rounded-lg border-0 bg-transparent text-inherit hover:bg-muted"
			onClick={onToggle}
			type="button"
		>
			{theme === 'dark' ? (
				<Sun aria-hidden size={17} />
			) : (
				<Moon aria-hidden size={17} />
			)}
		</button>
	)
}
