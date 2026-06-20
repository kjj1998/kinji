import { Alert } from "@mantine/core";

interface ErrorStatementProps {
	message: string;
	onDismiss: () => void;
}

export function ErrorStatement({ message, onDismiss }: ErrorStatementProps) {
	return (
		<Alert
			color="red"
			title="Couldn't process the statement"
			withCloseButton
			onClose={onDismiss}
		>
			{message || "Something went wrong while uploading."}
		</Alert>
	);
}
