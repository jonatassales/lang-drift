'use client'
import { Button } from '@langdrift/web/ui/button'
import { useFormatter, useTranslations } from 'next-intl'
import { useEffect, useRef, useState } from 'react'

import { useDateTime } from '@app/template/formatters/date-time.hook'
import { ErrorNotice } from '@components/error-notice/error-notice'
import { consoleRequest } from '@data/api'
import type { RuntimeLoadTestResult, RuntimeWidgetSnapshot } from '@data/types'

import { isLoadTestResult } from './load-test-result'

interface RuntimeLoadTestProps {
	snapshot: RuntimeWidgetSnapshot
}

export function RuntimeLoadTest(props: RuntimeLoadTestProps) {
	const { snapshot } = props

	const translate = useTranslations('ui')
	const format = useFormatter()
	const dateTime = useDateTime()
	const [result, setResult] = useState<RuntimeLoadTestResult | null>(null)
	const [error, setError] = useState<unknown>(null)
	const [pending, setPending] = useState(false)
	const controller = useRef<AbortController | null>(null)
	const saved = snapshot.runtime.loadTest
	const report =
		result && (!saved || result.completedAt > saved.completedAt)
			? result
			: saved
	const running = snapshot.executions.some(
		(execution) => execution.state === 'running',
	)
	const completedAt = report?.completedAt
	const [cooling, setCooling] = useState(!!completedAt)
	useEffect(() => {
		const remaining = completedAt
			? Date.parse(completedAt) + 10000 - Date.now()
			: 0
		setCooling(remaining > 0)
		if (remaining <= 0) return
		const timer = window.setTimeout(() => setCooling(false), remaining)

		return () => window.clearTimeout(timer)
	}, [completedAt])
	useEffect(() => () => controller.current?.abort(), [])
	const actionLabel = pending
		? 'loadTestRunning'
		: cooling
			? 'loadTestCooldown'
			: 'runLoadTest'

	async function run() {
		if (pending) return
		const request = new AbortController()
		controller.current = request
		setPending(true)
		setError(null)
		try {
			const response = await consoleRequest<{
				loadTest: RuntimeLoadTestResult
			}>(`runtimes/${snapshot.runtime.id}/load-test`, {
				method: 'POST',
				body: '{}',
				signal: request.signal,
			})

			if (!isLoadTestResult(response.loadTest))
				throw new Error('invalid_load_test')
			if (!request.signal.aborted) setResult(response.loadTest)
		} catch (failure) {
			if (!request.signal.aborted) setError(failure)
		} finally {
			if (!request.signal.aborted) setPending(false)
		}
	}

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
						onClick={run}
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
