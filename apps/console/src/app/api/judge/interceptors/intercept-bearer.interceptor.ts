import { judgeReply } from '@app/api/judge/judge-reply.map'
import { secretMatches } from '@app/api/judge/secret-matches.compute'

const bearerPrefix = 'Bearer '

// Only a request that carries the deployment's secret as a bearer reaches the
// judge. A deployment without a secret answers no one.
export function interceptBearer(request: Request): Response | undefined {
	const secret = process.env.RUNTIME_SECRET

	if (!secret) {
		return judgeReply(503, {
			code: 'UNAVAILABLE',
			reason: 'judge-not-configured',
		})
	}

	const header = request.headers.get('authorization') ?? ''
	const presented = header.startsWith(bearerPrefix)
		? header.slice(bearerPrefix.length)
		: ''
	const authorized = presented !== '' && secretMatches(presented, secret)

	if (authorized) return undefined

	return judgeReply(401, {
		code: 'UNAUTHORIZED',
	})
}
