import bookDemoJson from '@i18n/messages/book-a-demo.json'
import ciGatesJson from '@i18n/messages/ci-gates.json'
import cliJson from '@i18n/messages/cli.json'
import finopsJson from '@i18n/messages/finops.json'
import footerJson from '@i18n/messages/footer.json'
import headerJson from '@i18n/messages/header.json'
import llmopsJson from '@i18n/messages/llmops.json'
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
import type { CiGatesMessages } from '@routes/platform/ci-gates/ci-gates-messages.types'
import type { CliMessages } from '@routes/platform/cli/cli-messages.types'
import type { FinopsMessages } from '@routes/platform/finops/content/finops-content.types'
import type { LlmopsMessages } from '@routes/platform/llmops/content/llmops-content.types'
import type { MicroruntimesMessages } from '@routes/platform/microruntimes/content/microruntimes-content.types'
import type { SafeSeatMessages } from '@routes/platform/safe-seat/safe-seat-messages.types'

import { assertCatalogShape } from './assert-catalog-shape.compute'

export const headerMessages: Record<WebsiteLocale, HeaderMessages> =
	assertCatalogShape('header', headerJson)
export const privacyMessages: Record<WebsiteLocale, PrivacyMessages> =
	assertCatalogShape('privacy', privacyJson)
export const footerMessages: Record<WebsiteLocale, FooterMessages> =
	assertCatalogShape('footer', footerJson)
export const platformMessages: Record<WebsiteLocale, PlatformMessages> =
	assertCatalogShape('platform', platformJson)
export const llmopsMessages: Record<WebsiteLocale, LlmopsMessages> =
	assertCatalogShape('llmops', llmopsJson)
export const finopsMessages: Record<WebsiteLocale, FinopsMessages> =
	assertCatalogShape('finops', finopsJson)
export const microruntimesMessages: Record<
	WebsiteLocale,
	MicroruntimesMessages
> = assertCatalogShape('microruntimes', microruntimesJson)
export const safeSeatMessages: Record<WebsiteLocale, SafeSeatMessages> =
	assertCatalogShape('safe-seat', safeSeatJson)
export const cliMessages: Record<WebsiteLocale, CliMessages> =
	assertCatalogShape('cli', cliJson)
export const ciGatesMessages: Record<WebsiteLocale, CiGatesMessages> =
	assertCatalogShape('ci-gates', ciGatesJson)
export const bookDemoMessages: Record<WebsiteLocale, BookDemoMessages> =
	assertCatalogShape('book-a-demo', bookDemoJson)
