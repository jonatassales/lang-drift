import Image from 'next/image'

import type { CiGatesScreenshotCopy } from './ci-gates-messages.types'

interface ScreenshotSources {
	dark: string
	height: number
	light: string
	width: number
}

// A screenshot captured in both color schemes: the page shows the one that
// matches the theme, with its caption below.
interface ThemedScreenshotProps {
	copy: CiGatesScreenshotCopy
	screenshot: ScreenshotSources
}

export function ThemedScreenshot(props: ThemedScreenshotProps) {
	const { copy, screenshot } = props

	const imageClass = 'h-auto w-full rounded-2xl border border-border'

	return (
		<figure className="m-0">
			<Image
				alt={copy.alt}
				className={`${imageClass} dark:hidden`}
				height={screenshot.height}
				sizes="(min-width: 1180px) 1108px, 100vw"
				src={screenshot.light}
				width={screenshot.width}
			/>
			<Image
				alt={copy.alt}
				className={`${imageClass} hidden dark:block`}
				height={screenshot.height}
				sizes="(min-width: 1180px) 1108px, 100vw"
				src={screenshot.dark}
				width={screenshot.width}
			/>
			<figcaption className="mt-3 text-xs leading-relaxed text-muted-foreground">
				{copy.caption}
			</figcaption>
		</figure>
	)
}
