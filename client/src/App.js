import { useState, useEffect } from 'react'
import axios from 'axios'
import FormData from 'form-data'
import ReactGA from 'react-ga'
import { Form, TextField, Button, Text, Flex } from '@adobe/react-spectrum'
import UploadIcon from '@spectrum-icons/workflow/UploadToCloudOutline'
import Dropzone from './Dropzone'
import styles from './App.module.css'
import FileDataModal from './Modal'

function App() {
	const [progress, setProgress] = useState(-1)
	const [selectedFile, setSelectedFile] = useState()
	const [slug, setSlug] = useState('')
	const [fileData, setFileData] = useState(undefined)
	const [error, setError] = useState('')
	const [selectedFilename, setSelectedFilename] = useState('')
	const [showModal, setShowModal] = useState(false)

	useEffect(() => {
		ReactGA.initialize('G-S7NPY62JTS')
		ReactGA.pageview(window.location.pathname)
	}, [])

	const onDrop = (acceptedFiles, rejectedFiles) => {
		if (acceptedFiles.length === 1) {
			setError('')
			setSelectedFile(acceptedFiles[0])
			setSelectedFilename(acceptedFiles[0].name)
			
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
		const formdata = new FormData()
		formdata.append('file', selectedFile)

		const res = await axios.post('/api/file', formdata, {
			onUploadProgress: progressEvent => {
				const uploadPercent = Math.round(
					(progressEvent.loaded / progressEvent.total) * 100
				)
				setProgress(uploadPercent)
			},
			params: { slug }
		})
		setFileData(res.data)
		setShowModal(true)
		ReactGA.event({
			category: 'file',
			action: 'Upload file',
			value: selectedFile.size
		})
	}

	const closeDialogAndReset = () => {
		setShowModal(false)
		setProgress(-1)
		setSelectedFile()
		setSlug('')
		setFileData(undefined)
		setError('')
		setSelectedFilename('')
	}

	return (
		<div className={styles.App}>
			<div className={styles.AppContainer}>
				{/* <h1 className={styles.Heading}>CSCMS Temp Storage</h1> */}
				<Dropzone
					onDrop={onDrop}
					selectedFilename={selectedFilename}
					progress={progress}
					error={error}
				/>
				<div className={styles.FormContainer}>
					<Form className={styles.Form} onSubmit={handleSubmission}>
						<Flex direction="row" gap="size-300" alignItems="end" justifyContent="center">
							<TextField
								label="Custom slug for accessing the file (Optional)"
								placeholder="Slug"
								value={slug}
								onChange={e => setSlug(e)}
								width="300px"
								type="text"
								inputMode="text"
							/>
							<Button variant="primary" width="100px" type="submit">
								<UploadIcon />
								<Text>Upload</Text>
							</Button>
						</Flex>
					</Form>
				</div>
			</div>
			<FileDataModal
				show={showModal}
				closeDialog={closeDialogAndReset}
				fileData={fileData}
			/>
		</div>
	)
}

export default App
