import { Fragment } from 'react'
import SocialAuthButton from './SocialAuthButton'

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
			<div style={{ padding: '3rem 8rem', width: '360px' }}>
				<h2 style={{ textAlign: 'center', marginBottom: '64px', fontSize: '2.25rem' }}>
					{renderTitle()}
				</h2>
				<SocialAuthButton mode={mode} platform="Google" />
				<SocialAuthButton mode={mode} platform="Github" />
				<div style={{ textAlign: 'center', marginTop: '64px' }}>
					By creating an account, you agree to cscms{' '}
					<span role="link">Terms of Use, Privacy Policy</span> and to receive news and
					updates.
				</div>
			</div>
			<div
				style={{
					textAlign: 'center',
					padding: '3rem',
					background: '#EBEDF4',
					color: '#959595'
				}}
			>
				{renderSwitchText()}
			</div>
		</Fragment>
	)
}

export default AuthForm
