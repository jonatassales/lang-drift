'use client'

import { useAtomValue } from 'jotai'
import { createTranslator } from 'next-intl'

import type { WebsiteLocale } from '@i18n/routing'
import {
	dollars,
	type FinopsMessages,
	workItems,
} from '@routes/platform/finops/content/finops-content'
import { selectedWorkItemAtom } from '@state'
import { DetailList } from 'web/ui/detail-list'
import { Eyebrow } from 'web/ui/eyebrow'

// The receipt of the work item selected in the ledger: what it was estimated,
// what was reserved, what it consumed, and how it was verified and cleaned up.
// Its head stacks at every width, so it reads the same docked beside the
// ledger and under it.
interface DeliveryReceiptProps {
	copy: FinopsMessages['receipt']
	/** 2 under the FinOps page's title, 3 under a home section's heading. */
	headingLevel?: 2 | 3
	items: FinopsMessages['ledger']['items']
	locale: WebsiteLocale
}

export function DeliveryReceipt(props: DeliveryReceiptProps) {
	const { copy, headingLevel = 2, items, locale } = props

	const selectedId = useAtomValue(selectedWorkItemAtom)

	const Heading = headingLevel === 3 ? 'h3' : 'h2'
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
			className="rounded-2xl border border-border bg-card p-6 max-md:p-[18px]"
		>
			<div className="grid gap-5">
				<div className="grid gap-[7px]">
					<Eyebrow>{copy.selected}</Eyebrow>
					<Heading className="m-0 text-[28px] font-medium tracking-[-0.025em]">
						{itemCopy.title}
					</Heading>
					<p className="m-0 text-[12px] text-muted-foreground">
						{itemCopy.owner}
					</p>
				</div>
				<div className="grid gap-[7px]">
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
