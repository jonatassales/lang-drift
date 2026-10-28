// The two areas of the platform the console switches between.
export type Area = 'microruntimes' | 'finops'

export const areas: Area[] = ['microruntimes', 'finops']

// The route each area opens on.
export const areaHome: Record<Area, string> = {
	microruntimes: '/home',
	finops: '/finops',
}

// The area a console pathname belongs to: routes under /finops are FinOps,
// and every other route is Microruntimes. The URL is the only state.
export function areaOf(pathname: string): Area {
	if (pathname === '/finops' || pathname.startsWith('/finops/')) return 'finops'

	return 'microruntimes'
}
