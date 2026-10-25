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
		label: 'finops',
		description: 'finopsDescription',
		icon: FinopsIcon,
	},
	{
		label: 'microruntimes',
		description: 'microruntimesDescription',
		icon: MicroruntimesIcon,
	},
	{
		label: 'aiVoice',
		description: 'aiVoiceDescription',
		icon: AiVoiceIcon,
	},
] as const
