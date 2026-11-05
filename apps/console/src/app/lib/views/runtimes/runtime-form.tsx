'use client'

import { useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { mutation } from '@data/api'
import type { SavedRuntimeEnvelope } from '@data/read-envelopes.types'
import type { Component, Runtime, RuntimeConfiguration } from '@data/types'
import { Button } from 'web/ui/button'
import { Card } from 'web/ui/card'

import { ComponentFields } from './component-fields'
import { componentsFromForm } from './components-from-form'
import { RecommendationPanel } from './recommendation-panel'

interface RuntimeFormProps {
	productId: string
	components: Component[]
	runtime?: Runtime
	locked?: boolean
	onSaved: (runtime: Runtime, secret?: string) => void
}

export function RuntimeForm(props: RuntimeFormProps) {
	const { productId, components, runtime, locked = false, onSaved } = props

	const [error, setError] = useState<unknown>(null)
	const [invalidBuild, setInvalidBuild] = useState(false)
	const [pending, startTransition] = useTransition()
	const [selectedPaths, setSelectedPaths] = useState(
		(runtime?.components ?? components).map((component) => component.path),
	)
	const [configuration, setConfiguration] = useState({
		llmBudgetUsd: runtime?.configuration.llmBudgetUsd?.toString() ?? '',
		cpu: runtime?.configuration.cpu.toString() ?? '',
		memoryMiB: runtime?.configuration.memoryMiB.toString() ?? '',
		model: runtime?.configuration.model ?? '',
		effort: runtime?.configuration.effort ?? 'medium',
	})

	const translate = useTranslations('ui')

	function handleApply(values: Partial<RuntimeConfiguration>) {
		setConfiguration((current) => ({
			llmBudgetUsd: values.llmBudgetUsd?.toString() ?? current.llmBudgetUsd,
			cpu: values.cpu?.toString() ?? current.cpu,
			memoryMiB: values.memoryMiB?.toString() ?? current.memoryMiB,
			model: values.model ?? current.model,
			effort: values.effort ?? current.effort,
		}))
	}

	function handleToggleComponent(path: string, selected: boolean) {
		setSelectedPaths((current) =>
			selected
				? [...current, path]
				: current.filter((currentPath) => currentPath !== path),
		)
	}

	function handleLlmBudgetChange(event: React.ChangeEvent<HTMLInputElement>) {
		setConfiguration({
			...configuration,
			llmBudgetUsd: event.target.value,
		})
	}

	function handleCpuChange(event: React.ChangeEvent<HTMLInputElement>) {
		setConfiguration({
			...configuration,
			cpu: event.target.value,
		})
	}

	function handleMemoryChange(event: React.ChangeEvent<HTMLInputElement>) {
		setConfiguration({
			...configuration,
			memoryMiB: event.target.value,
		})
	}

	function handleModelChange(event: React.ChangeEvent<HTMLInputElement>) {
		setConfiguration({
			...configuration,
			model: event.target.value,
		})
	}

	function handleEffortChange(event: React.ChangeEvent<HTMLSelectElement>) {
		const value = event.target.value
		if (value === 'low' || value === 'medium' || value === 'high')
			setConfiguration({
				...configuration,
				effort: value,
			})
	}

	function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
		event.preventDefault()
		const form = new FormData(event.currentTarget)
		const selected = componentsFromForm(form, components)
		if (!selected) {
			setInvalidBuild(true)

			return
		}
		setInvalidBuild(false)
		setError(null)
		startTransition(async () => {
			try {
				const nativeRuntimeId = String(form.get('nativeRuntimeId') ?? '').trim()
				const config: RuntimeConfiguration = {
					...(configuration.llmBudgetUsd
						? {
								llmBudgetUsd: Number(configuration.llmBudgetUsd),
							}
						: {}),
					cpu: Number(configuration.cpu),
					memoryMiB: Number(configuration.memoryMiB),
					model: configuration.model,
					effort: configuration.effort,
					...(nativeRuntimeId
						? {
								nativeRuntimeId,
							}
						: {}),
				}

				const result = await mutation<SavedRuntimeEnvelope>(
					runtime ? `runtimes/${runtime.id}` : 'runtimes',
					runtime ? 'PATCH' : 'POST',
					{
						...(runtime
							? {
									revision: runtime.revision,
								}
							: {
									productId,
								}),
						name: form.get('name'),
						components: selected,
						configuration: config,
					},
				)

				onSaved(result.runtime, result.secret)
			} catch (failure) {
				setError(failure)
			}
		})
	}

	let submitLabel = translate('createRuntime')
	if (pending) {
		submitLabel = translate('saving')
	} else if (runtime) {
		submitLabel = translate('save')
	}
	const cannotRecommend = locked || pending || selectedPaths.length === 0

	return (
		<div className="grid items-start gap-5 xl:grid-cols-[minmax(0,1fr)_320px]">
			<Card className="block p-5">
				<form className="grid gap-5" onSubmit={handleSubmit}>
					<fieldset disabled={locked || pending} className="grid gap-5">
						<label className={labelClass}>
							{translate('name')}
							<input
								className={fieldClass}
								name="name"
								required
								maxLength={100}
								defaultValue={runtime?.name ?? ''}
							/>
						</label>
						<label className={labelClass}>
							{translate('llmBudget')}
							<input
								className={fieldClass}
								type="number"
								min="0.01"
								max="1000000"
								step="0.01"
								value={configuration.llmBudgetUsd}
								onChange={handleLlmBudgetChange}
							/>
							<span className="text-xs text-muted-foreground">
								{translate('llmBudgetHint')}
							</span>
						</label>
						<div className="grid gap-4 sm:grid-cols-2">
							<label className={labelClass}>
								{translate('cpu')}
								<input
									className={fieldClass}
									type="number"
									min="0.001"
									max="64"
									step="0.001"
									required
									value={configuration.cpu}
									onChange={handleCpuChange}
								/>
							</label>
							<label className={labelClass}>
								{translate('memory')}
								<input
									className={fieldClass}
									type="number"
									min="1"
									max="1048576"
									step="1"
									required
									value={configuration.memoryMiB}
									onChange={handleMemoryChange}
								/>
							</label>
							<label className={labelClass}>
								{translate('model')}
								<input
									className={fieldClass}
									required
									maxLength={160}
									value={configuration.model}
									onChange={handleModelChange}
								/>
							</label>
							<label className={labelClass}>
								{translate('effort')}
								<select
									className={fieldClass}
									value={configuration.effort}
									onChange={handleEffortChange}
								>
									{['low', 'medium', 'high'].map((effort) => (
										<option key={effort} value={effort}>
											{translate(effort)}
										</option>
									))}
								</select>
							</label>
						</div>
						<p className="text-xs leading-relaxed text-muted-foreground">
							{translate('agentSettingsHint')}
						</p>
						{runtime ? (
							<input
								type="hidden"
								name="nativeRuntimeId"
								value={runtime.configuration.nativeRuntimeId ?? ''}
							/>
						) : (
							<label className={labelClass}>
								{translate('migration')}
								<input
									className={fieldClass}
									name="nativeRuntimeId"
									pattern="[a-z][a-z0-9-]*"
									maxLength={63}
								/>
								<span className="font-normal text-muted-foreground">
									{translate('migrationHint')}
								</span>
							</label>
						)}
						<div className="grid gap-4">
							<div>
								<h2 className="text-sm font-medium">
									{translate('components')}
								</h2>
								<p className="mt-2 text-xs leading-relaxed text-muted-foreground">
									{translate('buildHint')}
								</p>
							</div>
							{components.map((component) => {
								const existing = runtime?.components.find(
									(selected) => selected.path === component.path,
								)

								return (
									<ComponentFields
										key={component.path}
										component={existing ?? component}
										selected={selectedPaths.includes(component.path)}
										onToggle={handleToggleComponent}
									/>
								)
							})}
						</div>
						{invalidBuild ? (
							<p role="alert" className="text-sm text-status-danger">
								{translate('invalidBuild')}
							</p>
						) : null}
						{error ? <ErrorNotice error={error} /> : null}
						<Button type="submit">{submitLabel}</Button>
					</fieldset>
					{locked ? (
						<p role="status" className="mt-4 text-sm text-muted-foreground">
							{translate('configurationLocked')}
						</p>
					) : null}
				</form>
			</Card>
			<RecommendationPanel
				key={JSON.stringify(selectedPaths)}
				productId={productId}
				components={components.filter((component) =>
					selectedPaths.includes(component.path),
				)}
				agent={{
					...(configuration.model.trim()
						? {
								model: configuration.model.trim(),
							}
						: {}),
					effort: configuration.effort,
				}}
				endpoint={
					runtime
						? `runtimes/${runtime.id}/recommendations`
						: `products/${productId}/recommendations`
				}
				disabled={cannotRecommend}
				onApply={handleApply}
			/>
		</div>
	)
}
