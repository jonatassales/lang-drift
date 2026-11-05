import { useId } from 'react'

import type { CommandLineLabels } from './command-line-labels.type'
import { CopyCommandButton } from './copy-command-button'

interface CommandLineProps {
	command: string
	labels: CommandLineLabels
}

// One command to type, set as code on a card so it reads as a terminal line,
// with a button that copies it and names the command it copies.
export function CommandLine(props: CommandLineProps) {
	const { command, labels } = props

	const commandId = useId()

	return (
		<div className="flex w-fit max-w-full items-center rounded-lg border border-border bg-card py-1.5 pr-1 pl-3">
			<code
				className="min-w-0 overflow-x-auto py-1.5 font-mono text-sm whitespace-nowrap text-foreground"
				id={commandId}
			>
				{command}
			</code>
			<CopyCommandButton
				command={command}
				commandId={commandId}
				labels={labels}
			/>
		</div>
	)
}
