import { useQuery } from "@tanstack/react-query";

import { fetchAllTransactions, fetchSummary } from "../services";

export function useTransactions(userId: string) {
	return useQuery({
		queryKey: ["transactions", userId],
		queryFn: () => fetchAllTransactions(userId),
	});
}

export function useSummary(userId: string, from?: string, to?: string) {
	return useQuery({
		queryKey: ["summary", userId, from, to],
		queryFn: () => fetchSummary(userId, from, to),
	});
}
