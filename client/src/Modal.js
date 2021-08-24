import dayjs from 'dayjs'
import localizedFormat from 'dayjs/plugin/localizedFormat'
import { Modal, Button } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faClipboard } from '@fortawesome/free-solid-svg-icons'
dayjs.extend(localizedFormat)

const FileDataModal = ({ show, onClose, fileData }) => {
	if (!show) {
		return null
	}
	const location = `${window.location.href}${fileData.token}`
	return (
		<Modal show={show} centered>
			<Modal.Header closeButton>
				<Modal.Title>Your File</Modal.Title>
			</Modal.Header>
			<Modal.Body>
				<a href={location}>{location}</a>
				<FontAwesomeIcon
					icon={faClipboard}
					onClick={() => navigator.clipboard.writeText(location)}
					cursor="pointer"
				/>
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
