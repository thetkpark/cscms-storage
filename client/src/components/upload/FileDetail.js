import Icon from '../util/Icon'
import { toTitleCase } from '../../utils/formatText'

const FileDetail = ({ type, file, onRemove }) => {
	const getType = () => {
		return file.type !== '' ? toTitleCase(file.type.split('/')[0]) : 'File'
	}
	const getFileSize = () => {
		if (file.size >= 1e6) {
			return `${parseInt(file.size / 1e6)} mb`
		} else if (file.size >= 1e3 && file.size < 1e6) {
			return `${parseInt(file.size / 1e3)} kb`
		} else {
			return `${file.size} b`
		}
	}
	const getExtension = () => {
		return file.name.split('.').pop().toUpperCase()
	}
	const getDisplay = () => {
		if (type === 'image') {
			return (
				<img src={URL.createObjectURL(file)} alt={file.name} style={{ width: '100%' }} />
			)
		} else {
			let ext = getExtension()
			if (ext === 'PDF') {
				return <Icon name="pdf-file" />
			} else if (file.type.split('/')[0] === 'image') {
				return <Icon name="img-file" />
			} else if (file.type.split('/')[0] === 'video') {
				return <Icon name="video-file" />
			} else if (file.type.split('/')[0] === 'audio') {
				return <Icon name="music-file" />
			} else if (ext === 'DOCX' || ext === 'DOC') {
				return <Icon name="document-file" />
			} else if (ext === 'XLSX' || ext === 'XLS') {
				return <Icon name="sheet-file" />
			} else if (ext === 'PPTX' || ext === 'PPT') {
				return <Icon name="slide-file" />
			} else if (ext === 'PSD') {
				return <Icon name="ps-file" />
			} else if (ext === 'AI') {
				return <Icon name="ai-file" />
			} else {
				return <Icon name="other-file" />
			}
		}
	}
	return (
		<div
			style={{ marginTop: '3rem', display: 'flex', width: '60%', alignItems: 'center' }}
		>
			<div
				style={{
					width: '80px',
					height: '80px',
					display: 'flex',
					alignItems: 'center',
					justifyContent: 'center'
				}}
			>
				{getDisplay()}
			</div>
			<div style={{ flex: '1', margin: '1.5rem' }}>
				<div>{file.name}</div>
				<div>
					{getType()} • {getFileSize()} • {getExtension()}
				</div>
			</div>
			<div onClick={onRemove}>
				<Icon name="close" />
			</div>
		</div>
	)
}

export default FileDetail
