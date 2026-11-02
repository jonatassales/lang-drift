import type { BookDemoErrorCode } from '@components/book-a-demo/form.types'
import type { PositionId } from '@components/book-a-demo/positions'
import type { EmailErrorCode } from '@components/early-access/form.types'

export interface HeaderMessages {
	home: string
	navigation: string
	platform: string
	platformMenu: string
	finops: string
	microruntimes: string
	ciGates: string
	comingSoon: string
	finopsDescription: string
	microruntimesDescription: string
	ciGatesDescription: string
	pricing: string
	docs: string
	github: string
	bookDemo: string
	language: string
	languages: {
		english: string
		portuguese: string
		chinese: string
		japanese: string
		vietnamese: string
	}
	lightTheme: string
	darkTheme: string
	openNavigation: string
	closeNavigation: string
}

export interface PrivacyMessages {
	title: string
	description: string
	updated: string
	skip: string
	introduction: string
	overviewTitle: string
	overview: string
	paused: string
	active: string
	turnstileTitle: string
	turnstile: string[]
	addendumIntroduction: string
	addendumLabel: string
	addendumUrl: string
	providerPolicy: string
	preferencesTitle: string
	preferences: string
	contactTitle: string
	contact: string
	changesTitle: string
	changes: string
}

export interface ConversionCopy {
	plans: {
		skip: string
		eyebrow: string
		title: string
		introduction: string
		cards: {
			trial: {
				name: string
				badge: string
				value: string
				support: string
				description: string
				listLabel: string
				items: [string, string, string]
				cta: string
			}
			plus: {
				name: string
				badge: string
				value: string
				support: string
				description: string
				listLabel: string
				items: [string, string, string]
				cta: string
			}
			pro: {
				name: string
				badge: string
				value: string
				support: string
				description: string
				listLabel: string
				items: [string, string, string]
				cta: string
			}
		}
		availability: string
		faqLink: string
	}
}

export interface EarlyAccessMessages {
	pausedHeading: string
	pausedDescription: string
	pricingHeading: string
	description: string
	helper: string
	emailLabel: string
	placeholder: string
	submit: string
	submitting: string
	success: string
	unavailable: string
	rateLimited: string
	challengeFailed: string
	errors: Record<EmailErrorCode, string>
}

export interface BookDemoMessages {
	skip: string
	title: string
	lead: string
	form: {
		email: string
		emailPlaceholder: string
		company: string
		companyPlaceholder: string
		position: string
		positionPlaceholder: string
		submit: string
		submitting: string
		failed: string
		unavailable: {
			title: string
			description: string
		}
	}
	positions: Record<PositionId, string>
	errors: Record<BookDemoErrorCode, string>
}

export interface EarlyAccessSentMessages {
	title: string
	description: string
	home: string
	image: string
}

export interface FooterMessages {
	description: string
	navigation: string
	faq: string
	social: string
	copyright: string
	security: string
	license: string
}

export interface PlatformMessages {
	skip: string
	hero: {
		titleStart: string
		words: string[]
		titleEnd: string
	}
	// The home sections' own copy; the headline FinOps reuses lives in its
	// catalog.
	sections: {
		finops: { description: string; link: string }
	}
	faq: {
		title: string
		description: string
		tabsLabel: string
	}
}
