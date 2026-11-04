import { sourceIds, sourceLinks } from './safe-seat-content.data'
import type { SafeSeatSourcesCopy } from './safe-seat-messages.types'

// The full citation of every source the page draws on, each linked to its
// primary record.
interface SourceCitationsProps {
	copy: SafeSeatSourcesCopy
}

export function SourceCitations(props: SourceCitationsProps) {
	const { copy } = props

	return (
		<ol className="m-0 mt-11 grid max-w-240 list-none gap-0 p-0">
			{sourceIds.map((id) => (
				<li className="border-t border-border py-4" key={id}>
					<a
						className="text-[13px] underline decoration-current/25 underline-offset-5 leading-[1.6] focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
						href={sourceLinks[id]}
						rel="noopener noreferrer"
						target="_blank"
					>
						{copy.items[id]}
					</a>
				</li>
			))}
		</ol>
	)
}
