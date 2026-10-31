import type { Metadata, ResolvingMetadata } from 'next'
import Image from 'next/image'
import { notFound } from 'next/navigation'
import { hasLocale } from 'next-intl'
import { setRequestLocale } from 'next-intl/server'

import { socialImages } from '@app/metadata/page-metadata'
import { WebsiteFooter } from '@components/website-footer'
import { WebsiteHeader } from '@components/website-header'
import { earlyAccessSentMessages } from '@i18n/messages/index'
import { Link } from '@i18n/navigation'
import { routing } from '@i18n/routing'

interface PageProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	props: PageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	const title = earlyAccessSentMessages[locale].title

	const images = await socialImages(parent)

	const openGraphImages = images.openGraph
	const twitterImages = images.twitter

	return {
		title: { absolute: title },
		robots: {
			index: false,
			follow: false,
		},
		openGraph: {
			title,
			...(openGraphImages ? { images: openGraphImages } : {}),
		},
		twitter: {
			card: 'summary_large_image',
			title,
			...(twitterImages ? { images: twitterImages } : {}),
		},
	}
}

export default async function EarlyAccessSentPage(props: PageProps) {
	const { params } = props

	const { locale } = await params

	if (!hasLocale(routing.locales, locale)) notFound()

	setRequestLocale(locale)

	const copy = earlyAccessSentMessages[locale]

	return (
		<>
			<WebsiteHeader />
			<main className="flex min-h-[calc(100svh-170px)] items-center justify-center px-5 pt-12 pb-16 text-center">
				<div className="max-w-2xl">
					<h1 className="text-4xl leading-tight tracking-tighter text-balance text-foreground sm:text-5xl">
						{copy.title}
					</h1>
					<p className="mx-auto mt-4 max-w-sm text-sm text-pretty text-primary-text">
						{copy.description}
					</p>
					<Image
						alt={copy.image}
						className="mx-auto my-7 block h-auto max-w-full"
						height={280}
						preload
						src="/early-access/early-bird.png"
						width={360}
					/>
					<Link
						className="inline-block p-2.5 text-foreground focus-visible:outline-2 focus-visible:outline-offset-3 focus-visible:outline-primary"
						href="/"
					>
						{copy.home}
					</Link>
				</div>
			</main>
			<WebsiteFooter locale={locale} />
		</>
	)
}
