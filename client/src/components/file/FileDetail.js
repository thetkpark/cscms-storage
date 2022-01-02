import { Fragment } from 'react'
import FileIcon from '../util/FileIcon'
import { formatFileSize } from '../../utils/formatFileSize'
import { formatDate } from '../../utils/formatText'
import Icon from '../util/Icon'

const FileDetail = ({ file }) => {
	return (
		<Fragment>
			<td width="35%">
				<FileIcon ext={file.filename.split('.')[1]} type={file.file_type.split("/")[0]} />{' '}
				{file.filename}
			</td>
			<td width="10%">{formatFileSize(file.file_size)}</td>
			<td width="20%">{formatDate(file.updated_at)}</td>
			<td width="35%">
				<Icon name="edit" />
				<Icon name="delete" />
				<Icon name="copy" />
			</td>
		</Fragment>
	)
}

export default FileDetail
