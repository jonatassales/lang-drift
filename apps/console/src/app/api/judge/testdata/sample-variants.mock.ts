// Values each field of a sample is replaced with, one at a time: absent and
// empty values, a value outside every enum, numbers on either side of the
// contract's bounds, a fraction, and values of every other JSON type.
const probes: unknown[] = [
	null,
	'',
	'outside-every-enum',
	-1,
	0,
	1.5,
	2,
	true,
	[],
	{},
]

// Every value one edit away from a contract sample: the sample replaced by
// each probe; for an object, each key removed, an unknown key added and each
// field's own variants; for an array, empty, four items long and its first
// item's variants. A schema and its mirror must agree on every one.
export function sampleVariants(sample: unknown): unknown[] {
	if (Array.isArray(sample)) {
		const [first, ...rest] = sample

		if (first === undefined) return probes

		const items = sampleVariants(first).map((variant) => [variant, ...rest])

		return [...probes, [], [first, first, first, first], ...items]
	}

	if (sample === null || typeof sample !== 'object') return probes

	const entries = Object.entries(sample)
	const removed = entries.map(([key]) =>
		Object.fromEntries(entries.filter(([other]) => other !== key)),
	)
	const replaced = entries.flatMap(([key, field]) =>
		sampleVariants(field).map((variant) => ({
			...sample,
			[key]: variant,
		})),
	)
	const unknownKey = {
		...sample,
		unexpected: 'field',
	}

	return [...probes, ...removed, unknownKey, ...replaced]
}
