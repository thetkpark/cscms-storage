import { Fragment } from 'react'
import FileIcon from '../util/FileIcon'
import { formatFileSize } from '../../utils/formatFileSize'
import { formatDate } from '../../utils/formatText'
import Icon from '../util/Icon'
import styles from '../../styles/file/Detail.module.css';
const FileDetail = ({ file }) => {
	const copyToClipboard = () => {
		var copyText = document.createElement('input')
		copyText.setAttribute('value', file.url)
		document.body.appendChild(copyText)
		copyText.select()
		copyText.setSelectionRange(0, 99999)
		navigator.clipboard.writeText(copyText.value)
		document.body.removeChild(copyText)
	}
	return (
		<Fragment>
			<td>
				<FileIcon ext={file.filename.split('.')[1]} type={file.file_type.split('/')[0]} />{' '}
				{file.filename}
			</td>
			<td>{formatFileSize(file.file_size)}</td>
			<td>{formatDate(file.updated_at)}</td>
			<td>
				<div className={styles.ActionList}>
					<div>
						<Icon name="edit" />
					</div>
					<div>
						<Icon name="delete" />
					</div>
					<div onClick={copyToClipboard}>
						<Icon name="copy" />
					</div>
				</div>
			</td>
		</Fragment>
	)
}

export default FileDetail
