'use client'

import { useAtomValue } from 'jotai'

import {
	type MicroruntimesMessages,
	stageEvidencePaths,
} from '@routes/platform/microruntimes/content/microruntimes-content'
import { selectedLifecycleStageAtom } from '@state'
import { Eyebrow } from 'web/ui/eyebrow'

// The record the lifecycle step selected in the hero leaves behind.
interface StageEvidenceProps {
	copy: MicroruntimesMessages['evidence']
	stageLabels: MicroruntimesMessages['lifecycle']['stages']
}

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
				<Eyebrow>{copy.selectedStage}</Eyebrow>
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
