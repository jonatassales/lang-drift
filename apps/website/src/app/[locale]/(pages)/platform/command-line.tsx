interface CommandLineProps {
	command: string
}

// One command to type, set as code on a card so it reads as a terminal line.
export function CommandLine(props: CommandLineProps) {
	const { command } = props

	return (
		<code className="block w-fit max-w-full overflow-x-auto rounded-lg border border-border bg-card px-4 py-3 font-mono text-sm whitespace-nowrap text-foreground">
			{command}
		</code>
	)
}
