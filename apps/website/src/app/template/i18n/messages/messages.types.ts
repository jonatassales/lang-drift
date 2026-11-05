import type { BookDemoErrorCode } from '@routes/book-a-demo/form.types'
import type { PositionId } from '@routes/book-a-demo/positions.data'

export interface HeaderMessages {
	home: string
	navigation: string
	platform: string
	platformMenu: string
	finops: string
	microruntimes: string
	ciGates: string
	safeSeat: string
	comingSoon: string
	finopsDescription: string
	microruntimesDescription: string
	cli: string
	cliDescription: string
	ciGatesDescription: string
	safeSeatDescription: string
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
	preferencesTitle: string
	preferences: string
	contactTitle: string
	contact: string
	changesTitle: string
	changes: string
}

export interface PlanCopy {
	name: string
	value: string
	support: string
	description: string
	listLabel: string
	items: [string, string, string]
	cta: string
}

export interface ConversionCopy {
	plans: {
		skip: string
		eyebrow: string
		comingSoon: string
		title: string
		introduction: string
		cards: {
			starter: PlanCopy
			plus: PlanCopy
			pro: PlanCopy
		}
		availability: string
		faqLink: string
	}
}

export interface NoticeCopy {
	title: string
	description: string
}

export interface PausedNoticeCopy extends NoticeCopy {
	home: string
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
		unavailable: NoticeCopy
		sent: NoticeCopy
	}
	paused: PausedNoticeCopy
	positions: Record<PositionId, string>
	errors: Record<BookDemoErrorCode, string>
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
