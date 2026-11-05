// Package judgment is the contract of the three answers the plan asks its
// judge for, as the simulation shows them: the guardrail's criteria, the
// verdict on each retrieved episode, and the suggested model and effort. No
// answer carries a price, since the plan pays for nothing. The JSON Schemas
// under schema/ define the same shapes; both sides test against them. It
// stands beside the older judge package until ticket 0325 removes that one.
package judgment
