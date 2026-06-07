import "@mantine/core/styles.css";
import "@mantine/dropzone/styles.css";
import "@mantine/notifications/styles.css";
import { MantineProvider } from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import "./index.css";
import App from "./App.tsx";

// biome-ignore lint/style/noNonNullAssertion: ignore
createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<MantineProvider>
			<Notifications position="top-right" autoClose={5000} />
			<BrowserRouter>
				<App />
			</BrowserRouter>
		</MantineProvider>
	</StrictMode>,
);
