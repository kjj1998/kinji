import { Box, Card, Group, Stack, Text, Title } from "@mantine/core";

export interface CategoryItem {
	name: string;
	amount: number;
	color?: string;
}

export interface SpendingByCategoryProps {
	categories: CategoryItem[];
	currency?: string;
}

export function SpendingByCategory({
	categories,
	currency = "$",
}: SpendingByCategoryProps) {
	const sorted = [...categories].sort((a, b) => b.amount - a.amount);
	const max = sorted.reduce((sum, cat) => sum + cat.amount, 0) || 1;

	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb="xs">
				Spending by Category
			</Title>
			<Stack gap="xs">
				{sorted.map((cat) => (
					<Box key={cat.name}>
						<Group justify="space-between" mb={4}>
							<Text size="sm">{cat.name}</Text>
							<Text size="sm" fw={500}>
								{currency}
								{cat.amount.toLocaleString("en-US", {
									minimumFractionDigits: 2,
									maximumFractionDigits: 2,
								})}
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
									width: `${(cat.amount / max) * 100}%`,
									borderRadius: 3,
									backgroundColor: cat.color ?? "#D4A853",
									transition: "width 0.4s ease",
								}}
							/>
						</Box>
					</Box>
				))}
			</Stack>
		</Card>
	);
}
