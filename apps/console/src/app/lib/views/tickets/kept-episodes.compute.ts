import type { TicketEpisode } from '@data/types'

export function keptEpisodes(episodes: readonly TicketEpisode[]) {
	return {
		kept: episodes.filter((episode) => episode.verdict === 'kept').length,
		total: episodes.length,
	}
}
