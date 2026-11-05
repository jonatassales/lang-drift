import { useTranslations } from 'next-intl'

import type { TicketStage } from '@data/types'
import { Eyebrow } from 'web/ui/eyebrow'
import { PencilSimple, UserCircle } from 'web/ui/icons/client'
import { cn } from 'web/utilities'

import { stageMarker } from './stage-marker.compute'

interface WorkflowStageListProps {
	stages: TicketStage[]
}

const progressTones: Record<TicketStage['progress'], string> = {
	pending: 'text-muted-foreground',
	active: 'text-primary',
	waiting: 'text-status-warning',
	done: 'text-foreground',
}

// The delivery stages the workflow runs, read-only: each stage's name, a
// read mark once done, whether a person writes it by hand or approves it,
// and its progress. Merge is always the owner's.
export function WorkflowStageList(props: WorkflowStageListProps) {
	const { stages } = props

	const translate = useTranslations('tickets')

	return (
		<div className="grid gap-1.5">
			<Eyebrow>{translate('workflow.humanInTheLoop')}</Eyebrow>
			<ul className="m-0 grid list-none gap-x-4 gap-y-1 p-0 sm:grid-cols-2">
				{stages.map((stage) => {
					const marker = stageMarker(stage)

					return (
						<li
							className={cn(
								'flex min-h-7 items-center gap-2 text-[12px]',
								progressTones[stage.progress],
							)}
							key={stage.id}
						>
							{stage.progress === 'done' ? (
								<span
									aria-hidden="true"
									className="font-mono text-status-success"
								>
									✓
								</span>
							) : (
								<span
									aria-hidden="true"
									className="inline-block size-1.5 rounded-full bg-current"
								/>
							)}
							<span>{translate(`workflow.stages.${stage.id}`)}</span>
							{marker === 'person' && (
								<PencilSimple
									aria-label={translate('safeSeat.personDriven')}
									className="shrink-0 text-primary"
									role="img"
									size={14}
								/>
							)}
							{marker === 'approval' && (
								<UserCircle
									aria-label={translate('workflow.humanInTheLoop')}
									className="shrink-0 text-muted-foreground"
									role="img"
									size={14}
								/>
							)}
							<span className="font-mono text-[10px] text-muted-foreground">
								{stage.id === 'merge'
									? translate('workflow.mergeLocked')
									: translate(`workflow.progress.${stage.progress}`)}
							</span>
							{stage.skills && (
								<span
									aria-hidden="true"
									className="font-mono text-[10px] text-muted-foreground max-sm:hidden"
								>
									{stage.skills}
								</span>
							)}
						</li>
					)
				})}
			</ul>
		</div>
	)
}
