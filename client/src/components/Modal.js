import dayjs from 'dayjs'
import localizedFormat from 'dayjs/plugin/localizedFormat'
import {
	DialogContainer,
	Dialog,
	Heading,
	Divider,
	Content,
	Text,
	ButtonGroup,
	Link,
	Flex,
	Button
} from '@adobe/react-spectrum'
import { formatBytes } from '../utils/formatBytes'
dayjs.extend(localizedFormat)

const FileDataModal = ({ show, fileData, closeDialog }) => {
	if (!show) {
		return null
	}

	const location = `${window.location.href}${fileData.token}`
	return (
		<DialogContainer onDismiss={closeDialog}>
			<Dialog>
				<Heading>Your File</Heading>
				<Divider />
				<Content>
					<Flex direction="column" justifyContent="start">
						<Flex direction="row" justifyContent="start">
							<Text>
								Download URL:{' '}
								<Link>
									<a href={location}>{location}</a>
								</Link>
							</Text>
						</Flex>
						<Text>File name: {fileData.file_name}</Text>
						<Text>File size: {formatBytes(fileData.file_size)}</Text>
						<Text>Valid Though: {dayjs(fileData.expired_at).format('LLL')}</Text>
					</Flex>
				</Content>
				<ButtonGroup>
					<Button
						variant="primary"
						onPress={() => navigator.clipboard.writeText(location)}
					>
						Copy URL
					</Button>
					<Button variant="primary" onPress={closeDialog}>
						Close
					</Button>
				</ButtonGroup>
			</Dialog>
		</DialogContainer>
		// <Modal show={show} centered>
		// 	<Modal.Header closeButton>
		// 		<Modal.Title>Your File</Modal.Title>
		// 	</Modal.Header>
		// 	<Modal.Body>
		// 		<a href={location}>{location}</a>
		// 		<FontAwesomeIcon
		// 			icon={faClipboard}
		// 			onClick={() => navigator.clipboard.writeText(location)}
		// 			cursor="pointer"
		// 		/>
		// 		<p>File Size: {fileData.file_size} bytes</p>
		// 		<p>Valid Though: {dayjs(fileData.created_at).add(1, 'month').format('LLL')}</p>
		// 	</Modal.Body>
		// 	<Modal.Footer>
		// 		<Button onClick={onClose}>Close</Button>
		// 	</Modal.Footer>
		// </Modal>
	)
}

export default FileDataModal