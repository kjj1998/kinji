import { Card, Divider, Group, Stack, Text, Title } from "@mantine/core";

export interface TransactionItem {
	merchant: string;
	date: string;
	amount: number;
	category: string;
}

export interface RecentTransactionsProps {
	transactions: TransactionItem[];
	currency?: string;
}

export function RecentTransactions({
	transactions,
	currency = "$",
}: RecentTransactionsProps) {
	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb="xs">
				Recent Transactions
			</Title>
			<Stack gap={0}>
				{transactions.map((tx, index) => {
					const isIncome = tx.amount > 0;
					return (
						<div key={`${tx.merchant}-${tx.date}`}>
							<Group justify="space-between" py={4} wrap="nowrap">
								<div>
									<Text size="sm" fw={500}>
										{tx.merchant}
									</Text>
									<Text size="xs" c="dimmed">
										{tx.date} · {tx.category}
									</Text>
								</div>
								<Text
									size="sm"
									fw={500}
									c={isIncome ? "green" : undefined}
									style={{ whiteSpace: "nowrap" }}
								>
									{isIncome ? "+" : ""}
									{currency}
									{Math.abs(tx.amount).toLocaleString("en-US", {
										minimumFractionDigits: 2,
										maximumFractionDigits: 2,
									})}
								</Text>
							</Group>
							{index < transactions.length - 1 && <Divider />}
						</div>
					);
				})}
			</Stack>
		</Card>
	);
}
