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

import './early-access-sent.css'

interface PageProps {
	params: Promise<{ locale: string }>
}

export async function generateMetadata(
	{ params }: PageProps,
	parent: ResolvingMetadata,
): Promise<Metadata> {
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

export default async function EarlyAccessSentPage({ params }: PageProps) {
	const { locale } = await params
	if (!hasLocale(routing.locales, locale)) notFound()
	setRequestLocale(locale)
	const copy = earlyAccessSentMessages[locale]

	return (
		<>
			<WebsiteHeader />
			<main className="early-access-sent">
				<div className="early-access-sent-content">
					<h1>{copy.title}</h1>
					<p>{copy.description}</p>
					<Image
						alt={copy.image}
						height={280}
						priority
						src="/early-access/early-bird.png"
						width={360}
					/>
					<Link href="/">{copy.home}</Link>
				</div>
			</main>
			<WebsiteFooter compact locale={locale} />
		</>
	)
}
