import type { Episode, LifecycleEvent } from './types'

// What the console API wraps a history read in: lifecycle events and episodes.

export interface LifecycleEventsEnvelope {
	events: LifecycleEvent[]
}

export interface EpisodesEnvelope {
	episodes: Episode[]
}
