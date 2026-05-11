import type { Transaction } from "./transaction";

export interface ValueAndChange {
	value: number;
	change: number;
}

export interface CategorySpending {
	category: string;
	amount: number;
}

export interface DateSpending {
	date: string;
	amount: number;
}

export interface CategorySpendingChange {
	category: string;
	amount: number;
	change: number;
	percentageChange: number;
}

export interface Merchant {
	name: string;
	category: string;
	amount: number;
}

export interface Summary {
	totalIncome: number;
	totalSpent: ValueAndChange;
	netSavings: ValueAndChange;
	savingsRate: ValueAndChange;
	lastMonthSpent: number;
	topCategory: CategorySpending;
	monthlySummary: string;
	spendingByCategory: CategorySpending[];
	monthlyTrend: DateSpending[];
	dailyTrend: DateSpending[];
	biggestChanges: CategorySpendingChange[];
	topMerchants: Merchant[];
	recentTransactions: Transaction[];
}
