import { useQuery } from "@tanstack/react-query";

import { fetchAllTransactions } from "../services";

export function useTransactions(userId: string) {
	return useQuery({
		queryKey: ["transactions", userId],
		queryFn: () => fetchAllTransactions(userId),
	});
}
