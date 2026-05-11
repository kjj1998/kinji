import type { Transaction } from "../types";

const BASE_URL = import.meta.env.VITE_API_URL;

export async function fetchAllTransactions(userId: string) {
	const url = `${BASE_URL}/api/v1/transactions/${userId}`;

	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`Response status: ${response.status}`);
	}

	const result = await response.json();

	return result as Transaction[];
}

export async function fetchSummary(userId: string, from?: string, to?: string) {
	const params = new URLSearchParams();
	if (from) params.set("from", from);
	if (to) params.set("to", to);

	const query = params.size > 0 ? `?${params}` : "";
	const url = `${BASE_URL}/api/v1/summary/${userId}${query}`;

	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`Response status: ${response.status}`);
	}

	return response.json();
}
