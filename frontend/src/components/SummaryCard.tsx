import { Badge, Card, Group, Text, Title } from "@mantine/core";

export interface SummaryCardProps {
	label: string;
	value: number;
	currency?: string;
	delta?: number;
	format?: "currency" | "percent";
	color?: string;
	invertDelta?: boolean;
}

export function SummaryCard({
	label,
	value,
	currency = "$",
	delta,
	format = "currency",
	color,
	invertDelta = false,
}: SummaryCardProps) {
	const isPositive = delta != null && delta > 0;
	const badgeColor = invertDelta
		? isPositive
			? "green"
			: "red"
		: isPositive
			? "red"
			: "green";

	const formattedValue =
		format === "currency" ? `${currency}${value / 100}` : `${value}%`;

	return (
		<Card withBorder radius="md" p="md">
			<Text size="sm" c="dimmed">
				{label}
			</Text>
			<Group justify="space-between" mt="xs">
				<Title order={2} c={color}>
					{formattedValue}
				</Title>
				{delta != null && delta !== 0 && (
					<Badge color={badgeColor}>
						{isPositive ? "↑" : "↓"} {Math.abs(delta).toFixed(1)}%
					</Badge>
				)}
			</Group>
		</Card>
	);
}
