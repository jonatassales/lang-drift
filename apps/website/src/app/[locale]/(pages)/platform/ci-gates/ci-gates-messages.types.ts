import type { CommandLineLabels } from '@routes/platform/command-line-labels.type'

import type { finopsGateIds, postureIds } from './ci-gates-content.data'

// The shape of the ci-gates catalog, keyed by the page's ids.

export type PostureId = (typeof postureIds)[number]
export type FinopsGateId = (typeof finopsGateIds)[number]

export interface CiGatesTitledCopy {
	title: string
	description: string
}

export interface CiGatesGateCopy extends CiGatesTitledCopy {
	eyebrow: string
}

export interface CiGatesSectionCopy {
	kicker: string
	title: string
	lead: string
}

export interface CiGatesHeroCopy {
	title: string
	lead: string
	note: string
}

export interface CiGatesInstallCopy extends CiGatesSectionCopy {
	command: string
}

export interface CiGatesPostureCopy extends CiGatesSectionCopy {
	items: Record<PostureId, CiGatesTitledCopy>
}

export interface CiGatesScreenshotCopy {
	alt: string
	caption: string
}

export interface CiGatesReportCopy extends CiGatesSectionCopy {
	summary: CiGatesScreenshotCopy
	dashboard: CiGatesScreenshotCopy
}

export interface CiGatesGatesCopy extends CiGatesSectionCopy {
	items: Record<FinopsGateId, CiGatesGateCopy>
}

export interface CiGatesMessages {
	skip: string
	hero: CiGatesHeroCopy
	install: CiGatesInstallCopy
	posture: CiGatesPostureCopy
	report: CiGatesReportCopy
	gates: CiGatesGatesCopy
	copyCommand: CommandLineLabels
}
