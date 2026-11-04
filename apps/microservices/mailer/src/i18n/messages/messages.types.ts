import type { demoRequestLocales } from 'streaming'

export type MailLocale = (typeof demoRequestLocales)[number]

export interface ConfirmationMessages {
	subject: string
	greeting: string
	thanks: string
	followUp: string
	signoff: string
}

export interface LeadNotificationMessages {
	subject: string
	intro: string
	company: string
	contact: string
	position: string
	locale: string
	submittedAt: string
}

export interface DemoEmailsMessages {
	confirmation: ConfirmationMessages
	leadNotification: LeadNotificationMessages
}

export type DemoEmailsCatalog = Record<MailLocale, DemoEmailsMessages>
