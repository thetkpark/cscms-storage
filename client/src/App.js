import { useState, useEffect } from 'react'
import axios from 'axios'
import FormData from 'form-data'
import ReactGA from 'react-ga'
import { TextField, Button, Text, Slider } from '@adobe/react-spectrum'
import UploadIcon from '@spectrum-icons/workflow/UploadToCloudOutline'
import Dropzone from './components/Dropzone'
import styles from './styles/App.module.css'
import FileDataModal from './components/Modal'
import ErrorDialog from './components/ErrorDialog'

function App() {
	const [progress, setProgress] = useState(-1)
	const [selectedFile, setSelectedFile] = useState(null)
	const [slug, setSlug] = useState('')
	const [fileData, setFileData] = useState(undefined)
	const [error, setError] = useState('')
	const [showModal, setShowModal] = useState(false)
	const [storeDuration, setStoreDuration] = useState(7)

	useEffect(() => {
		ReactGA.initialize('G-S7NPY62JTS')
		ReactGA.pageview(window.location.pathname)
	}, [])

	const onDrop = (acceptedFiles, rejectedFiles) => {
		if (acceptedFiles.length === 1) {
			setError('')
			setSelectedFile(acceptedFiles[0])
		} else {
			if (rejectedFiles[0].errors[0].code === 'too-many-files') {
				setError('Too many files. You can only upload one file at a time')
			} else if (rejectedFiles[0].errors[0].code === 'file-too-large') {
				setError('File too big. The size limit is 100MB')
			} else setError('File not accepted')
		}
	}

	const handleSubmission = async event => {
		event.preventDefault()
		if (!selectedFile) {
			setError('Please select a file')
			return
		}
		const formdata = new FormData()
		formdata.append('file', selectedFile)

		try {
			const res = await axios.post('/api/file', formdata, {
				onUploadProgress: progressEvent => {
					const uploadPercent = Math.round(
						(progressEvent.loaded / progressEvent.total) * 100
					)
					setProgress(uploadPercent)
				},
				params: { slug, duration: storeDuration }
			})

			setFileData(res.data)
			setShowModal(true)
			ReactGA.event({
				category: 'file',
				action: 'Upload file',
				value: selectedFile.size
			})
		} catch (err) {
			console.log(err)
			setError(err.response.data.message)
		}
	}

	const closeDialogAndReset = () => {
		setShowModal(false)
		setProgress(-1)
		setSelectedFile(null)
		setSlug('')
		setFileData(undefined)
		setError('')
	}

	return (
		<div className={styles.App}>
			<div className={styles.AppContainer}>
				{/* <h1 className={styles.Heading}>CSCMS Temp Storage</h1> */}
				<Dropzone
					onDrop={onDrop}
					selectedFilename={selectedFile ? selectedFile.name : ''}
					progress={progress}
					error={error}
				/>
				<div className={styles.FormContainer}>
					<form className={styles.Form} onSubmit={handleSubmission}>
						<TextField
							label="Custom slug for accessing the file (Optional)"
							placeholder="Slug"
							value={slug}
							onChange={e => setSlug(e)}
							width="300px"
							type="text"
							inputMode="text"
						/>
						<Slider
							label="Store Duration (Days)"
							value={storeDuration}
							onChange={setStoreDuration}
							minValue={1}
							maxValue={30}
						/>
						<Button variant="primary" type="submit">
							<UploadIcon />
							<Text>Upload</Text>
						</Button>
					</form>
				</div>
				<FileDataModal
					show={showModal}
					closeDialog={closeDialogAndReset}
					fileData={fileData}
				/>
				{error && <ErrorDialog errorMessage={error} closeDialog={() => setError('')} />}
			</div>
		</div>
	)
}

export default App
