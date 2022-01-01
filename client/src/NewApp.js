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
import axios from 'axios'
import Swal from 'sweetalert2'
import FileList from './components/file/FileList'
import UserProfile from './components/auth/UserProfile'
function App() {
	const [route, setRoute] = useState('file')
	const [auth, setAuth] = useRecoilState(authState)
	const [dialog, setDialog] = useState(null)
	const [progress, setProgress] = useState(0)
	const [error, setError] = useState(null)
	useEffect(() => {
		ReactGA.initialize('G-S7NPY62JTS')
		ReactGA.pageview(window.location.pathname)
	}, [])
	useEffect(() => {
		axios
			.get('https://storage.cscms.me/auth/user')
			.then(res => {
				setAuth({
					isAuthenticated: true,
					user: {
						name: res.data.username,
						image: res.data.avatar_url,
						email: res.data.email
					}
				})
			})
			.catch(err => {
				setAuth({
					isAuthenticated: false,
					user: null
				})
			})
	}, [])
	useEffect(() => {
		if (auth.isAuthenticated) {
			setDialog(null)
		}
	}, [auth])

	useEffect(() => {
		if (error != null) {
			Swal.fire({
				title: 'An error occured',
				text: error,
				icon: 'error',
				confirmButtonText: 'Ok'
			})
			setError(null)
		}
	}, [error])

	const handleAction = action => {
		switch (action) {
			case 'signup':
				setDialog('signup')
				break
			case 'login':
				setDialog('login')
				break
			case 'logout':
				axios
					.get('https://storage.cscms.me/auth/logout')
					.then(() => {
						setAuth(false)
					})
					.catch(err => {
						console.log(err)
					})
				break
			default:
				break
		}
	}
	const handleChangeRoute = newRoute => {
		if (newRoute === route) return
		setRoute(newRoute)
	}
	const handleUpload = data => {
		if (route === 'file') {
			handleUploadFile(data)
		} else if (route === 'image') {
			handleUploadImage(data)
		}
	}
	const handleUploadFile = async ({ selectedFile, slug, duration }) => {
		const formdata = new FormData()
		formdata.append('file', selectedFile)

		try {
			// setShowModal(true)
			const res = await axios.post('https://storage.cscms.me/api/file', formdata, {
				onUploadProgress: progressEvent => {
					const uploadPercent = Math.round(
						(progressEvent.loaded / progressEvent.total) * 100
					)
					setProgress(uploadPercent)
				},
				params: { slug, duration }
			})
			// setFileData(res.data)

			ReactGA.event({
				category: 'file',
				action: 'Upload file',
				value: selectedFile.size
			})
		} catch (err) {
			setError(err.response.data.message)
		}
	}
	const handleUploadImage = async ({ selectedFile }) => {
		const formdata = new FormData()
		formdata.append('image', selectedFile)

		try {
			// setShowModal(true)
			const res = await axios.post('https://storage.cscms.me/api/image', formdata, {
				onUploadProgress: progressEvent => {
					const uploadPercent = Math.round(
						(progressEvent.loaded / progressEvent.total) * 100
					)
					setProgress(uploadPercent)
				}
			})
			// setFileData(res.data)

			ReactGA.event({
				category: 'file',
				action: 'Upload file',
				value: selectedFile.size
			})
		} catch (err) {
			setError(err.response.data.message)
		}
	}
	const renderScreen = () => {
		switch (route) {
			case 'file':
			case 'image':
				return (
					<UploadContainer
						type={route}
						handleUpload={handleUpload}
						setError={setError}
						progress={progress}
					/>
				)
			case 'myfile':
				if (auth.isAuthenticated) return <FileList />
				setRoute('file')
				break
			default:
				setRoute('file')
				break
		}
	}
	return (
		<Fragment>
			<div className={styles.App}>
				<div className={styles.Wrapper}>
					<Navbar auth={auth.isAuthenticated} handleAction={handleAction} />
					<div style={{ padding: '2rem 8rem', display: 'flex', flexDirection: 'column' }}>
						{auth.isAuthenticated ? (
							<UserProfile user={auth.user} handleChangeRoute={handleChangeRoute} />
						) : null}

						{renderScreen()}
					</div>
					<Sidebar currentRoute={route} handleChangeRoute={handleChangeRoute} />
					{!auth.isAuthenticated && dialog ? (
						<Dialog open={dialog !== null} onClose={() => setDialog(null)}>
							<AuthForm mode={dialog} changeMode={mode => setDialog(mode)} />
						</Dialog>
					) : null}
				</div>
			</div>
		</Fragment>
	)
}

export default App
