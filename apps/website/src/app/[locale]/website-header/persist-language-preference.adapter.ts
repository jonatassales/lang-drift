const storageKey = 'langdrift-website-locale'

// Remembers a visitor's language for their next visit and answers whether
// the browser kept it. Routing never depends on it: the proxy's cookie
// decides the locale.
export function persistLanguagePreference(locale: string): boolean {
	try {
		window.localStorage.setItem(storageKey, locale)

		return true
	} catch {
		return false
	}
}
