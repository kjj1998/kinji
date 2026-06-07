import { Box, Group, Text } from "@mantine/core";
import { IconFileTypePdf } from "@tabler/icons-react";

interface FileChipProps {
	file: File | null;
}

function formatSize(bytes: number) {
	return bytes < 1024 ** 2
		? `${(bytes / 1024).toFixed(0)} KB`
		: `${(bytes / 1024 ** 2).toFixed(1)} MB`;
}

export function FileChip({ file }: FileChipProps) {
	return (
		<Group gap="sm" wrap="nowrap" style={{ minWidth: 0 }}>
			<IconFileTypePdf size={28} stroke={1.5} color="#1c1c1c" />
			<Box style={{ minWidth: 0 }}>
				<Text size="sm" fw={500} truncate>
					{file?.name}
				</Text>
				<Text size="xs" c="dimmed">
					{file && formatSize(file.size)}
				</Text>
			</Box>
		</Group>
	);
}
