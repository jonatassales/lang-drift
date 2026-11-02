import type { CSSProperties } from 'react'

import { LogoMark } from './logo-mark'

export interface BrandProps {
	className?: string
	compact?: boolean
	/** A fixed wordmark color; without it the wordmark takes the surrounding text color. */
	tone?: 'dark' | 'light'
}

// The fixed wordmark colors a caller can ask for by tone.
const toneColor: Record<NonNullable<BrandProps['tone']>, string> = {
	dark: '#f5f5f5',
	light: '#111111',
}

const wordmarkStyle: CSSProperties = {
	letterSpacing: '-0.055em',
	lineHeight: 1,
}

export function Brand(props: BrandProps) {
	const { className, compact = false, tone } = props

	const textColor = tone ? toneColor[tone] : undefined

	return (
		<span
			className={className}
			style={{
				alignItems: 'center',
				color: textColor,
				display: 'inline-flex',
				fontFamily: 'var(--lg-typography-font-body)',
				fontWeight: 650,
				gap: compact ? '8px' : '11px',
			}}
		>
			<LogoMark size={compact ? 23 : 30} />
			<span
				style={{
					...wordmarkStyle,
					fontSize: compact ? 18 : 24,
				}}
			>
				LangDrift
			</span>
		</span>
	)
}
