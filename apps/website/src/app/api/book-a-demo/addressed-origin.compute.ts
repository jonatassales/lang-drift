// The origin a request was addressed to. Next names every loopback host
// localhost in the request URL, so the Host header keeps 127.0.0.1 apart, as
// a browser does.
export function addressedOrigin(request: Request): string {
	const addressed = new URL(request.url)
	addressed.host = request.headers.get('host') ?? addressed.host

	return addressed.origin
}
