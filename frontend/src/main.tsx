import "@mantine/core/styles.css";
import { MantineProvider } from "@mantine/core";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";

// biome-ignore lint/style/noNonNullAssertion: ignore
createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<MantineProvider>
			<App />
		</MantineProvider>
	</StrictMode>,
);
