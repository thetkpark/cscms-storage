import React, { Fragment } from 'react'
import { useDropzone } from 'react-dropzone'
import styles from '../../styles/upload/Dropzone.module.css'
const DropZone = ({ onDrop, selectedFilename, type }) => {
	const { getRootProps, getInputProps, isDragActive } = useDropzone({
		onDrop,
		maxFiles: 1,
		multiple: false,
		maxSize: type === 'file' ? 100 << 20 : (100 << 20) / 20
	})

	const renderContainerText = () => {
		if (type === 'file') {
			return (
				<Fragment>
					<img src="file.png" alt="file" />
					<h2>Drag and Drop your file</h2>
					<p>Drop your file here, or click to select file</p>
					<p>The maximum file size is 100MB</p>
				</Fragment>
			)
		} else if (type === 'image') {
			return (
				<Fragment>
					<img src="image.png" alt="image" />
					<h2>Drag and Drop your image</h2>
					<p>Drop your image here, or click to select image</p>
					<p>The maximum file size is 5MB</p>
				</Fragment>
			)
		}
	}

	return (
		<div className={`${styles.Dropzone}`} {...getRootProps()}>
			<input className="dropzone-input" {...getInputProps()} />
			<div className={styles.DropZoneTextContainer}>{renderContainerText()}</div>
		</div>
	)
}

export default DropZone
