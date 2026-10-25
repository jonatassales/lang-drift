import Image from 'next/image'

import type { IntegrationId } from '@i18n/messages/messages.types'

const integrationLogos: Record<
	IntegrationId,
	{ light: string; dark?: string; width: number }
> = {
	github: {
		light: 'github-black.svg',
		dark: 'github-white.svg',
		width: 39,
	},
	obsidian: {
		light: 'obsidian.svg',
		width: 38,
	},
	linear: {
		light: 'linear-dark.svg',
		dark: 'linear-light.svg',
		width: 38,
	},
}

interface IntegrationLogoProps {
	id: IntegrationId
}

export function IntegrationLogo(props: IntegrationLogoProps) {
	const { id } = props

	const { light, dark, width } = integrationLogos[id]

	if (!dark) {
		return (
			<Image alt="" height={38} src={`/integrations/${light}`} width={width} />
		)
	}

	return (
		<>
			<Image
				alt=""
				className="integration-logo-light"
				height={38}
				src={`/integrations/${light}`}
				width={width}
			/>
			<Image
				alt=""
				className="integration-logo-dark"
				height={38}
				src={`/integrations/${dark}`}
				width={width}
			/>
		</>
	)
}
