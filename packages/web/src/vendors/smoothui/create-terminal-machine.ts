export type MachinePhase = 'delay' | 'typing'

export interface Machine {
	finished: boolean
	lineIndex: number
	phase: MachinePhase
	phaseElapsed: number
}

export function createTerminalMachine(): Machine {
	return {
		finished: false,
		lineIndex: 0,
		phase: 'delay',
		phaseElapsed: 0,
	}
}
