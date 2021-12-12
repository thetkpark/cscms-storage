import { Fragment } from 'react'
import SocialAuthButton from './SocialAuthButton'
import styles from '../../styles/auth/AuthForm.module.css'
const AuthForm = ({ mode, changeMode }) => {
	const renderTitle = () => {
		switch (mode) {
			case 'login':
				return 'Login'
			case 'signup':
				return 'Create Account'
			default:
				return ''
		}
	}
	const renderSwitchText = () => {
		switch (mode) {
			case 'login':
				return (
					<Fragment>
						Don't have an account?{' '}
						<span role="link" onClick={() => changeMode('signup')}>
							Signup Now.
						</span>
					</Fragment>
				)
			case 'signup':
				return (
					<Fragment>
						Already have an account?{' '}
						<span role="link" onClick={() => changeMode('login')}>
							Login.
						</span>
					</Fragment>
				)
			default:
				return null
		}
	}
	return (
		<Fragment>
			<div className={styles.AuthFormWrapper}>
				<h2 className={styles.TitleText}>{renderTitle()}</h2>
				<SocialAuthButton mode={mode} platform="Google" />
				<SocialAuthButton mode={mode} platform="Github" />
				<div className={styles.Description}>
					By creating an account, you agree to cscms{' '}
					<span role="link">Terms of Use, Privacy Policy</span> and to receive news and
					updates.
				</div>
			</div>
			<div className={styles.SwitchField}
			>
				{renderSwitchText()}
			</div>
		</Fragment>
	)
}

export default AuthForm
