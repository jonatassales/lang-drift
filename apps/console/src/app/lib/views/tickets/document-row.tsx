import { ReadMark } from './read-mark'
import { UnknownValue } from './unknown-value'

interface DocumentRowProps {
	/** What the read extracted; null while the document is unread. */
	extracted: string | null
	file: string
	/** The agent's judgment of the document, once it has made one. */
	judgment?: string
}

const note = 'text-[11px] text-muted-foreground'

// One linked document in the analysis: its file name, then a dash until
// the agent reads it, or the read mark and what the read extracted, with
// the agent's judgment beneath once made.
export function DocumentRow(props: DocumentRowProps) {
	const { extracted, file, judgment } = props

	return (
		<li className="grid gap-0.5 py-1.5">
			<span className="flex flex-wrap items-baseline gap-x-2">
				<span className="font-mono text-[12px] break-all">{file}</span>
				{extracted === null ? (
					<UnknownValue className="text-[11px]" />
				) : (
					<>
						<ReadMark />
						<span className={note}>{extracted}</span>
					</>
				)}
			</span>
			{judgment && <span className={note}>{judgment}</span>}
		</li>
	)
}
