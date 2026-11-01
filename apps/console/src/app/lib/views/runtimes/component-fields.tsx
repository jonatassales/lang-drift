'use client'
import { useTranslations } from 'next-intl'

import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import type { Component } from '@data/types'

interface ComponentFieldsProps {
	component: Component
	selected: boolean
	onSelectionChange: (selected: boolean) => void
}
export function ComponentFields(props: ComponentFieldsProps) {
	const { component, selected: enabled, onSelectionChange } = props

	const translate = useTranslations('ui')
	const prefix = `component:${component.path}:`

	return (
		<fieldset className="rounded-lg border border-border p-4">
			<legend className="px-1 text-sm">
				<label className="flex items-center gap-2">
					<input
						className="focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring"
						type="checkbox"
						name="selectedComponents"
						value={component.path}
						checked={enabled}
						onChange={(event) => onSelectionChange(event.target.checked)}
					/>
					{component.name}
					<code className="text-[11px] text-muted-foreground">
						{component.path}
					</code>
				</label>
			</legend>
			<div className="grid gap-4 sm:grid-cols-2">
				<label className={labelClass}>
					{translate('image')}
					<input
						className={fieldClass}
						name={`${prefix}image`}
						required={enabled}
						disabled={!enabled}
						defaultValue={component.image ?? ''}
						maxLength={200}
					/>
				</label>
				<label className={labelClass}>
					{translate('port')}
					<input
						className={fieldClass}
						name={`${prefix}port`}
						type="number"
						min={1}
						max={65535}
						required={enabled}
						disabled={!enabled}
						defaultValue={component.port ?? ''}
					/>
				</label>
				<label className={labelClass}>
					{translate('readinessPath')}
					<input
						className={fieldClass}
						name={`${prefix}readinessPath`}
						required={enabled}
						disabled={!enabled}
						defaultValue={component.readinessPath ?? ''}
						maxLength={200}
					/>
				</label>
				<label className={labelClass}>
					{translate('buildExecutable')}
					<input
						className={fieldClass}
						name={`${prefix}executable`}
						required={enabled}
						disabled={!enabled}
						defaultValue={component.build?.executable ?? ''}
						maxLength={200}
					/>
				</label>
				<label className={`${labelClass} sm:col-span-2`}>
					{translate('buildArguments')}
					<textarea
						className={`${fieldClass} font-mono`}
						name={`${prefix}args`}
						required={enabled}
						disabled={!enabled}
						defaultValue={
							component.build ? JSON.stringify(component.build.args) : ''
						}
						rows={2}
						maxLength={4096}
					/>
				</label>
			</div>
		</fieldset>
	)
}
