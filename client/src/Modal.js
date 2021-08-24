import dayjs from 'dayjs'
import localizedFormat from 'dayjs/plugin/localizedFormat'
import { Modal, Button } from 'react-bootstrap'
dayjs.extend(localizedFormat)

const FileDataModal = ({ show, onClose, fileData }) => {
	if (!show) {
		return null
	}
	return (
		<Modal show={show} centered>
			<Modal.Header closeButton>
				<Modal.Title>Your File</Modal.Title>
			</Modal.Header>
			<Modal.Body>
				<a
					href={`${window.location.href}${fileData.token}`}
				>{`${window.location.href}${fileData.token}`}</a>
				<p>Token: {fileData.token}</p>
				<p>File Size: {fileData.file_size} bytes</p>
				<p>Valid Though: {dayjs(fileData.created_at).add(1, 'month').format('LLL')}</p>
			</Modal.Body>
			<Modal.Footer>
				<Button onClick={onClose}>Close</Button>
			</Modal.Footer>
		</Modal>
	)
}

export default FileDataModal
