import { Box, Card, Group, Text, Title } from "@mantine/core";

export interface MonthlyDataPoint {
	month: string;
	amount: number;
	color?: string;
}

export interface MonthlyTrendProps {
	data: MonthlyDataPoint[];
	currency?: string;
	color?: string;
}

export function MonthlyTrend({
	data,
	currency = "$",
	color = "#D4A853",
}: MonthlyTrendProps) {
	const max = Math.max(...data.map((d) => d.amount)) || 1;

	return (
		<Card withBorder radius="md" p="md" mt="md">
			<Title order={5} mb="md">
				Monthly Trend
			</Title>
			<Group align="flex-end" justify="space-between" h={150} gap="xs">
				{data.map((point) => (
					<Box
						key={point.month}
						style={{
							flex: 1,
							display: "flex",
							flexDirection: "column",
							alignItems: "center",
							gap: 6,
						}}
					>
						<Text size="xs" c="dimmed">
							{currency}
							{(point.amount / 1000).toFixed(1)}k
						</Text>
						<Box
							style={{
								width: "100%",
								height: `${(point.amount / max) * 100}px`,
								backgroundColor: point.color ?? color,
								borderRadius: "3px 3px 0 0",
								transition: "height 0.4s ease",
							}}
						/>
						<Text size="xs" c="dimmed">
							{point.month}
						</Text>
					</Box>
				))}
			</Group>
		</Card>
	);
}
