import { atom } from 'jotai'

import { isRequiredHumanStage } from '@routes/platform/microruntimes/content/is-required-human-stage.compute'
import {
	type PlanEffort,
	type PlanPhaseId,
	planSession,
} from '@routes/platform/microruntimes/content/microruntimes-content.data'
import {
	type HumanStageId,
	initialWorkflowChoices,
	type WorkflowChoices,
	type WorkflowModelId,
} from '@routes/platform/microruntimes/content/microruntimes-workflow.data'
import { phaseOf } from '@routes/platform/microruntimes/content/phase-of.compute'
import { runSession } from '@routes/platform/microruntimes/content/run-session.compute'

// The Microruntimes hero's simulated session: the plan it types first, then,
// once the visitor plays the workflow, the run. The session holds how many
// terminal lines have completed, whether the visitor paused it, a run counter
// that remounts the terminal, the card the visitor opened by hand, which
// script plays, and the choices a run was launched with. The terminal writes
// the count from its own clock; the cards derive every field from it.
export type SessionStage = 'plan' | 'launching' | 'running'

export interface MicroruntimeSession {
	/** The card the visitor opened; null follows the session's act. */
	openedCard: PlanPhaseId | null
	playedCount: number
	paused: boolean
	/** The choices the playing run was launched with. */
	ran: WorkflowChoices | null
	run: number
	stage: SessionStage
}

export const initialMicroruntimeSession: MicroruntimeSession = {
	openedCard: null,
	playedCount: 0,
	paused: false,
	ran: null,
	run: 0,
	stage: 'plan',
}

export const microruntimeSessionAtom = atom<MicroruntimeSession>(
	initialMicroruntimeSession,
)

// The model, effort and human stages the Workflow card shows. A replay
// keeps them; a run reads its own snapshot instead.
export const workflowChoicesAtom = atom<WorkflowChoices>(initialWorkflowChoices)

// The lines of the run the session is playing, or none before a run.
export const runScriptAtom = atom((get) => {
	const { ran } = get(microruntimeSessionAtom)

	if (!ran) return []

	return runSession(ran)
})

// The plan's progress: the live count while it plays, complete once the
// session has moved on to the workflow, so the plan cards keep their end.
export const planPlayedCountAtom = atom((get) => {
	const { playedCount, stage } = get(microruntimeSessionAtom)

	if (stage === 'plan') return playedCount

	return planSession.length
})

export const runPlayedCountAtom = atom((get) => {
	const { playedCount, stage } = get(microruntimeSessionAtom)

	if (stage !== 'running') return 0

	return playedCount
})

// True once the terminal has completed the script it is playing.
export const sessionEndedAtom = atom((get) => {
	const { playedCount, stage } = get(microruntimeSessionAtom)

	switch (stage) {
		case 'plan':
			return playedCount >= planSession.length
		case 'launching':
			return false
		case 'running':
			return playedCount >= get(runScriptAtom).length
	}
})

// True while a run is starting or playing: the choices it reads are locked.
export const runInProgressAtom = atom((get) => {
	const { stage } = get(microruntimeSessionAtom)

	if (stage === 'plan') return false

	return !get(sessionEndedAtom)
})

// True when the visitor may play the workflow: the plan has been handed to
// the team, or the last run has ended.
export const canPlayWorkflowAtom = atom((get) => {
	const { stage } = get(microruntimeSessionAtom)

	if (stage === 'launching') return false

	return get(sessionEndedAtom)
})

// The card open beside the terminal: the one the visitor chose, else the
// one the session's act is filling; during a run, the Workflow card.
export const openCardAtom = atom((get) => {
	const { openedCard, playedCount, stage } = get(microruntimeSessionAtom)

	if (openedCard) return openedCard
	if (stage !== 'plan') return 'workflow'

	return phaseOf(playedCount)
})

// Records how many lines the terminal has completed. When that moves the
// plan into its next act, the act takes the cards back from the visitor.
export const recordPlayedCountAtom = atom(
	null,
	(get, set, playedCount: number) => {
		const session = get(microruntimeSessionAtom)
		const actChanged =
			session.stage === 'plan' &&
			phaseOf(playedCount) !== phaseOf(session.playedCount)

		set(microruntimeSessionAtom, {
			...session,
			openedCard: actChanged ? null : session.openedCard,
			playedCount,
		})
	},
)

// Opens a card by hand. The open card stays open: one card is always
// expanded, so the column beside the terminal never collapses to its
// headers.
export const openCardByHandAtom = atom(null, (get, set, card: PlanPhaseId) => {
	const session = get(microruntimeSessionAtom)

	set(microruntimeSessionAtom, {
		...session,
		openedCard: card,
	})
})

export const chooseModelAtom = atom(
	null,
	(get, set, model: WorkflowModelId) => {
		if (get(runInProgressAtom)) return

		set(workflowChoicesAtom, {
			...get(workflowChoicesAtom),
			model,
		})
	},
)

export const chooseEffortAtom = atom(null, (get, set, effort: PlanEffort) => {
	if (get(runInProgressAtom)) return

	set(workflowChoicesAtom, {
		...get(workflowChoicesAtom),
		effort,
	})
})

export const toggleHumanStageAtom = atom(
	null,
	(get, set, stage: HumanStageId) => {
		if (get(runInProgressAtom)) return
		if (isRequiredHumanStage(stage)) return

		const choices = get(workflowChoicesAtom)

		set(workflowChoicesAtom, {
			...choices,
			human: {
				...choices.human,
				[stage]: !choices.human[stage],
			},
		})
	},
)

// Plays the workflow: the session starts launching with a snapshot of the
// choices, resumed and with the Workflow card in front. The terminal keeps
// the finished transcript until the run starts.
export const playWorkflowAtom = atom(null, (get, set) => {
	if (!get(canPlayWorkflowAtom)) return

	set(microruntimeSessionAtom, {
		...get(microruntimeSessionAtom),
		openedCard: null,
		paused: false,
		ran: get(workflowChoicesAtom),
		stage: 'launching',
	})
})

// Starts the launched run: the count resets and the run bumps, so the
// terminal remounts empty and types the run.
export const startRunAtom = atom(null, (get, set) => {
	const session = get(microruntimeSessionAtom)

	if (session.stage !== 'launching') return

	set(microruntimeSessionAtom, {
		...session,
		playedCount: 0,
		run: session.run + 1,
		stage: 'running',
	})
})

// Replays the session from the plan: the count, the pause, the opened card
// and the run reset, and the run bumps so the terminal remounts with a fresh
// machine. The visitor's workflow choices stay.
export const replaySessionAtom = atom(null, (get, set) => {
	const { run } = get(microruntimeSessionAtom)

	set(microruntimeSessionAtom, {
		...initialMicroruntimeSession,
		run: run + 1,
	})
})

// Pauses the session where it is, or resumes it.
export const toggleSessionPauseAtom = atom(null, (get, set) => {
	const session = get(microruntimeSessionAtom)

	set(microruntimeSessionAtom, {
		...session,
		paused: !session.paused,
	})
})
