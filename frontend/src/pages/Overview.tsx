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
import type { DateSpending } from "../types";

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
	const now = new Date();
	const month = String(now.getMonth() + 1).padStart(2, "0");
	const year = String(now.getFullYear());
	const {
		data: summary,
		isLoading,
		isError,
	} = useSummary("gomez", month, year);

	if (isLoading) {
		return <>Loading...</>;
	}

	if (isError || !summary) {
		return <>No Data</>;
	}

	const formattedMonthlyTrendData = formatSpendingBarChartData(
		summary.monthlyExpenses,
	);
	const formattedDailyTrendData = formatSpendingBarChartData(
		summary.dailyTrend,
		true,
	);

	return (
		<>
			<Header text={`Good morning, ${userName}`} />
			<SimpleGrid cols={4} mt="md">
				<SummaryCard label="Total Income" value={summary.totalIncome.value} />
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
					value={summary.savingsRate}
					format={"percent"}
					invertDelta
				/>
			</SimpleGrid>
			<Box mt="xs">
				<MonthSummary monthlySummary={summary.monthlySummary} />
			</Box>
			<Grid mt="xs">
				<Grid.Col span={7}>
					<Stack gap="xs">
						<SpendingByCategory spendingByCategory={summary.topCategories} />
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
