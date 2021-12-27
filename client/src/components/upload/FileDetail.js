import Icon from '../util/Icon'
import { toTitleCase } from '../../utils/formatText'
import FileIcon from '../util/FileIcon'
import { formatFileSize } from '../../utils/formatFileSize'

const FileDetail = ({ type, file, onRemove }) => {
	const getType = () => {
		return file.type !== '' ? toTitleCase(file.type.split('/')[0]) : 'File'
	}
	const getExtension = () => {
		return file.name.split('.').pop().toUpperCase()
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
				{type === 'image' ? (
					<img
						src={URL.createObjectURL(file)}
						alt={file.name}
						style={{ width: '100%' }}
					/>
				) : (
					<FileIcon ext={getExtension()} type={file.type.split('/')[0]} />
				)}
			</div>
			<div style={{ flex: '1', margin: '1.5rem' }}>
				<div>{file.name}</div>
				<div>
					{getType()} • {formatFileSize(file.size)} • {getExtension()}
				</div>
			</div>
			<div onClick={onRemove}>
				<Icon name="close" />
			</div>
		</div>
	)
}

export default FileDetail
