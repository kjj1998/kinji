import { Card, Divider, Group, Stack, Text, Title } from "@mantine/core";
import type { Transaction } from "../types";
import { formatCurrency } from "../utils";

export interface RecentTransactionsProps {
	transactions: Transaction[];
}

export function RecentTransactions({ transactions }: RecentTransactionsProps) {
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
									{formatCurrency(Math.abs(tx.amount))}
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
