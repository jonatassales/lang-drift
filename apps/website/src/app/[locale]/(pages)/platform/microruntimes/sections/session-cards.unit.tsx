// @vitest-environment jsdom

import { act, cleanup, fireEvent, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'

import { launchWorkflowRun } from '@routes/platform/microruntimes/sections/testdata/launch-workflow-run.mock'
import { mountSessionCards } from '@routes/platform/microruntimes/sections/testdata/mount-session-cards.mock'
import { playRunTo } from '@routes/platform/microruntimes/sections/testdata/play-run-to.mock'
import { workflowPlayButton } from '@routes/platform/microruntimes/sections/testdata/workflow-play-button.mock'
import { workflowSelects } from '@routes/platform/microruntimes/sections/testdata/workflow-selects.mock'
import {
	initialMicroruntimeSession,
	microruntimeSessionAtom,
	replaySessionAtom,
	runScriptAtom,
	startRunAtom,
	workflowChoicesAtom,
} from '@state'

afterEach(cleanup)

describe('Microruntimes session cards', () => {
	it('stays hidden with the ticket open until the command is entered', () => {
		const { expanded, play } = mountSessionCards()
		const stack = screen.getByText('0031').closest('[aria-hidden]')
		expect(stack?.getAttribute('aria-hidden')).toBe('true')
		expect(expanded('Ticket')).toBe(true)
		expect(expanded('Analysis')).toBe(false)
		expect(screen.getAllByText('Not set yet').length).toBeGreaterThanOrEqual(4)

		play('plan')
		expect(stack?.getAttribute('aria-hidden')).toBeNull()

		play('links')
		expect(screen.getAllByText('0002-PRD-book-a-demo.prd.md')).toHaveLength(2)
	})

	it('opens the card of each act as the session reaches it', () => {
		const { expanded, play } = mountSessionCards()
		play('links')
		expect(expanded('Ticket')).toBe(true)

		play('analysis')
		expect(expanded('Ticket')).toBe(false)
		expect(expanded('Analysis')).toBe(true)

		play('planning')
		expect(expanded('Analysis')).toBe(false)
		expect(expanded('Plan')).toBe(true)
	})

	it('lets the visitor open a card by hand until the next act takes it back', () => {
		const { expanded, header, play, store } = mountSessionCards()
		play('readPrd')
		fireEvent.click(header('Ticket'))
		expect(expanded('Ticket')).toBe(true)
		expect(expanded('Analysis')).toBe(false)
		expect(store.get(microruntimeSessionAtom).openedCard).toBe('ticket')

		play('readSpec')
		expect(expanded('Ticket')).toBe(true)

		play('planning')
		expect(expanded('Plan')).toBe(true)
		expect(store.get(microruntimeSessionAtom).openedCard).toBeNull()
	})

	it('reads the documents and judges the episodes line by line', () => {
		const { play } = mountSessionCards()
		expect(screen.queryAllByText('Read')).toHaveLength(0)

		play('readSpec')
		expect(screen.getAllByText('Read')).toHaveLength(3)
		expect(
			screen.getByText('5 components · 3 acceptance criteria'),
		).toBeTruthy()

		play('judgeGuardrail')
		expect(screen.getAllByText('Read')).toHaveLength(4)
		expect(
			screen.getByText(
				'FIN-02: estimate before a paid call · FIN-03: a budget before a paid run · FIN-04: size cpu and memory to the load',
			),
		).toBeTruthy()

		play('vectorSearch')
		expect(screen.getAllByText('Resolving')).toHaveLength(8)

		play('candidates')
		expect(screen.getByText('checkout-preview-sizing')).toBeTruthy()
		expect(screen.getByText('0.81')).toBeTruthy()
		expect(screen.queryByText('kept')).toBeNull()

		play('drop')
		expect(screen.getAllByText('kept')).toHaveLength(3)
		expect(screen.getByText('dropped · docs only')).toBeTruthy()

		play('suggest')
		expect(screen.getByText('claude-sonnet-5-5')).toBeTruthy()
		expect(
			screen.getByText('recommendation · the team authorizes'),
		).toBeTruthy()
	})

	it('carries each act’s latest figure in its collapsed header', () => {
		const { header, play } = mountSessionCards()
		expect(header('Ticket').textContent).toContain('0031 · book-a-demo-runtime')

		play('reads')
		expect(header('Analysis').textContent).toContain('4 documents read')

		play('drop')
		expect(header('Analysis').textContent).toContain('3 of 6 episodes kept')

		play('suggest')
		expect(header('Analysis').textContent).toContain(
			'Claude Sonnet 5.5 · effort medium',
		)

		play('go')
		expect(header('Plan').textContent).toContain('250m · 256 MiB')

		play('draft')
		expect(header('Plan').textContent).toContain(
			'1,450m · 2,304 MiB · pending team authorization',
		)
	})

	it('steps the plan’s figures up with the components and closes the report', () => {
		const { play } = mountSessionCards()
		play('bounds')
		expect(screen.getByText('2,000m · 3,072 MiB')).toBeTruthy()

		play('resolve')
		expect(screen.getAllByText('Resolving').length).toBeGreaterThanOrEqual(8)

		play('go')
		expect(screen.getAllByText('250m · 256 MiB').length).toBeGreaterThanOrEqual(
			1,
		)
		expect(screen.getByText('250m of 2,000m')).toBeTruthy()
		expect(screen.getByText('1 of 5 resolved · 250m · 256 MiB')).toBeTruthy()

		play('crossCheck')
		expect(screen.getByText('5 components · 1,450m · 2,304 MiB')).toBeTruthy()
		expect(
			screen.getByText(
				'ports 8080 · 8081 · 4222 · 6379 · 3000 · distinct · no collisions',
			),
		).toBeTruthy()

		play('size')
		expect(screen.getByText('550m · 768 MiB')).toBeTruthy()
		expect(screen.getByText('5 one-replica Deployments')).toBeTruthy()

		play('price')
		expect(
			screen.getByText('runtime: no price source · units only'),
		).toBeTruthy()

		play('validate')
		expect(screen.getByText('schema · 0 errors')).toBeTruthy()

		play('draft')
		expect(screen.getAllByText('pending team authorization')).toHaveLength(1)
		expect(
			screen.getByText(
				'runtime.md draft · revision 1 · proposed to the console',
			),
		).toBeTruthy()
	})

	it('opens the ticket again when the session replays', () => {
		const { expanded, header, play, store } = mountSessionCards()
		play('handoff')
		fireEvent.click(header('Analysis'))
		expect(expanded('Analysis')).toBe(true)

		act(() => {
			store.set(replaySessionAtom)
		})
		expect(store.get(microruntimeSessionAtom)).toEqual({
			...initialMicroruntimeSession,
			run: 1,
		})
		expect(expanded('Ticket')).toBe(true)
		expect(expanded('Analysis')).toBe(false)
	})
})

describe('Microruntimes Workflow card', () => {
	it('opens when the plan is handed over, with the suggestion priced against the budget', () => {
		const { expanded, header, play } = mountSessionCards()
		expect(workflowPlayButton().getAttribute('aria-disabled')).toBe('true')
		expect(screen.getByText('available once the plan is drafted')).toBeTruthy()

		play('handoff')
		expect(expanded('Workflow')).toBe(true)
		expect(expanded('Plan')).toBe(false)
		expect(header('Workflow').textContent).toContain(
			'Sonnet 5.5 · medium · $9.50 of $25.00',
		)
		expect(screen.getByText('$9.50')).toBeTruthy()
		expect(
			screen.getByText('reserves $11.40 of the $25.00 budget'),
		).toBeTruthy()
		expect(workflowPlayButton().getAttribute('aria-disabled')).toBeNull()
	})

	it('keeps planning checked for a person and says why', () => {
		const { play, store } = mountSessionCards()
		play('handoff')
		const planning = screen.getByRole('checkbox', {
			hidden: true,
			name: 'Planning · human in the loop',
		})
		expect(planning.getAttribute('aria-disabled')).toBe('true')

		fireEvent.click(planning)
		expect(store.get(workflowChoicesAtom).human.planning).toBe(true)

		const note = screen.getByRole('tooltip', { hidden: true })
		expect(note.textContent).toBe(
			'Planning always has a person in this workflow',
		)
		expect(planning.getAttribute('aria-describedby')).toBe(note.id)
	})

	it('reprices the run as the visitor chooses and warns over the budget', () => {
		const { header, play } = mountSessionCards()
		play('handoff')
		const [model, effort] = workflowSelects()

		fireEvent.change(model, { target: { value: 'claude-haiku-5-5' } })
		expect(header('Workflow').textContent).toContain(
			'Haiku 5.5 · medium · $0.48',
		)

		fireEvent.change(model, { target: { value: 'claude-opus-5-5' } })
		fireEvent.change(effort, { target: { value: 'high' } })
		expect(screen.getByText('$34.00')).toBeTruthy()
		expect(
			screen.getByText('over the budget · the gate will pause this run'),
		).toBeTruthy()
	})

	it('launches with a snapshot, locks the choices and plays to a reviewed pull request', () => {
		const { header, play, store } = mountSessionCards()
		play('handoff')

		fireEvent.click(workflowPlayButton())
		expect(store.get(microruntimeSessionAtom).stage).toBe('launching')
		expect(workflowPlayButton().textContent).toContain('Starting…')

		fireEvent.click(
			screen.getByRole('checkbox', {
				hidden: true,
				name: 'Spec · human in the loop',
			}),
		)
		expect(store.get(workflowChoicesAtom).human.spec).toBe(false)

		act(() => {
			store.set(startRunAtom)
		})
		const session = store.get(microruntimeSessionAtom)
		expect(session).toMatchObject({
			playedCount: 0,
			ran: store.get(workflowChoicesAtom),
			run: 1,
			stage: 'running',
		})
		expect(screen.getByRole('status', { hidden: true }).textContent).toBe('')

		const script = store.get(runScriptAtom)
		playRunTo(store, script.findIndex((line) => line.id === 'gateOpened') + 1)
		expect(screen.getByText('$11.40')).toBeTruthy()

		playRunTo(
			store,
			script.findIndex((line) => line.message === 'devTicketWaiting') + 1,
		)
		expect(screen.getByText('the person writes the reserved part')).toBeTruthy()

		playRunTo(store, script.findIndex((line) => line.id === 'codeReview') + 1)
		expect(screen.queryByText('waiting for the owner')).toBeNull()

		playRunTo(store, script.findIndex((line) => line.id === 'completed') + 1)
		expect(screen.getByText('$9.80')).toBeTruthy()
		expect(screen.queryByText('waiting for the owner')).toBeNull()

		playRunTo(store, script.length)
		expect(screen.getByText('waiting for the owner')).toBeTruthy()
		expect(screen.getByText('driven by the person')).toBeTruthy()
		expect(screen.getByText('met')).toBeTruthy()
		const pullRequest = screen.getByText('#19 ·3f9c2e1')
		expect(pullRequest.querySelector('svg.text-status-success')).toBeTruthy()
		expect(header('Workflow').textContent).toContain(
			'#19 · ready for review · $9.80',
		)
		expect(screen.getByRole('status', { hidden: true }).textContent).toBe(
			'The simulated run ended: pull request #19 is ready for review, $9.80 observed of $11.40 reserved.',
		)
		expect(workflowPlayButton().textContent).toContain('Play again')
	})

	it('pauses at the gate when the reservation breaks the budget', () => {
		const { header, play, store } = mountSessionCards()
		play('handoff')
		const [model, effort] = workflowSelects()
		fireEvent.change(model, { target: { value: 'claude-opus-5-5' } })
		fireEvent.change(effort, { target: { value: 'high' } })

		const script = launchWorkflowRun(store)
		expect(script.some((line) => line.id === 'pullRequest')).toBe(false)

		playRunTo(store, script.length)
		expect(
			screen.getByText('paused at the gate · needs a decision'),
		).toBeTruthy()
		expect(header('Workflow').textContent).toContain('paused at the gate')
		expect(screen.getByRole('status', { hidden: true }).textContent).toBe(
			'The simulated run paused at the gate: the reservation of $40.80 exceeds the $25.00 budget, and nothing was paid for.',
		)
	})

	it('ends ready to merge when the review is left out of the loop, and keeps the choices on replay', () => {
		const { header, play, store } = mountSessionCards()
		play('handoff')
		fireEvent.click(
			screen.getByRole('checkbox', {
				hidden: true,
				name: 'Review · human in the loop',
			}),
		)
		expect(
			screen.queryByRole('checkbox', {
				hidden: true,
				name: 'Merge · human in the loop',
			}),
		).toBeNull()

		const script = launchWorkflowRun(store)
		playRunTo(store, script.length)
		expect(header('Workflow').textContent).toContain('#19 · ready to merge')

		act(() => {
			store.set(replaySessionAtom)
		})
		expect(store.get(microruntimeSessionAtom).stage).toBe('plan')
		expect(store.get(workflowChoicesAtom).human.review).toBe(false)
	})
})
