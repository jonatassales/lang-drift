// The paths excluded from the console proxy: API routes, Next internal assets,
// root static files, and any file with an extension.
export const proxyMatcher = [
	'/((?!api|_next|favicon\\.ico|icon\\.svg|apple-icon\\.png|robots\\.txt|.*\\..*).*)',
]
