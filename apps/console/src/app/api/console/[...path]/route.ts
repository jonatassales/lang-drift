import { consoleHandler } from '@data/console-handler'

export const runtime = 'nodejs'
export const dynamic = 'force-dynamic'

export async function GET(request: Request): Promise<Response> {
	return consoleHandler(request)
}

export async function POST(request: Request): Promise<Response> {
	return consoleHandler(request)
}

export async function PATCH(request: Request): Promise<Response> {
	return consoleHandler(request)
}

export async function DELETE(request: Request): Promise<Response> {
	return consoleHandler(request)
}
