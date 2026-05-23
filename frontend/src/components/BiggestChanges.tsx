import { Badge, Card, Divider, Group, Stack, Text, Title } from "@mantine/core";
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
				{changes.length === 0 ? (
					<Text size="sm" c="dimmed" ta="center" py="md">
						No changes to show yet
					</Text>
				) : (
					changes.map((item, index) => {
						const isNew = item.isNew;
						const isUp = item.change > 0;
						const noChange = item.percentageChange === 0;
						const color = noChange ? undefined : isUp ? "red" : "green";
						return (
							<div key={item.category}>
								<Group justify="space-between" py={4} wrap="nowrap">
									<Text size="sm" fw={500}>
										{item.category}
									</Text>
									<Group gap="xs" wrap="nowrap">
										{isNew ? (
											<Badge size="sm" variant="light" color="gray">
												New
											</Badge>
										) : (
											<Text size="sm" c={color} fw={500}>
												{noChange ? "" : isUp ? "▲" : "▼"}{" "}
												{Math.abs(item.percentageChange)}%
											</Text>
										)}
										<Text size="sm" c={isNew ? undefined : color} fw={500}>
											{isNew || noChange ? "" : isUp ? "+" : "-"}
											{formatCurrency(item.change / 100)}
										</Text>
									</Group>
								</Group>
								{index < changes.length - 1 && <Divider />}
							</div>
						);
					})
				)}
			</Stack>
		</Card>
	);
}
