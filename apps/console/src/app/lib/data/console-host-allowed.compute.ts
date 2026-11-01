import { hosted } from './hosted.compute.ts'

const loopbackHosts = ['localhost', '127.0.0.1', '[::1]']

export function consoleHostAllowed(
	hostname: string,
	token: string | undefined,
) {
	return hosted(token) || loopbackHosts.includes(hostname)
}
