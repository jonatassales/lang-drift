import { consoleStubsRequest } from '@data/console-stubs'

export const runtime = 'nodejs'
export const dynamic = 'force-dynamic'

async function forward(request: Request): Promise<Response> {
	const path = new URL(request.url).pathname.replace(
		/^\/api\/console-stubs\//u,
		'',
	)
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

	return consoleStubsRequest(request, path)
}

export async function GET(request: Request): Promise<Response> {
	return forward(request)
}

export async function POST(request: Request): Promise<Response> {
	return forward(request)
}
