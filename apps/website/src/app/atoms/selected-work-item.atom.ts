import { atom } from 'jotai'

import {
	initialWorkItemId,
	type WorkItemId,
} from '@routes/platform/finops/content/finops-content.data'

// The work item selected in the FinOps ledger; the delivery receipt follows it.
export const selectedWorkItemAtom = atom<WorkItemId>(initialWorkItemId)
