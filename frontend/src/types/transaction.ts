export const CATEGORIES = [
	"Entertainment",
	"Food",
	"Groceries",
	"Health",
	"Income",
	"Shopping",
	"Subscriptions",
	"Transport",
	"Utilities",
	"Credit",
] as const;

export type Category = (typeof CATEGORIES)[number];

export interface Transaction {
	userId: string;
	id?: string | null;
	date: string;
	merchant: string;
	category: Category;
	amount: number;
	direction: string;
	notes: string;
	split: number | null;
}

export interface Period {
	year: number;
	months: number[];
}

export type Stage = "uploaded" | "validating" | "parsing" | "checking_balances";
