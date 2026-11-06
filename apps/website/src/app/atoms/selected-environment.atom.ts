import { atom } from 'jotai'

import { initialEnvironmentId } from '@routes/platform/finops/content/finops-content.data'
import type { EnvironmentId } from '@routes/platform/finops/content/finops-content.types'

// The environment selected in the FinOps ledger; the environment receipt
// follows it.
export const selectedEnvironmentAtom = atom<EnvironmentId>(initialEnvironmentId)
