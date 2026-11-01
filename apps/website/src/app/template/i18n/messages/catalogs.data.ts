import type { FinopsMessages } from '@data/finops-content'
import type { MicroruntimesMessages } from '@data/microruntimes-content'
import type { WebsiteLocale } from '@i18n/routing'

import bookDemoJson from './book-a-demo.json'
import earlyAccessJson from './early-access.json'
import earlyAccessSentJson from './early-access-sent.json'
import finopsJson from './finops.json'
import footerJson from './footer.json'
import headerJson from './header.json'
import type {
	BookDemoMessages,
	EarlyAccessMessages,
	EarlyAccessSentMessages,
	FooterMessages,
	HeaderMessages,
	PlatformMessages,
	PrivacyMessages,
} from './messages.types'
import microruntimesJson from './microruntimes.json'
import platformJson from './platform.json'
import privacyJson from './privacy.json'

export const headerMessages: Record<WebsiteLocale, HeaderMessages> = headerJson
export const privacyMessages: Record<WebsiteLocale, PrivacyMessages> =
	privacyJson
export const earlyAccessMessages: Record<WebsiteLocale, EarlyAccessMessages> =
	earlyAccessJson
export const earlyAccessSentMessages: Record<
	WebsiteLocale,
	EarlyAccessSentMessages
> = earlyAccessSentJson
export const footerMessages: Record<WebsiteLocale, FooterMessages> = footerJson
export const platformMessages: Record<WebsiteLocale, PlatformMessages> =
	platformJson
export const finopsMessages: Record<WebsiteLocale, FinopsMessages> = finopsJson
export const microruntimesMessages: Record<
	WebsiteLocale,
	MicroruntimesMessages
> = microruntimesJson
export const bookDemoMessages: Record<WebsiteLocale, BookDemoMessages> =
	bookDemoJson
