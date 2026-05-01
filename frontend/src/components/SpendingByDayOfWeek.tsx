import { Box, Card, Group, Text, Title } from "@mantine/core";

export interface DayOfWeekData {
	day: string;
	amount: number;
}

export interface SpendingByDayOfWeekProps {
	data: DayOfWeekData[];
	currency?: string;
	color?: string;
}

export function SpendingByDayOfWeek({
	data,
	currency = "$",
	color = "#D4A853",
}: SpendingByDayOfWeekProps) {
	const max = Math.max(...data.map((d) => d.amount)) || 1;

	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb="xs">
				Spending by Day of Week
			</Title>
			<Group align="flex-end" justify="space-between" h={120} gap="xs">
				{data.map((point) => (
					<Box
						key={point.day}
						style={{
							flex: 1,
							display: "flex",
							flexDirection: "column",
							alignItems: "center",
							gap: 4,
						}}
					>
						<Text size="xs" c="dimmed">
							{currency}
							{point.amount}
						</Text>
						<Box
							style={{
								width: "100%",
								height: `${(point.amount / max) * 60}px`,
								backgroundColor: color,
								borderRadius: "3px 3px 0 0",
								transition: "height 0.4s ease",
							}}
						/>
						<Text size="xs" c="dimmed">
							{point.day}
						</Text>
					</Box>
				))}
			</Group>
		</Card>
	);
}
