import './App.css'
import { useState } from 'react'
import axios from 'axios'
import FormData from 'form-data'

function App() {
	const [timeUsed, setTimeUsed] = useState(0)
	const [progress, setProgress] = useState(0)
	const [selectedFile, setSelectedFile] = useState()
	const [isFileSelected, setIsFileSelected] = useState(false)

	const changeHandler = event => {
		setSelectedFile(event.target.files[0])
		setIsFileSelected(true)
	}

	const handleSubmission = async event => {
		event.preventDefault()
		const start = new Date()
		const formdata = new FormData()
		formdata.append('file', selectedFile)
		const res = await axios.post('/api/upload', formdata, {
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

	return (
		<div className="App">
			<input type="file" name="file" onChange={changeHandler} />
			<button onClick={handleSubmission}>Submit</button>
			<p>{progress}</p>
			<p>{timeUsed} ms</p>
		</div>
	)
}

export default App
