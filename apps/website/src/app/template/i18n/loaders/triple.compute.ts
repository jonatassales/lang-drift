export function triple(
	value: readonly string[],
	label: string,
): [string, string, string] {
	if (value.length !== 3) {
		throw new Error(`${label} must contain exactly three strings`)
	}

	const first = value[0]
	const second = value[1]
	const third = value[2]
	const allStrings =
		typeof first === 'string' &&
		typeof second === 'string' &&
		typeof third === 'string'

	if (!allStrings) {
		throw new Error(`${label} must contain exactly three strings`)
	}

	return [first, second, third]
}
