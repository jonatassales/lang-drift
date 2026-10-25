'use client'

import { useAtomValue } from 'jotai'

import {
	type MicroruntimesMessages,
	stageEvidencePaths,
} from '@data/microruntimes-content'
import { selectedLifecycleStageAtom } from '@state'

interface StageEvidenceProps {
	copy: MicroruntimesMessages['evidence']
	stageLabels: MicroruntimesMessages['lifecycle']['stages']
}

// The record the lifecycle step selected in the hero leaves behind.
export function StageEvidence(props: StageEvidenceProps) {
	const { copy, stageLabels } = props

	const selectedId = useAtomValue(selectedLifecycleStageAtom)
	const record = copy.stages[selectedId]

	return (
		<article
			aria-atomic="true"
			aria-live="polite"
			className="mt-3 grid gap-2 rounded-2xl border border-border bg-card p-6 sm:p-8"
		>
			<div className="flex flex-wrap items-center gap-3">
				<span className="font-mono text-[9px] tracking-[0.08em] text-muted-foreground uppercase">
					{copy.selectedStage}
				</span>
				<span className="rounded-full border border-border px-2 py-[3px] text-[11px]">
					{stageLabels[selectedId].label}
				</span>
				<span className="text-[11px] text-muted-foreground">{record.kind}</span>
			</div>
			<h3 className="m-0 text-xl font-medium tracking-[-0.02em]">
				{record.title}
			</h3>
			<code className="font-mono text-[12px] break-all text-muted-foreground">
				{stageEvidencePaths[selectedId]}
			</code>
		</article>
	)
}
