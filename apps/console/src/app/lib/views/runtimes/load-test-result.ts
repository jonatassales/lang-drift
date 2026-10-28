import type { RuntimeLoadTestResult } from '@data/types'

/** Validate measured results at both HTTP and WebRTC consumption boundaries. */
export function isLoadTestResult(
	value: unknown,
): value is RuntimeLoadTestResult {
	if (!value || typeof value !== 'object') return false
	const report = value as Record<string, unknown>
	const counts = [report.succeeded, report.failed]
	const times = [report.durationMs, report.p50Ms, report.p95Ms, report.maxMs]
	if (
		report.target !== 'extub-loopback' ||
		report.requests !== 50 ||
		report.concurrency !== 50 ||
		!counts.every(
			(count) =>
				typeof count === 'number' &&
				Number.isInteger(count) &&
				count >= 0 &&
				count <= 50,
		) ||
		!times.every(
			(time) =>
				typeof time === 'number' &&
				Number.isFinite(time) &&
				time >= 0 &&
				time <= 15000,
		) ||
		typeof report.startedAt !== 'string' ||
		typeof report.completedAt !== 'string'
	)
		return false

	return (
		(report.succeeded as number) + (report.failed as number) === 50 &&
		(report.p50Ms as number) <= (report.p95Ms as number) &&
		(report.p95Ms as number) <= (report.maxMs as number) &&
		(report.maxMs as number) <= (report.durationMs as number) &&
		Number.isFinite(Date.parse(report.startedAt)) &&
		Number.isFinite(Date.parse(report.completedAt)) &&
		Date.parse(report.completedAt) >= Date.parse(report.startedAt)
	)
}
