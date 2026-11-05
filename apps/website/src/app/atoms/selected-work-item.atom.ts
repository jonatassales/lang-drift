import { atom } from 'jotai'

import { initialWorkItemId } from '@routes/platform/finops/content/finops-content.data'
import type { WorkItemId } from '@routes/platform/finops/content/finops-content.types'

// The work item selected in the FinOps ledger; the delivery receipt follows it.
export const selectedWorkItemAtom = atom<WorkItemId>(initialWorkItemId)
