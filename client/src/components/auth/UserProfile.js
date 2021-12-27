import Icon from '../util/Icon'
import styles from '../../styles/auth/UserProfile.module.css'
const UserProfile = ({ user, handleChangeRoute }) => {
	return (
		<div className={styles.Box}>
			<div className={styles.ImageBox}>
				<img
                className={styles.Image}
					src={user.image}
					alt="user-profile"
					width={100}
					height={100}
				/>
			</div>
			<div className={styles.DetailBox}>
				<div className={styles.Name}>Hey {user.name}!</div>
				<div className={styles.MyFiles} onClick={() => handleChangeRoute('myfile')}
				>
					<Icon name="folder" />
					<span>My Files</span>
				</div>
			</div>
		</div>
	)
}
export default UserProfile
