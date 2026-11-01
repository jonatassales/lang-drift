'use client'
import { useTranslations } from 'next-intl'

import type { Recommendation } from '@data/types'

interface RecommendationAccountingProps {
	recommendation: Recommendation
}
export function RecommendationAccounting(props: RecommendationAccountingProps) {
	const { recommendation } = props

	const translate = useTranslations('ui')

	return (
		<dl className="grid grid-cols-2 gap-2 border-t border-border pt-3 text-[11px]">
			<dt className="text-muted-foreground">{translate('generationModel')}</dt>
			<dd className="break-all">
				{recommendation.generationModel ?? translate('unknown')}
			</dd>
			<dt className="text-muted-foreground">{translate('generationEffort')}</dt>
			<dd>
				{recommendation.generationEffort
					? translate(recommendation.generationEffort)
					: translate('unknown')}
			</dd>
			<dt className="text-muted-foreground">
				{translate('generationSettingsSource')}
			</dt>
			<dd>
				{recommendation.generationSettingsSource
					? translate(
							recommendation.generationSettingsSource === 'draft'
								? 'generationSourceDraft'
								: recommendation.generationSettingsSource === 'runtime'
									? 'generationSourceRuntime'
									: 'generationSourceOperator',
						)
					: translate('unknown')}
			</dd>
			<dt className="text-muted-foreground">{translate('providerUsage')}</dt>
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
					translate('unknown')
				)}
			</dd>
			<dt className="text-muted-foreground">{translate('cost')}</dt>
			<dd>{translate('unknown')}</dd>
			{recommendation.limits ? (
				<>
					<dt className="text-muted-foreground">{translate('outputLimit')}</dt>
					<dd>{recommendation.limits.outputTokens}</dd>
					<dt className="text-muted-foreground">
						{translate('requestTimeout')}
					</dt>
					<dd>{recommendation.limits.timeoutSeconds}</dd>
					<dt className="text-muted-foreground">{translate('retries')}</dt>
					<dd>{recommendation.limits.retries}</dd>
				</>
			) : null}
		</dl>
	)
}
