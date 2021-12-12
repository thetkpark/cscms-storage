import { Fragment } from 'react'
import SocialAuthButton from './SocialAuthButton'

const AuthForm = ({ mode }) => {
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
	return (
		<Fragment>
			<div style={{ padding: '3rem 8rem', width: '360px' }}>
				<h2 style={{ textAlign: 'center', marginBottom: '64px' }}>{renderTitle()}</h2>
				<SocialAuthButton mode={mode} platform="Google" />
				<SocialAuthButton mode={mode} platform="Github" />
				<div style={{ textAlign: 'center', marginTop: '64px' }}>
					By creating an account, you agree to cscms Terms of Use, Privacy Policy and to
					receive news and updates.
				</div>
			</div>
			<div
				style={{
					textAlign: 'center',
					padding: '3rem',
					background: '#ccc'
				}}
			>
				Don't have an account? <span>Sign Up</span>
			</div>
		</Fragment>
	)
}

export default AuthForm
