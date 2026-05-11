import { Box, Grid, SimpleGrid, Stack } from "@mantine/core";
import {
	BiggestChanges,
	Header,
	MonthSummary,
	RecentTransactions,
	SpendingBarChart,
	SpendingByCategory,
	SummaryCard,
	TopMerchants,
} from "../components";
import { useSummary } from "../hooks";
import type { DateSpending, Summary } from "../types";

interface OverviewProps {
	userName: string;
}

function formatSpendingBarChartData(
	spendings: DateSpending[],
	isDayOfWeek: boolean = false,
) {
	const data = spendings.map((spending) => {
		return { label: spending.date, amount: spending.amount, color: "#A8C5DA" };
	});

	if (isDayOfWeek) {
		const max = Math.max(...data.map((d) => d.amount));
		return data.map((d) => ({
			...d,
			color: d.amount === max ? "#D4A853" : "#A8C5DA",
		}));
	} else {
		return data.map((d, i) => ({
			...d,
			color: i === data.length - 1 ? "#D4A853" : d.color,
		}));
	}
}

export function Overview({ userName }: OverviewProps) {
	const { data: summary = {} as Summary, isLoading } = useSummary(
		"james",
		"2026-01-01",
		"2026-04-30",
	);

	if (isLoading) {
		return <>Loading...</>;
	}

	const formattedMonthlyTrendData = formatSpendingBarChartData(
		summary.monthlyTrend,
	);
	const formattedDailyTrendData = formatSpendingBarChartData(
		summary.dailyTrend,
		true,
	);

	return (
		<>
			<Header text={`Good morning, ${userName}`} />
			<SimpleGrid cols={4} mt="md">
				<SummaryCard label="Total Income" value={summary.totalIncome} />
				<SummaryCard
					label="Total Spent"
					value={summary.totalSpent.value}
					delta={summary.totalSpent.change}
				/>
				<SummaryCard
					label="Net Savings"
					value={summary.netSavings.value}
					delta={summary.netSavings.change}
					invertDelta
				/>
				<SummaryCard
					label="Savings Rate"
					value={summary.savingsRate.value}
					format={"percent"}
					delta={summary.savingsRate.change}
					invertDelta
				/>
			</SimpleGrid>
			<Box mt="xs">
				<MonthSummary
					currentSpend={summary.totalSpent.value}
					lastMonthSpend={summary.lastMonthSpent}
					topCategory={summary.topCategory.category}
					topCategoryAmount={summary.topCategory.amount}
					netSavings={summary.netSavings.value}
					savingsRate={summary.savingsRate.value}
				/>
			</Box>
			<Grid mt="xs">
				<Grid.Col span={7}>
					<Stack gap="xs">
						<SpendingByCategory
							spendingByCategory={summary.spendingByCategory}
						/>
						<SpendingBarChart
							title="Monthly Trend"
							data={formattedMonthlyTrendData}
						/>
						<SpendingBarChart
							title="Daily Trend"
							data={formattedDailyTrendData}
						/>
					</Stack>
				</Grid.Col>
				<Grid.Col span={5}>
					<Stack gap="xs">
						<BiggestChanges changes={summary.biggestChanges} />
						<TopMerchants merchants={summary.topMerchants} />
						<RecentTransactions transactions={summary.recentTransactions} />
					</Stack>
				</Grid.Col>
			</Grid>
		</>
	);
}
