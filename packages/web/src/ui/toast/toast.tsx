import type { ReactNode } from 'react'

import { cn } from 'web/utilities'

export interface ToastProps {
	children: ReactNode
	className?: string
	onClose?: () => void
	variant?: 'danger' | 'info' | 'success'
}

export function Toast(props: ToastProps) {
	const { children, className, onClose, variant = 'danger' } = props

	function handleClose() {
		onClose?.()
	}

	return (
		<div
			aria-live="assertive"
			className={cn(
				'fixed right-6 bottom-6 z-50 flex max-w-md items-center gap-3 rounded-xl border px-4 py-3 text-sm shadow-lg backdrop-blur-md transition-all',
				variant === 'danger' &&
					'border-status-danger/40 bg-card text-foreground',
				variant === 'success' &&
					'border-status-positive/40 bg-card text-foreground',
				variant === 'info' && 'border-border bg-card text-foreground',
				className,
			)}
			role="alert"
		>
			<div className="flex-1">{children}</div>
			{onClose ? (
				<button
					aria-label="Close notification"
					className="rounded-md p-1 text-muted-foreground transition hover:text-foreground focus-visible:outline-2 focus-visible:outline-primary"
					onClick={handleClose}
					type="button"
				>
					<span aria-hidden="true">&times;</span>
				</button>
			) : null}
		</div>
	)
}
