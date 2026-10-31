'use client'

import { Eyebrow } from '@langdrift/web/ui/eyebrow'
import { cn } from '@langdrift/web/utilities'
import { useAtom } from 'jotai'

import { ChoiceButton } from '@components/product-page/choice-button'
import {
	lifecycleStageIds,
	type MicroruntimesMessages,
} from '@data/microruntimes-content'
import { selectedLifecycleStageAtom } from '@state'

interface LifecycleStepsProps {
	copy: MicroruntimesMessages['lifecycle']
}

// The six lifecycle steps and what the selected one does and leaves behind.
// Selecting a step drives the evidence record right below.
export function LifecycleSteps(props: LifecycleStepsProps) {
	const { copy } = props

	const [selectedId, setSelectedId] = useAtom(selectedLifecycleStageAtom)
	const stage = copy.stages[selectedId]

	return (
		<div className="@container mt-8 grid gap-5">
			<ol
				aria-label={copy.stepperLabel}
				className="m-0 grid list-none grid-cols-2 gap-2 p-0 @md:grid-cols-3 @4xl:grid-cols-6"
			>
				{lifecycleStageIds.map((id, index) => {
					const selected = id === selectedId

					return (
						<li key={id}>
							<ChoiceButton
								className="flex items-center gap-2.5 border-border px-3 py-2.5"
								kind="step"
								onSelect={() => setSelectedId(id)}
								selected={selected}
							>
								<span
									className={cn(
										'inline-flex size-6 shrink-0 items-center justify-center rounded-full border border-border font-mono text-[10px] text-muted-foreground',
										selected &&
											'border-primary bg-primary text-primary-foreground',
									)}
								>
									{index + 1}
								</span>
								{copy.stages[id].label}
							</ChoiceButton>
						</li>
					)
				})}
			</ol>

			<div
				aria-atomic="true"
				aria-live="polite"
				className="grid gap-4 rounded-xl bg-foreground/[0.03] p-4 @lg:grid-cols-2 sm:p-5"
			>
				<div className="grid content-start gap-1.5">
					<Eyebrow>{copy.whatItDoes}</Eyebrow>
					<p className="m-0 text-[13px] leading-[1.6]">{stage.does}</p>
				</div>
				<div className="grid content-start gap-1.5">
					<Eyebrow>{copy.evidenceItLeaves}</Eyebrow>
					<p className="m-0 text-[13px] leading-[1.6]">{stage.leaves}</p>
				</div>
			</div>
		</div>
	)
}
