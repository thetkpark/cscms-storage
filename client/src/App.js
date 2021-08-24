import { useState } from 'react'
import axios from 'axios'
import FormData from 'form-data'
import Dropzone from './Dropzone'
import styles from './App.module.css'
import { Form } from 'react-bootstrap'
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

		const res = await axios.post('http://localhost:5000/api/file', formdata, {
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
				<Dropzone onDrop={onDrop} selectedFilename={selectedFilename} />
				<div className={styles.FormContainer}>
					<Form className={styles.Form} onSubmit={handleSubmission}>
						<Form.Group controlId="formBasicSlug">
							<Form.Control
								type="text"
								placeholder="Enter slug"
								value={slug}
								onChange={e => setSlug(e.target.value)}
							/>
						</Form.Group>
						<Form.Group controlId="formBasicSubmit">
							<Form.Control type="submit" value="Submit" />
						</Form.Group>
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
