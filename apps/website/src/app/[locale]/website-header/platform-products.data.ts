import {
	CiGatesIcon,
	CliIcon,
	FinopsIcon,
	MicroruntimesIcon,
	SafeSeatIcon,
} from 'web/ui/icons/platform'

// The products the header's Platform menu lists, by their message keys, in
// the menu's order: the engine first, then what runs on it.
export const platformProducts = [
	{
		label: 'microruntimes',
		description: 'microruntimesDescription',
		icon: MicroruntimesIcon,
		href: '/platform/microruntimes',
	},
	{
		label: 'finops',
		description: 'finopsDescription',
		icon: FinopsIcon,
		href: '/platform/finops',
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
	{
		label: 'safeSeat',
		description: 'safeSeatDescription',
		icon: SafeSeatIcon,
		href: '/platform/safe-seat',
		badge: 'comingSoon',
	},
] as const

export type PlatformProduct = (typeof platformProducts)[number]
