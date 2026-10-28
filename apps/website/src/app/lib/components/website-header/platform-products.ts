import { FinopsIcon, MicroruntimesIcon } from '@langdrift/web/ui/icons/platform'

// The products the header's Platform menu lists, by their message keys. Each
// one has a page of its own.
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
] as const

export type PlatformProduct = (typeof platformProducts)[number]

export function platformProduct(label: PlatformProduct['label']) {
	const product = platformProducts.find((entry) => entry.label === label)
	if (!product) throw new Error(`Unknown platform product: ${label}`)

	return product
}
