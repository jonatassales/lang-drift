import { extubRequest } from '@data/extub'

export async function extubProxyRoute(request: Request) {
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
