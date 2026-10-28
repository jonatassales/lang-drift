import { cn } from '@langdrift/web/utilities'
import type { ButtonHTMLAttributes } from 'react'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
	tone?: 'primary' | 'secondary' | 'danger'
}

export function Button(props: ButtonProps) {
	const {
		children,
		className,
		tone = 'secondary',
		type = 'button',
		...rest
	} = props

	return (
		<button
			type={type}
			className={cn(
				'inline-flex min-h-9 cursor-pointer items-center justify-center gap-2 rounded-md border px-3 py-2 text-xs font-medium transition-colors disabled:cursor-not-allowed disabled:opacity-50',
				tone === 'primary'
					? 'border-primary bg-primary text-primary-foreground hover:bg-primary-hover'
					: tone === 'danger'
						? 'border-status-danger/30 text-status-danger hover:bg-status-danger/5'
						: 'border-border bg-card text-foreground hover:bg-muted',
				className,
			)}
			{...rest}
		>
			{children}
		</button>
	)
}
