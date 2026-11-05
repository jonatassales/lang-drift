import type { CommandLineLabels } from '@app/[locale]/(pages)/platform/command-line-labels.type'

import type { cliFlowStepIds, cliGovernIds } from './cli-content.data'

// The shape of the cli catalog, keyed by the page's ids.

export type CliFlowStepId = (typeof cliFlowStepIds)[number]
export type CliGovernId = (typeof cliGovernIds)[number]

export interface CliTitledCopy {
	title: string
	description: string
}

export interface CliStepCopy extends CliTitledCopy {
	eyebrow: string
}

export interface CliSectionCopy {
	kicker: string
	title: string
	lead: string
}

export interface CliHeroCopy {
	title: string
	lead: string
	requirements: string
	command: string
}

export interface CliCommandsCopy {
	npx: string
	pnpm: string
}

export interface CliImageCopy {
	alt: string
	caption: string
}

export interface CliInstallCopy extends CliSectionCopy {
	commands: CliCommandsCopy
	image: CliImageCopy
}

export interface CliFlowCopy extends CliSectionCopy {
	steps: Record<CliFlowStepId, CliStepCopy>
}

export interface CliGovernsCopy extends CliSectionCopy {
	items: Record<CliGovernId, CliTitledCopy>
}

export interface CliLinksCopy {
	npm: string
	github: string
	license: string
}

export interface CliMessages {
	skip: string
	hero: CliHeroCopy
	install: CliInstallCopy
	flow: CliFlowCopy
	governs: CliGovernsCopy
	links: CliLinksCopy
	copyCommand: CommandLineLabels
}
