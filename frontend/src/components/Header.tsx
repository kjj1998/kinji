import { Group, Text } from "@mantine/core";

export interface HeaderProps {
	text: string;
}

export function Header({ text }: HeaderProps) {
	return (
		<Group h={60} px="md" style={{ borderBottom: "1px solid #D4A853" }}>
			<Text fw={400} size="xl">
				{text}
			</Text>
		</Group>
	);
}
