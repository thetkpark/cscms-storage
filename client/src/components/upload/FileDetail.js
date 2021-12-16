import Icon from '../util/Icon'
import { toTitleCase } from '../../utils/formatText'

const FileDetail = ({ type, file, onRemove }) => {
	const getType = () => {
		return toTitleCase(file.type.split('/')[0])
	}
	const getFileSize = () => {
		return file.size
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
			return <span>{getType()}</span>
		}
	}
	return (
		<div
			style={{ marginTop: '3rem', display: 'flex', width: '60%', alignItems: 'center' }}
		>
			<div style={{ width: '80px', height: '80px' }}>{getDisplay()}</div>
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
