'use client'

import { useRef, useState } from 'react'

import { Check, Copy } from 'web/ui/icons/client'
import { cn } from 'web/utilities'

import type { CommandLineLabels } from './command-line-labels.type'

type CopyState = 'idle' | 'copied' | 'failed'

const resetAfterMs = 2000

interface CopyCommandButtonProps {
	command: string
	commandId?: string
	labels: CommandLineLabels
}

// Copies a command to the clipboard and says so for a moment: a check when
// the copy took, the failure label when the browser refused it.
export function CopyCommandButton(props: CopyCommandButtonProps) {
	const { command, commandId, labels } = props

	const [state, setState] = useState<CopyState>('idle')
	const resetTimer = useRef<ReturnType<typeof setTimeout>>(undefined)

	async function handleClick() {
		clearTimeout(resetTimer.current)

		try {
			await navigator.clipboard.writeText(command)

			setState('copied')
		} catch {
			setState('failed')
		}

		resetTimer.current = setTimeout(() => {
			setState('idle')
		}, resetAfterMs)
	}

	const labelOf: Record<CopyState, string> = {
		idle: labels.copy,
		copied: labels.copied,
		failed: labels.failed,
	}
	const label = labelOf[state]

	return (
		<>
			<button
				aria-describedby={commandId}
				aria-label={label}
				className={cn(
					'inline-flex size-8 shrink-0 cursor-pointer items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary',
					state === 'copied' && 'text-status-success hover:text-status-success',
					state === 'failed' && 'text-status-danger hover:text-status-danger',
				)}
				onClick={handleClick}
				title={label}
				type="button"
			>
				{state === 'copied' ? (
					<Check aria-hidden="true" className="size-4" />
				) : (
					<Copy aria-hidden="true" className="size-4" />
				)}
			</button>
			<span aria-live="polite" className="sr-only">
				{state === 'idle' ? '' : label}
			</span>
		</>
	)
}
