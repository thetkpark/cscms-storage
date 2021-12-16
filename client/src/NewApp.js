import { useState, useEffect, Fragment } from 'react'
import ReactGA from 'react-ga'
import Navbar from './components/layout/Navbar'
import Sidebar from './components/layout/Sidebar'
import AuthForm from './components/auth/AuthForm'
import styles from './styles/NewApp.module.css'
import { Dialog } from '@material-ui/core'
import UploadContainer from './components/upload/UploadContainer'
import { useRecoilState } from 'recoil'
import { authState } from './store/auth'
function App() {
	const [route, setRoute] = useState('file')
	const [auth, setAuth] = useRecoilState(authState)
	const [dialog, setDialog] = useState(null)
	useEffect(() => {
		ReactGA.initialize('G-S7NPY62JTS')
		ReactGA.pageview(window.location.pathname)
	}, [])
	useEffect(() => {
		if (auth) {
			setDialog(null)
		}
	}, [auth])

	const handleAction = action => {
		switch (action) {
			case 'signup':
				setDialog('signup')
				break
			case 'login':
				setDialog('login')
				break
			case 'logout':
				setAuth(false)
				break
			default:
				break
		}
	}
	const handleChangeRoute = newRoute => {
		if (newRoute === route) return
		setRoute(newRoute)
	}
	const renderScreen = () => {
		switch (route) {
			case 'file':
			case 'image':
				return <UploadContainer type={route} />
			case 'myfile':
				if (auth) return <Fragment></Fragment>
				setRoute('file')
				break
			default:
				setRoute('file')
				break
		}
	}
	return (
		<div className={styles.App}>
			<div className={styles.Wrapper}>
				<Navbar auth={auth} handleAction={handleAction} />
				<div style={{ padding: '2rem 8rem', display: 'flex', flexDirection: 'column' }}>
					{auth ? (
						<Fragment>
							<div>Hey Wagyu!</div>
						</Fragment>
					) : null}

					{renderScreen()}
				</div>
				<Sidebar currentRoute={route} handleChangeRoute={handleChangeRoute} />
				{!auth && dialog ? (
					<Dialog open={dialog !== null} onClose={() => setDialog(null)}>
						<AuthForm mode={dialog} changeMode={mode => setDialog(mode)} />
					</Dialog>
				) : null}
			</div>
		</div>
	)
}

export default App
