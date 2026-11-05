// Package handoff is what the console receives from a plan and a run, as
// files: the plan's card data (plan.json), the recorded intent (intent.json)
// and the run's record (run.json). Every figure is in its unit, integers for
// millicores, MiB, cents and tokens and two decimals for similarities, and a
// value not known is null, never zero. A section absent from a record has
// not been revealed yet, as its card reads "Not set yet". The fields inside
// each section join these types as the element specs that fill them land.
package handoff
