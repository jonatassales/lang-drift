import type { FakeDialog } from '@layout/website-header/testdata/fake-dialog.mock'
import { useNavigationDrawer } from '@layout/website-header/use-navigation-drawer.hook'

// The drawer's trigger and close button around the hook, with the fake dialog
// attached where the header attaches its dialog element.
interface NavigationDrawerHarnessProps {
	dialog: FakeDialog
}

export function NavigationDrawerHarness(props: NavigationDrawerHarnessProps) {
	const { dialog } = props

	const drawer = useNavigationDrawer()

	function attachDialog(node: HTMLSpanElement | null) {
		drawer.drawerRef.current = node
			? (dialog as unknown as HTMLDialogElement)
			: null
	}

	return (
		<>
			<button onClick={drawer.openDrawer} ref={drawer.triggerRef} type="button">
				Open navigation
			</button>
			<button onClick={drawer.closeDrawer} ref={drawer.closeRef} type="button">
				Close navigation
			</button>
			<span ref={attachDialog} />
			<output>{drawer.open ? 'open' : 'closed'}</output>
		</>
	)
}
