import type { BookDemoErrorCode } from '@components/book-a-demo/form.schema'
import type { PositionId } from '@components/book-a-demo/positions'
import type { EmailErrorCode } from '@components/early-access/form.types'

export interface HomeMessages {
	skip: string
	hero: {
		title: string
		description: string
		seeMovement: string
		howItWorks: string
		caption: string
		summary: string
	}
	why: {
		kicker: string
		title: string
		description: string
		selected: string
		reason: string
		decision: string
		people: string
		product: string
	}
	attribution: {
		kicker: string
		title: string
		description: string
	}
	intent: {
		kicker: string
		title: string
		points: string
		events: string
		intentional: string
		unexplained: string
		review: string
	}
	voice: {
		title: string
		description: string
		question: string
		ask: string
		points: string
		largest: string
		answerTitle: string
		answer: string
		listening: string
		thinking: string
		speaking: string
		idle: string
	}
	report: {
		kicker: string
		title: string
		weekly: string
		vision: string
		largest: string
		movement: string
		why: string
		rationale: string
		who: string
		actor: string
		status: string
		attention: string
		action: string
	}
	pricing: {
		kicker: string
		title: string
		description: string
		cta: string
	}
	earlyAccess: {
		title: string
		description: string
	}
	footer: {
		description: string
		tagline: string
	}
}

export interface HeaderMessages {
	home: string
	navigation: string
	platform: string
	platformMenu: string
	vision: string
	finops: string
	microruntimes: string
	aiVoice: string
	visionDescription: string
	finopsDescription: string
	microruntimesDescription: string
	aiVoiceDescription: string
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

export interface DemoMessages {
	chart: {
		title: string
		movement: string
		movementDetail: string
		summary: string
		event: string
		why: string
		product: string
	}
	graph: {
		title: string
		description: string
		instructions: string
		summary: string
		closeLabel: string
		cardDescription: string
		cardBadge: string
		loading: string
		fallback: string
	}
	classifications: {
		baseline: string
		intentional: string
		review: string
		unexplained: string
	}
	teams: {
		product: string
		platform: string
		leadership: string
	}
	areas: {
		vision: string
		pricing: string
		authentication: string
		exports: string
	}
	events: {
		baseline: {
			action: string
			decision: string
			reason: string
			title: string
		}
		pricing: {
			action: string
			decision: string
			reason: string
			title: string
		}
		authentication: {
			action: string
			decision: string
			reason: string
			title: string
		}
		exports: {
			action: string
			decision: string
			reason: string
			title: string
		}
	}
	nodes: {
		vision: string
		target: string
		loop: string
		strategy: string
		baseline: string
		scope: string
		experience: string
		quality: string
		decision: string
		delivery: string
		feedback: string
		pricing: string
		onboarding: string
		evidence: string
		review: string
		exports: string
		accessibility: string
		tests: string
		reliability: string
		context: string
		intent: string
		release: string
		research: string
		authentication: string
		outcomes: string
		principles: string
		alignment: string
		roadmap: string
		segmentation: string
		packaging: string
		enterprise: string
		commercial: string
		discovery: string
		approval: string
		tradeoff: string
		identity: string
		session: string
		permissions: string
		security: string
		account: string
		roles: string
		migration: string
		verification: string
		formats: string
		reports: string
		retention: string
		workflow: string
		compatibility: string
	}
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
	roi: {
		eyebrow: string
		title: string
		description: string
		benefits: {
			context: {
				title: string
				description: string
			}
			tradeoff: {
				title: string
				description: string
			}
			opportunity: {
				title: string
				description: string
			}
		}
		earlyAccess: string
		explorePlans: string
		brief: {
			label: string
			context: string
			question: string
			evidence: [string, string, string]
		}
	}
	faq: {
		title: string
		description: string
		offer: string
		items: {
			meaning: {
				question: string
				answer: string
			}
			expected: {
				question: string
				answer: string
			}
			connection: {
				question: string
				answer: string
			}
			tools: {
				question: string
				answer: string
			}
			trial: {
				question: string
				answer: string
			}
			free: {
				question: string
				answer: string
			}
			paid: {
				question: string
				answer: string
			}
			availability: {
				question: string
				answer: string
			}
		}
	}
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
	homeHeading: string
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
	band: {
		title: string
		description: string
		cta: string
	}
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

export type IntegrationId = 'github' | 'obsidian' | 'linear'

export interface IntegrationsMessages {
	eyebrow: string
	title: string
	description: string
	project: {
		title: string
		label: string
		description: string
		connectionLabel: string
	}
	core: {
		title: string
		badge: string
		description: string
		supporting: string
	}
	tools: {
		github: {
			title: string
			badge: string
			description: string
			exampleLabel: string
			example: string
		}
		obsidian: {
			title: string
			badge: string
			description: string
			exampleLabel: string
			example: string
		}
		linear: {
			title: string
			badge: string
			description: string
			exampleLabel: string
			example: string
		}
	}
	technical: {
		title: string
		platform: string
		connectors: string
	}
	closing: {
		description: string
		cta: string
	}
}

export interface FooterMessages {
	navigation: string
	faq: string
	social: string
	copyright: string
}

export interface PlatformMessages {
	skip: string
	hero: {
		titleStart: string
		words: string[]
		titleEnd: string
		primaryCta: string
		secondaryCta: string
		// One line under each product in the hero, keyed by its header label.
		products: Record<'vision' | 'finops' | 'microruntimes' | 'aiVoice', string>
	}
	// The home sections' own copy; the headlines Vision and FinOps reuse live
	// in their catalogs.
	sections: {
		vision: { link: string }
		microruntimes: { title: string; description: string; link: string }
		finops: { description: string; link: string }
	}
	faq: {
		title: string
		description: string
		tabsLabel: string
	}
}
