import { useEffect, useState } from 'react'
import DropZone from './Dropzone'
import Button from '../util/Button'
import Icon from '../util/Icon'
import { Slider, TextField } from '@material-ui/core'
import FileDetail from './FileDetail'
const UploadContainer = ({ type, handleUpload, setError, progress }) => {
	const [selectedFile, setSelectedFile] = useState(null)
	const [duration, setDuration] = useState(7)
	const [slug, setSlug] = useState('')
	useEffect(() => {
		setSelectedFile(null)
	}, [type])
	const onDrop = (acceptedFiles, rejectedFiles) => {
		if (acceptedFiles.length === 1) {
			if (type === 'image') {
				if (acceptedFiles[0].type.includes('image')) {
					setSelectedFile(acceptedFiles[0])
				} else {
					setError('Please upload an image')
					return
				}
			}
			setSelectedFile(acceptedFiles[0])
		} else {
			if (rejectedFiles[0].errors[0].code === 'too-many-files') {
				setError('Too many files. You can only upload one file at a time')
			} else if (rejectedFiles[0].errors[0].code === 'file-too-large') {
				setError('File too big. The size limit is 100MB')
			} else setError('File not accepted')
			return
		}
	}
	const onClick = async e => {
		e.preventDefault()
		if (!selectedFile) {
			setError('Please select a file')
			return
		}
		try {
			await handleUpload({
				selectedFile,
				slug,
				duration,
				clearFile: () => setSelectedFile(null)
			})
		} catch (err) {
			setError(err.response.data.message)
		}
	}
	return (
		<div
			style={{
				minHeight: '75vh',
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
					progress={progress}
				/>
				{selectedFile ? (
					<FileDetail
						type={type}
						file={selectedFile}
						onRemove={() => setSelectedFile(null)}
					/>
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
					action={onClick}
				>
					<Icon name="upload" role="icon" /> Upload
				</Button>
			</div>
			{type === 'file' ? (
				<div style={{ marginTop: '1rem', display: 'flex', justifyContent: 'center' }}>
					<div style={{ width: '280px', padding: '0 1rem' }}>
						<div
							style={{
								display: 'flex',
								justifyContent: 'space-between',
								marginBottom: '1rem',
								fontWeight: '500'
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
								marginBottom: '1rem',
								fontWeight: '500'
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
