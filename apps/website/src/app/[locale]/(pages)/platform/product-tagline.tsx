interface ProductTaglineProps {
	end: string
	start: string
}

// The sentence a platform product page puts under its heading, in two lines.
export function ProductTagline(props: ProductTaglineProps) {
	const { end, start } = props

	return (
		<p className="mt-4 mb-0 text-xl leading-snug text-muted-foreground sm:text-2xl">
			<span className="block">{start}</span>{' '}
			<span className="block">{end}</span>
		</p>
	)
}
