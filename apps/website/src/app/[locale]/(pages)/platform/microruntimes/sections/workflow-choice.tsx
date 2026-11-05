import type { ChangeEvent } from 'react'

import { Eyebrow } from 'web/ui/eyebrow'
import { CaretDownIcon } from 'web/ui/icons/arrow'

interface WorkflowChoiceOption<Value extends string> {
	label: string
	value: Value
}

interface WorkflowChoiceProps<Value extends string> {
	label: string
	/** True while a run reads the choice; the select keeps focus but no change. */
	locked: boolean
	/** The field's name, so the browser can identify it. */
	name: string
	note?: string
	onChange: (value: Value) => void
	options: readonly WorkflowChoiceOption<Value>[]
	value: Value
}

// One labelled choice of the Workflow card, a native select drawn in the
// cards' roles. Its text is 16px on phones so focusing it does not zoom the
// page.
export function WorkflowChoice<Value extends string>(
	props: WorkflowChoiceProps<Value>,
) {
	const { label, locked, name, note, onChange, options, value } = props

	function handleChange(event: ChangeEvent<HTMLSelectElement>) {
		const option = options.find(
			(candidate) => candidate.value === event.target.value,
		)

		if (option) {
			onChange(option.value)
		}
	}

	return (
		<label className="grid min-w-0 content-start gap-1.5">
			<Eyebrow>{label}</Eyebrow>
			<span className="relative block">
				<select
					aria-disabled={locked || undefined}
					className="w-full cursor-pointer appearance-none rounded-lg border border-border bg-card py-1.5 pr-8 pl-2.5 font-mono text-base text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary aria-disabled:cursor-not-allowed aria-disabled:opacity-60 sm:text-[12px]"
					name={name}
					onChange={handleChange}
					value={value}
				>
					{options.map((option) => (
						<option key={option.value} value={option.value}>
							{option.label}
						</option>
					))}
				</select>
				<CaretDownIcon
					className="pointer-events-none absolute top-1/2 right-2 -translate-y-1/2 text-muted-foreground"
					size={14}
				/>
			</span>
			{note && (
				<span className="text-[11px] leading-snug text-muted-foreground">
					{note}
				</span>
			)}
		</label>
	)
}
