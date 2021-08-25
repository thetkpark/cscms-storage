import { useState } from 'react'
import axios from 'axios'
import FormData from 'form-data'
import { Form, TextField, Button, Text, Flex } from '@adobe/react-spectrum'
import UploadIcon from '@spectrum-icons/workflow/UploadToCloudOutline'
import Dropzone from './Dropzone'
import styles from './App.module.css'
import FileDataModal from './Modal'

function App() {
	const [progress, setProgress] = useState(-1)
	const [selectedFile, setSelectedFile] = useState()
	const [isFileSelected, setIsFileSelected] = useState(false)
	const [slug, setSlug] = useState('')
	const [fileData, setFileData] = useState(undefined)
	const [error, setError] = useState('')
	const [selectedFilename, setSelectedFilename] = useState('')
	const [showModal, setShowModal] = useState(false)

	const onDrop = (acceptedFiles, rejectedFiles) => {
		if (acceptedFiles.length === 1) {
			setError('')
			setIsFileSelected(true)
			setSelectedFile(acceptedFiles[0])
			setSelectedFilename(acceptedFiles[0].name)
		} else {
			if (rejectedFiles[0].errors[0].code === 'too-many-files') {
				setError('Too many files')
			} else if (rejectedFiles[0].errors[0].code === 'file-too-large') {
				setError('File too big')
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
	}

	const closeAndReset = () => {
		setShowModal(false)
		setProgress(-1)
		setSelectedFile()
		setIsFileSelected(false)
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
				/>
				<div className={styles.FormContainer}>
					<Form className={styles.Form} onSubmit={handleSubmission}>
						<Flex direction="row" gap="size-300" alignItems="end" justifyContent="center">
							<TextField
								label="Custom slug for accessing the file (Optional)"
								placeholder="Slug"
								value={slug}
								onChange={e => setSlug(e.target.value)}
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
				{progress < 0 ? null : <p>{progress}%</p>}
				{error.length > 0 ? <p>{error}</p> : null}
			</div>
			<FileDataModal show={showModal} onClose={closeAndReset} fileData={fileData} />
		</div>
	)
}

export default App
