import type { Period, Stage, Transaction } from "../types";

interface ImportHandlers {
	onProgress: (stage: Stage) => void;
	onDone: (transactions: Transaction[]) => void;
	onError: (message: string) => void;
}

const BASE_URL = import.meta.env.VITE_API_URL;

export async function fetchAllTransactions(
	userId: string,
	month?: string,
	year?: string,
) {
	const params = new URLSearchParams();
	if (month) params.set("month", month);
	if (year) params.set("year", year);

	const query = params.size > 0 ? `?${params}` : "";
	const url = `${BASE_URL}/api/v1/transactions/${userId}${query}`;

	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`Response status: ${response.status}`);
	}

	const result = await response.json();

	return result as Transaction[];
}

export async function fetchSummary(
	userId: string,
	month?: string,
	year?: string,
) {
	const params = new URLSearchParams();
	if (month) params.set("month", month);
	if (year) params.set("year", year);

	const query = params.size > 0 ? `?${params}` : "";
	const url = `${BASE_URL}/api/v1/summary/${userId}${query}`;

	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`Response status: ${response.status}`);
	}

	return response.json();
}

export async function saveTransactions(
	userId: string,
	transactions: Transaction[],
) {
	const url = `${BASE_URL}/api/v1/transactions/${userId}`;
	const response = await fetch(url, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(transactions),
	});

	if (!response.ok) {
		throw new Error(`Response status: ${response.status}`);
	}

	return response.json();
}

export async function importStatement(
	file: File,
	userId: string,
	password: string,
	handlers: ImportHandlers,
	signal?: AbortSignal,
) {
	const form = new FormData();
	form.append("statement", file);
	form.append("password", password);

	const url = `${BASE_URL}/api/v1/transactions/${userId}/import`;
	const response = await fetch(url, { method: "POST", body: form, signal });

	if (!response.ok) {
		let message = `Response status: ${response.status}`;
		try {
			message = (await response.json()).message ?? message;
		} catch {}
		handlers.onError(message);
		return;
	}

	if (!response.body) {
		handlers.onError("no response body");
		return;
	}
	const reader = response.body.getReader();
	const decoder = new TextDecoder();
	let buffer = "";

	while (true) {
		const { value, done } = await reader.read();
		if (done) break;
		buffer += decoder.decode(value, { stream: true });

		const frames = buffer.split("\n\n");
		buffer = frames.pop() ?? "";

		for (const frame of frames) {
			if (!frame.trim() || frame.startsWith(":")) continue; // skip blanks + heartbeats
			let event = "message";
			let data = "";
			for (const line of frame.split("\n")) {
				if (line.startsWith("event:")) event = line.slice(6).trim();
				else if (line.startsWith("data:")) data += line.slice(5).trim();
			}

			if (event === "progress") {
				handlers.onProgress(JSON.parse(data).stage);
			} else if (event === "done") {
				handlers.onDone(JSON.parse(data));
				return;
			} else if (event === "error") {
				handlers.onError(JSON.parse(data).message);
			}
		}
	}
}

export async function fetchPeriods(userId: string) {
	const url = `${BASE_URL}/api/v1/transactions/${userId}/periods`;
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`Response status: ${response.status}`);
	}

	const result = await response.json();

	return result as Period[];
}
