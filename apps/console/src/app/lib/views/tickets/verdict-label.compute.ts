import type { TicketEpisode } from '@data/types'

import type { TicketFormat } from './ticket-format.hook'

export function verdictLabel(
	episode: TicketEpisode,
	format: TicketFormat,
): string {
	const { translate } = format
	const label = translate(`analysis.verdicts.${episode.verdict}`)
	if (!episode.reason) return label

	return `${label} · ${translate(`analysis.reasons.${episode.reason}`)}`
}
