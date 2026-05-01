import "./App.css";
import { AppShell, Group, SimpleGrid, Text } from "@mantine/core";
import { Navbar } from "./components/Navbar";
import { SpendingByCategory } from "./components/SpendingByCategory";
import { SummaryCard } from "./components/SummaryCard";

const categories = [
	{ name: "Food & Dining", amount: 820, color: "#D4A853" },
	{ name: "Shopping", amount: 430, color: "#A8C5DA" },
	{ name: "Transport", amount: 310, color: "#B8D4A8" },
	{ name: "Subscriptions", amount: 145, color: "#C4A8D4" },
	{ name: "Entertainment", amount: 95, color: "#D4B8A8" },
];

function App() {
	return (
		<AppShell navbar={{ width: 250, breakpoint: "sm" }} padding="md">
			<Navbar
				current="overview"
				onNavigate={() => {}}
				onUpload={() => {}}
				monthStatus={{
					label: "April 2026",
					uploaded: 3,
					expected: 5,
					detail: "2 statements still missing",
				}}
				statementsMissing={2}
				user={{
					name: "James",
					email: "james@example.com",
				}}
			/>
			<AppShell.Main>
				<Group h={60} px="md" style={{ borderBottom: "1px solid #D4A853" }}>
					<Text fw={400} size="xl">
						Good morning, James
					</Text>
				</Group>
				<SimpleGrid cols={4} mt="md">
					<SummaryCard label="Total Income" value={4900} />
					<SummaryCard label="Total Spent" value={347} delta={2.5} />
					<SummaryCard
						label="Net Savings"
						value={1345.2}
						delta={8}
						invertDelta
					/>
					<SummaryCard
						label="Savings Rate"
						value={26}
						format={"percent"}
						delta={5}
						invertDelta
					/>
				</SimpleGrid>
				<SpendingByCategory categories={categories} />
			</AppShell.Main>
		</AppShell>
	);
}

export default App;
