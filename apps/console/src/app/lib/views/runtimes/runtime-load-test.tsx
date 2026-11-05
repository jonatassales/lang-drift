'use client'
import { useFormatter, useTranslations } from 'next-intl'
import { useEffect, useMemo, useRef, useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { consoleRequest } from '@data/api'
import type { LoadTestEnvelope } from '@data/runtime-envelopes.types'
import type { RuntimeLoadTestResult, RuntimeWidgetSnapshot } from '@data/types'
import { useDateTime } from '@i18n/use-date-time.hook'
import { Button } from 'web/ui/button'

import { isLoadTestResult } from './load-test-result'

interface RuntimeLoadTestProps {
	snapshot: RuntimeWidgetSnapshot
}

export function RuntimeLoadTest(props: RuntimeLoadTestProps) {
	const { snapshot } = props

	const [pending, startTransition] = useTransition()
	const [result, setResult] = useState<RuntimeLoadTestResult | null>(null)
	const [error, setError] = useState<unknown>(null)
	const [cooling, setCooling] = useState(
		Boolean(snapshot.runtime.loadTest?.completedAt),
	)
	const controller = useRef<AbortController | null>(null)
	const saved = snapshot.runtime.loadTest
	const report = useMemo(
		() =>
			result && (!saved || result.completedAt > saved.completedAt)
				? result
				: saved,
		[result, saved],
	)
	const completedAt = report?.completedAt

	const dateTime = useDateTime()

	const translate = useTranslations('ui')
	const format = useFormatter()

	// Manage cooldown timer after a load test run completes.
	useEffect(() => {
		const remaining = completedAt
			? Date.parse(completedAt) + 10000 - Date.now()
			: 0
		setCooling(remaining > 0)
		if (remaining <= 0) return
		const timer = window.setTimeout(() => setCooling(false), remaining)

		return () => window.clearTimeout(timer)
	}, [completedAt])

	// Abort in-flight load test request on unmount.
	useEffect(() => () => controller.current?.abort(), [])

	function handleRun() {
		if (pending) return
		const request = new AbortController()
		controller.current = request
		setError(null)
		startTransition(async () => {
			try {
				const response = await consoleRequest<LoadTestEnvelope>(
					`runtimes/${snapshot.runtime.id}/load-test`,
					{
						method: 'POST',
						body: '{}',
						signal: request.signal,
					},
				)

				if (!isLoadTestResult(response.loadTest))
					throw new Error('invalid_load_test')
				if (!request.signal.aborted) setResult(response.loadTest)
			} catch (failure) {
				if (!request.signal.aborted) setError(failure)
			}
		})
	}

	const running = snapshot.executions.some(
		(execution) => execution.state === 'running',
	)
	const actionLabel = pending
		? 'loadTestRunning'
		: cooling
			? 'loadTestCooldown'
			: 'runLoadTest'

	return (
		<section className="rounded-lg border border-border p-4" data-runtime-tests>
			<div className="flex flex-wrap items-center justify-between gap-3">
				<h2 className="text-sm font-medium">{translate('loadTestTitle')}</h2>
				{snapshot.source === 'stub' ? (
					<Button
						disabled={
							pending ||
							!running ||
							cooling ||
							snapshot.runtime.credentialRevoked
						}
						onClick={handleRun}
						type="button"
						variant="outline"
					>
						{translate(actionLabel)}
					</Button>
				) : null}
			</div>
			<p className="mt-3 text-xs leading-relaxed text-muted-foreground">
				{translate(
					snapshot.source === 'stub'
						? 'loadTestScope'
						: 'testEvidenceUnavailable',
				)}
			</p>
			{error ? (
				<div className="mt-3">
					<ErrorNotice error={error} />
				</div>
			) : null}
			{report ? (
				<div role="status" className="mt-4 text-xs">
					<p className="font-medium">
						{translate('loadTestOutcome', {
							succeeded: report.succeeded,
							failed: report.failed,
							requests: report.requests,
						})}
					</p>
					<dl className="mt-3 grid grid-cols-2 gap-x-5 gap-y-2">
						<dt>{translate('loadTestDuration')}</dt>
						<dd>
							{format.number(report.durationMs, { maximumFractionDigits: 2 })}{' '}
							ms
						</dd>
						<dt>P50</dt>
						<dd>
							{format.number(report.p50Ms, { maximumFractionDigits: 2 })} ms
						</dd>
						<dt>P95</dt>
						<dd>
							{format.number(report.p95Ms, { maximumFractionDigits: 2 })} ms
						</dd>
						<dt>{translate('loadTestMaximum')}</dt>
						<dd>
							{format.number(report.maxMs, { maximumFractionDigits: 2 })} ms
						</dd>
					</dl>
					<p className="mt-3 text-[10px] text-muted-foreground">
						{dateTime(new Date(report.completedAt), {
							dateStyle: 'short',
							timeStyle: 'medium',
						})}{' '}
						· {translate('loadTestObserved')}
					</p>
				</div>
			) : snapshot.source === 'stub' ? (
				<p className="mt-3 text-xs text-muted-foreground">
					{translate('loadTestEmpty')}
				</p>
			) : null}
		</section>
	)
}
