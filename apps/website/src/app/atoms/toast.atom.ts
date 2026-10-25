import { atom } from 'jotai'

export interface WebsiteToast {
	message: string
	open: boolean
	tone: 'info' | 'success' | 'warning'
}

export const toastAtom = atom<WebsiteToast>({
	message: '',
	open: false,
	tone: 'info',
})
