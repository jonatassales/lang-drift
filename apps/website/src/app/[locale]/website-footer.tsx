import { footerMessages, privacyMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'
import { AgentOrb } from 'web/ui/agent-orb'
import { LogoMark } from 'web/ui/brand'
import { GithubLogo, LinkedinLogo } from 'web/ui/icons/server'

import { NpmBadge } from './npm-badge'

const focusRing =
	'rounded focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-primary'
const brandLink = `inline-flex items-center gap-2 text-lg font-[650] tracking-[-0.055em] text-foreground hover:text-foreground ${focusRing}`
const footerLink = `inline-flex min-h-11 items-center text-[13px] text-muted-foreground hover:text-foreground ${focusRing}`
const creditLink = `inline-flex min-h-11 basis-full items-center hover:text-foreground md:basis-auto ${focusRing}`
const socialLink = `inline-flex size-11 items-center justify-center rounded-lg text-muted-foreground hover:text-foreground ${focusRing}`
const npmLink = `inline-flex h-11 items-center px-1 text-muted-foreground hover:text-foreground ${focusRing}`

interface WebsiteFooterProps {
	locale: WebsiteLocale
}

export function WebsiteFooter(props: WebsiteFooterProps) {
	const { locale } = props

	const copy = footerMessages[locale]

	return (
		<footer className="mx-auto max-w-[1240px] border-t border-border px-7 pt-13 pb-6 text-foreground [overflow-wrap:anywhere] md:[overflow-wrap:normal]">
			<div className="flex flex-col items-start justify-between gap-4.5 md:flex-row md:gap-8">
				<div className="max-w-[460px]">
					<Link className={brandLink} href="/" locale={locale}>
						<LogoMark size={23} />
						LangDrift
					</Link>
					<p className="mt-4 text-[15px] leading-[1.7] text-muted-foreground">
						{copy.description}
					</p>
				</div>
				<nav
					aria-label={copy.navigation}
					className="flex flex-wrap items-center gap-x-6 gap-y-2"
				>
					<Link className={footerLink} href="/#faq" locale={locale}>
						{copy.faq}
					</Link>
					<Link className={footerLink} href="/privacy">
						{privacyMessages[locale].title}
					</Link>
					<a
						className={footerLink}
						href="https://github.com/jonatassales/lang-drift/blob/main/SECURITY.md"
						rel="noopener noreferrer"
						target="_blank"
					>
						{copy.security}
					</a>
					<a
						className={footerLink}
						href="https://github.com/jonatassales/lang-drift/blob/main/LICENSE"
						rel="noopener noreferrer"
						target="_blank"
					>
						{copy.license}
					</a>
					<a className={footerLink} href="mailto:founders@langdrift.md">
						founders@langdrift.md
					</a>
				</nav>
			</div>
			<div className="mt-[34px] flex flex-col items-start justify-between gap-2.5 pt-4.5 md:flex-row md:items-center md:gap-5">
				<div className="flex flex-wrap items-center gap-2 text-xs leading-[1.6] text-muted-foreground">
					<span className="inline-flex size-7 shrink-0">
						<AgentOrb preset="neongate" size="28px" speed={0.35} />
					</span>
					<span>Neongate AI</span>
					<span className="basis-full md:basis-auto">{copy.copyright}</span>
					<a
						className={creditLink}
						href="https://langwatch.ai/"
						rel="noopener noreferrer"
						target="_blank"
					>
						{copy.acknowledgment}
					</a>
				</div>
				<nav
					aria-label={copy.social}
					className="flex flex-wrap items-center gap-1.5"
				>
					<a
						aria-label="LinkedIn"
						className={socialLink}
						href="https://www.linkedin.com/company/langdrift/"
						title="LinkedIn"
					>
						<LinkedinLogo aria-hidden="true" size={22} weight="fill" />
					</a>
					<a
						aria-label="GitHub"
						className={socialLink}
						href="https://github.com/jonatassales/lang-drift"
						title="GitHub"
					>
						<GithubLogo aria-hidden="true" size={22} weight="fill" />
					</a>
					<NpmBadge className={npmLink} packageName="@langdrift/cli" />
					<NpmBadge className={npmLink} packageName="@langdrift/core" />
				</nav>
			</div>
		</footer>
	)
}
