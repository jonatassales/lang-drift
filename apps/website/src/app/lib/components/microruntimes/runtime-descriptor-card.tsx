import { RuntimeFileEditor } from '@components/microruntimes/runtime-file-editor'
import {
	descriptorAnnotationIds,
	formatDescriptor,
	type MicroruntimesMessages,
	runtimeDescriptor,
	runtimeFilePath,
} from '@data/microruntimes-content'
import type { WebsiteLocale } from '@i18n/routing'

interface RuntimeDescriptorCardProps {
	copy: MicroruntimesMessages['descriptor']
	locale: WebsiteLocale
}

// The runtime file's langdrift JSON fence in the editor, beside numbered notes
// on what each part decides.
export function RuntimeDescriptorCard(props: RuntimeDescriptorCardProps) {
	const { copy, locale } = props

	return (
		<div className="mt-11 grid items-start gap-4 lg:grid-cols-[minmax(0,1.3fr)_minmax(0,0.7fr)]">
			<RuntimeFileEditor
				locale={locale}
				path={runtimeFilePath}
				value={formatDescriptor(runtimeDescriptor)}
			/>
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
