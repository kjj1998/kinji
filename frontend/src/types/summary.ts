import type { Transaction } from "./transaction";

export interface ValueAndChange {
	value: number;
	change: number;
}

export interface CategorySpending {
	category: string;
	value: number;
}

export interface DateSpending {
	dayOrMonth: string;
	value: number;
}

export interface CategorySpendingChange {
	categorySpending: CategorySpending;
	valueChange: number;
	percentageChange: number;
}

export interface Merchant {
	name: string;
	category: string;
	value: number;
}

export interface Summary {
	totalIncome: number;
	totalSpent: ValueAndChange;
	netSavings: ValueAndChange;
	savingsRate: ValueAndChange;
	monthlySummary: string;
	spendingByCategory: CategorySpending[];
	monthlyTrend: DateSpending[];
	spendingByDayOfWeek: DateSpending[];
	biggestChanges: CategorySpendingChange[];
	topMerchants: Merchant[];
	recentTransactions: Transaction[];
}
