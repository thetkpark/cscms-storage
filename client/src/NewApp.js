import { useState, useEffect } from 'react'
import ReactGA from 'react-ga'
import Navbar from './components/Navbar'
import Sidebar from './components/Sidebar'
import styles from './styles/NewApp.module.css'
function App() {
	const [route, setRoute] = useState('file')
	const [auth, setAuth] = useState(false)

	useEffect(() => {
		ReactGA.initialize('G-S7NPY62JTS')
		ReactGA.pageview(window.location.pathname)
	}, [])
	
	const handleAction = (action) => {
		switch (action) {
			case 'register':
				setAuth(true)
				break
			case 'login':
				setAuth(true)
				break
			case 'logout':
				setAuth(false)
				break
		}
	}

	return (
		<div className={styles.App}>
			<div className={styles.Wrapper}>
				<Navbar auth={auth} handleAction={handleAction} />
				<Sidebar currentRoute={route} handleChangeRoute={route => setRoute(route)} />.
			</div>
		</div>
	)
}

export default App
