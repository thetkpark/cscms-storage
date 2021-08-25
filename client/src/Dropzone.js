import React from 'react'
import { useDropzone } from 'react-dropzone'
import { ProgressBar } from '@adobe/react-spectrum'
import Upload from '@spectrum-icons/illustrations/Upload'
import FileShareIcon from '@spectrum-icons/workflow/FileShare'
import styles from './Dropzone.module.css'

const Dropzone = ({ onDrop, selectedFilename, progress }) => {
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

	const renderProgressBar = () => {
		return (
			<ProgressBar
				label={progress === 100 ? 'Encrypting...' : 'Uploading...'}
				value={progress}
				isIndeterminate={progress === 100}
			/>
		)
	}

	return (
		<div className={styles.Dropzone} {...getRootProps()}>
			<input className="dropzone-input" {...getInputProps()} />
			<div
				className={styles.DropZoneTextContainer}
				style={
					isDragActive || selectedFilename.length !== 0
						? { border: '2px dashed #3fa13f' }
						: { border: '2px dashed #3f3f3f' }
				}
			>
				<FileShareIcon
					size="XXL"
					color={isDragActive || selectedFilename.length !== 0 ? 'positive' : ''}
				/>
				<p className={styles.DropZoneText}>{renderText()}</p>
				{progress !== -1 ? renderProgressBar() : null}
			</div>
		</div>
	)
}

export default Dropzone
