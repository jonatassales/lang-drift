import { extubRequest } from '@data/extub'

export const runtime = 'nodejs'
export const dynamic = 'force-dynamic'

async function handle(request: Request) {
	const path = new URL(request.url).pathname.replace(/^\/api\/extub\//u, '')
	if (!['health', 'webrtc/offer'].includes(path))
		return Response.json(
			{
				error: {
					code: 'not_found',
				},
			},
			{
				status: 404,
			},
		)

	return extubRequest(request, path)
}

export const GET = handle
export const POST = handle
