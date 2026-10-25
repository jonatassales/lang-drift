import { LogoMark } from '@langdrift/web/ui/brand'

export function DocsFooter() {
	return (
		<footer className="flex w-full items-center justify-center gap-2 self-end text-center text-xs text-muted-foreground">
			<LogoMark size={24} />
			<span>© 2026 LangDrift</span>
		</footer>
	)
}
