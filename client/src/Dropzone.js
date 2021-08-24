import React from 'react'
import { useDropzone } from 'react-dropzone'

const Dropzone = ({ onDrop, selectedFilename }) => {
	const { getRootProps, getInputProps, isDragActive } = useDropzone({
		onDrop,
		maxFiles: 1,
		multiple: false,
		maxSize: 100 << 20 // 100MB
	})

	const renderText = text => {
		if (isDragActive)
			return <p className="dropzone-content">Release to drop the files here</p>
		if (selectedFilename.length !== 0)
			return <p className="dropzone-content">{selectedFilename}</p>
		return (
			<p className="dropzone-content">
				Drag 'n' drop some files here, or click to select files
			</p>
		)
	}

	return (
		<div {...getRootProps()}>
			<input className="dropzone-input" {...getInputProps()} />
			<div className="text-center" style={{ height: '50vh', border: 'dashed red' }}>
				{renderText()}
			</div>
		</div>
	)
}

export default Dropzone
