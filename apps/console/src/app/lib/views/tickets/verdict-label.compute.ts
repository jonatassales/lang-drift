import type { TicketEpisode } from '@data/types'

import type { TicketFormat } from './ticket-format.hook'

export function verdictLabel(
	episode: TicketEpisode,
	format: TicketFormat,
): string {
	const { t } = format
	const label = t(`analysis.verdicts.${episode.verdict}`)
	if (!episode.reason) return label

	return `${label} · ${t(`analysis.reasons.${episode.reason}`)}`
}
