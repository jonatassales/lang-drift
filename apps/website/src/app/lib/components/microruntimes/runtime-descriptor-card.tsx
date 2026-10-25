import { Eyebrow } from '@langdrift/web/ui/eyebrow'

import {
	descriptorAnnotationIds,
	formatDescriptor,
	type MicroruntimesMessages,
	runtimeDescriptor,
	runtimeFilePath,
} from '@data/microruntimes-content'

interface RuntimeDescriptorCardProps {
	copy: MicroruntimesMessages['descriptor']
}

// The runtime file as a reader would open it: one langdrift JSON fence,
// beside numbered notes on what each part decides.
export function RuntimeDescriptorCard(props: RuntimeDescriptorCardProps) {
	const { copy } = props

	return (
		<div className="mt-11 grid gap-4 lg:grid-cols-[minmax(0,1.3fr)_minmax(0,0.7fr)]">
			<figure className="m-0 min-w-0 overflow-hidden rounded-2xl border border-border bg-card">
				<figcaption className="flex items-center justify-between gap-4 border-b border-border px-5 py-3">
					<code className="truncate font-mono text-[12px] text-muted-foreground">
						{runtimeFilePath}
					</code>
					<Eyebrow className="shrink-0 rounded-full border border-border px-2 py-[3px]">
						{copy.illustrative}
					</Eyebrow>
				</figcaption>
				<pre className="m-0 max-h-[560px] overflow-auto p-5 font-mono text-[12px] leading-[1.65] text-foreground">
					<code>
						<span className="text-muted-foreground">```langdrift</span>
						{`\n${formatDescriptor(runtimeDescriptor)}\n`}
						<span className="text-muted-foreground">```</span>
					</code>
				</pre>
			</figure>
			<ol className="m-0 grid list-none content-start gap-3 p-0">
				{descriptorAnnotationIds.map((id, index) => (
					<li
						className="grid grid-cols-[auto_minmax(0,1fr)] gap-x-3 gap-y-1 rounded-2xl border border-border bg-card p-5"
						key={id}
					>
						<span className="row-span-2 inline-flex size-6 items-center justify-center rounded-full border border-border font-mono text-[10px] text-muted-foreground">
							{index + 1}
						</span>
						<h3 className="m-0 text-[15px] font-medium tracking-[-0.01em]">
							{copy.annotations[id].title}
						</h3>
						<p className="m-0 text-[13px] leading-[1.55] text-muted-foreground">
							{copy.annotations[id].description}
						</p>
					</li>
				))}
			</ol>
		</div>
	)
}
