import "./App.css";
import { AppShell, Group, SimpleGrid, Text } from "@mantine/core";
import { Navbar } from "./components/Navbar";
import { SummaryCard } from "./components/SummaryCard";

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
			</AppShell.Main>
		</AppShell>
	);
}

export default App;
