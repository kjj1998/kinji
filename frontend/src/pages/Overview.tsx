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

interface OverviewProps {
	userName: string;
}

const categories = [
	{ name: "Food & Dining", amount: 820, color: "#D4A853" },
	{ name: "Shopping", amount: 430, color: "#A8C5DA" },
	{ name: "Transport", amount: 310, color: "#B8D4A8" },
	{ name: "Subscriptions", amount: 145, color: "#C4A8D4" },
	{ name: "Entertainment", amount: 95, color: "#D4B8A8" },
];

const merchants = [
	{ name: "FairPrice", amount: 312.4, category: "Groceries" },
	{ name: "Grab", amount: 224.8, category: "Transport" },
	{ name: "Koufu", amount: 186.0, category: "Food" },
	{ name: "Uniqlo", amount: 149.9, category: "Shopping" },
	{ name: "Netflix", amount: 15.98, category: "Subscriptions" },
];

const trend = [
	{ label: "Nov", amount: 2900, color: "#A8C5DA" },
	{ label: "Dec", amount: 4100, color: "#A8C5DA" },
	{ label: "Jan", amount: 3200, color: "#A8C5DA" },
	{ label: "Feb", amount: 2750, color: "#A8C5DA" },
	{ label: "Mar", amount: 3600, color: "#A8C5DA" },
	{ label: "Apr", amount: 3470, color: "#D4A853" },
];

const transactions = [
	{ merchant: "Koufu", date: "30 Apr", amount: -8.5, category: "Food" },
	{ merchant: "Grab", date: "30 Apr", amount: -14.2, category: "Transport" },
	{
		merchant: "FairPrice",
		date: "29 Apr",
		amount: -63.8,
		category: "Groceries",
	},
	{
		merchant: "Netflix",
		date: "28 Apr",
		amount: -15.98,
		category: "Subscriptions",
	},
	{ merchant: "Salary", date: "25 Apr", amount: 4900.0, category: "Income" },
];

const changes = [
	{ category: "Entertainment", current: 95, previous: 180 },
	{ category: "Shopping", current: 430, previous: 320 },
	{ category: "Transport", current: 310, previous: 378 },
	{ category: "Food & Dining", current: 820, previous: 790 },
];

const dayOfWeek = [
	{ label: "Mon", amount: 45, color: "#A8C5DA" },
	{ label: "Tue", amount: 30, color: "#A8C5DA" },
	{ label: "Wed", amount: 80, color: "#A8C5DA" },
	{ label: "Thu", amount: 55, color: "#A8C5DA" },
	{ label: "Fri", amount: 120, color: "#A8C5DA" },
	{ label: "Sat", amount: 210, color: "#D4A853" },
	{ label: "Sun", amount: 95, color: "#A8C5DA" },
];

export function Overview({ userName }: OverviewProps) {
	return (
		<>
			<Header text={`Good morning, ${userName}`} />
			<SimpleGrid cols={4} mt="md">
				<SummaryCard label="Total Income" value={4900} />
				<SummaryCard label="Total Spent" value={347} delta={2.5} />
				<SummaryCard label="Net Savings" value={1345.2} delta={8} invertDelta />
				<SummaryCard
					label="Savings Rate"
					value={26}
					format={"percent"}
					delta={5}
					invertDelta
				/>
			</SimpleGrid>
			<Box mt="xs">
				<MonthSummary
					currentSpend={3470}
					lastMonthSpend={3900}
					topCategory="Food & Dining"
					topCategoryAmount={820}
					netSavings={1352}
					savingsRate={26}
				/>
			</Box>
			<Grid mt="xs">
				<Grid.Col span={7}>
					<Stack gap="xs">
						<SpendingByCategory categories={categories} />
						<SpendingBarChart
							title="Monthly Trend"
							data={trend}
							height={160}
							formatAmount={(v) => `${(v / 1000).toFixed(1)}k`} // TODO: fix the formatting when values are less than a thousand
						/>
						<SpendingBarChart
							title="Spending by Day of Week"
							data={dayOfWeek}
						/>
					</Stack>
				</Grid.Col>
				<Grid.Col span={5}>
					<Stack gap="xs">
						<BiggestChanges changes={changes} topN={3} />
						<TopMerchants merchants={merchants} />
						<RecentTransactions transactions={transactions} />
					</Stack>
				</Grid.Col>
			</Grid>
		</>
	);
}
