import { Box, Grid, Group, SimpleGrid, Stack, Text } from "@mantine/core";
import {
	BiggestChanges,
	MonthlyTrend,
	MonthSummary,
	RecentTransactions,
	SpendingByCategory,
	SpendingByDayOfWeek,
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
	{ month: "Nov", amount: 2900, color: "#A8C5DA" },
	{ month: "Dec", amount: 4100, color: "#A8C5DA" },
	{ month: "Jan", amount: 3200, color: "#A8C5DA" },
	{ month: "Feb", amount: 2750, color: "#A8C5DA" },
	{ month: "Mar", amount: 3600, color: "#A8C5DA" },
	{ month: "Apr", amount: 3470, color: "#D4A853" },
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
	{ day: "Mon", amount: 45 },
	{ day: "Tue", amount: 30 },
	{ day: "Wed", amount: 80 },
	{ day: "Thu", amount: 55 },
	{ day: "Fri", amount: 120 },
	{ day: "Sat", amount: 210 },
	{ day: "Sun", amount: 95 },
];

export function Overview({ userName }: OverviewProps) {
	return (
		<>
			<Group h={60} px="md" style={{ borderBottom: "1px solid #D4A853" }}>
				<Text fw={400} size="xl">
					Good morning, {userName}
				</Text>
			</Group>
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
						<MonthlyTrend data={trend} />
						<SpendingByDayOfWeek data={dayOfWeek} />
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
