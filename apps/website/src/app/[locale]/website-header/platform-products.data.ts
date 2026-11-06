import {
	CiGatesIcon,
	CliIcon,
	FinopsIcon,
	LlmopsIcon,
	MicroruntimesIcon,
	SafeSeatIcon,
} from 'web/ui/icons/platform'

// The products the header's Platform menu lists, by their message keys, in
// the menu's order: the products fill the left column, the utilities that
// run them the right one.
export const platformProducts = [
	{
		label: 'llmops',
		description: 'llmopsDescription',
		icon: LlmopsIcon,
		href: '/platform/llmops',
	},
	{
		label: 'finops',
		description: 'finopsDescription',
		icon: FinopsIcon,
		href: '/platform/finops',
	},
	{
		label: 'safeSeat',
		description: 'safeSeatDescription',
		icon: SafeSeatIcon,
		href: '/platform/safe-seat',
		badge: 'comingSoon',
	},
	{
		label: 'microruntimes',
		description: 'microruntimesDescription',
		icon: MicroruntimesIcon,
		href: '/platform/microruntimes',
	},
	{
		label: 'cli',
		description: 'cliDescription',
		icon: CliIcon,
		href: '/platform/cli',
	},
	{
		label: 'ciGates',
		description: 'ciGatesDescription',
		icon: CiGatesIcon,
		href: '/platform/ci-gates',
	},
] as const

export type PlatformProduct = (typeof platformProducts)[number]
