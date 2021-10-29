import { AlertDialog, DialogContainer } from '@adobe/react-spectrum'

const ErrorDialog = ({ errorMessage, closeDialog }) => {
	return (
		<DialogContainer>
			<AlertDialog
				variant="error"
				title="Error occur"
				primaryActionLabel="Close"
				onPrimaryAction={closeDialog}
			>
				{errorMessage}
			</AlertDialog>
		</DialogContainer>
	)
}

export default ErrorDialog
