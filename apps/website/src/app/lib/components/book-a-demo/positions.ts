// The positions a visitor can pick when asking for a demo; the labels live in
// the book-a-demo catalog under the same ids.
export const positionIds = [
	'founder',
	'ceo',
	'cto',
	'vpEngineering',
	'vpProduct',
	'engineeringManager',
	'productManager',
	'staffEngineer',
	'softwareEngineer',
	'platformEngineer',
	'finance',
	'other',
] as const

export type PositionId = (typeof positionIds)[number]
