import { timingSafeEqual } from 'node:crypto'
import type { CanActivate, ExecutionContext } from '@nestjs/common'
import { Injectable, UnauthorizedException } from '@nestjs/common'
import type { Request } from 'express'

const scheme = 'Bearer '

// Only the website may create or read demo requests: it presents the shared
// token as a bearer. The token is read at the check and compared in constant
// time; a missing token is a configuration error, never an open door.
@Injectable()
export class BearerTokenGuard implements CanActivate {
	canActivate(context: ExecutionContext): boolean {
		const token = process.env.DEMO_REQUESTS_TOKEN
		if (!token) throw new Error('Missing DEMO_REQUESTS_TOKEN')

		const request = context.switchToHttp().getRequest<Request>()
		const presented = this.presentedToken(request.headers.authorization)
		if (!presented) throw new UnauthorizedException()

		const expected = Buffer.from(token)
		const actual = Buffer.from(presented)
		const matches =
			expected.length === actual.length && timingSafeEqual(expected, actual)
		if (!matches) throw new UnauthorizedException()

		return true
	}

	private presentedToken(authorization: string | undefined): string | null {
		if (!authorization?.startsWith(scheme)) return null

		return authorization.slice(scheme.length)
	}
}
