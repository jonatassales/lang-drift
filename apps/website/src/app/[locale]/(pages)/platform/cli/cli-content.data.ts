// LangDrift CLI page content: the ids, links and image the page renders.
// Sentences live in the cli catalog; this module holds what a reader does not
// read as prose.

export const cliFlowStepIds = ['init', 'plan', 'run', 'status'] as const

export const cliGovernIds = ['budget', 'runtime', 'driver', 'evidence'] as const

export const cliLinks = {
	github: 'https://github.com/jonatassales/lang-drift/tree/main/lib/cli',
	license: 'https://github.com/jonatassales/lang-drift/blob/main/LICENSE',
	npm: 'https://www.npmjs.com/package/@langdrift/cli',
} as const

// The terminal screenshot of dft agent plan, served from public/.
export const cliPlanImage = {
	height: 1479,
	src: '/cli/dft-agent-plan.webp',
	width: 2000,
} as const
