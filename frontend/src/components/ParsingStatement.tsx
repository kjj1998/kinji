import {
	Box,
	Button,
	Group,
	Loader,
	Paper,
	Progress,
	Stack,
	Text,
} from "@mantine/core";
import { IconCheck } from "@tabler/icons-react";
import type { Stage } from "../types";
import { FileChip } from "./FileChip";

interface ParsingStatementProps {
	file: File | null;
	onCancel: () => void;
	stage: Stage | null;
}

const STEPS: { stage: Stage; label: string }[] = [
	{ stage: "uploaded", label: "Uploaded" },
	{ stage: "validating", label: "Validating PDF" },
	{ stage: "parsing", label: "Reading statement with Claude" },
	{ stage: "checking_balances", label: "Checking balances" },
];

export function ParsingStatement({
	file,
	stage,
	onCancel,
}: ParsingStatementProps) {
	const currentIndex = STEPS.findIndex((s) => s.stage === stage);

	return (
		<Paper withBorder p="lg" radius="md">
			<Stack gap="lg">
				<FileChip file={file} />
				<Stack align="center" gap="md" py="md">
					<Group gap="xs">
						<Loader size="sm" color="dark" />
						<Text size="sm" fw={500}>
							Extracting transactions
						</Text>
					</Group>
					<Progress
						value={100}
						striped
						animated
						color="dark"
						w="100%"
						radius="xl"
					/>
				</Stack>

				<Stack gap={6}>
					{STEPS.map((step, i) => {
						const status =
							i < currentIndex
								? "done"
								: i === currentIndex
									? "active"
									: "pending";
						return (
							<Group gap="xs" key={step.stage}>
								{status === "done" ? (
									<IconCheck size={16} color="#1c1c1c" />
								) : status === "active" ? (
									<Loader size={14} color="dark" />
								) : (
									<Box
										w={16}
										h={16}
										style={{ border: "2px solid #d8d4cc", borderRadius: "50%" }}
									/>
								)}
								<Text size="sm" c={status === "pending" ? "dimmed" : undefined}>
									{step.label}
								</Text>
							</Group>
						);
					})}
				</Stack>

				<Group justify="flex-end">
					<Button variant="default" onClick={onCancel}>
						Cancel
					</Button>
				</Group>
			</Stack>
		</Paper>
	);
}
