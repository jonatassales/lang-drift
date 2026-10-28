import {
	AiVoiceIcon,
	FinopsIcon,
	MicroruntimesIcon,
	VisionIcon,
} from '@langdrift/web/ui/icons/platform'

// The products the header's Platform menu lists, by their message keys. Only
// a product with a page links; the others are named with a line about them.
export const platformProducts = [
	{
		label: 'vision',
		description: 'visionDescription',
		icon: VisionIcon,
		href: '/platform/vision',
	},
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
		label: 'aiVoice',
		description: 'aiVoiceDescription',
		icon: AiVoiceIcon,
	},
] as const

export type PlatformProduct = (typeof platformProducts)[number]

export function platformProduct(label: PlatformProduct['label']) {
	const product = platformProducts.find((entry) => entry.label === label)
	if (!product) throw new Error(`Unknown platform product: ${label}`)

	return product
}
