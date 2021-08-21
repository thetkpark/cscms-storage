import './App.css'
import { useState } from 'react'
import axios from 'axios'
import FormData from 'form-data'

function App() {
	const [progress, setProgress] = useState(0)
	const [selectedFile, setSelectedFile] = useState()
	const [isFileSelected, setIsFileSelected] = useState(false)

	const changeHandler = event => {
		setSelectedFile(event.target.files[0])
		setIsFileSelected(true)
	}

	const handleSubmission = async event => {
		event.preventDefault()
		const formdata = new FormData()
		formdata.append('file', selectedFile)
		const res = await axios.post('http://localhost:5000/api/upload', formdata, {
			onUploadProgress: progressEvent => {
				const uploadPercent = Math.round(
					(progressEvent.loaded / progressEvent.total) * 100
				)
				setProgress(uploadPercent)
			}
		})
		console.log(res.data)
	}

	return (
		<div className="App">
			<input type="file" name="file" onChange={changeHandler} />
			<button onClick={handleSubmission}>Submit</button>
			<p>{progress}</p>
		</div>
	)
}

export default App
