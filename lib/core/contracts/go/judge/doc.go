// Package judge is the contract between the planner of the ex CLI and the
// judge the platform serves at <console-url>/api/judge: the question the
// planner asks, the response the console returns and the typed answer of
// each kind. The JSON Schemas under schema/ define the same shapes for the
// console; both sides test against them.
package judge
