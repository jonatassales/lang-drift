import { atom } from 'jotai'

import { initialWorkItemId, type WorkItemId } from '@data/finops-content'

// The work item selected in the FinOps ledger; the delivery receipt follows it.
export const selectedWorkItemAtom = atom<WorkItemId>(initialWorkItemId)
