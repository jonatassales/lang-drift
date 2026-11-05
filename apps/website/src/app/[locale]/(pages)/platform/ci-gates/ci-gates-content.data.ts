// CI Gates page content: the ids and screenshots the page renders. Sentences
// live in the ci-gates catalog; this module holds what a reader does not read
// as prose.

export const postureIds = ['token', 'pinned', 'bounded', 'unknown'] as const

export const finopsGateIds = [
	'fg-01',
	'fg-02',
	'fg-03',
	'fg-04',
	'fg-05',
	'fg-06',
] as const

// Renderings of pull request #122 of the LangDrift repository, one per color
// scheme, served from public/: the job summary dft gates report wrote, and the
// repository's own dashboard comment.
export const ciGatesScreenshots = {
	dashboard: {
		dark: '/ci-gates/dashboard-dark.png',
		height: 760,
		light: '/ci-gates/dashboard-light.png',
		width: 1800,
	},
	summary: {
		dark: '/ci-gates/report-dark.png',
		height: 1094,
		light: '/ci-gates/report-light.png',
		width: 1800,
	},
} as const
