import "./App.css";
import { AppShell } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useCallback } from "react";
import {
	Navigate,
	Route,
	Routes,
	useLocation,
	useNavigate,
} from "react-router-dom";
import type { NavbarItemKey } from "./components/Navbar";
import { Navbar } from "./components/Navbar";
import { Upload } from "./pages";
import { Overview } from "./pages/Overview";
import { Transactions } from "./pages/Transactions";

const queryClient = new QueryClient();

const pathToKey: Record<string, NavbarItemKey> = {
	"/overview": "overview",
	"/transactions": "transactions",
	"/statements": "statements",
};

const keyToPath: Record<NavbarItemKey, string> = {
	overview: "/overview",
	transactions: "/transactions",
	statements: "/statements",
	upload: "/upload",
};

function App() {
	const navigate = useNavigate();
	const location = useLocation();
	const current = pathToKey[location.pathname] ?? "overview";
	const handleNavigate = useCallback(
		(key: NavbarItemKey) => navigate(keyToPath[key]),
		[navigate],
	);

	return (
		<QueryClientProvider client={queryClient}>
			<AppShell navbar={{ width: 250, breakpoint: "sm" }} padding="md">
				<Navbar
					current={current}
					onNavigate={handleNavigate}
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
						<Route path="/" element={<Navigate to="/overview" replace />} />
						<Route path="/overview" element={<Overview userName="James" />} />
						<Route path="/transactions" element={<Transactions />} />
						<Route path="/upload" element={<Upload userId="james" />} />
					</Routes>
				</AppShell.Main>
			</AppShell>
		</QueryClientProvider>
	);
}

export default App;
