import { useEffect, useState } from 'react'
import DropZone from './Dropzone'
import Button from '../util/Button'
import Icon from '../util/Icon'
const UploadContainer = ({ type }) => {
	const [selectedFile, setSelectedFile] = useState(null)
	const [error, setError] = useState('')
	useEffect(() => {
		setSelectedFile(null)
		setError('')
	}, [type])
	const onDrop = (acceptedFiles, rejectedFiles) => {
		if (acceptedFiles.length === 1) {
			setError('')
			setSelectedFile(acceptedFiles[0])
			console.log(acceptedFiles[0])
		} else {
			if (rejectedFiles[0].errors[0].code === 'too-many-files') {
				setError('Too many files. You can only upload one file at a time')
			} else if (rejectedFiles[0].errors[0].code === 'file-too-large') {
				setError('File too big. The size limit is 100MB')
			} else setError('File not accepted')
		}
	}
	return (
		<div
			style={{
				background: 'white',
				width: '65vw',
				height: '60vh',
				margin: '1rem auto',
				borderRadius: '50px',
				padding: '3rem',
				display: 'flex',
				flexDirection: 'column',
				alignItems: 'center'
			}}
		>
			<DropZone
				type={type}
				selectedFilename={selectedFile ? selectedFile.name : ''}
				onDrop={onDrop}
			/>
			{selectedFile ? <div style={{ margin: '1rem' }}>{selectedFile.name}</div> : null}
			<Button
				bgColor={'#E9EEFF'}
				style={{
					border: 'none',
					fontSize: '.9rem',
					width: '170px',
					height: '50px',
					marginTop: '1rem'
				}}
			>
				<Icon name="upload" role="icon" /> Upload
			</Button>
		</div>
	)
}

export default UploadContainer
