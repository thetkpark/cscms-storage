import Icon from '../util/Icon'
import styles from '../../styles/auth/SocialAuth.module.css'
import { toTitleCase } from '../../utils/formatText'
const SocialAuthButton = ({ mode, platform }) => {
	const handleAction = (mode, platform) => {
		switch (mode) {
			case 'login':
                break
			case 'signup':
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
