// The request a server-rendered read hands the console's route boundary: the
// incoming host and the path, nothing else. The browser's Origin, fetch
// metadata and cookies describe its navigation, not this read, so a page
// reached from another site, or rendered over HTTPS, reads normally.
export function serverRead(host: string, path: string) {
	return new Request(new URL(`/api/console/${path}`, `http://${host}`))
}
