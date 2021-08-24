import React from 'react'
import { useDropzone } from 'react-dropzone'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faFileUpload } from '@fortawesome/free-solid-svg-icons'
import styles from './Dropzone.module.css'

const Dropzone = ({ onDrop, selectedFilename }) => {
	const { getRootProps, getInputProps, isDragActive } = useDropzone({
		onDrop,
		maxFiles: 1,
		multiple: false,
		maxSize: 100 << 20 // 100MB
	})

	const renderText = () => {
		if (isDragActive) return 'Release to drop your file here'
		if (selectedFilename.length !== 0) return selectedFilename
		return 'Drop your file here, or click to select file'
	}

	return (
		<div className={styles.Dropzone} {...getRootProps()}>
			<input className="dropzone-input" {...getInputProps()} />
			<div className={styles.DropZoneTextContainer}>
				<FontAwesomeIcon icon={faFileUpload} size="5x" color="#a4e1ef" />
				<p className={styles.DropZoneText}>{renderText()}</p>
			</div>
		</div>
	)
}

export default Dropzone
