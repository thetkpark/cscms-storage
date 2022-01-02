import { atom } from 'recoil'
export const authState = atom({
	key: 'auth',
	default: { isAuthenticated: true, user: {name:'',image:''} }
})
