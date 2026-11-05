import type { CiGatesMessages } from '@app/[locale]/(pages)/platform/ci-gates/ci-gates-messages.types'
import type { CliMessages } from '@app/[locale]/(pages)/platform/cli/cli-messages.types'
import type { SafeSeatMessages } from '@app/[locale]/(pages)/platform/safe-seat/safe-seat-messages.types'
import type { FinopsMessages } from '@data/finops-content'
import type { MicroruntimesMessages } from '@data/microruntimes-content'
import bookDemoJson from '@i18n/messages/book-a-demo.json'
import ciGatesJson from '@i18n/messages/ci-gates.json'
import cliJson from '@i18n/messages/cli.json'
import finopsJson from '@i18n/messages/finops.json'
import footerJson from '@i18n/messages/footer.json'
import headerJson from '@i18n/messages/header.json'
import type {
	BookDemoMessages,
	FooterMessages,
	HeaderMessages,
	PlatformMessages,
	PrivacyMessages,
} from '@i18n/messages/messages.types'
import microruntimesJson from '@i18n/messages/microruntimes.json'
import platformJson from '@i18n/messages/platform.json'
import privacyJson from '@i18n/messages/privacy.json'
import safeSeatJson from '@i18n/messages/safe-seat.json'
import type { WebsiteLocale } from '@i18n/routing'

export const headerMessages: Record<WebsiteLocale, HeaderMessages> = headerJson
export const privacyMessages: Record<WebsiteLocale, PrivacyMessages> =
	privacyJson
export const footerMessages: Record<WebsiteLocale, FooterMessages> = footerJson
export const platformMessages: Record<WebsiteLocale, PlatformMessages> =
	platformJson
export const finopsMessages: Record<WebsiteLocale, FinopsMessages> = finopsJson
export const microruntimesMessages: Record<
	WebsiteLocale,
	MicroruntimesMessages
> = microruntimesJson
export const safeSeatMessages: Record<WebsiteLocale, SafeSeatMessages> =
	safeSeatJson
export const cliMessages: Record<WebsiteLocale, CliMessages> = cliJson
export const ciGatesMessages: Record<WebsiteLocale, CiGatesMessages> =
	ciGatesJson
export const bookDemoMessages: Record<WebsiteLocale, BookDemoMessages> =
	bookDemoJson
