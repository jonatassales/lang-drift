import { atom } from 'jotai'

import { initialStageId } from '@routes/platform/microruntimes/content/microruntimes-content.data'
import type { LifecycleStageId } from '@routes/platform/microruntimes/content/microruntimes-content.types'

// The lifecycle step selected in the Microruntimes hero; the evidence section
// follows it.
export const selectedLifecycleStageAtom = atom<LifecycleStageId>(initialStageId)
