export interface Transaction {
	id: number;
	date: string;
	merchant: string;
	category: string;
	amount: number;
	notes: string;
	split: number | null;
}

export const transactions: Transaction[] = [
	{ id: 1,  date: "2026-04-30", merchant: "Koufu",            category: "Food",          amount: -8.50,   notes: "",                  split: null },
	{ id: 2,  date: "2026-04-30", merchant: "Grab",             category: "Transport",     amount: -14.20,  notes: "To airport",        split: null },
	{ id: 3,  date: "2026-04-29", merchant: "FairPrice",        category: "Groceries",     amount: -63.80,  notes: "",                  split: 31.90 },
	{ id: 4,  date: "2026-04-28", merchant: "Netflix",          category: "Subscriptions", amount: -15.98,  notes: "Family plan",       split: 5.33 },
	{ id: 5,  date: "2026-04-25", merchant: "Employer",         category: "Income",        amount: 4900.00, notes: "April salary",      split: null },
	{ id: 6,  date: "2026-04-24", merchant: "Uniqlo",           category: "Shopping",      amount: -149.90, notes: "Birthday gift",     split: null },
	{ id: 7,  date: "2026-04-23", merchant: "Grab",             category: "Transport",     amount: -9.40,   notes: "",                  split: null },
	{ id: 8,  date: "2026-04-22", merchant: "Koufu",            category: "Food",          amount: -7.20,   notes: "",                  split: null },
	{ id: 9,  date: "2026-04-21", merchant: "Spotify",          category: "Subscriptions", amount: -9.99,   notes: "",                  split: null },
	{ id: 10, date: "2026-04-20", merchant: "Wingstop",         category: "Food",          amount: -38.40,  notes: "Lunch with team",   split: 12.80 },
	{ id: 11, date: "2026-04-19", merchant: "SMRT",             category: "Transport",     amount: -48.00,  notes: "Monthly top-up",    split: null },
	{ id: 12, date: "2026-04-18", merchant: "Guardian",         category: "Health",        amount: -23.50,  notes: "",                  split: null },
	{ id: 13, date: "2026-04-17", merchant: "FairPrice",        category: "Groceries",     amount: -41.20,  notes: "",                  split: null },
	{ id: 14, date: "2026-04-16", merchant: "Grab",             category: "Transport",     amount: -11.60,  notes: "",                  split: null },
	{ id: 15, date: "2026-04-15", merchant: "Singtel",          category: "Utilities",     amount: -45.00,  notes: "Mobile bill",       split: null },
	{ id: 16, date: "2026-04-14", merchant: "Nando's",          category: "Food",          amount: -52.70,  notes: "Dinner with family",split: 26.35 },
	{ id: 17, date: "2026-04-13", merchant: "Amazon",           category: "Shopping",      amount: -34.90,  notes: "",                  split: null },
	{ id: 18, date: "2026-04-12", merchant: "Cold Storage",     category: "Groceries",     amount: -29.60,  notes: "",                  split: null },
	{ id: 19, date: "2026-04-11", merchant: "Grab",             category: "Transport",     amount: -16.80,  notes: "",                  split: null },
	{ id: 20, date: "2026-04-10", merchant: "Starbucks",        category: "Food",          amount: -12.50,  notes: "Client meeting",    split: null },
	{ id: 21, date: "2026-04-09", merchant: "SP Group",         category: "Utilities",     amount: -87.30,  notes: "Electricity & water",split: null },
	{ id: 22, date: "2026-04-08", merchant: "Cathay Cineplexes",category: "Entertainment", amount: -28.00,  notes: "",                  split: 14.00 },
	{ id: 23, date: "2026-04-07", merchant: "Koufu",            category: "Food",          amount: -6.80,   notes: "",                  split: null },
	{ id: 24, date: "2026-04-06", merchant: "Shopee",           category: "Shopping",      amount: -56.40,  notes: "",                  split: null },
	{ id: 25, date: "2026-04-05", merchant: "FairPrice",        category: "Groceries",     amount: -38.90,  notes: "",                  split: null },
	{ id: 26, date: "2026-04-04", merchant: "Grab",             category: "Transport",     amount: -13.20,  notes: "",                  split: null },
	{ id: 27, date: "2026-04-03", merchant: "Freelance",        category: "Income",        amount: 800.00,  notes: "Design project",    split: null },
	{ id: 28, date: "2026-04-02", merchant: "Watsons",          category: "Health",        amount: -18.70,  notes: "",                  split: null },
	{ id: 29, date: "2026-04-01", merchant: "McDonald's",       category: "Food",          amount: -15.30,  notes: "",                  split: null },
	{ id: 30, date: "2026-03-31", merchant: "Koufu",            category: "Food",          amount: -9.10,   notes: "",                  split: null },
	{ id: 31, date: "2026-03-30", merchant: "Grab",             category: "Transport",     amount: -22.40,  notes: "Late night ride",   split: null },
	{ id: 32, date: "2026-03-28", merchant: "FairPrice",        category: "Groceries",     amount: -72.50,  notes: "",                  split: 36.25 },
	{ id: 33, date: "2026-03-26", merchant: "Disney+",          category: "Subscriptions", amount: -11.98,  notes: "",                  split: null },
	{ id: 34, date: "2026-03-25", merchant: "Employer",         category: "Income",        amount: 4900.00, notes: "March salary",      split: null },
	{ id: 35, date: "2026-03-24", merchant: "IKEA",             category: "Shopping",      amount: -234.00, notes: "New desk lamp",     split: null },
	{ id: 36, date: "2026-03-22", merchant: "Grab",             category: "Transport",     amount: -8.90,   notes: "",                  split: null },
	{ id: 37, date: "2026-03-20", merchant: "Ya Kun",           category: "Food",          amount: -5.50,   notes: "",                  split: null },
	{ id: 38, date: "2026-03-18", merchant: "Singtel",          category: "Utilities",     amount: -45.00,  notes: "Mobile bill",       split: null },
	{ id: 39, date: "2026-03-15", merchant: "Decathlon",        category: "Shopping",      amount: -89.90,  notes: "Running shoes",     split: null },
	{ id: 40, date: "2026-03-12", merchant: "Cold Storage",     category: "Groceries",     amount: -44.30,  notes: "",                  split: null },
	{ id: 41, date: "2026-03-10", merchant: "Grab",             category: "Transport",     amount: -17.60,  notes: "",                  split: null },
	{ id: 42, date: "2026-03-08", merchant: "PS Cafe",          category: "Food",          amount: -96.40,  notes: "Anniversary dinner",split: 48.20 },
	{ id: 43, date: "2026-03-06", merchant: "SP Group",         category: "Utilities",     amount: -81.10,  notes: "Electricity & water",split: null },
	{ id: 44, date: "2026-03-04", merchant: "Shopee",           category: "Shopping",      amount: -23.80,  notes: "",                  split: null },
	{ id: 45, date: "2026-03-02", merchant: "FairPrice",        category: "Groceries",     amount: -55.60,  notes: "",                  split: null },
];
