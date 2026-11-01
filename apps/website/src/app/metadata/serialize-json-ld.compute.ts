export function serializeJsonLd(value: object) {
	return JSON.stringify(value).replaceAll('<', '\\u003c')
}
