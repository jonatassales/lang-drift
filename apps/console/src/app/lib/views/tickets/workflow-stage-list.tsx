import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { cn } from '@langdrift/web/utilities'
import { useTranslations } from 'next-intl'

import type { TicketStage } from '@data/types'

import { ReadMark } from './read-mark'

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
// read mark once done, whether a person is in the loop, and its progress.
// Merge is always the owner's.
export function WorkflowStageList(props: WorkflowStageListProps) {
	const { stages } = props

	const t = useTranslations('tickets')

	return (
		<div className="grid gap-1.5">
			<Eyebrow>{t('workflow.humanInTheLoop')}</Eyebrow>
			<ul className="m-0 grid list-none gap-x-4 gap-y-1 p-0 sm:grid-cols-2">
				{stages.map((stage) => (
					<li
						className={cn(
							'flex min-h-7 items-center gap-2 text-[12px]',
							progressTones[stage.progress],
						)}
						key={stage.id}
					>
						{stage.progress === 'done' ? (
							<ReadMark />
						) : (
							<span
								aria-hidden="true"
								className="inline-block size-1.5 rounded-full bg-current"
							/>
						)}
						<span>{t(`workflow.stages.${stage.id}`)}</span>
						<span className="font-mono text-[10px] text-muted-foreground">
							{stage.id === 'merge'
								? t('workflow.mergeLocked')
								: t(`workflow.progress.${stage.progress}`)}
						</span>
						{stage.human && stage.id !== 'merge' && (
							<span className="sr-only">{t('workflow.humanInTheLoop')}</span>
						)}
					</li>
				))}
			</ul>
		</div>
	)
}
