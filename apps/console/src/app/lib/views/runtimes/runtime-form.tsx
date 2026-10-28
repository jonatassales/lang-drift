'use client'

import { useTranslations } from 'next-intl'
import { useState, useTransition } from 'react'

import { ErrorNotice } from '@components/error-notice/error-notice'
import { fieldClass, labelClass } from '@components/form-styles/form-styles'
import { mutation } from '@data/api'
import type { Component, Runtime, RuntimeConfiguration } from '@data/types'
import { Button } from '@ui/button'
import { Card } from '@ui/card'

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

	const t = useTranslations('ui')
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
	function apply(values: Partial<RuntimeConfiguration>) {
		setConfiguration((current) => ({
			llmBudgetUsd: values.llmBudgetUsd?.toString() ?? current.llmBudgetUsd,
			cpu: values.cpu?.toString() ?? current.cpu,
			memoryMiB: values.memoryMiB?.toString() ?? current.memoryMiB,
			model: values.model ?? current.model,
			effort: values.effort ?? current.effort,
		}))
	}

	return (
		<div className="grid items-start gap-5 xl:grid-cols-[minmax(0,1fr)_320px]">
			<Card className="p-5">
				<form
					className="grid gap-5"
					onSubmit={(event) => {
						event.preventDefault()
						const form = new FormData(event.currentTarget)
						let selected: Component[]
						try {
							selected = componentsFromForm(form, components)
							setInvalidBuild(false)
						} catch {
							setInvalidBuild(true)

							return
						}
						setError(null)
						startTransition(async () => {
							try {
								const nativeRuntimeId = String(
									form.get('nativeRuntimeId') ?? '',
								).trim()
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

								const result = await mutation<{
									runtime: Runtime
									secret?: string
								}>(
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
					}}
				>
					<fieldset disabled={locked || pending} className="grid gap-5">
						<label className={labelClass}>
							{t('name')}
							<input
								className={fieldClass}
								name="name"
								required
								maxLength={100}
								defaultValue={runtime?.name ?? ''}
							/>
						</label>
						<label className={labelClass}>
							{t('llmBudget')}
							<input
								className={fieldClass}
								type="number"
								min="0.01"
								max="1000000"
								step="0.01"
								value={configuration.llmBudgetUsd}
								onChange={(event) =>
									setConfiguration({
										...configuration,
										llmBudgetUsd: event.target.value,
									})
								}
							/>
							<span className="text-xs text-muted-foreground">
								{t('llmBudgetHint')}
							</span>
						</label>
						<div className="grid gap-4 sm:grid-cols-2">
							<label className={labelClass}>
								{t('cpu')}
								<input
									className={fieldClass}
									type="number"
									min="0.001"
									max="64"
									step="0.001"
									required
									value={configuration.cpu}
									onChange={(event) =>
										setConfiguration({
											...configuration,
											cpu: event.target.value,
										})
									}
								/>
							</label>
							<label className={labelClass}>
								{t('memory')}
								<input
									className={fieldClass}
									type="number"
									min="1"
									max="1048576"
									step="1"
									required
									value={configuration.memoryMiB}
									onChange={(event) =>
										setConfiguration({
											...configuration,
											memoryMiB: event.target.value,
										})
									}
								/>
							</label>
							<label className={labelClass}>
								{t('model')}
								<input
									className={fieldClass}
									required
									maxLength={160}
									value={configuration.model}
									onChange={(event) =>
										setConfiguration({
											...configuration,
											model: event.target.value,
										})
									}
								/>
							</label>
							<label className={labelClass}>
								{t('effort')}
								<select
									className={fieldClass}
									value={configuration.effort}
									onChange={(event) => {
										const value = event.target.value
										if (
											value === 'low' ||
											value === 'medium' ||
											value === 'high'
										)
											setConfiguration({
												...configuration,
												effort: value,
											})
									}}
								>
									{['low', 'medium', 'high'].map((effort) => (
										<option key={effort} value={effort}>
											{t(effort)}
										</option>
									))}
								</select>
							</label>
						</div>
						<p className="text-xs leading-relaxed text-muted-foreground">
							{t('agentSettingsHint')}
						</p>
						{runtime ? (
							<input
								type="hidden"
								name="nativeRuntimeId"
								value={runtime.configuration.nativeRuntimeId ?? ''}
							/>
						) : (
							<label className={labelClass}>
								{t('migration')}
								<input
									className={fieldClass}
									name="nativeRuntimeId"
									pattern="[a-z][a-z0-9-]*"
									maxLength={63}
								/>
								<span className="font-normal text-muted-foreground">
									{t('migrationHint')}
								</span>
							</label>
						)}
						<div className="grid gap-4">
							<div>
								<h2 className="text-sm font-medium">{t('components')}</h2>
								<p className="mt-2 text-xs leading-relaxed text-muted-foreground">
									{t('buildHint')}
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
										onSelectionChange={(selected) => {
											setSelectedPaths((current) =>
												selected
													? [...current, component.path]
													: current.filter((path) => path !== component.path),
											)
										}}
									/>
								)
							})}
						</div>
						{invalidBuild ? (
							<p role="alert" className="text-sm text-status-danger">
								{t('invalidBuild')}
							</p>
						) : null}
						{error ? <ErrorNotice error={error} /> : null}
						<Button type="submit" tone="primary">
							{t(pending ? 'saving' : runtime ? 'save' : 'createRuntime')}
						</Button>
					</fieldset>
					{locked ? (
						<p role="status" className="mt-4 text-sm text-muted-foreground">
							{t('configurationLocked')}
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
				disabled={locked || pending || selectedPaths.length === 0}
				onApply={apply}
			/>
		</div>
	)
}
