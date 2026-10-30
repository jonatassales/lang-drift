import Image from 'next/image'
import type { ReactNode } from 'react'

import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { Link } from '@i18n/navigation'
import type { WebsiteLocale } from '@i18n/routing'

interface RouteStateProps {
	locale: WebsiteLocale
	title: string
	description: string
	imageSrc: string
	imageAlt: string
	homeLabel: string
	children?: ReactNode
}

export function RouteState(props: RouteStateProps) {
	const {
		locale,
		title,
		description,
		imageSrc,
		imageAlt,
		homeLabel,
		children,
	} = props

	return (
		<>
			<WebsiteHeader />
			<main className="flex min-h-[calc(100svh-170px)] items-center justify-center px-5 pt-12 pb-16 text-center">
				<div className="max-w-2xl">
					<h1 className="text-4xl leading-tight tracking-tighter text-balance text-foreground sm:text-5xl">
						{title}
					</h1>
					<p className="mx-auto mt-4 max-w-md text-sm leading-relaxed text-pretty text-muted-foreground sm:text-base">
						{description}
					</p>
					<Image
						alt={imageAlt}
						className="mx-auto my-6 block size-72 sm:size-80"
						height={320}
						preload
						src={imageSrc}
						width={320}
					/>
					<div className="flex flex-wrap items-center justify-center gap-3">
						<Link
							className="inline-flex min-h-10 items-center rounded-full bg-primary px-5 text-sm font-semibold text-primary-foreground transition hover:-translate-y-px hover:bg-primary-hover"
							href="/"
						>
							{homeLabel}
						</Link>
						{children}
					</div>
				</div>
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
