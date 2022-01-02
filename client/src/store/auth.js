import { atom } from 'recoil'
export const authState = atom({
	key: 'auth',
	default: { isAuthenticated: false, user: null }
})
