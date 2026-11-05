import type { Resources } from './microruntimes-content.types'

// What the envelope leaves once the plan's totals are taken from it.
export function headroom(bounded: Resources, suggested: Resources): Resources {
	return {
		cpuMillis: bounded.cpuMillis - suggested.cpuMillis,
		memoryMiB: bounded.memoryMiB - suggested.memoryMiB,
	}
}
