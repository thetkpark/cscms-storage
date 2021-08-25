import React from 'react'
import { useDropzone } from 'react-dropzone'
import { ProgressBar, Heading, Content } from '@adobe/react-spectrum'
import Upload from '@spectrum-icons/illustrations/Upload'
import FileShareIcon from '@spectrum-icons/workflow/FileShare'
import styles from './Dropzone.module.css'

const Dropzone = ({ onDrop, selectedFilename, progress, error }) => {
	const { getRootProps, getInputProps, isDragActive } = useDropzone({
		onDrop,
		maxFiles: 1,
		multiple: false,
		maxSize: 100 << 20 // 100MB
	})

	const renderText = () => {
		if (isDragActive) return 'Release to drop your file here'
		if (selectedFilename.length !== 0) return selectedFilename
		if (error.length !== 0) return error
		return 'Drop your file here, or click to select file'
	}

	const renderProgressBar = () => {
		return (
			<ProgressBar
				label={progress === 100 ? 'Encrypting...' : 'Uploading...'}
				value={progress}
				isIndeterminate={progress === 100}
			/>
		)
	}

	const containerStyle = () => {
		if (isDragActive || selectedFilename.length !== 0) {
			return { border: '2px dashed #3fa13f' }
		}
		if (error.length !== 0) {
			return { border: '2px dashed #a13f3f' }
		}
		return { border: '2px dashed #3f3f3f' }
	}

	const iconColor = () => {
		if (isDragActive || selectedFilename.length !== 0) {
			return 'positive'
		}
		if (error.length !== 0) {
			return 'negative'
		}
		return ''
	}

	return (
		<div className={styles.Dropzone} {...getRootProps()}>
			<input className="dropzone-input" {...getInputProps()} />
			<div className={styles.DropZoneTextContainer} style={containerStyle()}>
				<FileShareIcon size="XXL" color={iconColor()} />
				{/* <Upload /> */}
				<Heading>Drag and Drop your file</Heading>
				<Content>{renderText()}</Content>
				{progress !== -1 ? renderProgressBar() : null}
			</div>
		</div>
	)
}

export default Dropzone
