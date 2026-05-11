import { Card, Text, Title } from "@mantine/core";
import { calculateDelta, formatCurrency } from "../utils";

export interface MonthSummaryProps {
	currentSpend: number;
	lastMonthSpend: number;
	topCategory: string;
	topCategoryAmount: number;
	netSavings: number;
	savingsRate: number;
}

export function MonthSummary({
	currentSpend,
	lastMonthSpend,
	topCategory,
	topCategoryAmount,
	netSavings,
	savingsRate,
}: MonthSummaryProps) {
	const delta =
		lastMonthSpend === 0
			? null
			: calculateDelta({ current: currentSpend, previous: lastMonthSpend });

	const spendSentence =
		delta == null
			? null
			: delta > 0
				? `You spent ${Math.abs(delta).toFixed(0)}% more than last month.`
				: `You spent ${Math.abs(delta).toFixed(0)}% less than last month.`;

	const breakdownSentence = `Your biggest expense was ${topCategory} at ${formatCurrency(topCategoryAmount)}, and you saved ${formatCurrency(netSavings)} (${savingsRate.toFixed(0)}% of income).`;

	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb={6}>
				Monthly Summary
			</Title>
			<Text size="sm" c="dimmed" lh={1.6}>
				{spendSentence} {breakdownSentence}
			</Text>
		</Card>
	);
}
