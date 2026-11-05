import type { User } from './types'

// What the console API wraps its own reads in: who is signed in, and where
// its data comes from.

export interface SessionEnvelope {
	user: User | null
}

export interface SourceEnvelope {
	source: string
}
