import { Fragment, useEffect, useState } from 'react'
import FileIcon from '../util/FileIcon'
import { formatFileSize } from '../../utils/formatFileSize'
import styles from '../../styles/file/FileList.module.css'
import Icon from '../util/Icon'
const FileList = () => {
	const [sort, setSort] = useState({ name: '', asc: true })
	const [files, setFiles] = useState([
		{
			name: 'b.docx',
			ext: 'DOCX',
			type: 'document',
			size: 100000,
			lastModified: 'Dec, 13 2021'
		},
		{
			name: 'a.docx',
			ext: 'DOCX',
			type: 'document',
			size: 1500000,
			lastModified: 'Dec, 13 2021'
		},
		{
			name: 'c.docx',
			ext: 'DOCX',
			type: 'document',
			size: 1200000,
			lastModified: 'Dec, 13 2021'
		}
	])
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
									<div onClick={() => handleSort('name')}>
										Name{' '}
										{sort.name === 'name' && sort.asc ? (
											<Icon name="arrow-down" />
										) : (
											<Icon name="arrow-up" />
										)}
									</div>
								</th>
								<th>
									<div onClick={() => handleSort('size')}>
										Size{' '}
										{sort.name === 'size' && sort.asc ? (
											<Icon name="arrow-down" />
										) : (
											<Icon name="arrow-up" />
										)}
									</div>
								</th>
								<th>
									<div onClick={() => handleSort('lastModified')}>
										Last Modified{' '}
										{sort.name === 'lastModified' && sort.asc ? (
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
												<FileIcon ext={file.ext} type={file.type} /> {file.name}
											</td>
											<td>{formatFileSize(file.size)}</td>
											<td>{file.lastModified}</td>
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
