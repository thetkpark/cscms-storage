import { useEffect, useState } from 'react'
import DropZone from './Dropzone'
import Button from '../util/Button'
import Icon from '../util/Icon'
import { Slider, TextField } from '@material-ui/core'
const UploadContainer = ({ type }) => {
	const [selectedFile, setSelectedFile] = useState(null)
	const [error, setError] = useState('')
	const [duration, setDuration] = useState(7)
	const [slug, setSlug] = useState('')
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
				height: '75vh',
				display: 'flex',
				flexDirection: 'column',
				margin: '1rem auto'
			}}
		>
			<div
				style={{
					background: 'white',
					width: '65vw',
					flex: '1',
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
				{selectedFile ? (
					<div style={{ marginTop: '2rem' }}>{selectedFile.name}</div>
				) : null}
				<Button
					bgColor={'#E9EEFF'}
					style={{
						border: 'none',
						fontSize: '.9rem',
						width: '170px',
						height: '50px',
						marginTop: '2rem'
					}}
				>
					<Icon name="upload" role="icon" /> Upload
				</Button>
			</div>
			{type === 'file' ? (
				<div style={{ marginTop: '3rem', display: 'flex', justifyContent: 'center' }}>
					<div style={{ width: '280px', padding: '0 3rem' }}>
						<div
							style={{
								display: 'flex',
								justifyContent: 'space-between',
								marginBottom: '1rem'
							}}
						>
							<div>Storage Duration (Days)</div>
							<div>{duration}</div>
						</div>
						<div style={{ width: '85%' }}>
							<Slider
								value={duration}
								min={1}
								max={30}
								onChange={(event, val) => {
									setDuration(val)
								}}
							/>
						</div>
					</div>
					<div style={{ padding: '0 3rem' }}>
						<div
							style={{
								marginBottom: '1rem'
							}}
						>
							Custom Slug for accessing the file (Optional)
						</div>
						<div>
							<TextField
								variant="outlined"
								role="textbox"
								placeholder="slug"
								value={slug}
								onChange={(event, val) => setSlug(val)}
							/>
						</div>
					</div>
				</div>
			) : null}
		</div>
	)
}

export default UploadContainer
