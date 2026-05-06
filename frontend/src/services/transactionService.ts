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
