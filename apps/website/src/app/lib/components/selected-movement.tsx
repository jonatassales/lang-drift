'use client'

import { driftSign } from '@pack/react/utilities'
import { AnimatedAvatarGroup } from '@pack/react/vendors/smoothui/animated-avatar-group'
import { useAtomValue } from 'jotai'
import { useFormatter, useTranslations } from 'next-intl'

import { selectedPointAtom } from '@state'

import { useHeroPoints } from './hero-chart/use-hero-points.hook'

export function SelectedMovement() {
	const t = useTranslations('home.why')
	const demo = useTranslations('demo')
	const format = useFormatter()
	const points = useHeroPoints()
	const selectedPoint = useAtomValue(selectedPointAtom)
	const activeEvent = points[selectedPoint]?.event ?? points.at(-1)?.event
	if (!activeEvent) return null

	return (
		<article
			className="mt-11 block rounded-2xl border border-border bg-card p-6 max-md:p-[18px]"
			aria-live="polite"
			aria-atomic="true"
		>
			<div className="flex items-start justify-between gap-6 max-md:flex-col">
				<div className="grid gap-[7px]">
					<span className="font-mono text-[8px] tracking-[0.08em] text-muted-foreground uppercase">
						{t('selected')}
					</span>
					<h3 className="m-0 text-[28px] font-medium tracking-[-0.025em]">
						{activeEvent.title}
					</h3>
					<p className="m-0 text-[11px] text-muted-foreground">
						{activeEvent.date} · {activeEvent.productArea ?? t('product')}
					</p>
				</div>
				<div className="grid items-end justify-items-end gap-[7px] max-md:justify-items-start">
					<strong
						className="font-mono text-[28px] font-medium"
						data-drift-sign={driftSign(activeEvent.delta)}
					>
						{format.number(activeEvent.delta, { signDisplay: 'exceptZero' })}
					</strong>
					<span
						className="rounded-full border border-border px-2 py-[5px] text-[9px] text-muted-foreground data-[classification=intentional]:text-classification-text-intentional data-[classification=review]:text-classification-text-review data-[classification=unexplained]:text-classification-text-unexplained"
						data-classification={activeEvent.classification}
					>
						{demo(`classifications.${activeEvent.classification}`)}
					</span>
				</div>
			</div>
			<div className="mt-[22px] grid grid-cols-[minmax(0,1.6fr)_minmax(180px,0.7fr)_minmax(180px,0.7fr)] border-t border-border max-md:grid-cols-1">
				<div className="grid min-h-32 content-start gap-[9px] py-5 pr-[22px] max-md:min-h-0 max-md:px-0 max-md:py-4 max-md:[&+div]:border-t max-md:[&+div]:border-l-0 max-md:[&+div]:pl-0 [&+div]:border-l [&+div]:border-border [&+div]:pl-[22px]">
					<span className="font-mono text-[8px] tracking-[0.08em] text-muted-foreground uppercase">
						{t('reason')}
					</span>
					<p className="m-0 text-[13px] leading-[1.55]">{activeEvent.reason}</p>
				</div>
				<div className="grid min-h-32 content-start gap-[9px] py-5 pr-[22px] max-md:min-h-0 max-md:px-0 max-md:py-4 max-md:[&+div]:border-t max-md:[&+div]:border-l-0 max-md:[&+div]:pl-0 [&+div]:border-l [&+div]:border-border [&+div]:pl-[22px]">
					<span className="font-mono text-[8px] tracking-[0.08em] text-muted-foreground uppercase">
						{t('decision')}
					</span>
					<strong className="text-[13px] font-medium">
						{activeEvent.decision}
					</strong>
				</div>
				<div className="grid min-h-32 content-start gap-[9px] py-5 pr-[22px] max-md:min-h-0 max-md:px-0 max-md:py-4">
					<span className="font-mono text-[8px] tracking-[0.08em] text-muted-foreground uppercase">
						{t('people')}
					</span>
					<AnimatedAvatarGroup
						people={activeEvent.actors.map((actor) => ({
							initials: actor.initials,
							name: actor.name,
							role: actor.team,
							src: actor.src,
						}))}
						size={32}
					/>
					<small className="text-[9px] text-muted-foreground">
						{activeEvent.actors.map((actor) => actor.name).join(' + ')}
					</small>
				</div>
			</div>
		</article>
	)
}
