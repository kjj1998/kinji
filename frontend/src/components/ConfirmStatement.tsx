import {
	ActionIcon,
	Button,
	Group,
	Paper,
	PasswordInput,
	Stack,
} from "@mantine/core";
import { IconX } from "@tabler/icons-react";
import { FileChip } from "./FileChip";

interface ConfirmStatementProps {
	file: File;
	password: string;
	onPasswordChange: (password: string) => void;
	onCancel: () => void;
	onConfirm: () => void;
}

export function ConfirmStatement({
	file,
	password,
	onPasswordChange,
	onCancel,
	onConfirm,
}: ConfirmStatementProps) {
	return (
		<Paper withBorder p="lg" radius="md">
			<Stack gap="lg">
				<Group justify="space-between" wrap="nowrap">
					<FileChip file={file} />
					<ActionIcon
						variant="subtle"
						color="gray"
						onClick={onCancel}
						aria-label="Remove file"
					>
						<IconX size={18} />
					</ActionIcon>
				</Group>

				<PasswordInput
					label="Statement password"
					description="Only needed if your PDF is password-protected"
					placeholder="Leave blank if none"
					value={password}
					onChange={(e) => onPasswordChange(e.currentTarget.value)}
				/>

				<Group justify="flex-end">
					<Button variant="default" onClick={onCancel}>
						Cancel
					</Button>
					<Button color="dark" onClick={onConfirm}>
						Parse statement
					</Button>
				</Group>
			</Stack>
		</Paper>
	);
}
