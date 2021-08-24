import { useState } from 'react'
import axios from 'axios'
import FormData from 'form-data'
import dayjs from 'dayjs'
import localizedFormat from 'dayjs/plugin/localizedFormat'
import Dropzone from './Dropzone'
import styles from './App.module.css'
import { Form } from 'react-bootstrap'
dayjs.extend(localizedFormat)

function App() {
	const [timeUsed, setTimeUsed] = useState(-1)
	const [progress, setProgress] = useState(-1)
	const [selectedFile, setSelectedFile] = useState()
	const [isFileSelected, setIsFileSelected] = useState(false)
	const [slug, setSlug] = useState('')
	const [fileData, setFileData] = useState(undefined)
	const [error, setError] = useState('')
	const [selectedFilename, setSelectedFilename] = useState('')

	const onDrop = (acceptedFiles, rejectedFiles) => {
		if (acceptedFiles.length === 1) {
			console.log(acceptedFiles)
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
		const start = new Date()
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
		const end = new Date()
		setTimeUsed(end.getTime() - start.getTime())
		setFileData(res.data)
	}

	const renderFileData = () => {
		if (fileData === undefined) {
			return null
		}
		return (
			<div>
				<h3>Your File</h3>
				<a
					href={`${window.location.href}${fileData.token}`}
				>{`${window.location.href}${fileData.token}`}</a>
				<p>Token: {fileData.token}</p>
				<p>File Size: {fileData.file_size} bytes</p>
				<p>Valid Though: {dayjs(fileData.created_at).add(1, 'month').format('LLL')}</p>
			</div>
		)
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
				{timeUsed < 0 ? null : <p>{timeUsed}ms</p>}
				{renderFileData()}
				{error.length > 0 ? <p>{error}</p> : null}
			</div>
		</div>
	)
}

export default App
