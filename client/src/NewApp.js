import { useState, useEffect } from 'react'
import ReactGA from 'react-ga'
import Navbar from './components/layout/Navbar'
import Sidebar from './components/layout/Sidebar'
import AuthForm from './components/auth/AuthForm'
import styles from './styles/NewApp.module.css'
import { Dialog } from '@material-ui/core'
import DropZone from './components/upload/Dropzone'
import Button from './components/util/Button'
import Icon from './components/util/Icon'
function App() {
	const [route, setRoute] = useState('file')
	const [auth, setAuth] = useState(true)
	const [dialog, setDialog] = useState(null)
	const [selectedFile, setSelectedFile] = useState(null)
	const [error, setError] = useState('')
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
	const handleChangeRoute = newRoute => {
		if (newRoute === route) return
		setRoute(newRoute)
		setSelectedFile(null)
		setError('')
	}
	const onDrop = (acceptedFiles, rejectedFiles) => {
		if (acceptedFiles.length === 1) {
			setError('')
			setSelectedFile(acceptedFiles[0])
			console.log(acceptedFiles[0])
		} else {
			if (rejectedFiles[0].errors[0].code === 'too-many-files') {
				setError('Too many files. You can only upload one file at a time')
			} else if (rejectedFiles[0].errors[0].code === 'file-too-large') {
				setError('File too big. The size limit is 100MB')
			} else setError('File not accepted')
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
							padding: '3rem',
							display: 'flex',
							flexDirection: 'column',
							alignItems: 'center'
						}}
					>
						<DropZone
							type={route}
							selectedFilename={selectedFile ? selectedFile.name : ''}
							onDrop={onDrop}
						/>
						{selectedFile ? (
							<div style={{ margin: '1rem' }}>{selectedFile.name}</div>
						) : null}
						<Button
							bgColor={'#E9EEFF'}
							style={{
								border: 'none',
								fontSize: '.9rem',
								width: '170px',
								height: '50px',
								marginTop: '1rem'
							}}
						>
							<Icon name="upload" role="icon" /> Upload
						</Button>
					</div>
				</div>
				<Sidebar currentRoute={route} handleChangeRoute={handleChangeRoute} />
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
