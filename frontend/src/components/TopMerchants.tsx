import { Avatar, Card, Group, Stack, Text, Title } from "@mantine/core";

import { formatCurrency } from "../utils";

export interface MerchantItem {
	name: string;
	amount: number;
	category: string;
}

export interface TopMerchantsProps {
	merchants: MerchantItem[];
	currency?: string;
}

export function TopMerchants({ merchants }: TopMerchantsProps) {
	const sorted = [...merchants].sort((a, b) => b.amount - a.amount);

	return (
		<Card withBorder radius="md" p="sm">
			<Title order={5} mb="xs">
				Top Merchants
			</Title>
			{sorted.length === 0 ? (
				<Text size="sm" c="dimmed" ta="center" py="md">
					No data found
				</Text>
			) : (
				<Stack gap="xs">
					{sorted.map((merchant, index) => (
						<Group key={merchant.name} justify="space-between" wrap="nowrap">
							<Group gap="sm" wrap="nowrap">
								<Avatar size={28} radius="xl" color="gray">
									{index + 1}
								</Avatar>
								<div>
									<Text size="sm" fw={500} truncate>
										{merchant.name}
									</Text>
									<Text size="xs" c="dimmed">
										{merchant.category}
									</Text>
								</div>
							</Group>
							<Text size="sm" fw={500}>
								{formatCurrency(merchant.amount / 100)}
							</Text>
						</Group>
					))}
				</Stack>
			)}
		</Card>
	);
}
