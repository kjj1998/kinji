const sgdFormatter = new Intl.NumberFormat("en-SG", { style: "currency", currency: "SGD" });

export function formatCurrency(amount: number): string {
	return sgdFormatter.format(amount);
}
