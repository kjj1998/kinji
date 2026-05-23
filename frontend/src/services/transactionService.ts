import type { Transactions } from "../types";

const BASE_URL = import.meta.env.VITE_API_URL;

export async function fetchAllTransactions(userId: string) {
	const url = `${BASE_URL}/api/v1/transactions/${userId}`;

	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`Response status: ${response.status}`);
	}

	const result = await response.json();

	return result as Transactions;
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
