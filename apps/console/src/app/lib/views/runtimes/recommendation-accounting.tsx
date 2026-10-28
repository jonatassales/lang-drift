'use client'
import { useTranslations } from 'next-intl'

import type { Recommendation } from '@data/types'

interface RecommendationAccountingProps {
	recommendation: Recommendation
}
export function RecommendationAccounting(props: RecommendationAccountingProps) {
	const { recommendation } = props

	const t = useTranslations('ui')

	return (
		<dl className="grid grid-cols-2 gap-2 border-t border-border pt-3 text-[11px]">
			<dt className="text-muted-foreground">{t('generationModel')}</dt>
			<dd className="break-all">
				{recommendation.generationModel ?? t('unknown')}
			</dd>
			<dt className="text-muted-foreground">{t('generationEffort')}</dt>
			<dd>
				{recommendation.generationEffort
					? t(recommendation.generationEffort)
					: t('unknown')}
			</dd>
			<dt className="text-muted-foreground">{t('generationSettingsSource')}</dt>
			<dd>
				{recommendation.generationSettingsSource
					? t(
							recommendation.generationSettingsSource === 'draft'
								? 'generationSourceDraft'
								: recommendation.generationSettingsSource === 'runtime'
									? 'generationSourceRuntime'
									: 'generationSourceOperator',
						)
					: t('unknown')}
			</dd>
			<dt className="text-muted-foreground">{t('providerUsage')}</dt>
			<dd>
				{recommendation.usage ? (
					<ul className="grid gap-1">
						{Object.entries(recommendation.usage).map(([key, value]) => (
							<li key={key}>
								<code>{key}</code>: {value}
							</li>
						))}
					</ul>
				) : (
					t('unknown')
				)}
			</dd>
			<dt className="text-muted-foreground">{t('cost')}</dt>
			<dd>{t('unknown')}</dd>
			{recommendation.limits ? (
				<>
					<dt className="text-muted-foreground">{t('outputLimit')}</dt>
					<dd>{recommendation.limits.outputTokens}</dd>
					<dt className="text-muted-foreground">{t('requestTimeout')}</dt>
					<dd>{recommendation.limits.timeoutSeconds}</dd>
					<dt className="text-muted-foreground">{t('retries')}</dt>
					<dd>{recommendation.limits.retries}</dd>
				</>
			) : null}
		</dl>
	)
}
