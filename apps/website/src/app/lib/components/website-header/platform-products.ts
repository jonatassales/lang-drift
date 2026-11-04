import {
	CiGatesIcon,
	FinopsIcon,
	MicroruntimesIcon,
	SafeSeatIcon,
} from '@langdrift/web/ui/icons/platform'

// The products the header's Platform menu lists, by their message keys. Only
// a product with a page links; the others are named with a line about them.
export const platformProducts = [
	{
		label: 'finops',
		description: 'finopsDescription',
		icon: FinopsIcon,
		href: '/platform/finops',
	},
	{
		label: 'microruntimes',
		description: 'microruntimesDescription',
		icon: MicroruntimesIcon,
		href: '/platform/microruntimes',
	},
	{
		label: 'ciGates',
		description: 'ciGatesDescription',
		icon: CiGatesIcon,
		badge: 'comingSoon',
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

export function platformProduct(label: PlatformProduct['label']) {
	const product = platformProducts.find((entry) => entry.label === label)
	if (!product) throw new Error(`Unknown platform product: ${label}`)

	return product
}
