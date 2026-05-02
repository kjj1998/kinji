import {
	Box,
	Button,
	Group,
	NumberInput,
	Select,
	SimpleGrid,
	TextInput,
} from "@mantine/core";
import { DateInput } from "@mantine/dates";
import { useForm } from "@mantine/form";

import type { Transaction } from "../types";

export interface AddTransactionFormProps {
	categories: string[];
	onAdd: (transaction: Transaction) => void;
}

export function AddTransactionForm({
	categories,
	onAdd,
}: AddTransactionFormProps) {
	const form = useForm({
		initialValues: {
			date: "",
			merchant: "",
			category: "",
			amount: 0,
			notes: "",
			split: null,
		},
		validate: {
			date: (v) => (v ? null : "Required"),
			merchant: (v) => (v.trim() ? null : "Required"),
			category: (v) => (v ? null : "Required"),
			amount: (v) => (v !== 0 ? null : "Required"),
			split: (v, values) =>
				v == null || Math.abs(v) <= Math.abs(values.amount)
					? null
					: "Cannot exceed amount",
		},
	});

	return (
		<Box mb="sm">
			<form
				onSubmit={form.onSubmit((values) => {
					const newTx: Transaction = {
						id: Date.now(),
						date: values.date,
						merchant: values.merchant,
						category: values.category,
						amount: values.amount,
						notes: values.notes,
						split: values.split ?? null,
					};
					onAdd(newTx);
					form.reset();
				})}
			>
				<SimpleGrid cols={3}>
					<DateInput label="Date" {...form.getInputProps("date")} />
					<TextInput label="Merchant" {...form.getInputProps("merchant")} />
					<Select
						label="Category"
						data={categories}
						{...form.getInputProps("category")}
					/>
					<NumberInput label="Amount" {...form.getInputProps("amount")} />
					<TextInput label="Notes" {...form.getInputProps("notes")} />
					<NumberInput
						label="Split (optional)"
						{...form.getInputProps("split")}
					/>
				</SimpleGrid>
				<Group justify="flex-end" mt="sm">
					<Button type="submit">Add</Button>
				</Group>
			</form>
		</Box>
	);
}
