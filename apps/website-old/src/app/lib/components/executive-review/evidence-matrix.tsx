import {
	type ExecutiveReviewMessages,
	evidenceKinds,
	getEventSources,
	type ReviewEvent,
} from '@data/executive-review-data'

type PanelTranslate = (
	key: 'labels.week' | 'matrix.linkedCellLabel',
	values: Record<string, string | number>,
) => string

interface EvidenceMatrixProps {
	events: readonly ReviewEvent[]
	copy: ExecutiveReviewMessages
	prefix: string
	selectedSourceId: string | null
	inspectorId: string
	onOpen: (sourceId: string, trigger: HTMLButtonElement) => void
	t: PanelTranslate
}

export function EvidenceMatrix(props: EvidenceMatrixProps) {
	const { events, copy, prefix, selectedSourceId, inspectorId, onOpen, t } =
		props

	return (
		<section
			className="review-matrix"
			aria-labelledby={`${prefix}-matrix-title`}
		>
			<div className="review-panel-heading">
				<h3 id={`${prefix}-matrix-title`}>{copy.matrix.title}</h3>
				<p>{copy.matrix.description}</p>
			</div>
			<table>
				<caption className="website-sr-only">{copy.matrix.title}</caption>
				<thead>
					<tr>
						<th scope="col" id={`${prefix}-event-column`}>
							{copy.matrix.eventHeading}
						</th>
						{evidenceKinds.map((kind) => (
							<th key={kind} scope="col" id={`${prefix}-column-${kind}`}>
								{copy.kinds[kind]}
							</th>
						))}
					</tr>
				</thead>
				<tbody>
					{events.map((event) => (
						<tr
							key={event.id}
							data-review-event={event.id}
							data-inspected={
								event.sourceIds.includes(selectedSourceId ?? '') || undefined
							}
						>
							<th scope="row" id={`${prefix}-event-${event.id}`} tabIndex={-1}>
								<span className="review-event-id">
									{event.id} · {t('labels.week', { week: event.week })}
								</span>
								<strong>{copy.events[event.id].title}</strong>
								<span className="review-event-meta">
									{copy.targets[event.targetId].label} ·{' '}
									{copy.classifications[event.classification ?? 'unclassified']}
								</span>
							</th>
							{evidenceKinds.map((kind) => {
								const sources = getEventSources(event, kind)
								const source = sources[0]

								return (
									<td
										key={kind}
										headers={`${prefix}-event-${event.id} ${prefix}-column-${kind}`}
									>
										<span className="review-mobile-kind" aria-hidden="true">
											{copy.kinds[kind]}
										</span>
										{source ? (
											<button
												type="button"
												className="review-linked-cell"
												data-review-source={source.id}
												aria-label={t('matrix.linkedCellLabel', {
													eventTitle: copy.events[event.id].title,
													sourceKind: copy.kinds[kind],
													count: sources.length,
												})}
												aria-controls={inspectorId}
												aria-expanded={selectedSourceId === source.id}
												onClick={(e) => onOpen(source.id, e.currentTarget)}
											>
												<span aria-hidden="true">{sources.length} ↗</span>
											</button>
										) : (
											<span className="review-not-linked">
												{copy.matrix.notLinked}
											</span>
										)}
									</td>
								)
							})}
						</tr>
					))}
				</tbody>
			</table>
		</section>
	)
}
