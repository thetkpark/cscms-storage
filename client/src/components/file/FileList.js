import { Fragment, useEffect, useState } from 'react'
import styles from '../../styles/file/FileList.module.css'
import Icon from '../util/Icon'
import axios from 'axios'
import FileDetail from './FileDetail'
import ImageDetail from './ImageDetail'
const FileList = ({setError}) => {
	const [sort, setSort] = useState({ name: '', asc: true })
	const [files, setFiles] = useState([])
	const [displayFile, setDisplayFile] = useState([])
	useEffect(() => {
		fetchFiles()
	}, [])
	const fetchFiles = async () => {
		// const fileRes = await axios.get('https://storage.cscms.me/api/file')
		let fileRes = {}
		fileRes.data = [
			{
				id: 'IQWKoJx_dYBBmL69d3xxvmgea3Wvjm',
				created_at: '2022-01-01T12:38:04.953Z',
				updated_at: '2022-01-02T08:38:11.937Z',
				expired_at: '2022-01-08T12:38:04.953Z',
				token: 'rrnutv',
				nonce: 'b77c1799b006e7fb2eaeac0473c9568e326c04f05f68850538af17be80bc53c7',
				filename: 'cactus.png',
				file_size: 96158787,
				visited: 1,
				UserID: 3,
				file_type: 'image/png',
				encrypted: true,
				DeletedAt: null
			},
			{
				id: 'KKMyT_XdwGFSRvqt20iw0rIY-c-Vec',
				created_at: '2021-12-31T05:34:17.36Z',
				updated_at: '2021-12-31T05:34:17.36Z',
				expired_at: '2022-01-07T05:34:17.36Z',
				token: '22u4ie',
				nonce: '93972fe4109056888b5d99d015f234c428b745f2a6f7ea67499915f9fa98fbbb',
				filename: 'main.py',
				file_size: 0,
				visited: 0,
				UserID: 3,
				file_type: 'text/x-python-script',
				encrypted: true,
				DeletedAt: null
			},
			{
				id: 'RD_q9z3MNZ-fCfGv_6Hrz2C30q0ytl',
				created_at: '2022-01-02T07:52:48.109Z',
				updated_at: '2022-01-02T07:52:54.526Z',
				expired_at: '2022-01-09T07:52:48.109Z',
				token: 'laby8f',
				nonce: '445ecbd1058900a0df2e6add140f941136e5bf48cb7a5a7979de2e11785f7de5',
				filename: 'Screen Shot 2565-01-02 at 14.06.44.png',
				file_size: 61266,
				visited: 1,
				UserID: 3,
				file_type: 'image/png',
				encrypted: true,
				DeletedAt: null
			},
			{
				id: 'smA2dYM6uu6mrcMgONmi10o3tZVZYB',
				created_at: '2021-12-30T13:00:59.415Z',
				updated_at: '2021-12-30T13:00:59.415Z',
				expired_at: '2022-01-29T13:00:59.415Z',
				token: 'ztnnvk',
				nonce: 'af5a7f5a02864c48189e9eb2abcbacd0dc35c1ceede035f70bb9e6d525229d5e',
				filename: 'index.html',
				file_size: 941,
				visited: 0,
				UserID: 3,
				file_type: 'text/html',
				encrypted: true,
				DeletedAt: null
			},
			{
				id: 'vFnIP_s_4irDGgt-i1-XvafNi9rdIV',
				created_at: '2022-01-01T12:54:54.977Z',
				updated_at: '2022-01-01T12:54:54.977Z',
				expired_at: '2022-01-08T12:54:54.977Z',
				token: 'zck3z1',
				nonce: '17468ede101307fd9cf708bfbdf7b929bc394f6a2afaedb637dc6f45fe22754d',
				filename: 'poster firebase.png',
				file_size: 821766,
				visited: 0,
				UserID: 3,
				file_type: 'image/png',
				encrypted: true,
				DeletedAt: null
			},
			{
				id: 'wJqxzRVpvgJ8-GSsIJulEFKGaBHkuv',
				created_at: '2021-12-30T13:01:00.625Z',
				updated_at: '2021-12-30T13:01:00.625Z',
				expired_at: '2022-01-29T13:01:00.625Z',
				token: 'v1ddtb',
				nonce: 'ce072bb3a9aefd70df9dab4e03290c2c7460745838a546600a7e3e01b46227bf',
				filename: 'index.html',
				file_size: 941,
				visited: 0,
				UserID: 3,
				file_type: 'text/html',
				encrypted: true,
				DeletedAt: null
			},
			{
				id: 'X6mRPurR6VQrnPP1bKvbTvO8u2Ucww',
				created_at: '2022-01-01T12:58:10.772Z',
				updated_at: '2022-01-01T12:58:10.772Z',
				expired_at: '2022-01-08T12:58:10.772Z',
				token: 'j2rypv',
				nonce: '582650c983fd25a763102d293fc3ad709312af5de9de43d9dbb30231d5393520',
				filename: 'poster firebase.png',
				file_size: 821766,
				visited: 0,
				UserID: 3,
				file_type: 'image/png',
				encrypted: true,
				DeletedAt: null
			}
		]
		const fileData = fileRes.data.map(file => ({
			...file,
			type: 'file',
			url: 'https://storage.cscms.me/' + file.token
		}))
		// const imageRes = await axios.get('https://storage.cscms.me/api/image')
		let imageRes = {}
		imageRes.data = [
			{
				id: 19,
				created_at: '2022-01-01T12:39:09.019Z',
				updated_at: '2022-01-01T12:39:09.019Z',
				original_filename: 'poster firebase.png',
				file_size: 821766,
				file_path: 'x6HOTDxi8ijWlt9IJzgh.png',
				user_id: 3,
				DeletedAt: null
			},
			{
				id: 20,
				created_at: '2022-01-01T12:59:34.504Z',
				updated_at: '2022-01-01T12:59:34.504Z',
				original_filename: 'poster firebase.png',
				file_size: 821766,
				file_path: '4qNdVbUOPJleYH2e4gml.png',
				user_id: 3,
				DeletedAt: null
			},
			{
				id: 22,
				created_at: '2022-01-02T07:53:01.785Z',
				updated_at: '2022-01-02T07:53:01.785Z',
				original_filename: 'Screen Shot 2565-01-02 at 14.06.44.png',
				file_size: 61266,
				file_path: '8xGx6bkbuN11rBvKO0QL.png',
				user_id: 3,
				DeletedAt: null
			}
		]
		const imageData = imageRes.data.map(image => ({
			...image,
			type: 'image',
			url: 'https://img.cscms.me/' + image.file_path,
			file_type: 'image',
			filename: image.original_filename
		}))
		setFiles([...fileData, ...imageData])
	}
	useEffect(() => {
		setDisplayFile(files)
		setSort({ name: '', asc: true })
	}, [files])
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
								<th style={{ width: '45%' }}>
									<div onClick={() => handleSort('filename')}>
										Name{' '}
										{sort.name === 'filename' && sort.asc ? (
											<Icon name="arrow-down" />
										) : (
											<Icon name="arrow-up" />
										)}
									</div>
								</th>
								<th style={{ width: '15%' }}>
									<div onClick={() => handleSort('file_size')}>
										Size{' '}
										{sort.name === 'file_size' && sort.asc ? (
											<Icon name="arrow-down" />
										) : (
											<Icon name="arrow-up" />
										)}
									</div>
								</th>
								<th style={{ width: '25%' }}>
									<div onClick={() => handleSort('updated_at')}>
										Last Modified{' '}
										{sort.name === 'updated_at' && sort.asc ? (
											<Icon name="arrow-down" />
										) : (
											<Icon name="arrow-up" />
										)}
									</div>
								</th>
								<th style={{ width: '15%' }}></th>
							</tr>
						</thead>
						<tbody>
							{displayFile.length === 0 ? (
								<Fragment>
									<tr>
										<td className={styles.Empty} colSpan={4}>
											No files found
										</td>
									</tr>
								</Fragment>
							) : (
								displayFile.map((file, index) => {
									return (
										<tr key={index} className={styles.Row}>
											{file.type === 'file' ? (
												<FileDetail setError={setError} file={file} />
											) : (
												<ImageDetail setError={setError} file={file} />
											)}
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
