import { Card, Text, Title } from "@mantine/core";

export interface MonthSummaryProps {
	monthlySummary: string;
}

export function MonthSummary({ monthlySummary }: MonthSummaryProps) {
	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb={6}>
				Monthly Summary
			</Title>
			<Text size="sm" c="dimmed" lh={1.6}>
				{monthlySummary}
			</Text>
		</Card>
	);
}
