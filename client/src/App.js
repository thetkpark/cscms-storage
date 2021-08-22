import './App.css'
import { useState } from 'react'
import axios from 'axios'
import FormData from 'form-data'

function App() {
	const [timeUsed, setTimeUsed] = useState(0)
	const [downloadTime, setDownloadTime] = useState(0)
	const [progress, setProgress] = useState(0)
	const [selectedFile, setSelectedFile] = useState()
	const [isFileSelected, setIsFileSelected] = useState(false)
	const [fileId, setFileId] = useState('')

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
			}
		})
		const end = new Date()
		setTimeUsed(end.getTime() - start.getTime())
		console.log(res.data)
	}

	const handleDownload = async event => {
		event.preventDefault()
		const start = new Date()
		const { data, headers } = await axios.get(`/${fileId}`, {
			responseType: 'blob'
		})
		console.log(headers)
		const downloadUrl = window.URL.createObjectURL(new Blob([data]))
		const link = document.createElement('a')
		link.href = downloadUrl
		link.setAttribute('download', headers['File-Name']) //any other extension
		document.body.appendChild(link)
		link.click()
		link.remove()

		const end = new Date()
		setDownloadTime(end.getTime() - start.getTime())
	}

	return (
		<div className="App">
			<input type="file" name="file" onChange={changeHandler} />
			<button onClick={handleSubmission}>Submit</button>
			<p>{progress}</p>
			<p>{timeUsed} ms</p>
			<br />
			<input type="text" placeholder="fileId" onChange={e => setFileId(e.target.value)} />
			<button onClick={handleDownload}>Submit</button>
			<p>{downloadTime}</p>
		</div>
	)
}

export default App
