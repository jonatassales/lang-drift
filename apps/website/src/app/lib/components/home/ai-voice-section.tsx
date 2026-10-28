import { HomeProductSection } from '@components/home/home-product-section'
import { VoicePreview } from '@components/voice-preview/voice-preview'
import { platformProduct } from '@components/website-header/platform-products'
import { homeMessages } from '@i18n/messages/index'
import type { WebsiteLocale } from '@i18n/routing'

interface AiVoiceSectionProps {
	locale: WebsiteLocale
}

// AI Voice on the home: the orb playing the scripted sample answer under the
// voice title and its disclosure. No link, since the product has no page yet.
export function AiVoiceSection(props: AiVoiceSectionProps) {
	const { locale } = props

	const copy = homeMessages[locale].voice

	return (
		<HomeProductSection
			description={copy.description}
			headline={copy.title}
			id="ai-voice"
			locale={locale}
			product={platformProduct('aiVoice')}
		>
			<VoicePreview />
		</HomeProductSection>
	)
}
