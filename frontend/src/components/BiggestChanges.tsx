import { Card, Divider, Group, Stack, Text, Title } from "@mantine/core";
import { calculateDelta } from "../utils";

export interface CategoryChange {
	category: string;
	current: number;
	previous: number;
}

export interface BiggestChangesProps {
	changes: CategoryChange[];
	topN?: number;
	currency?: string;
}

export function BiggestChanges({
	changes,
	topN = 3,
	currency = "$",
}: BiggestChangesProps) {
	const withDelta = changes
		.filter((c) => c.previous > 0)
		.map((c) => ({
			...c,
			delta: calculateDelta({ current: c.current, previous: c.previous }),
			diff: c.current - c.previous,
		}))
		.filter((c): c is typeof c & { delta: number } => c.delta !== null)
		.sort((a, b) => Math.abs(b.delta) - Math.abs(a.delta))
		.slice(0, topN);

	const fmt = (amount: number) =>
		`${currency}${Math.abs(amount).toLocaleString("en-US", {
			minimumFractionDigits: 0,
			maximumFractionDigits: 0,
		})}`;

	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb="xs">
				Biggest Changes vs Last Month
			</Title>
			<Stack gap={0}>
				{withDelta.map((item, index) => {
					const isUp = item.delta > 0;
					return (
						<div key={item.category}>
							<Group justify="space-between" py={4} wrap="nowrap">
								<Text size="sm" fw={500}>
									{item.category}
								</Text>
								<Group gap="xs" wrap="nowrap">
									<Text size="sm" c={isUp ? "red" : "green"} fw={500}>
										{isUp ? "▲" : "▼"} {Math.abs(item.delta).toFixed(0)}%
									</Text>
									<Text size="xs" c="dimmed">
										{fmt(item.previous)} → {fmt(item.current)}
									</Text>
									<Text size="sm" c={isUp ? "red" : "green"} fw={500}>
										{isUp ? "+" : "-"}
										{fmt(item.diff)}
									</Text>
								</Group>
							</Group>
							{index < withDelta.length - 1 && <Divider />}
						</div>
					);
				})}
			</Stack>
		</Card>
	);
}
