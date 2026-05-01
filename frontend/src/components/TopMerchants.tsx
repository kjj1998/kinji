import { Avatar, Card, Group, Stack, Text, Title } from "@mantine/core";

export interface MerchantItem {
	name: string;
	amount: number;
	category: string;
}

export interface TopMerchantsProps {
	merchants: MerchantItem[];
	currency?: string;
}

export function TopMerchants({
	merchants,
	currency = "$",
}: TopMerchantsProps) {
	const sorted = [...merchants].sort((a, b) => b.amount - a.amount);

	return (
		<Card withBorder radius="md" p="md">
			<Title order={5} mb="md">
				Top Merchants
			</Title>
			<Stack gap="sm">
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
							{currency}
							{merchant.amount.toLocaleString("en-US", {
								minimumFractionDigits: 2,
								maximumFractionDigits: 2,
							})}
						</Text>
					</Group>
				))}
			</Stack>
		</Card>
	);
}
