import { Button, Stack, Text } from "@mantine/core";
import { Dropzone } from "@mantine/dropzone";
import { IconFileTypePdf, IconUpload, IconX } from "@tabler/icons-react";

interface SelectStatementProps {
	maxSize: number;
	onFileSelect: (file: File) => void;
}

export function SelectStatement({
	maxSize,
	onFileSelect,
}: SelectStatementProps) {
	return (
		<Dropzone
			accept={["application/pdf"]}
			maxSize={maxSize}
			multiple={false}
			onDrop={(files) => onFileSelect(files[0])}
		>
			<Stack
				align="center"
				justify="center"
				gap="sm"
				mih={220}
				style={{ pointerEvents: "none" }}
			>
				<Dropzone.Idle>
					<IconUpload size={40} stroke={1.5} color="#1c1c1c" />
				</Dropzone.Idle>
				<Dropzone.Accept>
					<IconFileTypePdf size={40} stroke={1.5} color="#1c1c1c" />
				</Dropzone.Accept>
				<Dropzone.Reject>
					<IconX size={40} stroke={1.5} color="#8a3a1f" />
				</Dropzone.Reject>

				<Text size="lg" ta="center">
					Drag & drop your statement
				</Text>
				<Text size="sm" c="dimmed" ta="center">
					or
				</Text>
				<Button variant="default" size="sm" style={{ pointerEvents: "auto" }}>
					Browse files
				</Button>
				<Text size="sm" c="dimmed" ta="center">
					PDF only · up to 10 MB
				</Text>
			</Stack>
		</Dropzone>
	);
}
