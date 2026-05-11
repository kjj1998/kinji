import { Box, Card, Group, Text, Title } from "@mantine/core";
import { formatCurrency } from "../utils";

export interface BarChartDataPoint {
	label: string;
	amount: number;
	color?: string;
}

export interface BarChartProps {
	title: string;
	data: BarChartDataPoint[];
	height?: number;
}

export function SpendingBarChart({ title, data, height = 160 }: BarChartProps) {
	const max = Math.max(...data.map((d) => d.amount)) || 1;

	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb="xs">
				{title}
			</Title>
			<Group align="flex-end" justify="space-between" h={height} gap="xs">
				{data.map((point) => (
					<Box
						key={point.label}
						style={{
							flex: 1,
							display: "flex",
							flexDirection: "column",
							alignItems: "center",
							gap: 6,
						}}
					>
						<Text size="xs" c="dimmed">
							{formatCurrency(point.amount)}
						</Text>
						<Box
							style={{
								width: "100%",
								height: `${(point.amount / max) * 100}px`,
								backgroundColor: `${point.color}`,
								borderRadius: "3px 3px 0 0",
								transition: "height 0.4s ease",
							}}
						/>
						<Text size="xs" c="dimmed">
							{point.label}
						</Text>
					</Box>
				))}
			</Group>
		</Card>
	);
}
