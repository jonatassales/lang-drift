import { atom } from 'jotai'

import {
	initialStageId,
	type LifecycleStageId,
} from '@routes/platform/microruntimes/content/microruntimes-content'

// The lifecycle step selected in the Microruntimes hero; the evidence section
// follows it.
export const selectedLifecycleStageAtom = atom<LifecycleStageId>(initialStageId)
