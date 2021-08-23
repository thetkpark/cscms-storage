import './App.css'
import { useState } from 'react'
import axios from 'axios'
import FormData from 'form-data'
import dayjs from 'dayjs'
import localizedFormat from 'dayjs/plugin/localizedFormat'
dayjs.extend(localizedFormat)

function App() {
	const [timeUsed, setTimeUsed] = useState(-1)
	const [progress, setProgress] = useState(-1)
	const [selectedFile, setSelectedFile] = useState()
	const [isFileSelected, setIsFileSelected] = useState(false)
	const [slug, setSlug] = useState('')
	const [fileData, setFileData] = useState(undefined)

	const changeHandler = event => {
		setSelectedFile(event.target.files[0])
		setIsFileSelected(true)
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
		<div className="App">
			<input type="file" name="file" onChange={changeHandler} />
			<input
				type="text"
				name="slug"
				placeholder="slug"
				onChange={e => setSlug(e.target.value)}
			/>
			<button onClick={handleSubmission}>Submit</button>
			{progress < 0 ? null : <p>{progress}%</p>}
			{timeUsed < 0 ? null : <p>{timeUsed}ms</p>}
			{renderFileData()}
		</div>
	)
}

export default App
