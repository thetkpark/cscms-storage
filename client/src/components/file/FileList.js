import { Fragment, useEffect, useState } from 'react'
import FileIcon from '../util/FileIcon'
import { formatFileSize } from '../../utils/formatFileSize'
import styles from '../../styles/file/FileList.module.css'
import Icon from '../util/Icon'
import axios from 'axios'
const FileList = () => {
	const [sort, setSort] = useState({ name: '', asc: true })
	const [files, setFiles] = useState([])
	useEffect(() => {
		fetchFiles()
	}, [])
	const fetchFiles = async () => {
		const fileRes = await axios.get('https://storage.cscms.me/api/file')
		const fileData = fileRes.data
		const imageRes = await axios.get('https://storage.cscms.me/api/image')
		const imageData = imageRes.data
		setFiles([...fileData, ...imageData])
	}
	const [displayFile, setDisplayFile] = useState(files)
	useEffect(() => {
		if (sort.name === '') {
			setDisplayFile(files)
		} else {
			let temp = [...files].sort((a, b) => {
				if (sort.name === 'size') {
					if (sort.asc) return a.size - b.size
					return b.size - a.size
				}
				if (sort.asc) {
					return a[sort.name] > b[sort.name] ? 1 : -1
				} else {
					return b[sort.name] > a[sort.name] ? 1 : -1
				}
			})
			setDisplayFile(temp)
		}
	}, [sort])
	const handleSort = type => {
		if (type === sort.name) {
			if (sort.asc) {
				setSort({ name: type, asc: false })
			} else {
				setSort({ name: '', asc: true })
			}
		} else {
			setSort({ name: type, asc: true })
		}
	}

	return (
		<Fragment>
			<div className={styles.FileListWrapper}>
				<h3>My Files</h3>
				<div>
					<table className={styles.FileList}>
						<thead>
							<tr>
								<th>
									<div onClick={() => handleSort('filename')}>
										Name{' '}
										{sort.name === 'filename' && sort.asc ? (
											<Icon name="arrow-down" />
										) : (
											<Icon name="arrow-up" />
										)}
									</div>
								</th>
								<th>
									<div onClick={() => handleSort('file_size')}>
										Size{' '}
										{sort.name === 'file_size' && sort.asc ? (
											<Icon name="arrow-down" />
										) : (
											<Icon name="arrow-up" />
										)}
									</div>
								</th>
								<th>
									<div onClick={() => handleSort('updated_at')}>
										Last Modified{' '}
										{sort.name === 'updated_at' && sort.asc ? (
											<Icon name="arrow-down" />
										) : (
											<Icon name="arrow-up" />
										)}
									</div>
								</th>
								<th></th>
							</tr>
						</thead>
						<tbody>
							{displayFile.length === 0 ? (
								<Fragment>
									<tr>
										<td colSpan={4}>No files found</td>
									</tr>
								</Fragment>
							) : (
								displayFile.map((file, index) => {
									return (
										<tr key={index}>
											<td>
												<FileIcon
													ext={file.filename.split('.')[1]}
													type={file.file_type}
												/>{' '}
												{file.filename}
											</td>
											<td>{formatFileSize(file.file_size)}</td>
											<td>{file.updated_at}</td>
											<td>
												<div className={styles.EditIcon}>
													<Icon name="edit" />
												</div>
											</td>
										</tr>
									)
								})
							)}
						</tbody>
					</table>
				</div>
			</div>
		</Fragment>
	)
}

export default FileList
