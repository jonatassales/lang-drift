'use client'

import { AnimatedTabs } from '@pack/react/vendors/smoothui/animated-tabs'
import { useId, useState } from 'react'

import {
	type ScenarioId,
	type VisionLoopMessages,
	visionLoopScenarios,
} from '@data/vision-loop-content'

import { EvidenceRecord } from './vision-loop-evidence/evidence-record'
import { PillarIcon } from './vision-loop-evidence/pillar-icon'

interface DemoProps {
	labels: VisionLoopMessages['labels']
	scenarios: VisionLoopMessages['scenarios']
	selectorLabel: VisionLoopMessages['selectorLabel']
}

export function VisionLoopEvidenceDemo(props: DemoProps) {
	const { labels, scenarios, selectorLabel } = props

	const instanceId = useId()
	const [scenarioId, setScenarioId] = useState<ScenarioId>('priority-shift')
	const tabs = visionLoopScenarios.map(({ id }) => ({
		id,
		label: scenarios[id].tabLabel,
	}))

	return (
		<div className="vision-loop-demo">
			<span className="vle-selector-label">{selectorLabel}</span>
			<AnimatedTabs
				tabs={tabs}
				activeTab={scenarioId}
				onChange={(id) => {
					const nextScenario = visionLoopScenarios.find(
						(scenario) => scenario.id === id,
					)
					if (nextScenario) setScenarioId(nextScenario.id)
				}}
				variant="underline"
				layoutId={instanceId}
				ariaLabel={selectorLabel}
				className="vle-scenario-tabs"
			/>
			{visionLoopScenarios.map((scenario) => {
				const copy = scenarios[scenario.id]
				const active = scenario.id === scenarioId

				return (
					<div
						key={scenario.id}
						role="tabpanel"
						id={`${instanceId}-panel-${scenario.id}`}
						aria-labelledby={`${instanceId}-tab-${scenario.id}`}
						tabIndex={active ? 0 : -1}
						hidden={!active}
						className="vle-panel"
					>
						{active ? (
							<div className="vle-scenario-content" data-scenario={scenario.id}>
								<p className="vle-scenario-context">{copy.context}</p>
								<div className="vle-surface">
									<div className="vle-pillars">
										<div className="vle-pillar vle-vision">
											<div className="vle-pillar-heading min-w-0 [&_h3]:min-w-0 [&_h3]:[overflow-wrap:anywhere]">
												<PillarIcon type="vision" />
												<h3>{labels.vision}</h3>
												<span className="vle-connector" aria-hidden="true">
													→
												</span>
											</div>
											<dl className="vle-vision-target">
												<div>
													<dt>{labels.target}</dt>
													<dd>{copy.vision.target}</dd>
												</div>
												<div>
													<dt>{labels.audience}</dt>
													<dd>{copy.vision.audience}</dd>
												</div>
												<div>
													<dt>{labels.successCriterion}</dt>
													<dd>{copy.vision.successCriterion}</dd>
												</div>
											</dl>
										</div>
										<div className="vle-pillar vle-loop">
											<div className="vle-pillar-heading min-w-0 [&_h3]:min-w-0 [&_h3]:[overflow-wrap:anywhere]">
												<PillarIcon type="loop" />
												<h3>{labels.loop}</h3>
												<span className="vle-connector" aria-hidden="true">
													→
												</span>
											</div>
											<ol className="vle-loop-records">
												{scenario.loop.map((id) => (
													<li key={id}>
														<h4>{copy.loop[id].label}</h4>
														<p>{copy.loop[id].description}</p>
													</li>
												))}
											</ol>
										</div>
										<div className="vle-pillar vle-evidence">
											<div className="vle-pillar-heading min-w-0 [&_h3]:min-w-0 [&_h3]:[overflow-wrap:anywhere]">
												<PillarIcon type="evidence" />
												<h3>{labels.evidence}</h3>
											</div>
											<p className="vle-records-label">
												{labels.supportingRecords}
											</p>
											<ul className="vle-evidence-records">
												{scenario.evidence.map((id) => (
													<EvidenceRecord
														key={id}
														id={`${instanceId}-${scenario.id}-${id}`}
														copy={copy.evidence[id]}
														sourceLabel={labels.sourceContext}
													/>
												))}
											</ul>
										</div>
									</div>
									<div className="vle-interpretation">
										<div>
											<span className="vle-label">{labels.interpretation}</span>
											<h3>{copy.interpretation.title}</h3>
										</div>
										<div className="vle-interpretation-copy">
											<p>{copy.interpretation.description}</p>
											<div className="vle-open-question">
												<span className="vle-label">{labels.openQuestion}</span>
												<p>{copy.interpretation.openQuestion}</p>
											</div>
										</div>
									</div>
								</div>
							</div>
						) : null}
					</div>
				)
			})}
		</div>
	)
}
