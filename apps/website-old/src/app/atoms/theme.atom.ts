import {
	atomWithStorage,
	createJSONStorage,
	unstable_withStorageValidator as withStorageValidator,
} from 'jotai/utils'

export type WebsiteTheme = 'dark' | 'light'

export const themeAtom = atomWithStorage<WebsiteTheme>(
	'langdrift.website.theme',
	'dark',
	withStorageValidator(
		(value): value is WebsiteTheme => value === 'dark' || value === 'light',
	)(createJSONStorage()),
)
