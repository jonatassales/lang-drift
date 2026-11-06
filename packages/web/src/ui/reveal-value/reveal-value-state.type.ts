// Where a value stands as it is revealed: not known yet, being resolved,
// or known. A value not known is unset, never zero.
export type RevealValueState = 'unset' | 'resolving' | 'set'
