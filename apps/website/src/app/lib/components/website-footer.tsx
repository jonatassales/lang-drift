import { AgentOrb } from '@pack/react/ui/agent-orb'
import { LogoMark } from '@pack/react/ui/brand'
import { GithubLogo, LinkedinLogo } from '@pack/react/ui/icons/server'

import {
	footerMessages,
	homeMessages,
	privacyMessages,
} from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

import './website-footer.css'

interface WebsiteFooterProps {
	locale: WebsiteLocale
	compact?: boolean
}

export function WebsiteFooter(props: WebsiteFooterProps) {
	const { locale, compact = false } = props

	const copy = footerMessages[locale]

	return (
		<footer
			className="smooth-footer website-footer"
			data-compact={compact || undefined}
		>
			<div className="website-footer-top">
				{!compact ? (
					<div className="website-footer-intro">
						<Link className="website-footer-brand" href="/" locale={locale}>
							<LogoMark size={23} />
							LangDrift
						</Link>
						<p>{homeMessages[locale].footer.description}</p>
					</div>
				) : null}
				<nav aria-label={copy.navigation} className="website-footer-links">
					{/* <Link href="/pricing" locale={locale}>
            {headerMessages[locale].plans}
          </Link> */}
					<Link href="/#faq" locale={locale}>
						{copy.faq}
					</Link>
					<Link href="/privacy">{privacyMessages[locale].title}</Link>
					<a href="mailto:founders@langdrift.md">founders@langdrift.md</a>
				</nav>
			</div>
			<div className="website-footer-bottom">
				<div className="website-footer-credit">
					<span className="website-footer-orb">
						<AgentOrb preset="neongate" size="28px" speed={0.35} />
					</span>
					<span>Neongate AI</span>
					<span>{copy.copyright}</span>
				</div>
				<nav aria-label={copy.social} className="website-footer-social">
					<a
						href="https://www.linkedin.com/company/langdrift/"
						aria-label="LinkedIn"
						title="LinkedIn"
					>
						<LinkedinLogo aria-hidden="true" size={22} weight="fill" />
					</a>
					<a
						href="https://github.com/jonatassales/langdrift"
						aria-label="GitHub"
						title="GitHub"
					>
						<GithubLogo aria-hidden="true" size={22} weight="fill" />
					</a>
				</nav>
			</div>
		</footer>
	)
}
