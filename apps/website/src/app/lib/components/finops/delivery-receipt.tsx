'use client'

import { DetailList } from '@langdrift/web/ui/detail-list'
import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { useAtomValue } from 'jotai'
import { createTranslator } from 'next-intl'

import { dollars, type FinopsMessages, workItems } from '@data/finops-content'
import type { WebsiteLocale } from '@i18n/routing'
import { selectedWorkItemAtom } from '@state'

interface DeliveryReceiptProps {
	copy: FinopsMessages['receipt']
	items: FinopsMessages['ledger']['items']
	locale: WebsiteLocale
}

// The receipt of the work item selected in the ledger: what it was estimated,
// what was reserved, what it consumed, and how it was verified and cleaned up.
export function DeliveryReceipt(props: DeliveryReceiptProps) {
	const { copy, items, locale } = props

	const selectedId = useAtomValue(selectedWorkItemAtom)
	const item = workItems.find((workItem) => workItem.id === selectedId)
	if (!item) return null

	const receipt = createTranslator({
		locale,
		messages: copy,
	})
	const money = dollars(locale)
	const itemCopy = items[item.id]

	return (
		<article
			aria-atomic="true"
			aria-live="polite"
			className="mt-11 rounded-2xl border border-border bg-card p-6 max-md:p-[18px]"
		>
			<div className="flex items-start justify-between gap-6 max-md:flex-col">
				<div className="grid gap-[7px]">
					<Eyebrow>{copy.selected}</Eyebrow>
					<h3 className="m-0 text-[28px] font-medium tracking-[-0.025em]">
						{itemCopy.title}
					</h3>
					<p className="m-0 text-[12px] text-muted-foreground">
						{itemCopy.owner}
					</p>
				</div>
				<div className="grid justify-items-end gap-[7px] max-md:justify-items-start">
					<Eyebrow>{copy.labels.observed}</Eyebrow>
					{item.observed === null ? (
						<strong className="text-[18px] font-medium text-status-warning">
							{copy.unknownObserved}
						</strong>
					) : (
						<strong className="font-mono text-[28px] font-medium tabular-nums">
							{money.format(item.observed)}
						</strong>
					)}
				</div>
			</div>
			<DetailList
				className="mt-[22px] border-t border-border pt-1"
				items={[
					{
						detail: money.format(item.estimate),
						id: 'estimate',
						term: copy.labels.estimate,
					},
					{
						detail: money.format(item.reserved),
						id: 'reserved',
						term: copy.labels.reserved,
					},
					{
						detail: receipt('modelValue', {
							effort: copy.efforts[item.effort],
							tier: copy.tiers[item.tier],
						}),
						id: 'model',
						term: copy.labels.modelEffort,
					},
					{
						detail: receipt('retries', { count: item.retries }),
						id: 'retries',
						term: copy.labels.retries,
					},
					{
						detail: item.verified
							? copy.verification.verified
							: copy.verification.pending,
						id: 'verification',
						term: copy.labels.verification,
						tone: item.verified ? 'success' : 'warning',
					},
					{
						detail: copy.cleanup[item.cleanup],
						id: 'cleanup',
						term: copy.labels.cleanup,
						tone: item.cleanup === 'done' ? 'default' : 'warning',
					},
					{
						detail: (
							<code className="font-mono text-[12px] break-all">
								{item.evidencePath}
							</code>
						),
						id: 'evidence',
						term: copy.labels.evidence,
					},
				]}
			/>
		</article>
	)
}
