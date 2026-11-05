// Package stages is the contract of the report each stage's skill writes, so
// that the run's stage lines print their figures: "console › grilling   3
// open questions settled", "console › to-spec    0011-SPEC revised · 3
// acceptance criteria" and the others. The JSON Schemas under schema/ define
// the same shapes, so a skill's author can validate a report before a run.
package stages
