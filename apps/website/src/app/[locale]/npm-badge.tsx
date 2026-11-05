// A link to one package on npm, shown as its shields.io version badge in the
// npm red. The badge reads the latest published version, so the link never
// pins one by hand.
interface NpmBadgeProps {
	className?: string
	/** The scoped package name, e.g. `@langdrift/cli`. */
	packageName: string
}

export function NpmBadge(props: NpmBadgeProps) {
	const { className, packageName } = props

	const name = `npm: ${packageName}`
	const label = encodeURIComponent(packageName)
	const badgeSource = `https://img.shields.io/npm/v/${packageName}?logo=npm&logoColor=white&label=${label}&color=CB3837&labelColor=18181B`

	return (
		<a
			aria-label={name}
			className={className}
			href={`https://www.npmjs.com/package/${packageName}`}
			title={name}
		>
			<img
				alt={name}
				className="block h-5 w-auto"
				height={20}
				src={badgeSource}
			/>
		</a>
	)
}
