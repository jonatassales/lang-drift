import { CardValue, type CardValueLabels } from './card-value'
import { ReadMark } from './read-mark'

const note =
	'text-[11px] text-muted-foreground motion-safe:transition-opacity motion-safe:duration-300 starting:opacity-0'

// One linked document in the analysis, on one line: its file name, a dash
// until the agent reads it, then the read mark and what the read
// extracted, with the agent's judgment of it beneath once made.
interface DocumentRowProps {
	/** What the read extracted from the document. */
	extracted: string
	file: string
	/** The agent's judgment of the document, once it has made one. */
	judgment?: string
	labels: CardValueLabels & { readMark: string }
	read: boolean
}

export function DocumentRow(props: DocumentRowProps) {
	const { extracted, file, judgment, labels, read } = props

	return (
		<li className="grid gap-0.5 py-1.5">
			<span className="flex flex-wrap items-baseline gap-x-2">
				<span className="font-mono text-[12px]">{file}</span>
				{read ? (
					<>
						<ReadMark label={labels.readMark} />
						<span className={note}>{extracted}</span>
					</>
				) : (
					<CardValue className="text-[11px]" labels={labels} state="unset" />
				)}
			</span>
			{judgment && <span className={note}>{judgment}</span>}
		</li>
	)
}
