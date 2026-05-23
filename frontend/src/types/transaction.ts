export interface Transaction {
	id: number;
	date: string;
	merchant: string;
	category: string;
	amount: number;
	direction: string;
	notes: string;
	split: number | null;
}

export interface TransactionsAvailability {
	year: number;
	months: number[];
}

export interface Transactions {
	transactions: Transaction[];
	availabilities: TransactionsAvailability[];
}
