export interface Transaction {
	id: number;
	date: string;
	merchant: string;
	category: string;
	amount: number;
	notes: string;
	split: number | null;
}
