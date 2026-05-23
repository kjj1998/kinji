import { useQuery } from "@tanstack/react-query";

import { fetchAllTransactions, fetchSummary } from "../services";

export function useTransactions(userId: string) {
	return useQuery({
		queryKey: ["transactions", userId],
		queryFn: () => fetchAllTransactions(userId),
	});
}

export function useSummary(userId: string, month?: string, year?: string) {
	return useQuery({
		queryKey: ["summary", userId, month, year],
		queryFn: () => fetchSummary(userId, month, year),
	});
}
