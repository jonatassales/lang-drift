import { atom } from 'jotai'

import { initialWorkItemId } from '@routes/platform/llmops/content/llmops-content.data'
import type { WorkItemId } from '@routes/platform/llmops/content/llmops-content.types'

// The work item selected in the LLMOps ledger; the delivery receipt follows it.
export const selectedWorkItemAtom = atom<WorkItemId>(initialWorkItemId)
