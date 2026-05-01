import "./App.css";
import { AppShell } from "@mantine/core";
import { Route, Routes } from "react-router-dom";
import { Navbar } from "./components/Navbar";
import { Overview } from "./pages/Overview";
import { Transactions } from "./pages/Transactions";

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
				<Routes>
					<Route path="/" element={<Overview userName="James" />} />
					<Route path="/overview" element={<Overview userName="James" />} />
					<Route path="/transactions" element={<Transactions />} />
				</Routes>
			</AppShell.Main>
		</AppShell>
	);
}

export default App;
