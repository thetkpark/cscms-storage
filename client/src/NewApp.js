import { useState, useEffect } from 'react'
import ReactGA from 'react-ga'
import Navbar from './components/Navbar'
import Sidebar from './components/Sidebar'
import styles from './styles/NewApp.module.css'
function App() {
	const [route, setRoute] = useState('file')

	useEffect(() => {
		ReactGA.initialize('G-S7NPY62JTS')
		ReactGA.pageview(window.location.pathname)
	}, [])

	return (
		<div className={styles.App}>
			<div className={styles.Wrapper}>
				<Navbar />
				<Sidebar currentRoute={route} handleChangeRoute={route => setRoute(route)} />.
			</div>
		</div>
	)
}

export default App
