import Icon from '../util/Icon'
import styles from '../../styles/auth/SocialAuth.module.css'
import { toTitleCase } from '../../utils/formatText'
import { useRecoilState } from 'recoil'
import { authState } from '../../store/auth'
const SocialAuthButton = ({ mode, platform }) => {
    const [auth,setAuth] = useRecoilState(authState)
	const handleAction = (mode, platform) => {
		switch (mode) {
			case 'login':
				setAuth(true)
                break
			case 'signup':
				setAuth(true)
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
