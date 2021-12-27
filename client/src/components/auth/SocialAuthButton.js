import Icon from '../util/Icon'
import styles from '../../styles/auth/SocialAuth.module.css'
import { toTitleCase } from '../../utils/formatText'
import { useSetRecoilState } from 'recoil'
import { authState } from '../../store/auth'
const SocialAuthButton = ({ mode, platform }) => {
	const setAuth = useSetRecoilState(authState)
	const handleAction = (mode, platform) => {
		switch (mode) {
			case 'login':
				setAuth({
					isAuthenticated: true,
					user: {
						name: 'Wagyu',
						image: 'https://tmp.cscms.me/9nwmgr',
						email: 'wagyu@wagyu.com'
					}
				})
				break
			case 'signup':
				setAuth({
					isAuthenticated: true,
					user: {
						name: 'Momo',
						image: 'https://tmp.cscms.me/9nwmgr',
						email: 'momo@momo.com'
					}
				})
				break
			default:
				return null
		}
	}
	return (
		<div
			className={`${styles.SocialBtn} ${styles[platform]}`}
			onClick={() => handleAction(mode, platform)}
		>
			<Icon name={platform} />
			<span>
				{toTitleCase(mode)} with {platform}
			</span>
		</div>
	)
}
export default SocialAuthButton
