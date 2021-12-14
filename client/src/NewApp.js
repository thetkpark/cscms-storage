import { useState, useEffect } from 'react'
import ReactGA from 'react-ga'
import Navbar from './components/layout/Navbar'
import Sidebar from './components/layout/Sidebar'
import AuthForm from './components/auth/AuthForm'
import styles from './styles/NewApp.module.css'
import { Dialog } from '@material-ui/core'
import DropZone from './components/upload/Dropzone'
function App() {
	const [route, setRoute] = useState('file')
	const [auth, setAuth] = useState(false)
	const [dialog, setDialog] = useState(null)
	useEffect(() => {
		ReactGA.initialize('G-S7NPY62JTS')
		ReactGA.pageview(window.location.pathname)
	}, [])

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

	return (
		<div className={styles.App}>
			<div className={styles.Wrapper}>
				<Navbar auth={auth} handleAction={handleAction} />
				<div style={{ padding: '2rem 8rem', display: 'flex', flexDirection: 'column' }}>
					<div>Hey Wagyu!</div>
					<div
						style={{
							background: 'white',
							width: '65vw',
							height: '60vh',
							margin: '1rem auto',
							borderRadius: '50px',
							padding:'3rem'
						}}
					>

						<DropZone type={route}/>
					</div>
				</div>
				<Sidebar currentRoute={route} handleChangeRoute={route => setRoute(route)} />
				{dialog ? (
					<Dialog open={dialog !== null} onClose={() => setDialog(null)}>
						<AuthForm mode={dialog} changeMode={mode => setDialog(mode)} />
					</Dialog>
				) : null}
			</div>
		</div>
	)
}

export default App
