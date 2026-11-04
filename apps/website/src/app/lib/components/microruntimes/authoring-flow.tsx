import { Fragment } from 'react'

import {
	flowGuaranteeIds,
	flowTargetIds,
	type MicroruntimesMessages,
} from '@data/microruntimes-content'
import { LogoMark } from 'web/ui/brand'
import { Eyebrow } from 'web/ui/eyebrow'
import { ArrowRightIcon } from 'web/ui/icons/arrow'
import { MicroruntimesIcon } from 'web/ui/icons/platform'
import { InfoCards } from 'web/ui/info-cards'

interface AuthoringFlowProps {
	copy: MicroruntimesMessages['flow']
}

const flowCard =
	'flex flex-col gap-2 rounded-2xl border border-border bg-card p-6'
const flowArrow = (
	<li
		aria-hidden="true"
		className="flex items-center justify-center text-muted-foreground"
	>
		<ArrowRightIcon className="rotate-90 lg:rotate-0" size={20} />
	</li>
)

// Console → connected agent → local execution, then the three rules that
// keep a request from becoming an unbounded action.
export function AuthoringFlow(props: AuthoringFlowProps) {
	const { copy } = props

	const steps = [
		{
			...copy.console,
			icon: <LogoMark size={20} />,
			id: 'console',
		},
		{
			...copy.agent,
			icon: <MicroruntimesIcon className="size-5" />,
			id: 'agent',
		},
	]

	return (
		<>
			<ol className="m-0 mt-11 grid list-none gap-3 p-0 lg:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)_auto_minmax(0,1.2fr)]">
				{steps.map((step) => (
					<Fragment key={step.id}>
						<li className={flowCard}>
							<span className="mb-2 inline-flex size-9 items-center justify-center rounded-lg bg-foreground/5 text-foreground">
								{step.icon}
							</span>
							<h3 className="m-0 text-lg font-medium tracking-[-0.02em]">
								{step.title}
							</h3>
							<p className="m-0 text-[13px] leading-[1.55] text-muted-foreground">
								{step.description}
							</p>
						</li>
						{flowArrow}
					</Fragment>
				))}
				<li className={flowCard}>
					<Eyebrow>{copy.targetsTitle}</Eyebrow>
					<ul className="m-0 mt-1 grid list-none gap-3 p-0">
						{flowTargetIds.map((id) => (
							<li
								className="grid gap-0.5 border-t border-border pt-3 first:border-t-0 first:pt-0"
								key={id}
							>
								<strong className="font-mono text-[14px] font-medium">
									{copy.targets[id].title}
								</strong>
								<span className="text-[13px] text-muted-foreground">
									{copy.targets[id].description}
								</span>
							</li>
						))}
					</ul>
				</li>
			</ol>
			<InfoCards
				className="mt-3"
				items={flowGuaranteeIds.map((id) => ({
					description: copy.guarantees[id].description,
					id,
					title: copy.guarantees[id].title,
				}))}
			/>
		</>
	)
}
