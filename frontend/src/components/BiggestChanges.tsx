import { Card, Divider, Group, Stack, Text, Title } from "@mantine/core";
import type { CategorySpendingChange } from "../types";
import { formatCurrency } from "../utils";

export interface BiggestChangesProps {
	changes: CategorySpendingChange[];
}

export function BiggestChanges({ changes }: BiggestChangesProps) {
	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb="xs">
				Biggest Changes vs Last Month
			</Title>
			<Stack gap={0}>
				{changes.map((item, index) => {
					const isUp = item.change > 0;
					return (
						<div key={item.category}>
							<Group justify="space-between" py={4} wrap="nowrap">
								<Text size="sm" fw={500}>
									{item.category}
								</Text>
								<Group gap="xs" wrap="nowrap">
									<Text size="sm" c={isUp ? "red" : "green"} fw={500}>
										{isUp ? "▲" : "▼"} {Math.abs(item.percentageChange)}%
									</Text>
									<Text size="sm" c={isUp ? "red" : "green"} fw={500}>
										{isUp ? "+" : "-"}
										{formatCurrency(Math.abs(item.change))}
									</Text>
								</Group>
							</Group>
							{index < changes.length - 1 && <Divider />}
						</div>
					);
				})}
			</Stack>
		</Card>
	);
}
