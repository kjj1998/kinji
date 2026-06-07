import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
	fetchAllTransactions,
	fetchSummary,
	saveTransactions,
} from "../services";
import type { Transaction } from "../types";

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

export function useSaveTransactions(userId: string) {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: (transactions: Transaction[]) =>
			saveTransactions(userId, transactions),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["transactions", userId] });
			queryClient.invalidateQueries({ queryKey: ["summary", userId] });
		},
	});
}
