// Sign-up is paused until feature flags have a proper home. To open it,
// restore the import and the POST export below and delete the paused POST.
// import { postEarlyAccess } from '@components/early-access/submission.route'

export const runtime = 'nodejs'

// export const POST = postEarlyAccess

export function POST() {
	return Response.json(
		{
			ok: false,
			code: 'UNAVAILABLE',
		},
		{
			status: 503,
			headers: {
				'Cache-Control': 'no-store',
				'Retry-After': '30',
				'X-Content-Type-Options': 'nosniff',
			},
		},
	)
}
