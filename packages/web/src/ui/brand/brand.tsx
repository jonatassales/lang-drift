import type { CSSProperties } from 'react'

export interface BrandProps {
	className?: string
	compact?: boolean
	tone?: 'dark' | 'light'
}

export interface LogoMarkProps {
	size?: number
}

const wordmarkStyle: CSSProperties = {
	letterSpacing: '-0.055em',
	lineHeight: 1,
}

// The drift mark, in black and white: a ring with one segment carried outward.
// The ring is the vision; the piece that moved is the drift whose why
// LangDrift keeps. Both strokes take the foreground, so the mark reads in
// either color scheme.
export function LogoMark(props: LogoMarkProps) {
	const { size = 31 } = props

	return (
		<svg
			aria-hidden="true"
			fill="none"
			height={size}
			strokeLinecap="round"
			strokeWidth="9"
			viewBox="0 0 64 64"
			width={size}
		>
			<path
				className="stroke-foreground"
				d="M45.54 49.34A22 22 0 1 1 45.54 14.66"
			/>
			<path
				className="stroke-foreground"
				d="M56.72 18.86A28 28 0 0 1 56.72 45.14"
			/>
		</svg>
	)
}

export function Brand(props: BrandProps) {
	const { className, compact = false, tone = 'light' } = props

	const textColor = tone === 'dark' ? '#f5f5f5' : '#111111'

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
