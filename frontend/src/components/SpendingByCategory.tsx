import { Box, Card, Group, Stack, Text, Title } from "@mantine/core";
import type { CategorySpending } from "../types";
import { formatCurrency } from "../utils";

const COLORS = ["#D4A853", "#A8C5DA", "#B8D4A8", "#C4A8D4", "#D4B8A8"];

export interface SpendingByCategoryProps {
	spendingByCategory: CategorySpending[];
}

export function SpendingByCategory({
	spendingByCategory,
}: SpendingByCategoryProps) {
	const data = (spendingByCategory ?? []).map(
		(c: CategorySpending, i: number) => ({
			name: c.category,
			amount: c.amount,
			color: COLORS[i % COLORS.length],
		}),
	);
	const total = data.reduce((sum, cat) => sum + cat.amount, 0);

	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb="xs">
				Spending by Category
			</Title>
			{data.length === 0 ? (
				<Text ta="center">No data found</Text>
			) : (
				<Stack gap="xs">
					{data.map((cat) => (
						<Box key={cat.name}>
							<Group justify="space-between" mb={4}>
								<Text size="sm">{cat.name}</Text>
								<Text size="sm" fw={500}>
									{formatCurrency(cat.amount)}
								</Text>
							</Group>
							<Box
								style={{
									height: 6,
									borderRadius: 3,
									backgroundColor: "var(--mantine-color-gray-2)",
									overflow: "hidden",
								}}
							>
								<Box
									style={{
										height: "100%",
										width: `${(cat.amount / total) * 100}%`,
										borderRadius: 3,
										backgroundColor: cat.color,
										transition: "width 0.4s ease",
									}}
								/>
							</Box>
						</Box>
					))}
				</Stack>
			)}
		</Card>
	);
}
